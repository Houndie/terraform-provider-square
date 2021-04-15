package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const catalogTimePeriodEvent = "event"

var catalogTimePeriodSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogTimePeriodEvent: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func catalogTimePeriodSchemaToObject(input map[string]interface{}) *objects.CatalogTimePeriod {
	return &objects.CatalogTimePeriod{Event: input[catalogTimePeriodEvent].(string)}
}

func catalogTimePeriodObjectToSchema(input *objects.CatalogTimePeriod) map[string]interface{} {
	return map[string]interface{}{catalogTimePeriodEvent: input.Event}
}
