#!/bin/bash

# Create CA Key
openssl genrsa -out ca.key 2048

# Create CA Cert
openssl req -new -key ca.key -x509 -out ca.crt -days 3650 -subj "/CN=ca"

# Create Server Key and Signing Request
openssl req -new -nodes -newkey rsa:2048 -keyout server.key -out server.req -batch -subj "/CN=custom-metrics-apiserver.custom-metrics.svc" 

# Create Signed Server Cert
openssl x509 -req -in server.req -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650 -sha256

rm server.req
rm ca.srl
rm ca.key
