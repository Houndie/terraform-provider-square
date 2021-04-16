package main

import (
	"errors"
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogCustomAttributeDefinitionNumberConfigPrecision = "precision"
)

var catalogCustomAttributeDefinitionNumberConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogCustomAttributeDefinitionNumberConfigPrecision: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5, //nolint: gomnd
		},
	},
}

func catalogCustomAttributeDefinitionNumberConfigSchemaToObject(input map[string]interface{}) *objects.CatalogCustomAttributeDefinitionNumberConfig {
	precision := input[catalogCustomAttributeDefinitionNumberConfigPrecision].(int)

	return &objects.CatalogCustomAttributeDefinitionNumberConfig{
		Precision: &precision,
	}
}

func catalogCustomAttributeDefinitionNumberConfigObjectToSchema(input *objects.CatalogCustomAttributeDefinitionNumberConfig) map[string]interface{} {
	return map[string]interface{}{
		catalogCustomAttributeDefinitionNumberConfigPrecision: *input.Precision,
	}
}

const (
	catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionName = "name"
	catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionUID  = "uid"
)

var catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionUID: &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

func catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionSchemaToObject(input map[string]interface{}) *objects.CatalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelection {
	return &objects.CatalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelection{
		Name: input[catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionName].(string),
		UID:  input[catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionUID].(string),
	}
}

func catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionObjectToSchema(input *objects.CatalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelection) map[string]interface{} {
	return map[string]interface{}{
		catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionName: input.Name,
		catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionUID:  input.UID,
	}
}

const (
	catalogCustomAttributeDefinitionSelectionConfigAllowedSelections    = "allowed_selections"
	catalogCustomAttributeDefinitionSelectionConfigMaxAllowedSelections = "max_allowed_selections"
)

var catalogCustomAttributeDefinitionSelectionConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogCustomAttributeDefinitionSelectionConfigAllowedSelections: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionSchema,
		},
		catalogCustomAttributeDefinitionSelectionConfigMaxAllowedSelections: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
	},
}

func catalogCustomAttributeDefinitionSelectionConfigSchemaToObject(input map[string]interface{}) *objects.CatalogCustomAttributeDefinitionSelectionConfig {
	maxInt := input[catalogCustomAttributeDefinitionSelectionConfigMaxAllowedSelections].(int)
	result := &objects.CatalogCustomAttributeDefinitionSelectionConfig{
		MaxAllowedSelections: &maxInt,
	}

	if selections := input[catalogCustomAttributeDefinitionSelectionConfigAllowedSelections].([]interface{}); len(selections) > 0 {
		result.AllowedSelections = make([]*objects.CatalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelection, len(selections))
		for i, selection := range selections {
			result.AllowedSelections[i] = catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionSchemaToObject(selection.(map[string]interface{}))
		}
	}

	return result
}

func catalogCustomAttributeDefinitionSelectionConfigObjectToSchema(input *objects.CatalogCustomAttributeDefinitionSelectionConfig) map[string]interface{} {
	result := map[string]interface{}{
		catalogCustomAttributeDefinitionSelectionConfigMaxAllowedSelections: *input.MaxAllowedSelections,
	}

	if input.AllowedSelections != nil {
		selectionsList := make([]interface{}, len(input.AllowedSelections))
		for i, selection := range input.AllowedSelections {
			selectionsList[i] = catalogCustomAttributeDefinitionSelectionConfigCustomAttributeSelectionObjectToSchema(selection)
		}

		result[catalogCustomAttributeDefinitionSelectionConfigAllowedSelections] = selectionsList
	}

	return result
}

const catalogCustomAttributeDefinitionStringConfigEnforceUniqueness = "enforce_uniqueness"

var catalogCustomAttributeDefinitionStringConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogCustomAttributeDefinitionStringConfigEnforceUniqueness: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	},
}

func catalogCustomAttributeDefinitionStringConfigSchemaToObject(input map[string]interface{}) *objects.CatalogCustomAttributeDefinitionStringConfig {
	return &objects.CatalogCustomAttributeDefinitionStringConfig{
		EnforceUniqueness: input[catalogCustomAttributeDefinitionStringConfigEnforceUniqueness].(bool),
	}
}

func catalogCustomAttributeDefinitionStringConfigObjectToSchema(input *objects.CatalogCustomAttributeDefinitionStringConfig) map[string]interface{} {
	return map[string]interface{}{
		catalogCustomAttributeDefinitionStringConfigEnforceUniqueness: input.EnforceUniqueness,
	}
}

