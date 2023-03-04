terraform {
  required_providers {
    confluent-cloud-datasource-connectors = {
      source  = "terraform.local/local/confluent-cloud-datasource-connectors"
      version = "0.0.1"
      # Other parameters...
    }
  }
}


data "confluent-cloud-datasource-connectors" "mycluster" {
  environment_id = "anenv"
  kafka_cluster_id = "deadbeef"
}

output "debug" {
  value = data.confluent-cloud-datasource-connectors.mycluster.connectors
}

output "avalue" {
  value = data.confluent-cloud-datasource-connectors.mycluster.avalue
}

output "dump" {
  value = data.confluent-cloud-datasource-connectors.mycluster
}

output "bump" {
  value = "bump"
}
