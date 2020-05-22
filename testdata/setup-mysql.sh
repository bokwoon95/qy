#!/bin/bash
docker exec northwind-mysql mysql --user=user --password=password northwind -e 'source /testdata/northwind.mysql.sql'
docker exec northwind-mysql mysql --user=user --password=password northwind -e 'show tables;'
