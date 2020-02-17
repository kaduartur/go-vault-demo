#!/bin/bash

./create-vault-approle.sh . http://0.0.0.0:8200

export VAULT_ADDR=http://localhost:8200
export VAULT_AUTHENTICATION=APPROLE
export VAULT_ROLE_ID=$(cat /tmp/vault/role-id.txt)
export VAULT_SECRET_ID=$(cat /tmp/vault/secret-id.txt)
export VAULT_DB_BACKEND=database
export VAULT_DB_ROLE=demo_role

export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_NAME=vault_demo

export KAFKA_BOOTSTRAP_SERVERS=0.0.0.0:9092

go run server/cmd/server/main.go