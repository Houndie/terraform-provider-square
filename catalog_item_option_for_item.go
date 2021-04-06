package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const catalogItemOptionForItemItemOptionID = "item_option_id"

var catalogItemOptionForItemSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemOptionForItemItemOptionID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func catalogItemOptionForItemSchemaToObject(input map[string]interface{}) *objects.CatalogItemOptionForItem {
	return &objects.CatalogItemOptionForItem{
		ItemOptionID: input[catalogItemOptionForItemItemOptionID].(string),
	}
}

func catalogItemOptionForItemObjectToSchema(input *objects.CatalogItemOptionForItem) map[string]interface{} {
	return map[string]interface{}{
		catalogItemOptionForItemItemOptionID: input.ItemOptionID,
	}
}
