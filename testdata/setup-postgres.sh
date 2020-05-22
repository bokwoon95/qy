#!/bin/bash
docker exec northwind-postgres psql --username=user northwind -f /testdata/northwind.postgres.sql
docker exec northwind-postgres psql --username=user northwind -c '\dt'
