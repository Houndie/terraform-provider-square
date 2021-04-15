package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogImageName    = "name"
	catalogImageURL     = "url"
	catalogImageCaption = "caption"
)

var catalogImageSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogImageName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogImageURL: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogImageCaption: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
	},
}

func catalogImageSchemaToObject(input map[string]interface{}) *objects.CatalogImage {
	return &objects.CatalogImage{
		Name:    input[catalogImageName].(string),
		URL:     input[catalogImageURL].(string),
		Caption: input[catalogImageCaption].(string),
	}
}

func catalogImageObjectToSchema(input *objects.CatalogImage) map[string]interface{} {
	return map[string]interface{}{
		catalogImageName:    input.Name,
		catalogImageURL:     input.URL,
		catalogImageCaption: input.Caption,
	}
}
