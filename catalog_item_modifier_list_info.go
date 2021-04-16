package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogItemModifierListInfoModifierListID       = "modifer_list_id"
	catalogItemModifierListInfoModifierOverrides    = "modifier_overrides"
	catalogItemModifierListInfoMinSelectedModifiers = "min_selected_modifiers"
	catalogItemModifierListInfoMaxSelectedModifiers = "max_selected_modifiers"
	catalogItemModifierListInfoEnabled              = "enabled"

	catalogModifierOverrideModifierID  = "modifier_id"
	catalogModifierOverrideOnByDefault = "on_by_default"
)

var catalogModifierOverrideSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogModifierOverrideModifierID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogModifierOverrideOnByDefault: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	},
}

func catalogModifierOverrideSchemaToObject(input map[string]interface{}) *objects.CatalogModifierOverride {
	return &objects.CatalogModifierOverride{
		ModifierID:  input[catalogModifierOverrideModifierID].(string),
		OnByDefault: input[catalogModifierOverrideOnByDefault].(bool),
	}
}

func catalogModifierOverrideObjectToSchema(input *objects.CatalogModifierOverride) map[string]interface{} {
	return map[string]interface{}{
		catalogModifierOverrideModifierID:  input.ModifierID,
		catalogModifierOverrideOnByDefault: input.OnByDefault,
	}
}

var catalogItemModifierListInfoSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemModifierListInfoModifierListID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemModifierListInfoModifierOverrides: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     catalogModifierOverrideSchema,
		},
		catalogItemModifierListInfoMinSelectedModifiers: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		catalogItemModifierListInfoMaxSelectedModifiers: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		catalogItemModifierListInfoEnabled: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
	},
}

func catalogItemModifierListInfoSchemaToObject(input map[string]interface{}) *objects.CatalogItemModifierListInfo {
	result := &objects.CatalogItemModifierListInfo{
		ModifierListID:       input[catalogItemModifierListInfoModifierListID].(string),
		MinSelectedModifiers: input[catalogItemModifierListInfoMinSelectedModifiers].(int),
		MaxSelectedModifiers: input[catalogItemModifierListInfoMaxSelectedModifiers].(int),
	}

	if overrides := input[catalogItemModifierListInfoModifierOverrides].([]interface{}); len(overrides) > 0 {
		result.ModifierOverrides = make([]*objects.CatalogModifierOverride, len(overrides))

		for i, override := range overrides {
			result.ModifierOverrides[i] = catalogModifierOverrideSchemaToObject(override.(map[string]interface{}))
		}
	}

	enabledType := input[catalogItemModifierListInfoEnabled].(bool)
	result.Enabled = &enabledType

	return result
}

func catalogItemModifierListInfoObjectToSchema(input *objects.CatalogItemModifierListInfo) map[string]interface{} {
	result := map[string]interface{}{
		catalogItemModifierListInfoModifierListID:       input.ModifierListID,
		catalogItemModifierListInfoMinSelectedModifiers: input.MinSelectedModifiers,
		catalogItemModifierListInfoMaxSelectedModifiers: input.MaxSelectedModifiers,
	}

	if input.ModifierOverrides != nil {
		overrides := make([]interface{}, len(input.ModifierOverrides))
		for i, override := range input.ModifierOverrides {
			overrides[i] = catalogModifierOverrideObjectToSchema(override)
		}

		result[catalogItemModifierListInfoModifierOverrides] = overrides
	}

	if input.Enabled != nil {
		result[catalogItemModifierListInfoEnabled] = *input.Enabled
	}

	return result
}
