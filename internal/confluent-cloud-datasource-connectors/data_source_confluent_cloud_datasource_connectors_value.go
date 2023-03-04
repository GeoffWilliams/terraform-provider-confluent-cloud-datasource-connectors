package confluentCloudDatasourceConnectors

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	paramConnectors     = "connectors"
	paramKafkaClusterId = "kafka_cluster_id"
	paramEnvironmentId  = "environment_id"

	paramAvalue = "avalue"
)

func kafkaTopicDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: confluentCloudConnectorsRead,
		Schema: map[string]*schema.Schema{
			paramKafkaClusterId: {
				Type:     schema.TypeString,
				Required: true,
			},
			paramEnvironmentId: {
				Type:     schema.TypeString,
				Required: true,
			},

			paramConnectors: {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},

			paramAvalue: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		SchemaVersion: 2,
	}
}

func confluentCloudConnectorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// must be set to something or all data values be null
	d.SetId("must")

	connectors := make(map[string]string)
	connectors["x"] = "z"
	if err := d.Set(paramConnectors, connectors); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(paramAvalue, "hello"); err != nil {
		return diag.FromErr(err)
	}

	//tflog.Debug(fmt.Sprintf("============================================================= Connectors found: %d", len(connectors)))
	return nil
}

// func kafkaTopicDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	restEndpoint, err := extractRestEndpoint(meta.(*Client), d, false)
// 	if err != nil {
// 		return diag.Errorf("error reading Kafka Topic: %s", createDescriptiveError(err))
// 	}
// 	clusterId, err := extractKafkaClusterId(meta.(*Client), d, false)
// 	if err != nil {
// 		return diag.Errorf("error reading Kafka Topic: %s", createDescriptiveError(err))
// 	}
// 	clusterApiKey, clusterApiSecret, err := extractClusterApiKeyAndApiSecret(meta.(*Client), d, false)
// 	if err != nil {
// 		return diag.Errorf("error reading Kafka Topic: %s", createDescriptiveError(err))
// 	}
// 	kafkaRestClient := meta.(*Client).kafkaRestClientFactory.CreateKafkaRestClient(restEndpoint, clusterId, clusterApiKey, clusterApiSecret, meta.(*Client).isKafkaMetadataSet, meta.(*Client).isKafkaClusterIdSet)
// 	topicName := d.Get(paramTopicName).(string)
// 	tflog.Debug(ctx, fmt.Sprintf("Reading Kafka Topic %q", topicName))

// 	// Mark resource as new to avoid d.Set("") when getting 404
// 	d.MarkNewResource()

// 	if _, err := readTopicAndSetAttributes(ctx, d, kafkaRestClient, topicName); err != nil {
// 		return diag.Errorf("error reading Kafka Topic %q: %s", topicName, createDescriptiveError(err))
// 	}

// 	tflog.Debug(ctx, fmt.Sprintf("Finished reading Kafka Topic %q", topicName))

// 	return nil
// }
