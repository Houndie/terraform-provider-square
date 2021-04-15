package main

import (
	"errors"
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogItemOptionValueForItemVariationItemOptionID      = "item_option_id"
	catalogItemOptionValueForItemVariationItemOptionValueID = "item_option_value_id"
)

var catalogItemOptionValueForItemVariationSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemOptionValueForItemVariationItemOptionID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemOptionValueForItemVariationItemOptionValueID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func catalogItemOptionValueForItemVariationSchemaToObject(input map[string]interface{}) *objects.CatalogItemOptionValueForItemVariation {
	return &objects.CatalogItemOptionValueForItemVariation{
		ItemOptionID:      input[catalogItemOptionValueForItemVariationItemOptionID].(string),
		ItemOptionValueID: input[catalogItemOptionValueForItemVariationItemOptionValueID].(string),
	}
}

func catalogItemOptionValueForItemVariationObjectToSchema(input *objects.CatalogItemOptionValueForItemVariation) map[string]interface{} {
	return map[string]interface{}{
		catalogItemOptionValueForItemVariationItemOptionID:      input.ItemOptionID,
		catalogItemOptionValueForItemVariationItemOptionValueID: input.ItemOptionValueID,
	}
}

const (
	itemVariationLocationOverridesLocationID              = "location_id"
	itemVariationLocationOverridesPriceMoney              = "price_money"
	itemVariationLocationOverridesPricingType             = "pricing_type"
	itemVariationLocationOverridesTrackInventory          = "track_inventory"
	itemVariationLocationOverridesInventoryAlertType      = "inventory_alert_type"
	itemVariationLocationOverridesInventoryAlertThreshold = "inventory_alert_threshold"

	inventoryAlertTypeNone        = "NONE"
	inventoryAlertTypeLowQuantity = "LOW_QUANTITY"

	catalogPricingTypeFixed    = "FIXED_PRICING"
	catalogPricingTypeVariable = "VARIAIBLE_PRICING"
)

var (
	catalogPricingTypeStrToEnum = map[string]objects.CatalogPricingType{
		catalogPricingTypeFixed:    objects.CatalogPricingTypeFixed,
		catalogPricingTypeVariable: objects.CatalogPricingTypeVariable,
	}

	catalogPricingTypeEnumToStr = map[objects.CatalogPricingType]string{
		objects.CatalogPricingTypeFixed:    catalogPricingTypeFixed,
		objects.CatalogPricingTypeVariable: catalogPricingTypeVariable,
	}

	catalogPricingTypeValidate = stringInSlice([]string{catalogPricingTypeFixed, catalogPricingTypeVariable}, false)

	inventoryAlertType = stringInSlice([]string{inventoryAlertTypeNone, inventoryAlertTypeLowQuantity}, false)
)

var itemVariationLocationOverridesSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		itemVariationLocationOverridesLocationID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		itemVariationLocationOverridesPriceMoney: &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem:     moneySchema,
		},
		itemVariationLocationOverridesPricingType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: catalogPricingTypeValidate,
		},
		itemVariationLocationOverridesTrackInventory: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		itemVariationLocationOverridesInventoryAlertType: &schema.Schema{
			Type:             schema.TypeString,
			Optional:         true,
			Default:          inventoryAlertTypeNone,
			ValidateDiagFunc: catalogPricingTypeValidate,
		},
		itemVariationLocationOverridesInventoryAlertThreshold: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
	},
}

func inventoryAlertSchemaToObject(alert string, threshold *int) (objects.InventoryAlertType, error) {
	switch alert {
	case inventoryAlertTypeNone:
		return &objects.InventoryAlertTypeNone{}, nil
	case inventoryAlertTypeLowQuantity:
		if threshold == nil {
			return nil, errors.New("alert type set to \"low quantity\" but threshold not set")
		}
		return &objects.InventoryAlertTypeLowQuantity{
			Threshold: *threshold,
		}, nil
	}

	return nil, fmt.Errorf("unknown inventory alert type found: %s", alert)

}

