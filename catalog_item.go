package main

import (
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	catalogItemName                    = "name"
	catalogItemDescription             = "description"
	catalogItemAbbreviation            = "abbreviation"
	catalogItemLabelColor              = "label_color"
	catalogItemAvailableOnline         = "available_online"
	catalogItemAvailableForPickup      = "available_for_pickup"
	catalogItemAvailableElectronically = "available_electronically"
	catalogItemCategoryID              = "category_id"
	catalogItemTaxIDs                  = "tax_ids"
	catalogItemModifierListInfo        = "modifier_list_info"
	catalogItemProductType             = "product_type"
	catalogItemSkipModifierScreen      = "skip_modifier_screen"
	catalogItemItemOptions             = "item_options"

	catalogItemProductTypeRegular             = "REGULAR"
	catalogItemProductTypeAppointmentsService = "APPOINTMENTS_SERVICE"
)

var catalogItemSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemDescription: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemAbbreviation: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemLabelColor: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemAvailableOnline: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogItemAvailableForPickup: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogItemAvailableElectronically: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogItemCategoryID: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemTaxIDs: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		catalogItemModifierListInfo: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     catalogItemModifierListInfoSchema,
		},
		catalogItemProductType: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ValidateDiagFunc: func(input interface{}, path cty.Path) diag.Diagnostics {
				return toDiag(validation.StringInSlice([]string{
					catalogItemProductTypeRegular,
					catalogItemProductTypeAppointmentsService,
				}, false)(input, ""))
			},
			Default: catalogItemProductTypeRegular,
		},
		catalogItemSkipModifierScreen: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogItemItemOptions: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     catalogItemOptionForItemSchema,
		},
	},
}

func toDiag(warnings []string, errors []error) diag.Diagnostics {
	if len(warnings) == 0 && len(errors) == 0 {
		return nil
	}

	result := make([]diag.Diagnostic, 0, len(warnings)+len(errors))
	for _, warning := range warnings {
		result = append(result, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  warning,
		})
	}
	for _, err := range errors {
		result = append(result, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
	}

	return result
}

func catalogItemSchemaToObject(input map[string]interface{}) (*objects.CatalogItem, error) {
	result := &objects.CatalogItem{
		Name:                    input[catalogItemName].(string),
		Description:             input[catalogItemDescription].(string),
		Abbreviation:            input[catalogItemAbbreviation].(string),
		LabelColor:              input[catalogItemLabelColor].(string),
		AvailableOnline:         input[catalogItemAvailableOnline].(bool),
		AvailableForPickup:      input[catalogItemAvailableForPickup].(bool),
		AvailableElectronically: input[catalogItemAvailableElectronically].(bool),
		CategoryID:              input[catalogItemCategoryID].(string),
		SkipModifierScreen:      input[catalogItemSkipModifierScreen].(bool),
	}

	if taxIDs, ok := input[catalogItemTaxIDs]; ok {
		result.TaxIDs = taxIDs.([]string)
	}

	if modifierListInfo, ok := input[catalogItemModifierListInfo]; ok {
		modifierListInfoType := modifierListInfo.([]map[string]interface{})
		result.ModifierListInfo = make([]*objects.CatalogItemModifierListInfo, len(modifierListInfoType))
		for i, info := range modifierListInfoType {
			result.ModifierListInfo[i] = catalogItemModifierListInfoSchemaToObject(info)
		}
	}

	switch input[catalogItemProductType].(string) {
	case catalogItemProductTypeRegular:
		result.ProductType = objects.CatalogItemProductTypeRegular
	case catalogItemProductTypeAppointmentsService:
		result.ProductType = objects.CatalogItemProductTypeAppointmentsService
	default:
		return nil, fmt.Errorf("unknown value for product type: %s", input[catalogItemProductType].(string))
	}

	if options, ok := input[catalogItemItemOptions]; ok {
		optionsType := options.([]map[string]interface{})
		result.ItemOptions = make([]*objects.CatalogItemOptionForItem, len(optionsType))
		for i, option := range optionsType {
			result.ItemOptions[i] = catalogItemOptionForItemSchemaToObject(option)
		}
	}

	return result, nil
}

func catalogItemObjectToSchema(input *objects.CatalogItem) (map[string]interface{}, error) {
	result := map[string]interface{}{
		catalogItemName:                    input.Name,
		catalogItemDescription:             input.Description,
		catalogItemAbbreviation:            input.Abbreviation,
		catalogItemLabelColor:              input.LabelColor,
		catalogItemAvailableOnline:         input.AvailableOnline,
		catalogItemAvailableForPickup:      input.AvailableForPickup,
		catalogItemAvailableElectronically: input.AvailableElectronically,
		catalogItemCategoryID:              input.CategoryID,
		catalogItemSkipModifierScreen:      input.SkipModifierScreen,
	}

	if input.TaxIDs != nil {
		result[catalogItemTaxIDs] = input.TaxIDs
	}

	if input.ModifierListInfo != nil {
		resultModifierListInfo := make([]map[string]interface{}, len(input.ModifierListInfo))
		for i, info := range input.ModifierListInfo {
			resultModifierListInfo[i] = catalogItemModifierListInfoObjectToSchema(info)
		}
		result[catalogItemModifierListInfo] = resultModifierListInfo
	}

	switch input.ProductType {
	case objects.CatalogItemProductTypeRegular:
		result[catalogItemProductType] = catalogItemProductTypeRegular
	case objects.CatalogItemProductTypeAppointmentsService:
		result[catalogItemProductType] = catalogItemProductTypeAppointmentsService
	default:
		return nil, fmt.Errorf("cannot store product type %s", input.ProductType)
	}

	if input.ItemOptions != nil {
		resultOptions := make([]map[string]interface{}, len(input.ItemOptions))
		for i, o := range input.ItemOptions {
			resultOptions[i] = catalogItemOptionForItemObjectToSchema(o)
		}
		result[catalogItemItemOptions] = resultOptions
	}

	return result, nil
}
