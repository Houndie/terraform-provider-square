package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogV1IDCatalogV1ID = "catalog_v1_id"
	catalogV1IDLocationID  = "location_id"
)

var catalogV1IDSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogV1IDCatalogV1ID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogV1IDLocationID: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
	},
}

func catalogV1IDObjectToSchema(input *objects.CatalogV1ID) map[string]interface{} {
	return map[string]interface{}{
		catalogV1IDCatalogV1ID: input.CatalogV1ID,
		catalogV1IDLocationID:  input.LocationID,
	}
}

func catalogV1IDSchemaToObject(input map[string]interface{}) *objects.CatalogV1ID {
	return &objects.CatalogV1ID{
		CatalogV1ID: input[catalogV1IDCatalogV1ID].(string),
		LocationID:  input[catalogV1IDLocationID].(string),
	}
}
