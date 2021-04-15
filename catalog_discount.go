package main

import (
	"errors"
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogDiscountName           = "name"
	catalogDiscountPinRequired    = "pin_required"
	catalogDiscountLabelColor     = "label_color"
	catalogDiscountModifyTaxBasis = "modify_tax_basis"
	catalogDiscountDiscountType   = "discount_type"
	catalogDiscountPercentage     = "percentage"
	catalogDiscountAmountMoney    = "amount_money"

	catalogDiscountTypeFixedPercentage    = "FIXED_PERCENTAGE"
	catalogDiscountTypeFixedAmount        = "FIXED_AMOUNT"
	catalogDiscountTypeVariablePercentage = "VARIABLE_PERCENTAGE"
	catalogDiscountTypeVariableAmount     = "VARIABLE_AMOUNT"
)

var catalogDiscountSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogDiscountName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogDiscountPinRequired: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogDiscountLabelColor: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogDiscountModifyTaxBasis: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogDiscountDiscountType: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogDiscountPercentage: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogDiscountAmountMoney: &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem:     moneySchema,
		},
	},
}

func catalogDiscountSchemaToObject(input map[string]interface{}) (*objects.CatalogDiscount, error) {
	result := &objects.CatalogDiscount{
		Name:           input[catalogDiscountName].(string),
		PinRequired:    input[catalogDiscountPinRequired].(bool),
		LabelColor:     input[catalogDiscountLabelColor].(string),
		ModifyTaxBasis: input[catalogDiscountModifyTaxBasis].(string),
	}

	switch input[catalogDiscountDiscountType] {
	case catalogDiscountTypeFixedPercentage:
		percentage, ok := input[catalogDiscountPercentage]
		if !ok {
			return nil, errors.New("fixed percentage chosen, but percentage field empty")
		}
		result.DiscountType = &objects.CatalogDiscountFixedPercentage{Percentage: percentage.(string)}
	case catalogDiscountTypeVariablePercentage:
		percentage, ok := input[catalogDiscountPercentage]
		if !ok {
			return nil, errors.New("variable percentage chosen, but percentage field empty")
		}
		result.DiscountType = &objects.CatalogDiscountVariablePercentage{Percentage: percentage.(string)}
	case catalogDiscountTypeFixedAmount:
		amount, ok := input[catalogDiscountAmountMoney]
		if !ok {
			return nil, errors.New("fixed amount chosen, but amount field empty")
		}
		result.DiscountType = &objects.CatalogDiscountFixedAmount{AmountMoney: moneySchemaToObject(amount.(map[string]interface{}))}
	case catalogDiscountTypeVariableAmount:
		amount, ok := input[catalogDiscountAmountMoney]
		if !ok {
			return nil, errors.New("variable amount chosen, but amount field empty")
		}
		result.DiscountType = &objects.CatalogDiscountVariableAmount{AmountMoney: moneySchemaToObject(amount.(map[string]interface{}))}
	default:
		return nil, fmt.Errorf("unknown discount type: %s", input[catalogDiscountDiscountType])
	}

	return result, nil
}

func catalogDiscountObjectToSchema(input *objects.CatalogDiscount) (map[string]interface{}, error) {
	result := map[string]interface{}{
		catalogDiscountName:           input.Name,
		catalogDiscountPinRequired:    input.PinRequired,
		catalogDiscountLabelColor:     input.LabelColor,
		catalogDiscountModifyTaxBasis: input.ModifyTaxBasis,
	}

	switch t := input.DiscountType.(type) {
	case *objects.CatalogDiscountFixedPercentage:
		result[catalogDiscountDiscountType] = catalogDiscountTypeFixedPercentage
		result[catalogDiscountPercentage] = t.Percentage
	case *objects.CatalogDiscountVariablePercentage:
		result[catalogDiscountDiscountType] = catalogDiscountTypeVariablePercentage
		result[catalogDiscountPercentage] = t.Percentage
	case *objects.CatalogDiscountFixedAmount:
		result[catalogDiscountDiscountType] = catalogDiscountTypeFixedAmount
		result[catalogDiscountAmountMoney] = t.AmountMoney
	case *objects.CatalogDiscountVariableAmount:
		result[catalogDiscountDiscountType] = catalogDiscountTypeVariableAmount
		result[catalogDiscountAmountMoney] = t.AmountMoney
	default:
		return nil, fmt.Errorf("unknown discount type: %s", input.DiscountType)
	}

	return result, nil
}
