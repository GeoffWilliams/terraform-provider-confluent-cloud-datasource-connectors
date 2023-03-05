package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	connect "github.com/confluentinc/ccloud-sdk-go-v2/connect/v1"
)

const (
	terraformProviderUserAgent = "terraform-provider-confluent"
)

// // Provider builds and returns the Env provider schema. It is consumed only by
// // the main function.
// func Provider() *schema.Provider {
// 	return &schema.Provider{
// 		DataSourcesMap: map[string]*schema.Resource{
// 			"confluent-cloud-datasource-connectors": confluentConnectorsDataSource(),
// 		},
// 	}
// }

// steal the client setup code from https://github.com/confluentinc/terraform-provider-confluent/blob/master/internal/provider/provider.go
type Client struct {
	connectClient *connect.APIClient

	userAgent                   string
	cloudApiKey                 string
	cloudApiSecret              string
	kafkaClusterId              string
	kafkaApiKey                 string
	kafkaApiSecret              string
	kafkaRestEndpoint           string
	isKafkaClusterIdSet         bool
	isKafkaMetadataSet          bool
	schemaRegistryClusterId     string
	schemaRegistryApiKey        string
	schemaRegistryApiSecret     string
	schemaRegistryRestEndpoint  string
	isSchemaRegistryMetadataSet bool
}

