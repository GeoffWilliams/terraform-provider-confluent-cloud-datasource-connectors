package main

import (
	confluentCloudDatastore "github.com/geoffwilliams/terraform-provider-confluent-cloud-datasource-connectors/internal/confluent-cloud-datasource-connectors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: confluentCloudDatastore.Provider,
	})
}
