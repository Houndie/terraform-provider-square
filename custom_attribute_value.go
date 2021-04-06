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
			Required: false,
		},
		catalogCustomAttributeValueNumberValue: &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
		},
		catalogCustomAttributeValueSelectionUIDValues: &schema.Schema{
			Type:     schema.TypeList,
			Required: false,
			Elem:     schema.TypeString,
		},
		catalogCustomAttributeValueStringValue: &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
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
	case objects.CatalogCustomAttributeValueString:
		result[catalogCustomAttributeValueStringValue] = string(t)
	case objects.CatalogCustomAttributeValueSelection:
		result[catalogCustomAttributeValueSelectionUIDValues] = []string(t)
	case objects.CatalogCustomAttributeValueNumber:
		result[catalogCustomAttributeValueNumberValue] = string(t)
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

	if booleanValue, ok := input[catalogCustomAttributeValueBooleanValue]; ok {
		result.Type = objects.CatalogCustomAttributeValueBoolean(booleanValue.(bool))
	} else if numberValue, ok := input[catalogCustomAttributeValueNumberValue]; ok {
		result.Type = objects.CatalogCustomAttributeValueNumber(numberValue.(string))
	} else if selectionValue, ok := input[catalogCustomAttributeValueSelectionUIDValues]; ok {
		result.Type = objects.CatalogCustomAttributeValueSelection(selectionValue.([]string))
	} else if stringValue, ok := input[catalogCustomAttributeValueStringValue]; ok {
		result.Type = objects.CatalogCustomAttributeValueString(stringValue.(string))
	} else {
		return nil, errors.New("no *_type set in schema")
	}

	return result, nil
}

func catalogCustomAttributeValueValidate(input map[string]interface{}) error {
	_, hasBool := input[catalogCustomAttributeValueBooleanValue]
	_, hasString := input[catalogCustomAttributeValueStringValue]
	_, hasNumber := input[catalogCustomAttributeValueNumberValue]
	_, hasSelection := input[catalogCustomAttributeValueSelectionUIDValues]

	if !hasBool && !hasString && !hasNumber && !hasSelection {
		return fmt.Errorf("at least one *_type field must be set")
	}

	return nil
}
