package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogItemOptionValueColor        = "color"
	catalogItemOptionValueDescription  = "description"
	catalogItemOptionValueItemOptionID = "item_option_id"
	catalogItemOptionValueName         = "name"
	catalogItemOptionValueOrdinal      = "ordinal"
)

var catalogItemOptionValueSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemOptionValueColor: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemOptionValueDescription: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemOptionValueItemOptionID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemOptionValueName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemOptionValueOrdinal: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
	},
}

func catalogItemOptionValueSchemaToObject(input map[string]interface{}) *objects.CatalogItemOptionValue {
	return &objects.CatalogItemOptionValue{
		Color:        input[catalogItemOptionValueColor].(string),
		Description:  input[catalogItemOptionValueDescription].(string),
		ItemOptionID: input[catalogItemOptionValueItemOptionID].(string),
		Name:         input[catalogItemOptionValueName].(string),
		Ordinal:      input[catalogItemOptionValueOrdinal].(int),
	}
}

func catalogItemOptionValueObjectToSchema(input *objects.CatalogItemOptionValue) map[string]interface{} {
	return map[string]interface{}{
		catalogItemOptionValueColor:        input.Color,
		catalogItemOptionValueDescription:  input.Description,
		catalogItemOptionValueItemOptionID: input.ItemOptionID,
		catalogItemOptionValueName:         input.Name,
		catalogItemOptionValueOrdinal:      input.Ordinal,
	}
}
