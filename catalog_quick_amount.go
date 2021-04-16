package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogQuickAmountAmount  = "amount"
	catalogQuickAmountType    = "type"
	catalogQuickAmountOrdinal = "ordinal"

	catalogQuickAmountTypeManual = "QUICK_AMOUNT_TYPE_MANUAL"
	catalogQuickAmountTypeAuto   = "QUICK_AMOUNT_TYPE_AUTO"
)

var (
	catalogQuickAmountTypeStrToEnum = map[string]objects.CatalogQuickAmountType{
		catalogQuickAmountTypeManual: objects.CatalogQuickAmountTypeManual,
		catalogQuickAmountTypeAuto:   objects.CatalogQuickAmountTypeAuto,
	}

	catalogQuickAmountTypeEnumToStr = map[objects.CatalogQuickAmountType]string{
		objects.CatalogQuickAmountTypeManual: catalogQuickAmountTypeManual,
		objects.CatalogQuickAmountTypeAuto:   catalogQuickAmountTypeAuto,
	}

	catalogQuickAmountTypeValidate = stringInSlice([]string{
		catalogQuickAmountTypeManual,
		catalogQuickAmountTypeAuto,
	}, false)
)

var catalogQuickAmountSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogQuickAmountAmount: &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem:     moneySchema,
		},
		catalogQuickAmountType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: catalogQuickAmountTypeValidate,
		},
		catalogQuickAmountOrdinal: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
	},
}

func catalogQuickAmountSchemaToObject(input map[string]interface{}) *objects.CatalogQuickAmount {
	return &objects.CatalogQuickAmount{
		Type:    catalogQuickAmountTypeStrToEnum[input[catalogQuickAmountType].(string)],
		Ordinal: input[catalogQuickAmountOrdinal].(int),
		Amount:  moneySchemaToObject(input[catalogQuickAmountAmount].(*schema.Set).List()[0].(map[string]interface{})),
	}
}

func catalogQuickAmountObjectToSchema(input *objects.CatalogQuickAmount) map[string]interface{} {
	return map[string]interface{}{
		catalogQuickAmountType:    catalogQuickAmountTypeEnumToStr[input.Type],
		catalogQuickAmountOrdinal: input.Ordinal,
		catalogQuickAmountAmount:  schema.NewSet(schema.HashResource(moneySchema), []interface{}{moneyObjectToSchema(input.Amount)}),
	}
}
