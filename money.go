package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	moneyAmount   = "amount"
	moneyCurrency = "currency"
)

var moneySchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		moneyAmount: &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		moneyCurrency: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func moneySchemaToObject(input map[string]interface{}) *objects.Money {
	return &objects.Money{
		Amount:   input[moneyAmount].(int),
		Currency: input[moneyCurrency].(string),
	}
}

func moneyObjectToSchema(input *objects.Money) map[string]interface{} {
	return map[string]interface{}{
		moneyAmount:   input.Amount,
		moneyCurrency: input.Currency,
	}
}
