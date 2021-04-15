package main

import (
	"errors"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogProductSetAllProducts   = "all_products"
	catalogProductSetName          = "name"
	catalogProductSetProductIDsAll = "product_ids_all"
	catalogProductSetProductIDsAny = "product_ids_any"
	catalogProductSetQuantityExact = "quantity_exact"
	catalogProductSetQuantityMax   = "quantity_max"
	catalogProductSetQuantityMin   = "quantity_min"
)

var catalogProductSetSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogProductSetAllProducts: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogProductSetName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogProductSetProductIDsAll: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		catalogProductSetProductIDsAny: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		catalogProductSetQuantityExact: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		catalogProductSetQuantityMax: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		catalogProductSetQuantityMin: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
	},
}

func catalogProductSetSchemaToObject(input map[string]interface{}) (*objects.CatalogProductSet, error) {
	result := &objects.CatalogProductSet{
		Name: input[catalogProductSetName].(string),
	}

	if input[catalogProductSetAllProducts].(bool) {
		result.Products = &objects.CatalogProductSetAllProducts{}
	} else if products, ok := input[catalogProductSetProductIDsAll]; ok {
		result.Products = &objects.CatalogProductSetAllIDs{
			IDs: products.([]string),
		}
	} else if products, ok := input[catalogProductSetProductIDsAny]; ok {
		result.Products = &objects.CatalogProductSetAnyIDs{
			IDs: products.([]string),
		}
	} else {
		return nil, errors.New("one of all products, product ids all, or product ids any must be set")
	}

	if exact := input[catalogProductSetQuantityExact].(int); exact != 0 {
		result.Quantity = &objects.CatalogProductSetQuantityExact{
			Amount: exact,
		}
	} else if min, max := input[catalogProductSetQuantityMin].(int), input[catalogProductSetQuantityMax].(int); min != 0 || max != 0 {
		result.Quantity = &objects.CatalogProductSetQuantityRange{
			Min: min,
			Max: max,
		}
	} else {
		return nil, errors.New("one of quantity exact, quantity min, or quantity max must be set")
	}

	return result, nil
}

func catalogProductSetObjectToSchema(input *objects.CatalogProductSet) (map[string]interface{}, error) {
	result := map[string]interface{}{
		catalogProductSetName: input.Name,
	}

	switch t := input.Products.(type) {
	case *objects.CatalogProductSetAllProducts:
		result[catalogProductSetAllProducts] = true
	case *objects.CatalogProductSetAllIDs:
		result[catalogProductSetProductIDsAll] = t.IDs
	case *objects.CatalogProductSetAnyIDs:
		result[catalogProductSetProductIDsAny] = t.IDs
	default:
		return nil, errors.New("unknown product selection found")
	}

	switch t := input.Quantity.(type) {
	case *objects.CatalogProductSetQuantityExact:
		result[catalogProductSetQuantityExact] = t.Amount
	case *objects.CatalogProductSetQuantityRange:
		result[catalogProductSetQuantityMin] = t.Min
		result[catalogProductSetQuantityMax] = t.Max
	default:
		return nil, errors.New("unknown product set quantity found")
	}

	return result, nil
}
