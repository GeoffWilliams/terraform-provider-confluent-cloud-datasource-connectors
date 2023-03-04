package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	paramConfig         = "config"
	paramStatus         = "status"
	paramKafkaClusterId = "kafka_cluster_id"
	paramEnvironmentId  = "environment_id"

	paramConnectorName = "connector_name"
)

func confluentConnectorsDataSource() *schema.Resource {
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

			paramConnectorName: {
				Type:     schema.TypeString,
				Required: true,
			},

			paramConfig: {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			paramStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		SchemaVersion: 2,
	}
}

func confluentCloudConnectorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// must be set to something or all data values be null!
	c := meta.(*Client)

	// connector info
	req1 := c.connectClient.ConnectorsV1Api.ReadConnectv1Connector(
		c.connectApiContext(ctx),
		d.Get(paramConnectorName).(string), d.Get(paramEnvironmentId).(string), d.Get(paramKafkaClusterId).(string),
	)
	apiError, something, err := req1.Execute()
	tflog.Debug(ctx, fmt.Sprintf("apiError: %s", apiError))
	tflog.Debug(ctx, fmt.Sprintf("something: %s", something))
	tflog.Debug(ctx, fmt.Sprintf("err: %s", err))

	// connector config
	req2 := c.connectClient.ConnectorsV1Api.GetConnectv1ConnectorConfig(
		c.connectApiContext(ctx),
		d.Get(paramConnectorName).(string), d.Get(paramEnvironmentId).(string), d.Get(paramKafkaClusterId).(string),
	)
	apiError1, something1, err1 := req2.Execute()
	tflog.Debug(ctx, fmt.Sprintf("apiError1: %s", apiError1))
	tflog.Debug(ctx, fmt.Sprintf("something1: %s", something1))
	tflog.Debug(ctx, fmt.Sprintf("err1: %s", err1))

	d.SetId(d.Get(paramConnectorName).(string))

	// config := make(map[string]string)
	// config["x"] = "z"
	if err := d.Set(paramConfig, something1); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(paramStatus, "NOT_DEFINED"); err != nil {
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
