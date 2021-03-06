// Code generated by qygen-postgres; DO NOT EDIT.
package qx // modified to break import cycle

type TABLE_ACTOR struct {
	*TableInfo
	ACTOR_ID    NumberField
	FIRST_NAME  StringField
	LAST_NAME   StringField
	LAST_UPDATE TimeField
}

func ACTOR() TABLE_ACTOR {
	tbl := TABLE_ACTOR{TableInfo: NewTableInfo("public", "actor")}
	tbl.ACTOR_ID = NewNumberField("actor_id", tbl.TableInfo)
	tbl.FIRST_NAME = NewStringField("first_name", tbl.TableInfo)
	tbl.LAST_NAME = NewStringField("last_name", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_ACTOR) As(alias string) TABLE_ACTOR {
	tbl2 := ACTOR()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_ADDRESS struct {
	*TableInfo
	ADDRESS     StringField
	ADDRESS2    StringField
	ADDRESS_ID  NumberField
	CITY_ID     NumberField
	DISTRICT    StringField
	LAST_UPDATE TimeField
	PHONE       StringField
	POSTAL_CODE StringField
}

func ADDRESS() TABLE_ADDRESS {
	tbl := TABLE_ADDRESS{TableInfo: NewTableInfo("public", "address")}
	tbl.ADDRESS = NewStringField("address", tbl.TableInfo)
	tbl.ADDRESS2 = NewStringField("address2", tbl.TableInfo)
	tbl.ADDRESS_ID = NewNumberField("address_id", tbl.TableInfo)
	tbl.CITY_ID = NewNumberField("city_id", tbl.TableInfo)
	tbl.DISTRICT = NewStringField("district", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.PHONE = NewStringField("phone", tbl.TableInfo)
	tbl.POSTAL_CODE = NewStringField("postal_code", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_ADDRESS) As(alias string) TABLE_ADDRESS {
	tbl2 := ADDRESS()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_CATEGORY struct {
	*TableInfo
	CATEGORY_ID NumberField
	LAST_UPDATE TimeField
	NAME        StringField
}

func CATEGORY() TABLE_CATEGORY {
	tbl := TABLE_CATEGORY{TableInfo: NewTableInfo("public", "category")}
	tbl.CATEGORY_ID = NewNumberField("category_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.NAME = NewStringField("name", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_CATEGORY) As(alias string) TABLE_CATEGORY {
	tbl2 := CATEGORY()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_CITY struct {
	*TableInfo
	CITY        StringField
	CITY_ID     NumberField
	COUNTRY_ID  NumberField
	LAST_UPDATE TimeField
}

func CITY() TABLE_CITY {
	tbl := TABLE_CITY{TableInfo: NewTableInfo("public", "city")}
	tbl.CITY = NewStringField("city", tbl.TableInfo)
	tbl.CITY_ID = NewNumberField("city_id", tbl.TableInfo)
	tbl.COUNTRY_ID = NewNumberField("country_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_CITY) As(alias string) TABLE_CITY {
	tbl2 := CITY()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_COUNTRY struct {
	*TableInfo
	COUNTRY     StringField
	COUNTRY_ID  NumberField
	LAST_UPDATE TimeField
}

func COUNTRY() TABLE_COUNTRY {
	tbl := TABLE_COUNTRY{TableInfo: NewTableInfo("public", "country")}
	tbl.COUNTRY = NewStringField("country", tbl.TableInfo)
	tbl.COUNTRY_ID = NewNumberField("country_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_COUNTRY) As(alias string) TABLE_COUNTRY {
	tbl2 := COUNTRY()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_CUSTOMER struct {
	*TableInfo
	ACTIVE      NumberField
	ACTIVEBOOL  BooleanField
	ADDRESS_ID  NumberField
	CREATE_DATE TimeField
	CUSTOMER_ID NumberField
	EMAIL       StringField
	FIRST_NAME  StringField
	LAST_NAME   StringField
	LAST_UPDATE TimeField
	STORE_ID    NumberField
}

func CUSTOMER() TABLE_CUSTOMER {
	tbl := TABLE_CUSTOMER{TableInfo: NewTableInfo("public", "customer")}
	tbl.ACTIVE = NewNumberField("active", tbl.TableInfo)
	tbl.ACTIVEBOOL = NewBooleanField("activebool", tbl.TableInfo)
	tbl.ADDRESS_ID = NewNumberField("address_id", tbl.TableInfo)
	tbl.CREATE_DATE = NewTimeField("create_date", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.EMAIL = NewStringField("email", tbl.TableInfo)
	tbl.FIRST_NAME = NewStringField("first_name", tbl.TableInfo)
	tbl.LAST_NAME = NewStringField("last_name", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.STORE_ID = NewNumberField("store_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_CUSTOMER) As(alias string) TABLE_CUSTOMER {
	tbl2 := CUSTOMER()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_FILM struct {
	*TableInfo
	DESCRIPTION          StringField
	FILM_ID              NumberField
	LANGUAGE_ID          NumberField
	LAST_UPDATE          TimeField
	LENGTH               NumberField
	ORIGINAL_LANGUAGE_ID NumberField
	RATING               EnumField
	RELEASE_YEAR         NumberField
	RENTAL_DURATION      NumberField
	RENTAL_RATE          NumberField
	REPLACEMENT_COST     NumberField
	SPECIAL_FEATURES     ArrayField
	TITLE                StringField
}

func FILM() TABLE_FILM {
	tbl := TABLE_FILM{TableInfo: NewTableInfo("public", "film")}
	tbl.DESCRIPTION = NewStringField("description", tbl.TableInfo)
	tbl.FILM_ID = NewNumberField("film_id", tbl.TableInfo)
	tbl.LANGUAGE_ID = NewNumberField("language_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.LENGTH = NewNumberField("length", tbl.TableInfo)
	tbl.ORIGINAL_LANGUAGE_ID = NewNumberField("original_language_id", tbl.TableInfo)
	tbl.RATING = NewEnumField("rating", tbl.TableInfo)
	tbl.RELEASE_YEAR = NewNumberField("release_year", tbl.TableInfo)
	tbl.RENTAL_DURATION = NewNumberField("rental_duration", tbl.TableInfo)
	tbl.RENTAL_RATE = NewNumberField("rental_rate", tbl.TableInfo)
	tbl.REPLACEMENT_COST = NewNumberField("replacement_cost", tbl.TableInfo)
	tbl.SPECIAL_FEATURES = NewArrayField("special_features", tbl.TableInfo)
	tbl.TITLE = NewStringField("title", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_FILM) As(alias string) TABLE_FILM {
	tbl2 := FILM()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_FILM_ACTOR struct {
	*TableInfo
	ACTOR_ID    NumberField
	FILM_ID     NumberField
	LAST_UPDATE TimeField
}

func FILM_ACTOR() TABLE_FILM_ACTOR {
	tbl := TABLE_FILM_ACTOR{TableInfo: NewTableInfo("public", "film_actor")}
	tbl.ACTOR_ID = NewNumberField("actor_id", tbl.TableInfo)
	tbl.FILM_ID = NewNumberField("film_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_FILM_ACTOR) As(alias string) TABLE_FILM_ACTOR {
	tbl2 := FILM_ACTOR()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_FILM_CATEGORY struct {
	*TableInfo
	CATEGORY_ID NumberField
	FILM_ID     NumberField
	LAST_UPDATE TimeField
}

func FILM_CATEGORY() TABLE_FILM_CATEGORY {
	tbl := TABLE_FILM_CATEGORY{TableInfo: NewTableInfo("public", "film_category")}
	tbl.CATEGORY_ID = NewNumberField("category_id", tbl.TableInfo)
	tbl.FILM_ID = NewNumberField("film_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_FILM_CATEGORY) As(alias string) TABLE_FILM_CATEGORY {
	tbl2 := FILM_CATEGORY()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_INVENTORY struct {
	*TableInfo
	FILM_ID      NumberField
	INVENTORY_ID NumberField
	LAST_UPDATE  TimeField
	STORE_ID     NumberField
}

func INVENTORY() TABLE_INVENTORY {
	tbl := TABLE_INVENTORY{TableInfo: NewTableInfo("public", "inventory")}
	tbl.FILM_ID = NewNumberField("film_id", tbl.TableInfo)
	tbl.INVENTORY_ID = NewNumberField("inventory_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.STORE_ID = NewNumberField("store_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_INVENTORY) As(alias string) TABLE_INVENTORY {
	tbl2 := INVENTORY()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_LANGUAGE struct {
	*TableInfo
	LANGUAGE_ID NumberField
	LAST_UPDATE TimeField
	NAME        StringField
}

func LANGUAGE() TABLE_LANGUAGE {
	tbl := TABLE_LANGUAGE{TableInfo: NewTableInfo("public", "language")}
	tbl.LANGUAGE_ID = NewNumberField("language_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.NAME = NewStringField("name", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_LANGUAGE) As(alias string) TABLE_LANGUAGE {
	tbl2 := LANGUAGE()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_PAYMENT struct {
	*TableInfo
	AMOUNT       NumberField
	CUSTOMER_ID  NumberField
	PAYMENT_DATE TimeField
	PAYMENT_ID   NumberField
	RENTAL_ID    NumberField
	STAFF_ID     NumberField
}

func PAYMENT() TABLE_PAYMENT {
	tbl := TABLE_PAYMENT{TableInfo: NewTableInfo("public", "payment")}
	tbl.AMOUNT = NewNumberField("amount", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.PAYMENT_DATE = NewTimeField("payment_date", tbl.TableInfo)
	tbl.PAYMENT_ID = NewNumberField("payment_id", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_PAYMENT) As(alias string) TABLE_PAYMENT {
	tbl2 := PAYMENT()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_PAYMENT_P2007_01 struct {
	*TableInfo
	AMOUNT       NumberField
	CUSTOMER_ID  NumberField
	PAYMENT_DATE TimeField
	PAYMENT_ID   NumberField
	RENTAL_ID    NumberField
	STAFF_ID     NumberField
}

func PAYMENT_P2007_01() TABLE_PAYMENT_P2007_01 {
	tbl := TABLE_PAYMENT_P2007_01{TableInfo: NewTableInfo("public", "payment_p2007_01")}
	tbl.AMOUNT = NewNumberField("amount", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.PAYMENT_DATE = NewTimeField("payment_date", tbl.TableInfo)
	tbl.PAYMENT_ID = NewNumberField("payment_id", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_PAYMENT_P2007_01) As(alias string) TABLE_PAYMENT_P2007_01 {
	tbl2 := PAYMENT_P2007_01()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_PAYMENT_P2007_02 struct {
	*TableInfo
	AMOUNT       NumberField
	CUSTOMER_ID  NumberField
	PAYMENT_DATE TimeField
	PAYMENT_ID   NumberField
	RENTAL_ID    NumberField
	STAFF_ID     NumberField
}

func PAYMENT_P2007_02() TABLE_PAYMENT_P2007_02 {
	tbl := TABLE_PAYMENT_P2007_02{TableInfo: NewTableInfo("public", "payment_p2007_02")}
	tbl.AMOUNT = NewNumberField("amount", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.PAYMENT_DATE = NewTimeField("payment_date", tbl.TableInfo)
	tbl.PAYMENT_ID = NewNumberField("payment_id", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_PAYMENT_P2007_02) As(alias string) TABLE_PAYMENT_P2007_02 {
	tbl2 := PAYMENT_P2007_02()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_PAYMENT_P2007_03 struct {
	*TableInfo
	AMOUNT       NumberField
	CUSTOMER_ID  NumberField
	PAYMENT_DATE TimeField
	PAYMENT_ID   NumberField
	RENTAL_ID    NumberField
	STAFF_ID     NumberField
}

func PAYMENT_P2007_03() TABLE_PAYMENT_P2007_03 {
	tbl := TABLE_PAYMENT_P2007_03{TableInfo: NewTableInfo("public", "payment_p2007_03")}
	tbl.AMOUNT = NewNumberField("amount", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.PAYMENT_DATE = NewTimeField("payment_date", tbl.TableInfo)
	tbl.PAYMENT_ID = NewNumberField("payment_id", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_PAYMENT_P2007_03) As(alias string) TABLE_PAYMENT_P2007_03 {
	tbl2 := PAYMENT_P2007_03()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_PAYMENT_P2007_04 struct {
	*TableInfo
	AMOUNT       NumberField
	CUSTOMER_ID  NumberField
	PAYMENT_DATE TimeField
	PAYMENT_ID   NumberField
	RENTAL_ID    NumberField
	STAFF_ID     NumberField
}

func PAYMENT_P2007_04() TABLE_PAYMENT_P2007_04 {
	tbl := TABLE_PAYMENT_P2007_04{TableInfo: NewTableInfo("public", "payment_p2007_04")}
	tbl.AMOUNT = NewNumberField("amount", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.PAYMENT_DATE = NewTimeField("payment_date", tbl.TableInfo)
	tbl.PAYMENT_ID = NewNumberField("payment_id", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_PAYMENT_P2007_04) As(alias string) TABLE_PAYMENT_P2007_04 {
	tbl2 := PAYMENT_P2007_04()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_PAYMENT_P2007_05 struct {
	*TableInfo
	AMOUNT       NumberField
	CUSTOMER_ID  NumberField
	PAYMENT_DATE TimeField
	PAYMENT_ID   NumberField
	RENTAL_ID    NumberField
	STAFF_ID     NumberField
}

func PAYMENT_P2007_05() TABLE_PAYMENT_P2007_05 {
	tbl := TABLE_PAYMENT_P2007_05{TableInfo: NewTableInfo("public", "payment_p2007_05")}
	tbl.AMOUNT = NewNumberField("amount", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.PAYMENT_DATE = NewTimeField("payment_date", tbl.TableInfo)
	tbl.PAYMENT_ID = NewNumberField("payment_id", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_PAYMENT_P2007_05) As(alias string) TABLE_PAYMENT_P2007_05 {
	tbl2 := PAYMENT_P2007_05()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_PAYMENT_P2007_06 struct {
	*TableInfo
	AMOUNT       NumberField
	CUSTOMER_ID  NumberField
	PAYMENT_DATE TimeField
	PAYMENT_ID   NumberField
	RENTAL_ID    NumberField
	STAFF_ID     NumberField
}

func PAYMENT_P2007_06() TABLE_PAYMENT_P2007_06 {
	tbl := TABLE_PAYMENT_P2007_06{TableInfo: NewTableInfo("public", "payment_p2007_06")}
	tbl.AMOUNT = NewNumberField("amount", tbl.TableInfo)
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.PAYMENT_DATE = NewTimeField("payment_date", tbl.TableInfo)
	tbl.PAYMENT_ID = NewNumberField("payment_id", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_PAYMENT_P2007_06) As(alias string) TABLE_PAYMENT_P2007_06 {
	tbl2 := PAYMENT_P2007_06()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_RENTAL struct {
	*TableInfo
	CUSTOMER_ID  NumberField
	INVENTORY_ID NumberField
	LAST_UPDATE  TimeField
	RENTAL_DATE  TimeField
	RENTAL_ID    NumberField
	RETURN_DATE  TimeField
	STAFF_ID     NumberField
}

func RENTAL() TABLE_RENTAL {
	tbl := TABLE_RENTAL{TableInfo: NewTableInfo("public", "rental")}
	tbl.CUSTOMER_ID = NewNumberField("customer_id", tbl.TableInfo)
	tbl.INVENTORY_ID = NewNumberField("inventory_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.RENTAL_DATE = NewTimeField("rental_date", tbl.TableInfo)
	tbl.RENTAL_ID = NewNumberField("rental_id", tbl.TableInfo)
	tbl.RETURN_DATE = NewTimeField("return_date", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_RENTAL) As(alias string) TABLE_RENTAL {
	tbl2 := RENTAL()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_STAFF struct {
	*TableInfo
	ACTIVE      BooleanField
	ADDRESS_ID  NumberField
	EMAIL       StringField
	FIRST_NAME  StringField
	LAST_NAME   StringField
	LAST_UPDATE TimeField
	PASSWORD    StringField
	STAFF_ID    NumberField
	STORE_ID    NumberField
	USERNAME    StringField
}

func STAFF() TABLE_STAFF {
	tbl := TABLE_STAFF{TableInfo: NewTableInfo("public", "staff")}
	tbl.ACTIVE = NewBooleanField("active", tbl.TableInfo)
	tbl.ADDRESS_ID = NewNumberField("address_id", tbl.TableInfo)
	tbl.EMAIL = NewStringField("email", tbl.TableInfo)
	tbl.FIRST_NAME = NewStringField("first_name", tbl.TableInfo)
	tbl.LAST_NAME = NewStringField("last_name", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.PASSWORD = NewStringField("password", tbl.TableInfo)
	tbl.STAFF_ID = NewNumberField("staff_id", tbl.TableInfo)
	tbl.STORE_ID = NewNumberField("store_id", tbl.TableInfo)
	tbl.USERNAME = NewStringField("username", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_STAFF) As(alias string) TABLE_STAFF {
	tbl2 := STAFF()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type TABLE_STORE struct {
	*TableInfo
	ADDRESS_ID       NumberField
	LAST_UPDATE      TimeField
	MANAGER_STAFF_ID NumberField
	STORE_ID         NumberField
}

func STORE() TABLE_STORE {
	tbl := TABLE_STORE{TableInfo: NewTableInfo("public", "store")}
	tbl.ADDRESS_ID = NewNumberField("address_id", tbl.TableInfo)
	tbl.LAST_UPDATE = NewTimeField("last_update", tbl.TableInfo)
	tbl.MANAGER_STAFF_ID = NewNumberField("manager_staff_id", tbl.TableInfo)
	tbl.STORE_ID = NewNumberField("store_id", tbl.TableInfo)
	return tbl
}

func (tbl TABLE_STORE) As(alias string) TABLE_STORE {
	tbl2 := STORE()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type VIEW_ACTOR_INFO struct {
	*TableInfo
	ACTOR_ID   NumberField
	FILM_INFO  StringField
	FIRST_NAME StringField
	LAST_NAME  StringField
}

func ACTOR_INFO() VIEW_ACTOR_INFO {
	tbl := VIEW_ACTOR_INFO{TableInfo: NewTableInfo("public", "actor_info")}
	tbl.ACTOR_ID = NewNumberField("actor_id", tbl.TableInfo)
	tbl.FILM_INFO = NewStringField("film_info", tbl.TableInfo)
	tbl.FIRST_NAME = NewStringField("first_name", tbl.TableInfo)
	tbl.LAST_NAME = NewStringField("last_name", tbl.TableInfo)
	return tbl
}

func (tbl VIEW_ACTOR_INFO) As(alias string) VIEW_ACTOR_INFO {
	tbl2 := ACTOR_INFO()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type VIEW_CUSTOMER_LIST struct {
	*TableInfo
	ADDRESS  StringField
	CITY     StringField
	COUNTRY  StringField
	ID       NumberField
	NAME     StringField
	NOTES    StringField
	PHONE    StringField
	SID      NumberField
	ZIP_CODE StringField
}

func CUSTOMER_LIST() VIEW_CUSTOMER_LIST {
	tbl := VIEW_CUSTOMER_LIST{TableInfo: NewTableInfo("public", "customer_list")}
	tbl.ADDRESS = NewStringField("address", tbl.TableInfo)
	tbl.CITY = NewStringField("city", tbl.TableInfo)
	tbl.COUNTRY = NewStringField("country", tbl.TableInfo)
	tbl.ID = NewNumberField("id", tbl.TableInfo)
	tbl.NAME = NewStringField("name", tbl.TableInfo)
	tbl.NOTES = NewStringField("notes", tbl.TableInfo)
	tbl.PHONE = NewStringField("phone", tbl.TableInfo)
	tbl.SID = NewNumberField("sid", tbl.TableInfo)
	tbl.ZIP_CODE = NewStringField("zip_code", tbl.TableInfo)
	return tbl
}

func (tbl VIEW_CUSTOMER_LIST) As(alias string) VIEW_CUSTOMER_LIST {
	tbl2 := CUSTOMER_LIST()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type VIEW_FILM_LIST struct {
	*TableInfo
	ACTORS      StringField
	CATEGORY    StringField
	DESCRIPTION StringField
	FID         NumberField
	LENGTH      NumberField
	PRICE       NumberField
	RATING      EnumField
	TITLE       StringField
}

func FILM_LIST() VIEW_FILM_LIST {
	tbl := VIEW_FILM_LIST{TableInfo: NewTableInfo("public", "film_list")}
	tbl.ACTORS = NewStringField("actors", tbl.TableInfo)
	tbl.CATEGORY = NewStringField("category", tbl.TableInfo)
	tbl.DESCRIPTION = NewStringField("description", tbl.TableInfo)
	tbl.FID = NewNumberField("fid", tbl.TableInfo)
	tbl.LENGTH = NewNumberField("length", tbl.TableInfo)
	tbl.PRICE = NewNumberField("price", tbl.TableInfo)
	tbl.RATING = NewEnumField("rating", tbl.TableInfo)
	tbl.TITLE = NewStringField("title", tbl.TableInfo)
	return tbl
}

func (tbl VIEW_FILM_LIST) As(alias string) VIEW_FILM_LIST {
	tbl2 := FILM_LIST()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type VIEW_NICER_BUT_SLOWER_FILM_LIST struct {
	*TableInfo
	ACTORS      StringField
	CATEGORY    StringField
	DESCRIPTION StringField
	FID         NumberField
	LENGTH      NumberField
	PRICE       NumberField
	RATING      EnumField
	TITLE       StringField
}

func NICER_BUT_SLOWER_FILM_LIST() VIEW_NICER_BUT_SLOWER_FILM_LIST {
	tbl := VIEW_NICER_BUT_SLOWER_FILM_LIST{TableInfo: NewTableInfo("public", "nicer_but_slower_film_list")}
	tbl.ACTORS = NewStringField("actors", tbl.TableInfo)
	tbl.CATEGORY = NewStringField("category", tbl.TableInfo)
	tbl.DESCRIPTION = NewStringField("description", tbl.TableInfo)
	tbl.FID = NewNumberField("fid", tbl.TableInfo)
	tbl.LENGTH = NewNumberField("length", tbl.TableInfo)
	tbl.PRICE = NewNumberField("price", tbl.TableInfo)
	tbl.RATING = NewEnumField("rating", tbl.TableInfo)
	tbl.TITLE = NewStringField("title", tbl.TableInfo)
	return tbl
}

func (tbl VIEW_NICER_BUT_SLOWER_FILM_LIST) As(alias string) VIEW_NICER_BUT_SLOWER_FILM_LIST {
	tbl2 := NICER_BUT_SLOWER_FILM_LIST()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type VIEW_SALES_BY_FILM_CATEGORY struct {
	*TableInfo
	CATEGORY    StringField
	TOTAL_SALES NumberField
}

func SALES_BY_FILM_CATEGORY() VIEW_SALES_BY_FILM_CATEGORY {
	tbl := VIEW_SALES_BY_FILM_CATEGORY{TableInfo: NewTableInfo("public", "sales_by_film_category")}
	tbl.CATEGORY = NewStringField("category", tbl.TableInfo)
	tbl.TOTAL_SALES = NewNumberField("total_sales", tbl.TableInfo)
	return tbl
}

func (tbl VIEW_SALES_BY_FILM_CATEGORY) As(alias string) VIEW_SALES_BY_FILM_CATEGORY {
	tbl2 := SALES_BY_FILM_CATEGORY()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type VIEW_SALES_BY_STORE struct {
	*TableInfo
	MANAGER     StringField
	STORE       StringField
	TOTAL_SALES NumberField
}

func SALES_BY_STORE() VIEW_SALES_BY_STORE {
	tbl := VIEW_SALES_BY_STORE{TableInfo: NewTableInfo("public", "sales_by_store")}
	tbl.MANAGER = NewStringField("manager", tbl.TableInfo)
	tbl.STORE = NewStringField("store", tbl.TableInfo)
	tbl.TOTAL_SALES = NewNumberField("total_sales", tbl.TableInfo)
	return tbl
}

func (tbl VIEW_SALES_BY_STORE) As(alias string) VIEW_SALES_BY_STORE {
	tbl2 := SALES_BY_STORE()
	tbl2.TableInfo.Alias = alias
	return tbl2
}

type VIEW_STAFF_LIST struct {
	*TableInfo
	ADDRESS  StringField
	CITY     StringField
	COUNTRY  StringField
	ID       NumberField
	NAME     StringField
	PHONE    StringField
	SID      NumberField
	ZIP_CODE StringField
}

func STAFF_LIST() VIEW_STAFF_LIST {
	tbl := VIEW_STAFF_LIST{TableInfo: NewTableInfo("public", "staff_list")}
	tbl.ADDRESS = NewStringField("address", tbl.TableInfo)
	tbl.CITY = NewStringField("city", tbl.TableInfo)
	tbl.COUNTRY = NewStringField("country", tbl.TableInfo)
	tbl.ID = NewNumberField("id", tbl.TableInfo)
	tbl.NAME = NewStringField("name", tbl.TableInfo)
	tbl.PHONE = NewStringField("phone", tbl.TableInfo)
	tbl.SID = NewNumberField("sid", tbl.TableInfo)
	tbl.ZIP_CODE = NewStringField("zip_code", tbl.TableInfo)
	return tbl
}

func (tbl VIEW_STAFF_LIST) As(alias string) VIEW_STAFF_LIST {
	tbl2 := STAFF_LIST()
	tbl2.TableInfo.Alias = alias
	return tbl2
}
