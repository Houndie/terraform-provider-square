package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogQuickAmountsSettingsOption                 = "option"
	catalogQuickAmountsSettingsAmounts                = "amounts"
	catalogQuickAmountsSettingsEligibleForAutoAmounts = "eligible_for_auto_amounts"

	catalogQuickAmountsSettingsOptionDisabled = "DISABLED"
	catalogQuickAmountsSettingsOptionAuto     = "AUTO"
	catalogQuickAmountsSettingsOptionManual   = "MANUAL"
)

var (
	catalogQuickAmountsSettingsOptionStrToEnum = map[string]objects.CatalogQuickAmountsSettingsOption{
		catalogQuickAmountsSettingsOptionDisabled: objects.CatalogQuickAmountsSettingsOptionDisabled,
		catalogQuickAmountsSettingsOptionAuto:     objects.CatalogQuickAmountsSettingsOptionAuto,
		catalogQuickAmountsSettingsOptionManual:   objects.CatalogQuickAmountsSettingsOptionManual,
	}

	catalogQuickAmountsSettingsOptionEnumToStr = map[objects.CatalogQuickAmountsSettingsOption]string{
		objects.CatalogQuickAmountsSettingsOptionDisabled: catalogQuickAmountsSettingsOptionDisabled,
		objects.CatalogQuickAmountsSettingsOptionAuto:     catalogQuickAmountsSettingsOptionAuto,
		objects.CatalogQuickAmountsSettingsOptionManual:   catalogQuickAmountsSettingsOptionManual,
	}

	catalogQuickAmountsSettingsOptionValidate = stringInSlice([]string{
		catalogQuickAmountsSettingsOptionDisabled,
		catalogQuickAmountsSettingsOptionAuto,
		catalogQuickAmountsSettingsOptionManual,
	}, false)
)

var catalogQuickAmountsSettingsSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogQuickAmountsSettingsOption: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: catalogQuickAmountsSettingsOptionValidate,
		},
		catalogQuickAmountsSettingsAmounts: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     catalogQuickAmountSchema,
		},
		catalogQuickAmountsSettingsEligibleForAutoAmounts: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	},
}

func catalogQuickAmountsSettingsSchemaToObject(input map[string]interface{}) *objects.CatalogQuickAmountsSettings {
	result := &objects.CatalogQuickAmountsSettings{
		Option:                 catalogQuickAmountsSettingsOptionStrToEnum[input[catalogQuickAmountsSettingsOption].(string)],
		EligibleForAutoAmounts: input[catalogQuickAmountsSettingsEligibleForAutoAmounts].(bool),
	}

	if amounts, ok := input[catalogQuickAmountsSettingsAmounts]; ok {
		amountsList := amounts.([]map[string]interface{})
		result.Amounts = make([]*objects.CatalogQuickAmount, len(amountsList))
		for i, amount := range amountsList {
			result.Amounts[i] = catalogQuickAmountSchemaToObject(amount)
		}
	}

	return result
}

func catalogQuickAmountsSettingsObjectToSchema(input *objects.CatalogQuickAmountsSettings) map[string]interface{} {
	result := map[string]interface{}{
		catalogQuickAmountsSettingsOption:                 catalogQuickAmountsSettingsOptionEnumToStr[input.Option],
		catalogQuickAmountsSettingsEligibleForAutoAmounts: input.EligibleForAutoAmounts,
	}

	if input.Amounts != nil {
		amounts := make([]map[string]interface{}, len(input.Amounts))
		for i, amount := range input.Amounts {
			amounts[i] = catalogQuickAmountObjectToSchema(amount)
		}
		result[catalogQuickAmountsSettingsAmounts] = amounts
	}

	return result
}