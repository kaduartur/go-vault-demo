path "/sys/leases/renew" {
  capabilities = [ "update" ]
}

path "auth/approle/role/demo_role/role-id" {
  capabilities = [ "read" ]
}

path "auth/approle/role/demo_role/secret-id" {
  capabilities = ["create", "read", "update"]
}

path "database/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

path "database/creds/demo_role" {
  capabilities = [ "read" ]
}

path "secret/*" {
  capabilities = ["read"]
}

path "demo/credit-card/*" {
  capabilities = ["create", "update", "delete", "read", "list"]
}

path "demo/transit/encrypt/*" {
  capabilities = ["update"]
}

path "demo/transit/decrypt/*" {
  capabilities = ["update"]
}