func New(version, userAgent string) func() *schema.Provider {
	return func() *schema.Provider {
		provider := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"cloud_api_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("CONFLUENT_CLOUD_API_KEY", ""),
					Description: "The Confluent Cloud API Key.",
				},
				"cloud_api_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("CONFLUENT_CLOUD_API_SECRET", ""),
					Description: "The Confluent Cloud API Secret.",
				},
				"kafka_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("KAFKA_ID", ""),
					Description: "The Kafka Cluster ID.",
				},
				"kafka_api_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("KAFKA_API_KEY", ""),
					Description: "The Kafka Cluster API Key.",
				},
				"kafka_api_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("KAFKA_API_SECRET", ""),
					Description: "The Kafka Cluster API Secret.",
				},
				"kafka_rest_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("KAFKA_REST_ENDPOINT", ""),
					Description: "The Kafka Cluster REST Endpoint.",
				},
				"schema_registry_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SCHEMA_REGISTRY_ID", ""),
					Description: "The Schema Registry Cluster ID.",
				},
				"schema_registry_api_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SCHEMA_REGISTRY_API_KEY", ""),
					Description: "The Schema Registry Cluster API Key.",
				},
				"schema_registry_api_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SCHEMA_REGISTRY_API_SECRET", ""),
					Description: "The Schema Registry Cluster API Secret.",
				},
				"schema_registry_rest_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SCHEMA_REGISTRY_REST_ENDPOINT", ""),
					Description: "The Schema Registry Cluster REST Endpoint.",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "https://api.confluent.cloud",
					Description: "The base endpoint of Confluent Cloud API.",
				},
				"max_retries": {
					Type:     schema.TypeInt,
					Optional: true,
					//ValidateFunc: validation.IntAtLeast(5),
					Description: "Maximum number of retries of HTTP client. Defaults to 4.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"confluent-cloud-datasource-connectors": confluentConnectorsDataSource(),
			},
		}

		provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return providerConfigure(ctx, d, provider, version, userAgent)
		}

		return provider
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider, providerVersion, additionalUserAgent string) (interface{}, diag.Diagnostics) {
	tflog.Info(ctx, "Initializing Terraform Provider for Confluent Cloud")
	endpoint := d.Get("endpoint").(string)
	cloudApiKey := d.Get("cloud_api_key").(string)
	cloudApiSecret := d.Get("cloud_api_secret").(string)
	kafkaClusterId := d.Get("kafka_id").(string)
	kafkaApiKey := d.Get("kafka_api_key").(string)
	kafkaApiSecret := d.Get("kafka_api_secret").(string)
	kafkaRestEndpoint := d.Get("kafka_rest_endpoint").(string)
	schemaRegistryClusterId := d.Get("schema_registry_id").(string)
	schemaRegistryApiKey := d.Get("schema_registry_api_key").(string)
	schemaRegistryApiSecret := d.Get("schema_registry_api_secret").(string)
	schemaRegistryRestEndpoint := d.Get("schema_registry_rest_endpoint").(string)
	// If the max_retries doesn't exist in the configuration, then 0 will be returned.
	maxRetries := d.Get("max_retries").(int)

	// 3 or 4 attributes should be set or not set at the same time
	// Option #2: (kafka_api_key, kafka_api_secret, kafka_rest_endpoint)
	// Option #3 (primary): (kafka_api_key, kafka_api_secret, kafka_rest_endpoint, kafka_id)
	allKafkaAttributesAreSet := (kafkaApiKey != "") && (kafkaApiSecret != "") && (kafkaRestEndpoint != "")
	allKafkaAttributesAreNotSet := (kafkaApiKey == "") && (kafkaApiSecret == "") && (kafkaRestEndpoint == "")
	justOneOrTwoKafkaAttributesAreSet := !(allKafkaAttributesAreSet || allKafkaAttributesAreNotSet)
	if justOneOrTwoKafkaAttributesAreSet {
		return nil, diag.Errorf("(kafka_api_key, kafka_api_secret, kafka_rest_endpoint) or (kafka_api_key, kafka_api_secret, kafka_rest_endpoint, kafka_id) attributes should be set or not set in the provider block at the same time")
	}

	// All 4 attributes should be set or not set at the same time
	allSchemaRegistryAttributesAreSet := (schemaRegistryApiKey != "") && (schemaRegistryApiSecret != "") && (schemaRegistryRestEndpoint != "") && (schemaRegistryClusterId != "")
	allSchemaRegistryAttributesAreNotSet := (schemaRegistryApiKey == "") && (schemaRegistryApiSecret == "") && (schemaRegistryRestEndpoint == "") && (schemaRegistryClusterId == "")
	justSubsetOfSchemaRegistryAttributesAreSet := !(allSchemaRegistryAttributesAreSet || allSchemaRegistryAttributesAreNotSet)
	if justSubsetOfSchemaRegistryAttributesAreSet {
		return nil, diag.Errorf("All 4 schema_registry_api_key, schema_registry_api_secret, schema_registry_rest_endpoint, schema_registry_id attributes should be set or not set in the provider block at the same time")
	}

	userAgent := p.UserAgent(terraformProviderUserAgent, fmt.Sprintf("%s (https://confluent.cloud; support@confluent.io)", providerVersion))
	if additionalUserAgent != "" {
		userAgent = fmt.Sprintf("%s %s", additionalUserAgent, userAgent)
	}

	connectCfg := connect.NewConfiguration()

	connectCfg.Servers[0].URL = endpoint

	connectCfg.UserAgent = userAgent

	if maxRetries != 0 {
		connectCfg.HTTPClient = NewRetryableClientFactory(WithMaxRetries(maxRetries)).CreateRetryableClient()
	} else {
		connectCfg.HTTPClient = NewRetryableClientFactory().CreateRetryableClient()
	}

	client := Client{
		connectClient:              connect.NewAPIClient(connectCfg),
		userAgent:                  userAgent,
		cloudApiKey:                cloudApiKey,
		cloudApiSecret:             cloudApiSecret,
		kafkaClusterId:             kafkaClusterId,
		kafkaApiKey:                kafkaApiKey,
		kafkaApiSecret:             kafkaApiSecret,
		kafkaRestEndpoint:          kafkaRestEndpoint,
		schemaRegistryClusterId:    schemaRegistryClusterId,
		schemaRegistryApiKey:       schemaRegistryApiKey,
		schemaRegistryApiSecret:    schemaRegistryApiSecret,
		schemaRegistryRestEndpoint: schemaRegistryRestEndpoint,
		// For simplicity, treat 3 (for Kafka) and 4 (for SR) variables as a "single" one
		isKafkaMetadataSet:          allKafkaAttributesAreSet,
		isKafkaClusterIdSet:         kafkaClusterId != "",
		isSchemaRegistryMetadataSet: allSchemaRegistryAttributesAreSet,
	}

	return &client, nil
}

//