func itemVariationLocationOverridesSchemaToObject(input map[string]interface{}) (*objects.ItemVariationLocationOverrides, error) {
	result := &objects.ItemVariationLocationOverrides{
		LocationID:     input[itemVariationLocationOverridesLocationID].(string),
		TrackInventory: input[itemVariationLocationOverridesTrackInventory].(bool),
		PriceMoney:     moneySchemaToObject(input[itemVariationLocationOverridesPriceMoney].([]map[string]interface{})[0]),
		PricingType:    catalogPricingTypeStrToEnum[input[itemVariationLocationOverridesPricingType].(string)],
	}

	var (
		err       error
		threshold *int
	)

	if t, ok := input[itemVariationLocationOverridesInventoryAlertThreshold]; ok {
		tInt := t.(int)
		threshold = &tInt
	}

	if result.InventoryAlertType, err = inventoryAlertSchemaToObject(input[itemVariationLocationOverridesInventoryAlertType].(string), threshold); err != nil {
		return nil, err
	}

	return result, nil
}

func inventoryAlertObjectToSchema(input objects.InventoryAlertType) (string, *int, error) {
	switch t := input.(type) {
	case *objects.InventoryAlertTypeNone:
		return inventoryAlertTypeNone, nil, nil
	case *objects.InventoryAlertTypeLowQuantity:
		return inventoryAlertTypeLowQuantity, &t.Threshold, nil
	default:
		return "", nil, errors.New("unknown inventory alert type found")
	}
}

func itemVariationLocationOverridesObjectToSchema(input *objects.ItemVariationLocationOverrides) (map[string]interface{}, error) {
	result := map[string]interface{}{
		itemVariationLocationOverridesLocationID:     input.LocationID,
		itemVariationLocationOverridesTrackInventory: input.TrackInventory,
		itemVariationLocationOverridesPriceMoney:     []map[string]interface{}{moneyObjectToSchema(input.PriceMoney)},
		itemVariationLocationOverridesPricingType:    catalogPricingTypeEnumToStr[input.PricingType],
	}

	var err error
	if result[itemVariationLocationOverridesInventoryAlertType], result[itemVariationLocationOverridesInventoryAlertThreshold], err = inventoryAlertObjectToSchema(input.InventoryAlertType); err != nil {
		return nil, err
	}

	return result, nil
}

const (
	catalogItemVariationItemID                  = "item_id"
	catalogItemVariationName                    = "name"
	catalogItemVariationSKU                     = "sku"
	catalogItemVariationUPC                     = "upc"
	catalogItemVariationOrdinal                 = "ordinal"
	catalogItemVariationPricingType             = "pricing_type"
	catalogItemVariationPriceMoney              = "price_money"
	catalogItemVariationLocationOverrides       = "location_overrides"
	catalogItemVariationTrackInventory          = "track_inventory"
	catalogItemVariationInventoryAlertType      = "inventory_alert_type"
	catalogItemVariationInventoryAlertThreshold = "inventory_alert_threshold"
	catalogItemVariationUserData                = "user_data"
	catalogItemVariationServiceDuration         = "service_duration"
	catalogItemVariationAvailableForBooking     = "available_for_booking"
	catalogItemVariationItemOptionValues        = "item_option_values"
	catalogItemVariationMeasurementUnitID       = "measurement_unit_id"
	catalogItemVariationTeamMemberIDs           = "team_member_ids"
)

var catalogItemVariationSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogItemVariationItemID: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemVariationName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogItemVariationSKU: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemVariationUPC: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemVariationOrdinal: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		catalogItemVariationPricingType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: catalogPricingTypeValidate,
		},
		catalogItemVariationPriceMoney: &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem:     moneySchema,
		},
		catalogItemVariationLocationOverrides: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     itemVariationLocationOverridesSchema,
		},
		catalogItemVariationTrackInventory: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogItemVariationInventoryAlertType: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  inventoryAlertTypeNone,
		},
		catalogItemVariationInventoryAlertThreshold: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		catalogItemVariationUserData: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemVariationServiceDuration: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		catalogItemVariationAvailableForBooking: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		catalogItemVariationItemOptionValues: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     catalogItemOptionValueForItemVariationSchema,
		},
		catalogItemVariationMeasurementUnitID: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		catalogItemVariationTeamMemberIDs: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Default:  []string{},
		},
	},
}

