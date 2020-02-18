storage "consul" {
  address = "consul:8500"
  path    = "vault"
}

max_lease_ttl = "1h"
default_lease_ttl = "1h"