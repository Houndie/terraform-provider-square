package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogItemOptionDescription = "description"
	catalogItemOptionDisplayName = "display_name"
	catalogItemOptionName        = "name"
	catalogItemOptionShowColors  = "show_colors"
)

var catalogItemOptionSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemOptionDescription: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemOptionDisplayName: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemOptionName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemOptionShowColors: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	},
}

func catalogItemOptionSchemaToObject(input map[string]interface{}) *objects.CatalogItemOption {
	return &objects.CatalogItemOption{
		Description: input[catalogItemOptionDescription].(string),
		DisplayName: input[catalogItemOptionDisplayName].(string),
		Name:        input[catalogItemOptionName].(string),
		ShowColors:  input[catalogItemOptionShowColors].(bool),
	}
}

func catalogItemOptionObjectToSchema(input *objects.CatalogItemOption) map[string]interface{} {
	return map[string]interface{}{
		catalogItemOptionDescription: input.Description,
		catalogItemOptionDisplayName: input.DisplayName,
		catalogItemOptionName:        input.Name,
		catalogItemOptionShowColors:  input.ShowColors,
	}
}
