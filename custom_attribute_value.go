package main

import (
	"errors"
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogCustomAttributeValueCustomAttributeDefinitionID = "custom_attribute_definition_id"
	catalogCustomAttributeValueKey                         = "key"
	catalogCustomAttributeValueType                        = "type"
	catalogCustomAttributeValueName                        = "name"
	catalogCustomAttributeValueBooleanValue                = "boolean_value"
	catalogCustomAttributeValueNumberValue                 = "number_value"
	catalogCustomAttributeValueSelectionUIDValues          = "selection_uid_values"
	catalogCustomAttributeValueStringValue                 = "string_value"
)

var catalogCustomAttributeValueSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogCustomAttributeValueCustomAttributeDefinitionID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogCustomAttributeValueKey: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogCustomAttributeValueName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogCustomAttributeValueBooleanValue: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		catalogCustomAttributeValueNumberValue: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogCustomAttributeValueSelectionUIDValues: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		catalogCustomAttributeValueStringValue: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogCustomAttributeValueType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: catalogCustomAttributeDefinitionTypeValidate,
		},
	},
}

func catalogCustomAttributeValueObjectToSchema(input *objects.CatalogCustomAttributeValue) (map[string]interface{}, error) {
	result := map[string]interface{}{
		catalogCustomAttributeValueCustomAttributeDefinitionID: input.CustomAttributeDefinitionID,
		catalogCustomAttributeValueKey:                         input.Key,
		catalogCustomAttributeValueName:                        input.Name,
	}

	switch t := input.Type.(type) {
	case objects.CatalogCustomAttributeValueBoolean:
		result[catalogCustomAttributeValueBooleanValue] = bool(t)
		result[catalogCustomAttributeValueType] = catalogCustomAttributeDefinitionTypeBoolean
	case objects.CatalogCustomAttributeValueString:
		result[catalogCustomAttributeValueStringValue] = string(t)
		result[catalogCustomAttributeValueType] = catalogCustomAttributeDefinitionTypeString
	case objects.CatalogCustomAttributeValueSelection:
		result[catalogCustomAttributeValueSelectionUIDValues] = []string(t)
		result[catalogCustomAttributeValueType] = catalogCustomAttributeDefinitionTypeSelection
	case objects.CatalogCustomAttributeValueNumber:
		result[catalogCustomAttributeValueNumberValue] = string(t)
		result[catalogCustomAttributeValueType] = catalogCustomAttributeDefinitionTypeNumber
	default:
		return nil, fmt.Errorf("no Type found on input")
	}

	return result, nil
}

func catalogCustomAttributeValueSchemaToObject(input map[string]interface{}) (*objects.CatalogCustomAttributeValue, error) {
	result := &objects.CatalogCustomAttributeValue{
		CustomAttributeDefinitionID: input[catalogCustomAttributeValueCustomAttributeDefinitionID].(string),
		Key:                         input[catalogCustomAttributeValueKey].(string),
		Name:                        input[catalogCustomAttributeValueName].(string),
	}

	switch input[catalogCustomAttributeValueType].(string) {
	case catalogCustomAttributeDefinitionTypeBoolean:
		result.Type = objects.CatalogCustomAttributeValueBoolean(input[catalogCustomAttributeValueBooleanValue].(bool))
	case catalogCustomAttributeDefinitionTypeNumber:
		result.Type = objects.CatalogCustomAttributeValueNumber(input[catalogCustomAttributeValueNumberValue].(string))
	case catalogCustomAttributeDefinitionTypeSelection:
		result.Type = objects.CatalogCustomAttributeValueSelection(input[catalogCustomAttributeValueSelectionUIDValues].([]string))
	case catalogCustomAttributeDefinitionTypeString:
		result.Type = objects.CatalogCustomAttributeValueString(input[catalogCustomAttributeValueStringValue].(string))
	default:
		return nil, errors.New("no *_type set in schema")
	}

	return result, nil
}