const (
	catalogCustomAttributeDefinitionAllowedObjectTypes = "allowed_object_types"
	catalogCustomAttributeDefinitionName               = "name"
	catalogCustomAttributeDefinitionAppVisibility      = "app_visibility"
	catalogCustomAttributeDefinitionDescription        = "description"
	catalogCustomAttributeDefinitionKey                = "key"
	catalogCustomAttributeDefinitionSellerVisibility   = "seller_visibility"
	catalogCustomAttributeDefinitionSourceApplication  = "source_application"
	catalogCustomAttributeDefinitionType               = "type"
	catalogCustomAttributeDefinitionNumberConfig       = "number_config"
	catalogCustomAttributeDefinitionSelectionConfig    = "selection_config"
	catalogCustomAttributeDefinitionStringConfig       = "string_config"

	catalogCustomAttributeDefinitionTypeString    = "STRING"
	catalogCustomAttributeDefinitionTypeBoolean   = "BOOLEAN"
	catalogCustomAttributeDefinitionTypeNumber    = "NUMBER"
	catalogCustomAttributeDefinitionTypeSelection = "SELECTION"

	catalogCustomAttributeDefinitionAppVisibilityHidden          = "APP_VISIBILITY_HIDDEN"
	catalogCustomAttributeDefinitionAppVisibilityReadOnly        = "APP_VISIBILITY_READ_ONLY"
	catalogCustomAttributeDefinitionAppVisibilityReadWriteValues = "APP_VISIBILITY_READ_WRITE_VALUES"

	catalogCustomAttributeDefinitionSellerVisibilityHidden          = "APP_VISIBILITY_HIDDEN"
	catalogCustomAttributeDefinitionSellerVisibilityReadWriteValues = "APP_VISIBILITY_READ_WRITE_VALUES"
)

var (
	catalogCustomAttributeDefinitionTypeValidate = stringInSlice([]string{
		catalogCustomAttributeDefinitionTypeString,
		catalogCustomAttributeDefinitionTypeBoolean,
		catalogCustomAttributeDefinitionTypeNumber,
		catalogCustomAttributeDefinitionTypeSelection,
	}, false)

	catalogCustomAttributeDefinitionAppVisibilityStrToEnum = map[string]objects.CatalogCustomAttributeDefinitionAppVisibility{
		catalogCustomAttributeDefinitionAppVisibilityHidden:          objects.CatalogCustomAttributeDefinitionAppVisibilityHidden,
		catalogCustomAttributeDefinitionAppVisibilityReadOnly:        objects.CatalogCustomAttributeDefinitionAppVisibilityReadOnly,
		catalogCustomAttributeDefinitionAppVisibilityReadWriteValues: objects.CatalogCustomAttributeDefinitionAppVisibilityReadWriteValues,
		"": "",
	}

	catalogCustomAttributeDefinitionAppVisibilityEnumToStr = map[objects.CatalogCustomAttributeDefinitionAppVisibility]string{
		objects.CatalogCustomAttributeDefinitionAppVisibilityHidden:          catalogCustomAttributeDefinitionAppVisibilityHidden,
		objects.CatalogCustomAttributeDefinitionAppVisibilityReadOnly:        catalogCustomAttributeDefinitionAppVisibilityReadOnly,
		objects.CatalogCustomAttributeDefinitionAppVisibilityReadWriteValues: catalogCustomAttributeDefinitionAppVisibilityReadWriteValues,
		"": "",
	}

	catalogCustomAttributeDefinitionAppVisibilityValidate = stringInSlice([]string{
		catalogCustomAttributeDefinitionAppVisibilityHidden,
		catalogCustomAttributeDefinitionAppVisibilityReadOnly,
		catalogCustomAttributeDefinitionAppVisibilityReadWriteValues,
		"",
	}, false)

	catalogCustomAttributeDefinitionSellerVisibilityStrToEnum = map[string]objects.CatalogCustomAttributeDefinitionSellerVisibility{
		catalogCustomAttributeDefinitionSellerVisibilityHidden:          objects.CatalogCustomAttributeDefinitionSellerVisibilityHidden,
		catalogCustomAttributeDefinitionSellerVisibilityReadWriteValues: objects.CatalogCustomAttributeDefinitionSellerVisibilityReadWriteValues,
		"": "",
	}

	catalogCustomAttributeDefinitionSellerVisibilityEnumToStr = map[objects.CatalogCustomAttributeDefinitionSellerVisibility]string{
		objects.CatalogCustomAttributeDefinitionSellerVisibilityHidden:          catalogCustomAttributeDefinitionSellerVisibilityHidden,
		objects.CatalogCustomAttributeDefinitionSellerVisibilityReadWriteValues: catalogCustomAttributeDefinitionSellerVisibilityReadWriteValues,
		"": "",
	}

	catalogCustomAttributeDefinitionSellerVisibilityValidate = stringInSlice([]string{
		catalogCustomAttributeDefinitionSellerVisibilityHidden,
		catalogCustomAttributeDefinitionSellerVisibilityReadWriteValues,
		"",
	}, false)
)

var catalogCustomAttributeDefinitionSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogCustomAttributeDefinitionName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogCustomAttributeDefinitionAppVisibility: &schema.Schema{
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: catalogCustomAttributeDefinitionAppVisibilityValidate,
		},
		catalogCustomAttributeDefinitionDescription: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogCustomAttributeDefinitionKey: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogCustomAttributeDefinitionSellerVisibility: &schema.Schema{
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: catalogCustomAttributeDefinitionSellerVisibilityValidate,
		},
		catalogCustomAttributeDefinitionType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: catalogCustomAttributeDefinitionTypeValidate,
		},
		catalogCustomAttributeDefinitionNumberConfig: &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem:     catalogCustomAttributeDefinitionNumberConfigSchema,
		},
		catalogCustomAttributeDefinitionSelectionConfig: &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem:     catalogCustomAttributeDefinitionSelectionConfigSchema,
		},
		catalogCustomAttributeDefinitionStringConfig: &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem:     catalogCustomAttributeDefinitionStringConfigSchema,
		},
	},
}

