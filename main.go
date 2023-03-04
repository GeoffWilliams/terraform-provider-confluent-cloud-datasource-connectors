package main

import (
	"github.com/geoffwilliams/terraform-provider-confluent-cloud-datasource-connectors/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var (
	// Injected from linker flags like `go build -ldflags "-X main.version=$VERSION" -X ...`
	version   = ""
	userAgent = ""
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.New(version, userAgent)})
}
