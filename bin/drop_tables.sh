#!/bin/bash

cat <<EOF > /tmp/drop_tables.sql
USE rest_server;

SHOW tables;

DROP TABLE IF EXISTS config;

DROP TABLE IF EXISTS players;

SHOW tables;

EOF

mysql -urest_api_user -prest_api_pw < /tmp/drop_tables.sql
