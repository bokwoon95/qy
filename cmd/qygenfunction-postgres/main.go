package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	_ "github.com/lib/pq"
)

func main() {
	currdir, err := os.Getwd()
	if err != nil {
		printAndExit(err)
	}
	var (
		databaseFlag  = flag.String("database", "", "(required) Database URL")
		directoryFlag = flag.String("directory", filepath.Join(currdir, "tables"), "(optional) Directory to place the generated file. Can be absolute or relative filepath")
		dryrunFlag    = flag.Bool("dryrun", false, "(optional) Print the list of tables to be generated without generating the file")
		fileFlag      = flag.String("file", "functions.go", "(optional) Name of the file to be generated. If file already exists, -overwrite flag must be specified to overwrite the file")
		overwriteFlag = flag.Bool("overwrite", false, "(optional) Overwrite any files that already exist")
		packageFlag   = flag.String("package", "tables", "(optional) Package name of the file to be generated")
		schemasFlag   = flag.String("schema", "", "(required) A comma separated list of schemas that you want to generate tables for. Please don't include any spaces")
	)
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	if len(os.Args[1:]) == 0 {
		flag.PrintDefaults()
		return
	}
	config, err := prepConfig(*databaseFlag, *schemasFlag, *directoryFlag, *fileFlag, *packageFlag, *dryrunFlag, *overwriteFlag)
	if err != nil {
		printAndExit(err)
	}
	db, err := sql.Open("postgres", config.databaseURL)
	if err != nil {
		printAndExit(err)
	}
	err = db.Ping()
	if err != nil {
		printAndExit("Could not ping the database, is the database reachable via " + config.databaseURL + "? " + err.Error())
	}

	//--------------------------------------------------------------------------------//
	tables, err := processFunctions(getFunctions(db, config.databaseURL, config.schemas))
	if err != nil {
		printAndExit(err)
	}
	if config.dryrun {
		// Here is where you can print the list of filtered tables before you
		// actually write it out to a file.
		for _, table := range tables {
			fmt.Println(table)
		}
		return
	}
	//--------------------------------------------------------------------------------//

	err = writeFunctionsToFile(tables, config.directory, config.file, config.packageName)
	if err != nil {
		printAndExit(err)
	}
	fmt.Println("Result:      ", strconv.Itoa(len(tables)), "functions written into", filepath.Join(config.directory, config.file))
}

type Function struct {
	Schema       string
	Name         string
	RawResults   string
	RawArguments string
	StructName   string
	Constructor  string
	Results      []Field
	Arguments    []Field
}

type Field struct {
	RawField    string
	Name        string
	FieldType   string
	GoType      string
	Constructor string
}

func getFunctions(db *sql.DB, databaseURL string, schemas []string) ([]Function, error) {
	var functions []Function
	query := replacePlaceholders(
		"SELECT n.nspname, p.proname" +
			", pg_catalog.pg_get_function_result(p.oid) AS result" +
			", pg_catalog.pg_get_function_identity_arguments(p.oid) as arguments" +
			" FROM pg_catalog.pg_proc AS p" +
			" LEFT JOIN pg_catalog.pg_namespace AS n ON n.oid = p.pronamespace" +
			" WHERE n.nspname IN (?" + strings.Repeat(", ?", len(schemas)-1) + ") AND p.prokind = 'f'",
	)
	args := make([]interface{}, len(schemas))
	for i := range schemas {
		args[i] = schemas[i]
	}
	rows, err := db.Query(query, args...)
	fmt.Println("Query:       ", query, args)
	if err != nil {
		return functions, err
	}
	defer rows.Close()
NEXT_ROW:
	for rows.Next() {
		var functionSchema, functionName, functionResults, functionArguments string
		err := rows.Scan(&functionSchema, &functionName, &functionResults, &functionArguments)
		if err != nil {
			return functions, err
		}
		function := Function{
			Schema:       functionSchema,
			Name:         functionName,
			RawResults:   functionResults,
			RawArguments: functionArguments,
		}
		if functionSchema == "public" {
			function.StructName = "FUNCTION_" + strings.ToUpper(functionName)
			function.Constructor = strings.ToUpper(functionName)
		} else {
			function.StructName = "FUNCTION_" + strings.ToUpper(functionSchema+"__"+functionName)
			function.Constructor = strings.ToUpper(functionSchema + "__" + functionName)
		}
		{ // function.Arguments
			switch {
			case functionArguments == "":
				// Do nothing
			default:
				rawFields := strings.Split(functionArguments, ",")
				for i := range rawFields {
					field := extractNameAndType(strings.TrimSpace(rawFields[i]))
					if field.RawField != "void" && field.FieldType == "" {
						fmt.Printf("Skipping %s.%s: unable to process argument '%s'\n", function.Schema, function.Name, field.RawField)
						continue NEXT_ROW
					}
					if field.Name == "" {
						field.Name = "_arg" + strconv.Itoa(i+1)
					}
					function.Arguments = append(function.Arguments, field)
				}
			}
		}
		{ // function.Results
			var fields []Field
			switch {
			case function.RawResults == "void":
				// Do nothing
			case strings.HasPrefix(functionResults, "TABLE(") && strings.HasSuffix(functionResults, ")"):
				functionResults = functionResults[6 : len(functionResults)-1]
				rawFields := strings.Split(functionResults, ",")
				for i := range rawFields {
					field := extractNameAndType(strings.TrimSpace(rawFields[i]))
					fields = append(fields, field)
				}
			default:
				if strings.HasPrefix(functionResults, "SETOF ") {
					functionResults = functionResults[6:]
				}
				field := extractNameAndType(strings.TrimSpace(functionResults))
				field.Name = "Result"
				fields = append(fields, field)
			}
			for i := range fields {
				switch {
				case fields[i].FieldType == "":
					fmt.Printf("Skipping %s.%s: unable to process return value '%s'\n", function.Schema, function.Name, fields[i].RawField)
					continue NEXT_ROW
				default:
					if fields[i].Name == "" {
						fields[i].Name = "_res" + strconv.Itoa(i)
					}
					function.Results = append(function.Results, fields[i])
				}
			}
		}
		functions = append(functions, function)
	}
	return functions, nil
}