func catalogCustomAttributeDefinitionSchemaToObject(input map[string]interface{}) (*objects.CatalogCustomAttributeDefinition, error) {
	result := &objects.CatalogCustomAttributeDefinition{
		Name:             input[catalogCustomAttributeDefinitionName].(string),
		AppVisibility:    catalogCustomAttributeDefinitionAppVisibilityStrToEnum[input[catalogCustomAttributeDefinitionAppVisibility].(string)],
		SellerVisibility: catalogCustomAttributeDefinitionSellerVisibilityStrToEnum[input[catalogCustomAttributeDefinitionSellerVisibility].(string)],
		Description:      input[catalogCustomAttributeDefinitionDescription].(string),
		Key:              input[catalogCustomAttributeDefinitionKey].(string),
	}

	switch input[catalogCustomAttributeDefinitionType].(string) {
	case catalogCustomAttributeDefinitionTypeBoolean:
		result.Type = &objects.CatalogCustomAttributeDefinitionTypeBoolean{}
	case catalogCustomAttributeDefinitionTypeString:
		config := input[catalogCustomAttributeDefinitionStringConfig].(*schema.Set).List()
		if len(config) < 1 {
			return nil, errors.New("string attribute definition set without string config provided")
		}

		result.Type = &objects.CatalogCustomAttributeDefinitionTypeString{
			Config: catalogCustomAttributeDefinitionStringConfigSchemaToObject(config[0].(map[string]interface{})),
		}
	case catalogCustomAttributeDefinitionTypeNumber:
		config := input[catalogCustomAttributeDefinitionNumberConfig].(*schema.Set).List()
		if len(config) < 1 {
			return nil, errors.New("number attribute definition set without number config provided")
		}

		result.Type = &objects.CatalogCustomAttributeDefinitionTypeNumber{
			Config: catalogCustomAttributeDefinitionNumberConfigSchemaToObject(config[0].(map[string]interface{})),
		}
	case catalogCustomAttributeDefinitionTypeSelection:
		config := input[catalogCustomAttributeDefinitionSelectionConfig].(*schema.Set).List()
		if len(config) < 1 {
			return nil, errors.New("selection attribute definition set without selection config provided")
		}

		result.Type = &objects.CatalogCustomAttributeDefinitionTypeSelection{
			Config: catalogCustomAttributeDefinitionSelectionConfigSchemaToObject(config[0].(map[string]interface{})),
		}
	default:
		return nil, fmt.Errorf("unknown type provided: %s", input[catalogCustomAttributeDefinitionType].(string))
	}

	return result, nil
}

func catalogCustomAttributeDefinitionObjectToSchema(input *objects.CatalogCustomAttributeDefinition) (map[string]interface{}, error) {
	result := map[string]interface{}{
		catalogCustomAttributeDefinitionName:             input.Name,
		catalogCustomAttributeDefinitionAppVisibility:    input.AppVisibility,
		catalogCustomAttributeDefinitionSellerVisibility: input.SellerVisibility,
		catalogCustomAttributeDefinitionDescription:      input.Description,
		catalogCustomAttributeDefinitionKey:              input.Key,
	}

	switch t := input.Type.(type) {
	case *objects.CatalogCustomAttributeDefinitionTypeBoolean:
		result[catalogCustomAttributeDefinitionType] = catalogCustomAttributeDefinitionTypeBoolean
	case *objects.CatalogCustomAttributeDefinitionTypeString:
		result[catalogCustomAttributeDefinitionType] = catalogCustomAttributeDefinitionTypeString
		result[catalogCustomAttributeDefinitionStringConfig] = schema.NewSet(schema.HashResource(catalogCustomAttributeDefinitionStringConfigSchema), []interface{}{catalogCustomAttributeDefinitionStringConfigObjectToSchema(t.Config)})
	case *objects.CatalogCustomAttributeDefinitionTypeNumber:
		result[catalogCustomAttributeDefinitionType] = catalogCustomAttributeDefinitionTypeNumber
		result[catalogCustomAttributeDefinitionNumberConfig] = schema.NewSet(schema.HashResource(catalogCustomAttributeDefinitionNumberConfigSchema), []interface{}{catalogCustomAttributeDefinitionNumberConfigObjectToSchema(t.Config)})
	case *objects.CatalogCustomAttributeDefinitionTypeSelection:
		result[catalogCustomAttributeDefinitionType] = catalogCustomAttributeDefinitionTypeSelection
		result[catalogCustomAttributeDefinitionSelectionConfig] = schema.NewSet(schema.HashResource(catalogCustomAttributeDefinitionSelectionConfigSchema), []interface{}{catalogCustomAttributeDefinitionSelectionConfigObjectToSchema(t.Config)})
	default:
		return nil, errors.New("unknown definition type found")
	}

	return result, nil
}
