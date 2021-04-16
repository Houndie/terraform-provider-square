package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogTaxName                   = "name"
	catalogTaxCalculationPhase       = "calculation_phase"
	catalogTaxInclusionType          = "inclusion_type"
	catalogTaxPercentage             = "percentage"
	catalogTaxAppliesToCustomAmounts = "applies_to_custom_amounts"
	catalogTaxEnabled                = "enabled"

	taxCalculationPhaseSubtotalPhase = "TAX_SUBTOTAL_PHASE"
	taxCalculationPhaseTotalPhase    = "TAX_TOTAL_PHASE"

	taxInclusionTypeAdditive  = "ADDITIVE"
	taxInclusionTypeInclusive = "INCLUSIVE"
)

var (
	taxCalculationPhaseStrToEnum = map[string]objects.TaxCalculationPhase{
		taxCalculationPhaseSubtotalPhase: objects.TaxCalculationPhaseSubtotalPhase,
		taxCalculationPhaseTotalPhase:    objects.TaxCalculationPhaseTotalPhase,
	}

	taxCalculationPhaseEnumToStr = map[objects.TaxCalculationPhase]string{
		objects.TaxCalculationPhaseSubtotalPhase: taxCalculationPhaseSubtotalPhase,
		objects.TaxCalculationPhaseTotalPhase:    taxCalculationPhaseTotalPhase,
	}

	taxCalculationPhaseValidate = stringInSlice([]string{taxCalculationPhaseSubtotalPhase, taxCalculationPhaseTotalPhase}, false)

	taxInclusionTypeStrToEnum = map[string]objects.TaxInclusionType{
		taxInclusionTypeAdditive:  objects.TaxInclusionTypeAdditive,
		taxInclusionTypeInclusive: objects.TaxInclusionTypeInclusive,
	}

	taxInclusionTypeEnumToStr = map[objects.TaxInclusionType]string{
		objects.TaxInclusionTypeAdditive:  taxInclusionTypeAdditive,
		objects.TaxInclusionTypeInclusive: taxInclusionTypeInclusive,
	}

	taxInclusionTypeValidate = stringInSlice([]string{taxInclusionTypeAdditive, taxInclusionTypeInclusive}, false)
)

var catalogTaxSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogTaxName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogTaxCalculationPhase: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: taxCalculationPhaseValidate,
		},
		catalogTaxInclusionType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: taxInclusionTypeValidate,
		},
		catalogTaxPercentage: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogTaxAppliesToCustomAmounts: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		catalogTaxEnabled: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	},
}

func catalogTaxSchemaToObject(input map[string]interface{}) *objects.CatalogTax {
	return &objects.CatalogTax{
		Name:                   input[catalogTaxName].(string),
		Percentage:             input[catalogTaxPercentage].(string),
		AppliesToCustomAmounts: input[catalogTaxAppliesToCustomAmounts].(bool),
		Enabled:                input[catalogTaxEnabled].(bool),
		CalculationPhase:       taxCalculationPhaseStrToEnum[input[catalogTaxCalculationPhase].(string)],
		InclusionType:          taxInclusionTypeStrToEnum[input[catalogTaxInclusionType].(string)],
	}
}

func catalogTaxObjectToSchema(input *objects.CatalogTax) map[string]interface{} {
	return map[string]interface{}{
		catalogTaxName:                   input.Name,
		catalogTaxPercentage:             input.Percentage,
		catalogTaxAppliesToCustomAmounts: input.AppliesToCustomAmounts,
		catalogTaxEnabled:                input.Enabled,
		catalogTaxCalculationPhase:       taxCalculationPhaseEnumToStr[input.CalculationPhase],
		catalogTaxInclusionType:          taxInclusionTypeEnumToStr[input.InclusionType],
	}
}
