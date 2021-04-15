package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogModifierListName          = "name"
	catalogModifierListOrdinal       = "ordinal"
	catalogModifierListSelectionType = "selection_type"

	catalogModifierListSelectionTypeSingle   = "SINGLE"
	catalogModifierListSelectionTypeMultiple = "MULTIPLE"
)

var (
	catalogModifierListSelectionTypeStrToEnum = map[string]objects.CatalogModifierListSelectionType{
		catalogModifierListSelectionTypeSingle:   objects.CatalogModifierListSelectionTypeSingle,
		catalogModifierListSelectionTypeMultiple: objects.CatalogModifierListSelectionTypeMultiple,
	}

	catalogModifierListSelectionTypeEnumToStr = map[objects.CatalogModifierListSelectionType]string{
		objects.CatalogModifierListSelectionTypeSingle:   catalogModifierListSelectionTypeSingle,
		objects.CatalogModifierListSelectionTypeMultiple: catalogModifierListSelectionTypeMultiple,
	}

	catalogModifierListSelectionTypeValidate = stringInSlice([]string{catalogModifierListSelectionTypeSingle, catalogModifierListSelectionTypeMultiple}, false)
)

var catalogModifierListSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogModifierListName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogModifierListOrdinal: &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		catalogModifierListSelectionType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: catalogModifierListSelectionTypeValidate,
		},
	},
}

func catalogModifierListSchemaToObject(input map[string]interface{}) *objects.CatalogModifierList {
	return &objects.CatalogModifierList{
		Name:          input[catalogModifierListName].(string),
		Ordinal:       input[catalogModifierListOrdinal].(int),
		SelectionType: catalogModifierListSelectionTypeStrToEnum[input[catalogModifierListSelectionType].(string)],
	}
}

func catalogModifierListObjectToSchema(input *objects.CatalogModifierList) map[string]interface{} {
	return map[string]interface{}{
		catalogModifierListName:          input.Name,
		catalogModifierListOrdinal:       input.Ordinal,
		catalogModifierListSelectionType: catalogModifierListSelectionTypeEnumToStr[input.SelectionType],
	}
}
