package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogModifierName           = "name"
	catalogModifierPriceMoney     = "price_money"
	catalogModifierOrdinal        = "ordinal"
	catalogModifierModifierListID = "modifier_list_id"
)

var catalogModifierSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogModifierName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogModifierPriceMoney: &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem:     moneySchema,
		},
		catalogModifierOrdinal: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		catalogModifierModifierListID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func catalogModifierSchemaToObject(input map[string]interface{}) *objects.CatalogModifier {
	return &objects.CatalogModifier{
		Name:           input[catalogModifierName].(string),
		PriceMoney:     moneySchemaToObject(input[catalogModifierPriceMoney].(*schema.Set).List()[0].(map[string]interface{})),
		Ordinal:        input[catalogModifierOrdinal].(int),
		ModifierListID: input[catalogModifierModifierListID].(string),
	}
}

func catalogModifierObjectToSchema(input *objects.CatalogModifier) map[string]interface{} {
	return map[string]interface{}{
		catalogModifierName:           input.Name,
		catalogModifierPriceMoney:     schema.NewSet(schema.HashResource(moneySchema), []interface{}{moneyObjectToSchema(input.PriceMoney)}),
		catalogModifierOrdinal:        input.Ordinal,
		catalogModifierModifierListID: input.ModifierListID,
	}
}
