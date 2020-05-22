#!/bin/bash
# Prerequisite steps:
# cd testdata/
# export $(cat .env | xargs)
# docker volume rm testdata_postgres_data testdata_mysql_data # get rid of stale data
# docker-compose up -d
shopt -s extglob
set -o allexport; source .env; set +o allexport

SAKILA='Sakila'

if [ ! -d "$SAKILA" ]; then
  echo 'Sakila/ directory not found, please download it from https://github.com/jOOQ/jOOQ/tree/master/jOOQ-examples/Sakila and place it in the testdata directory'
  exit 1
fi

HELP="   Usage: $0 ehhh
"

declare -a Pre
declare -a Dirs
declare -a Post
declare -a Finish

# Unpack script arguments
argc="$#";: "$((i=0))"
while [ "$((i))" -lt "$argc" ]; do
  case "$1" in
    --help|-h) Help='true'
    ;; -postgres)
        Pre+=("docker exec ${POSTGRES_NAME}-postgres psql --username=${POSTGRES_USER} ${POSTGRES_NAME} -q -f '")
        Dirs+=("${SAKILA}/postgres-sakila-db/")
        Post+=("'")
        Finish+=("docker exec ${POSTGRES_NAME}-postgres psql --username=${POSTGRES_USER} ${POSTGRES_NAME} -c '\dt'")
    ;; -mysql)
        Pre+=("docker exec ${MYSQL_NAME}-mysql mysql --user=root --password=root ${MYSQL_NAME} -e 'source ")
        Dirs+=("${SAKILA}/mysql-sakila-db/")
        Post+=("'")
        Finish+=("docker exec ${MYSQL_NAME}-mysql mysql --user=root --password=root ${MYSQL_NAME} -e 'show tables;'")
    ;; -sqlite)
        Pre+=("echo ':^)'")
        Dirs+=("sqlite-sakila-db/")
    ;; -drop) Drop='true'
    ;; -delete) Delete='true'
    ;; -schema) Schema='true'
    ;; -insert) Insert='true'
    ;; *) :
  esac
  shift;: "$((i=i+1))"
done

[ "$Help" ] && echo "$HELP" && exit 0

len="${#Pre[@]}"
for (( i = 0; i < len; i++ )); do
  if [ "$Drop" ]; then
    files=( "${Dirs[i]}"*drop* )
    echo "${Pre[i]}/testdata/${files[0]}${Post[i]}"
  fi
  if [ "$Delete" ]; then
    files=( "${Dirs[i]}"*delete* )
    echo "${Pre[i]}/testdata/${files[0]}${Post[i]}"
  fi
  if [ "$Schema" ]; then
    files=( "${Dirs[i]}"*schema* )
    echo "${Pre[i]}/testdata/${files[0]}${Post[i]}"
  fi
  if [ "$Insert" ]; then
    files=( "${Dirs[i]}"*insert* )
    echo "${Pre[i]}/testdata/${files[0]}${Post[i]}"
  fi
  echo "${Finish[i]}"
done
