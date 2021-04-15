package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const catalogCategoryName = "name"

var catalogCategorySchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogCategoryName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func catalogCategorySchemaToObject(input map[string]interface{}) *objects.CatalogCategory {
	return &objects.CatalogCategory{
		Name: input[catalogCategoryName].(string),
	}
}

func catalogCategoryObjectToSchema(input *objects.CatalogCategory) map[string]interface{} {
	return map[string]interface{}{
		catalogCategoryName: input.Name,
	}
}