func extractNameAndType(rawField string) Field {
	var field Field
	field.RawField = rawField
	if matches := regexp.
		MustCompile(`boolean` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(rawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Boolean
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeBoolSlice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeBoolean
			field.GoType = GoTypeBool
			field.Constructor = FieldConstructorBoolean
		}

	} else if matches := regexp.
		MustCompile(`json` + `(?:b)?` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(rawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// JSON
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeInterface
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeJSON
			field.GoType = GoTypeInterface
			field.Constructor = FieldConstructorJSON
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `smallint` +
			`|` + `oid` +
			`|` + `integer` +
			`|` + `bigint` +
			`|` + `smallserial` +
			`|` + `serial` +
			`|` + `bigserial` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(rawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Integer
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeIntSlice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeNumber
			field.GoType = GoTypeInt
			field.Constructor = FieldConstructorNumber
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `decimal` +
			`|` + `numeric` +
			`|` + `real` +
			`|` + `double precision` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(rawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Float
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeFloat64Slice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeNumber
			field.GoType = GoTypeFloat64
			field.Constructor = FieldConstructorNumber
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `text` +
			`|` + `name` +
			`|` + `char` + `(?:\(\d+\))?` +
			`|` + `character` + `(?:\(\d+\))?` +
			`|` + `varchar` + `(?:\(\d+\))?` +
			`|` + `character varying` + `(?:\(\d+\))?` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(rawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// String
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			field.FieldType = FieldTypeArray
			field.GoType = GoTypeStringSlice
			field.Constructor = FieldConstructorArray
		} else {
			field.FieldType = FieldTypeString
			field.GoType = GoTypeString
			field.Constructor = FieldConstructorString
		}

	} else if matches := regexp.
		MustCompile(`(?:` + `date` +
			`|` + `(?:time|timestamp)` +
			`(?: \(\d+\))?` +
			`(?: without time zone| with time zone)?` + `)` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(rawField); len(matches) == 2 {
		// fmt.Println(rawField, matches)
		// Time
		field.Name = strings.TrimSpace(rawField[:len(rawField)-len(matches[0])])
		if matches[1] == "[]" {
			// Do nothing
		} else {
			field.FieldType = FieldTypeTime
			field.GoType = GoTypeTime
			field.Constructor = FieldConstructorTime
		}

	} else if matches := regexp.
		MustCompile(`bytea` +
			`(\[\])?` +
			`$`).
		FindStringSubmatch(rawField); len(matches) == 2 {
		if matches[1] == "[]" {
			// Do nothing
		} else {
			field.FieldType = FieldTypeBinary
			field.GoType = GoTypeBytes
			field.Constructor = FieldConstructorBinary
		}
	}

	return field
}

const (
	GoTypeInterface    = "interface{}"
	GoTypeBool         = "bool"
	GoTypeInt          = "int"
	GoTypeFloat64      = "float64"
	GoTypeString       = "string"
	GoTypeTime         = "time.Time"
	GoTypeBoolSlice    = "[]bool"
	GoTypeIntSlice     = "[]int"
	GoTypeFloat64Slice = "[]float64"
	GoTypeStringSlice  = "[]string"
	GoTypeTimeSlice    = "[]time.Time"
	GoTypeBytes        = "[]byte"

	FieldTypeBoolean = "qx.BooleanField"
	FieldTypeJSON    = "qx.JSONField"
	FieldTypeNumber  = "qx.NumberField"
	FieldTypeString  = "qx.StringField"
	FieldTypeTime    = "qx.TimeField"
	FieldTypeEnum    = "qx.EnumField"
	FieldTypeArray   = "qx.ArrayField"
	FieldTypeBinary  = "qx.BinaryField"

	FieldConstructorBoolean = "qx.NewBooleanField"
	FieldConstructorJSON    = "qx.NewJSONField"
	FieldConstructorNumber  = "qx.NewNumberField"
	FieldConstructorString  = "qx.NewStringField"
	FieldConstructorTime    = "qx.NewTimeField"
	FieldConstructorEnum    = "qx.NewEnumField"
	FieldConstructorArray   = "qx.NewArrayField"
	FieldConstructorBinary  = "qx.NewBinaryField"
)

func processFunctions(inputFunctions []Function, err error) ([]Function, error) {
	if err != nil {
		return inputFunctions, err
	}
	var outputFunctions []Function
	for _, function := range inputFunctions {
		_ = function
	}
	outputFunctions = inputFunctions
	return outputFunctions, nil
}

type FileData struct {
	PackageName string
	Imports     []string
	Functions   []Function
}

var Imports = []string{
	"github.com/bokwoon95/qy/qx",
}

var qygenfunctionTemplate = `// Code generated by qygenfunction-postgres; DO NOT EDIT.
package {{$.PackageName}}

import (
	{{- range $_, $import := $.Imports}}
	"{{$import}}"
	{{- end}}
)
{{- range $_, $function := $.Functions}}
{{template "function_struct_definition" $function}}
{{template "function_constructor" $function}}
{{template "function_as" $function}}
{{- end}}

{{- define "function_struct_definition"}}
{{- with $function := .}}
// {{$function.StructName | uppercase }} references the {{$function.Schema}}.{{$function.Name}} function
type {{$function.StructName | uppercase}} struct {
	*qx.FunctionInfo
	{{- range $_, $field := $function.Results}}
	{{$field.Name | uppercase | trimUnderscorePrefix}} {{$field.FieldType}}
	{{- end}}
}
{{- end}}
{{- end}}

{{- define "function_constructor"}}
{{- with $function := .}}
// {{$function.Constructor | trimUnderscorePrefix}} creates an instance of the {{$function.Schema}}.{{$function.Name}} function
func {{$function.Constructor | trimUnderscorePrefix}}(
	{{- range $i, $arg := $function.Arguments}}
	{{$arg.Name}} {{$arg.GoType}},
	{{- end}}
	) {{$function.StructName}} {
	return {{$function.Constructor}}_({{range $i, $arg := $function.Arguments}}{{if not $i}}{{$arg.Name}}{{else}}, {{$arg.Name}}{{end}}{{end}})
}

// {{$function.Constructor | trimUnderscorePrefix}}_ creates an instance of the {{$function.Schema}}.{{$function.Name}} function
func {{$function.Constructor | trimUnderscorePrefix}}_(
	{{- range $i, $arg := $function.Arguments}}
	{{$arg.Name}} interface{},
	{{- end}}
	) {{$function.StructName}} {
	f := {{$function.StructName}}{FunctionInfo: &qx.FunctionInfo{
		Schema: "{{$function.Schema}}",
		Name: "{{$function.Name}}",
		Arguments: []interface{}{{"{"}}{{range $i, $arg := $function.Arguments}}{{if not $i}}{{$arg.Name}}{{else}}, {{$arg.Name}}{{end}}{{end}}{{"}"}},
	},}
	{{- range $_, $field := $function.Results}}
	f.{{$field.Name | uppercase | trimUnderscorePrefix}} = {{$field.Constructor}}("{{$field.Name}}", f.FunctionInfo)
	{{- end}}
	return f
}
{{- end}}
{{- end}}

{{- define "function_as"}}
{{- with $function := .}}
func (f {{$function.StructName}}) As(alias string) {{$function.StructName}} {
	f.FunctionInfo.Alias = alias
	return f
}
{{- end}}
{{- end}}`

// writeFunctionsToFile will write the functions into a file specified by
// filepath.Join(directory, file).
func writeFunctionsToFile(functions []Function, directory, file, packageName string) error {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return fmt.Errorf("Could not create directory %s: %w", directory, err)
	}
	filename := filepath.Join(directory, file)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	t, err := template.New("").Funcs(template.FuncMap{
		"uppercase":            strings.ToUpper,
		"trimUnderscorePrefix": func(s string) string { return strings.TrimPrefix(s, "_") },
	}).Parse(qygenfunctionTemplate)
	if err != nil {
		return err
	}
	data := FileData{
		PackageName: packageName,
		Imports:     Imports,
		Functions:   functions,
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	if _, err := exec.LookPath("goimports"); err == nil {
		_ = exec.Command("goimports", "-w", filename).Run()
	} else if _, err := exec.LookPath("gofmt"); err == nil {
		_ = exec.Command("gofmt", "-w", filename).Run()
	}
	return nil
}

type Config struct {
	databaseURL string
	directory   string
	dryrun      bool
	file        string
	overwrite   bool
	packageName string
	schemas     []string
}

// prepConfig will process the incoming data and initialize the config object
// accordingly
func prepConfig(database, schemas, directory, file, packageName string, dryrun, overwrite bool) (cfg Config, err error) {
	// databaseURL
	cfg.databaseURL = database
	fmt.Println("Database URL:", cfg.databaseURL)
	if cfg.databaseURL == "" {
		return cfg, fmt.Errorf("Database URL is either empty or not passed in. You need to specify a database URL with the -database option.")
	}

	// schemas
	if schemas == "" {
		return cfg, fmt.Errorf("A database schema needs to be specified")
	}
	cfg.schemas = strings.FieldsFunc(schemas, func(r rune) bool { return r == ',' || unicode.IsSpace(r) })
	fmt.Println("Schemas:     ", cfg.schemas)

	// directory
	cfg.directory = directory
	fmt.Println("Directory:   ", cfg.directory)
	if cfg.directory == "" {
		return cfg, fmt.Errorf("-directory was not specified. You need to provide a directory to place the generated file in")
	}
	cfg.directory, err = filepath.Abs(cfg.directory)
	if err != nil {
		return cfg, err
	}

	// file
	cfg.file = file
	fmt.Println("File:        ", cfg.file)
	if cfg.file == "" {
		return cfg, fmt.Errorf("-file was not specified. You need to provide a file name (e.g. tables.go) for the generated file")
	}
	if !strings.HasSuffix(cfg.file, ".go") {
		cfg.file = cfg.file + ".go"
	}
	asboluteFilePath := filepath.Join(directory, file)
	if _, err := os.Stat(asboluteFilePath); err == nil && !overwrite {
		return cfg, fmt.Errorf("Specified file %s already exists. If you wish to overwrite it, provide the -overwrite flag", asboluteFilePath)
	}

	// packageName
	cfg.packageName = packageName
	fmt.Println("Package Name:", cfg.packageName)
	if cfg.packageName == "" {
		return cfg, fmt.Errorf("-package name was not provided. You need to provide the package name for the generated file")
	}

	// dryrun
	cfg.dryrun = dryrun
	return cfg, nil
}

/* Utility functions */

func printAndExit(v ...interface{}) {
	fmt.Println("======================================== ERROR! ========================================")
	log.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// replacePlaceholders will replace all ? placeholders in a query string to the
// postgres-valid incrementing placeholders $1, $2, $3 etc.
func replacePlaceholders(query string) string {
	buf := &strings.Builder{}
	i := 0
	for pos := strings.Index(query, "?"); pos >= 0; pos = strings.Index(query, "?") {
		i++
		buf.WriteString(query[:pos] + "$" + strconv.Itoa(i))
		query = query[pos+1:]
	}
	buf.WriteString(query)
	return buf.String()
}

// String implements fmt.Stringer for type Table, allowing you to call
// fmt.Println(table) and have it formatted accordingly.
func (f Function) String() string {
	var output string
	if f.Constructor != "" && f.StructName != "" {
		output += fmt.Sprintf("%s.%s => func %s() %s\n", f.Schema, f.Name, f.Constructor, f.StructName)
	} else {
		output += fmt.Sprintf("%s.%s\n", f.Schema, f.Name)
	}
	output += fmt.Sprintf("    Arguments\n")
	for _, field := range f.Arguments {
		output += fmt.Sprintf("        %#v\n", field)
		// if field.Constructor != "" && field.GoType != "" {
		// 	output += fmt.Sprintf("        %s: %s\n", field.Name, field.GoType)
		// } else {
		// 	output += fmt.Sprintf("        %s\n", field.RawField)
		// }
	}
	output += fmt.Sprintf("    Results\n")
	for _, field := range f.Results {
		output += fmt.Sprintf("        %#v\n", field)
		// if field.Constructor != "" && field.FieldType != "" {
		// 	output += fmt.Sprintf("        %s: %s\n", field.Name, field.FieldType)
		// } else {
		// 	output += fmt.Sprintf("        %s\n", field.RawField)
		// }
	}
	return output
}
