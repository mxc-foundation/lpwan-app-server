#!/usr/bin/env bash

mkdir -p certs

# make key and self-signed certificate for master server
openssl req -x509 -newkey rsa:4096 -keyout certs/master.key -out certs/master.cert -nodes -days 365 -subj "/CN=appserver"
# make key and cert signing request for slave (regional server)
openssl req -newkey rsa:4096 -keyout certs/slave.key -out certs/slave.csr -nodes -days 365 -subj "/CN=hax.matchx.io"
# sign slave's csr with master server acting as certificate authority
openssl x509 -req -in certs/slave.csr -CA certs/master.cert -CAkey certs/master.key -out certs/slave.cert -set_serial 01 -days 365
