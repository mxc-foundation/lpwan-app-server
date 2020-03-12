#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    create role mxprotocol_m2m_wallet with login password 'mxprotocol_m2m_wallet';
    create database mxprotocol_m2m_wallet with owner mxprotocol_m2m_wallet;

    create role test with login password 'test';
    create database test with owner test;
EOSQL
