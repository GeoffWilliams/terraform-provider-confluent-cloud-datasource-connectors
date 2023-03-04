terraform {
  required_providers {
    confluent-cloud-datasource-connectors = {
      source  = "terraform.local/local/confluent-cloud-datasource-connectors"
      version = "0.0.1"
      # Other parameters...
    }
  }
}

locals {
  connectors_to_lookup = {"aconnector":{}, "anotherconnector":{}}
}

data "confluent-cloud-datasource-connectors" "connectors" {
  for_each = local.connectors_to_lookup
  environment_id = "anenv"
  kafka_cluster_id = "deadbeef"
  connector_name = each.key
}

output "aconnector" {
  value = data.confluent-cloud-datasource-connectors.connectors["aconnector"]
}

output "anotherconnector" {
  value = data.confluent-cloud-datasource-connectors.connectors["anotherconnector"]
}

output "bump" {
  value = "bump"
}
