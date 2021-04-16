package main

import (
	"fmt"
	"time"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogPricingRuleDiscountID          = "discount_id"
	catalogPricingRuleExcludeProductsID   = "exclude_products_id"
	catalogPricingRuleExcludeStrategy     = "exclude_strategy"
	catalogPricingRuleMatchProductsID     = "match_products_id"
	catalogPricingRuleName                = "name"
	catalogPricingRuleTimePeriodIDs       = "time_period_ids"
	catalogPricingRuleValidFromDate       = "valid_from_date"
	catalogPricingRuleValidFromLocalTime  = "valid_from_local_time"
	catalogPricingRuleValidUntilDate      = "valid_until_date"
	catalogPricingRuleValidUntilLocalTime = "valid_until_local_time"

	excludeStrategyLeastExpensive = "LEAST_EXPENSIVE"
	excludeStrategyMostExpensive  = "MOST_EXPENSIVE"
)

var (
	excludeStrategyStrToEnum = map[string]objects.ExcludeStrategy{
		excludeStrategyLeastExpensive: objects.ExcludeStrategyLeastExpensive,
		excludeStrategyMostExpensive:  objects.ExcludeStrategyMostExpensive,
	}

	excludeStrategyEnumToStr = map[objects.ExcludeStrategy]string{
		objects.ExcludeStrategyLeastExpensive: excludeStrategyLeastExpensive,
		objects.ExcludeStrategyMostExpensive:  excludeStrategyMostExpensive,
	}

	excludeStrategyValidate = stringInSlice([]string{excludeStrategyLeastExpensive, excludeStrategyMostExpensive}, false)
)

var catalogPricingRuleSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogPricingRuleDiscountID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogPricingRuleExcludeProductsID: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogPricingRuleExcludeStrategy: &schema.Schema{
			Type:             schema.TypeString,
			Optional:         true,
			Default:          excludeStrategyLeastExpensive,
			ValidateDiagFunc: excludeStrategyValidate,
		},
		catalogPricingRuleMatchProductsID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogPricingRuleName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogPricingRuleTimePeriodIDs: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		catalogPricingRuleValidFromDate: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogPricingRuleValidFromLocalTime: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogPricingRuleValidUntilDate: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogPricingRuleValidUntilLocalTime: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

func catalogPricingRuleSchemaToObject(input map[string]interface{}) (*objects.CatalogPricingRule, error) {
	result := &objects.CatalogPricingRule{
		DiscountID:          input[catalogPricingRuleDiscountID].(string),
		ExcludeProductsID:   input[catalogPricingRuleExcludeProductsID].(string),
		ExcludeStrategy:     excludeStrategyStrToEnum[input[catalogPricingRuleExcludeStrategy].(string)],
		MatchProductsID:     input[catalogPricingRuleMatchProductsID].(string),
		Name:                input[catalogPricingRuleName].(string),
		ValidFromLocalTime:  input[catalogPricingRuleValidFromLocalTime].(string),
		ValidUntilLocalTime: input[catalogPricingRuleValidUntilLocalTime].(string),
	}

	if ids := input[catalogPricingRuleTimePeriodIDs].([]interface{}); len(ids) > 0 {
		result.TimePeriodIDs = make([]string, len(ids))
		for i, id := range ids {
			result.TimePeriodIDs[i] = id.(string)
		}
	}

	if date := input[catalogPricingRuleValidFromDate].(string); date == "" {
		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			return nil, fmt.Errorf("error parsing valid from date: %w", err)
		}

		result.ValidFromDate = &t
	}

	if date := input[catalogPricingRuleValidUntilDate].(string); date == "" {
		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			return nil, fmt.Errorf("error parsing valid until date: %w", err)
		}

		result.ValidUntilDate = &t
	}

	return result, nil
}

func catalogPricingRuleObjectToSchema(input *objects.CatalogPricingRule) map[string]interface{} {
	result := map[string]interface{}{
		catalogPricingRuleDiscountID:          input.DiscountID,
		catalogPricingRuleExcludeProductsID:   input.ExcludeProductsID,
		catalogPricingRuleExcludeStrategy:     input.ExcludeStrategy,
		catalogPricingRuleMatchProductsID:     input.MatchProductsID,
		catalogPricingRuleName:                input.Name,
		catalogPricingRuleValidFromLocalTime:  input.ValidFromLocalTime,
		catalogPricingRuleValidUntilLocalTime: input.ValidUntilLocalTime,
	}

	if input.TimePeriodIDs != nil {
		ids := make([]interface{}, len(input.TimePeriodIDs))
		for i, id := range input.TimePeriodIDs {
			ids[i] = id
		}

		result[catalogPricingRuleTimePeriodIDs] = ids
	}

	if input.ValidFromDate != nil {
		result[catalogPricingRuleValidFromDate] = input.ValidFromDate.Format(time.RFC3339)
	}

	if input.ValidUntilDate != nil {
		result[catalogPricingRuleValidUntilDate] = input.ValidUntilDate.Format(time.RFC3339)
	}

	return result
}
