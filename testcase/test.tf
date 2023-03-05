terraform {
  required_providers {
    confluent-cloud-datasource-connectors = {
      source  = "terraform.local/local/confluent-cloud-datasource-connectors"
      version = "0.0.1"
      # Other parameters...
    }
  }
}

provider "confluent-cloud-datasource-connectors" {
#  cloud_api_key    = var.confluent_cloud_api_key    # optionally use CONFLUENT_CLOUD_API_KEY env var
#  cloud_api_secret = var.confluent_cloud_api_secret # optionally use CONFLUENT_CLOUD_API_SECRET env var
}

locals {
  connectors_to_lookup = {"gwilliams-cpsconnect-discountCodes":{}, "gwilliams-mssql-jdbc":{}, "nonexistantconnector": {}}
}

data "confluent-cloud-datasource-connectors" "confluent_connectors" {
  for_each = local.connectors_to_lookup
  environment_id = "env-w5502w"
  kafka_cluster_id = "lkc-7n1ny1"
  connector_name = each.key
}

output "gwilliams-cpsconnect-discountCodes" {
  value = data.confluent-cloud-datasource-connectors.confluent_connectors["gwilliams-cpsconnect-discountCodes"]
}

output "gwilliams-mssql-jdbc" {
  value = data.confluent-cloud-datasource-connectors.confluent_connectors["gwilliams-mssql-jdbc"]
}

output "nonexistantconnector" {
  value = data.confluent-cloud-datasource-connectors.confluent_connectors["nonexistantconnector"]
}

output "bump" {
  value = "bump"
}
