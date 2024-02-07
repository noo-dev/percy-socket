#!/bin/bash

echo "Creating server.key"
openssl genrsa -output server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
echo "creating server.crt"
openssl req -new -x509 -sha256 -key server.key -out server.crt -batch -days 365