func catalogItemVariationSchemaToObject(input map[string]interface{}) (*objects.CatalogItemVariation, error) {
	result := &objects.CatalogItemVariation{
		ItemID:              input[catalogItemVariationItemID].(string),
		Name:                input[catalogItemVariationName].(string),
		SKU:                 input[catalogItemVariationSKU].(string),
		UPC:                 input[catalogItemVariationUPC].(string),
		Ordinal:             input[catalogItemVariationOrdinal].(int),
		PriceMoney:          moneySchemaToObject(input[catalogItemVariationPriceMoney].([]map[string]interface{})[0]),
		TrackInventory:      input[catalogItemVariationTrackInventory].(bool),
		UserData:            input[catalogItemVariationUserData].(string),
		ServiceDuration:     input[catalogItemVariationServiceDuration].(int),
		AvailableForBooking: input[catalogItemVariationAvailableForBooking].(bool),
		MeasurementUnitID:   input[catalogItemVariationMeasurementUnitID].(string),
		TeamMemberIDs:       input[catalogItemVariationTeamMemberIDs].([]string),
		PricingType:         catalogPricingTypeStrToEnum[input[catalogItemVariationPricingType].(string)],
	}

	if overrides, ok := input[catalogItemVariationLocationOverrides]; ok {
		overrideSchema := overrides.([]map[string]interface{})
		result.LocationOverrides = make([]*objects.ItemVariationLocationOverrides, len(overrideSchema))
		for i, override := range overrideSchema {
			var err error
			result.LocationOverrides[i], err = itemVariationLocationOverridesSchemaToObject(override)
			if err != nil {
				return nil, fmt.Errorf("error parsing location override: %w", err)
			}
		}
	}

	var (
		err       error
		threshold *int
	)

	if t, ok := input[itemVariationLocationOverridesInventoryAlertThreshold]; ok {
		tInt := t.(int)
		threshold = &tInt
	}

	if result.InventoryAlertType, err = inventoryAlertSchemaToObject(input[itemVariationLocationOverridesInventoryAlertType].(string), threshold); err != nil {
		return nil, err
	}

	if values, ok := input[catalogItemVariationItemOptionValues]; ok {
		valuesSchema := values.([]map[string]interface{})
		result.ItemOptionValues = make([]*objects.CatalogItemOptionValueForItemVariation, len(valuesSchema))
		for i, value := range valuesSchema {
			result.ItemOptionValues[i] = catalogItemOptionValueForItemVariationSchemaToObject(value)
		}
	}

	return result, nil
}

func catalogItemVariationObjectToSchema(input *objects.CatalogItemVariation) (map[string]interface{}, error) {
	result := map[string]interface{}{
		catalogItemVariationItemID:              input.ItemID,
		catalogItemVariationName:                input.Name,
		catalogItemVariationSKU:                 input.SKU,
		catalogItemVariationUPC:                 input.UPC,
		catalogItemVariationOrdinal:             input.Ordinal,
		catalogItemVariationPriceMoney:          []map[string]interface{}{moneyObjectToSchema(input.PriceMoney)},
		catalogItemVariationTrackInventory:      input.TrackInventory,
		catalogItemVariationUserData:            input.UserData,
		catalogItemVariationServiceDuration:     input.ServiceDuration,
		catalogItemVariationAvailableForBooking: input.AvailableForBooking,
		catalogItemVariationMeasurementUnitID:   input.MeasurementUnitID,
		catalogItemVariationTeamMemberIDs:       input.TeamMemberIDs,
		catalogItemVariationPricingType:         catalogPricingTypeEnumToStr[input.PricingType],
	}

	if input.LocationOverrides != nil {
		overrides := make([]map[string]interface{}, len(input.LocationOverrides))
		for i, override := range input.LocationOverrides {
			var err error
			overrides[i], err = itemVariationLocationOverridesObjectToSchema(override)
			if err != nil {
				return nil, fmt.Errorf("error calculating item variation override: %w", err)
			}
		}

		result[catalogItemVariationLocationOverrides] = overrides
	}

	var err error
	if result[itemVariationLocationOverridesInventoryAlertType], result[itemVariationLocationOverridesInventoryAlertThreshold], err = inventoryAlertObjectToSchema(input.InventoryAlertType); err != nil {
		return nil, err
	}

	if input.ItemOptionValues != nil {
		values := make([]map[string]interface{}, len(input.ItemOptionValues))
		for i, value := range input.ItemOptionValues {
			values[i] = catalogItemOptionValueForItemVariationObjectToSchema(value)
		}

		result[catalogItemVariationItemOptionValues] = values
	}

	return result, nil
}
