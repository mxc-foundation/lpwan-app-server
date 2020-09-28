#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    create role mxprotocol with login password 'mxprotocol';
    create database mxprotocol with owner mxprotocol;

    create role test with login password 'test';
    create database test with owner test;
EOSQL
