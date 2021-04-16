package main

import (
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

var (
	catalogItemProductTypeStrToEnum = map[string]objects.CatalogItemProductType{
		catalogItemProductTypeRegular:             objects.CatalogItemProductTypeRegular,
		catalogItemProductTypeAppointmentsService: objects.CatalogItemProductTypeAppointmentsService,
	}

	catalogItemProductTypeEnumToStr = map[objects.CatalogItemProductType]string{
		objects.CatalogItemProductTypeRegular:             catalogItemProductTypeRegular,
		objects.CatalogItemProductTypeAppointmentsService: catalogItemProductTypeAppointmentsService,
	}

	catalogItemProductTypeValidate = stringInSlice([]string{catalogItemProductTypeRegular, catalogItemProductTypeAppointmentsService}, false)
)

func stringInSlice(slice []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(input interface{}, path cty.Path) diag.Diagnostics {
		return toDiag(validation.StringInSlice(slice, ignoreCase)(input, ""))
	}
}

var catalogItemSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemDescription: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemAbbreviation: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemLabelColor: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemAvailableOnline: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		catalogItemAvailableForPickup: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		catalogItemAvailableElectronically: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		catalogItemCategoryID: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
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
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: catalogItemProductTypeValidate,
			Default:          catalogItemProductTypeRegular,
		},
		catalogItemSkipModifierScreen: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
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

func catalogItemSchemaToObject(input map[string]interface{}) *objects.CatalogItem {
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
		ProductType:             catalogItemProductTypeStrToEnum[input[catalogItemProductType].(string)],
	}

	if idList := input[catalogItemTaxIDs].([]interface{}); len(idList) > 0 {
		result.TaxIDs = make([]string, len(idList))
		for i, id := range idList {
			result.TaxIDs[i] = id.(string)
		}
	}

	if modifierListInfo := input[catalogItemModifierListInfo].([]interface{}); len(modifierListInfo) > 0 {
		result.ModifierListInfo = make([]*objects.CatalogItemModifierListInfo, len(modifierListInfo))
		for i, info := range modifierListInfo {
			result.ModifierListInfo[i] = catalogItemModifierListInfoSchemaToObject(info.(map[string]interface{}))
		}
	}

	if options := input[catalogItemItemOptions].([]interface{}); len(options) > 0 {
		result.ItemOptions = make([]*objects.CatalogItemOptionForItem, len(options))
		for i, option := range options {
			result.ItemOptions[i] = catalogItemOptionForItemSchemaToObject(option.(map[string]interface{}))
		}
	}

	return result
}

func catalogItemObjectToSchema(input *objects.CatalogItem) map[string]interface{} {
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
		catalogItemProductType:             catalogItemProductTypeEnumToStr[input.ProductType],
	}

	if input.TaxIDs != nil {
		result[catalogItemTaxIDs] = input.TaxIDs
	}

	if input.ModifierListInfo != nil {
		resultModifierListInfo := make([]interface{}, len(input.ModifierListInfo))
		for i, info := range input.ModifierListInfo {
			resultModifierListInfo[i] = catalogItemModifierListInfoObjectToSchema(info)
		}

		result[catalogItemModifierListInfo] = resultModifierListInfo
	}

	if input.ItemOptions != nil {
		resultOptions := make([]interface{}, len(input.ItemOptions))
		for i, o := range input.ItemOptions {
			resultOptions[i] = catalogItemOptionForItemObjectToSchema(o)
		}

		result[catalogItemItemOptions] = resultOptions
	}

	return result
}
