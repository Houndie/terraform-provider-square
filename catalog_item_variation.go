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
		},
	},
}

func itemVariationLocationOverridesSchemaToObject(input map[string]interface{}) (*objects.ItemVariationLocationOverrides, error) {
	result := &objects.ItemVariationLocationOverrides{
		LocationID:     input[itemVariationLocationOverridesLocationID].(string),
		TrackInventory: input[itemVariationLocationOverridesTrackInventory].(bool),
		PriceMoney:     moneySchemaToObject(input[itemVariationLocationOverridesPriceMoney].(*schema.Set).List()[0].(map[string]interface{})),
		PricingType:    catalogPricingTypeStrToEnum[input[itemVariationLocationOverridesPricingType].(string)],
	}

	switch input[itemVariationLocationOverridesInventoryAlertType].(string) {
	case inventoryAlertTypeNone:
		result.InventoryAlertType = &objects.InventoryAlertTypeNone{}
	case inventoryAlertTypeLowQuantity:
		result.InventoryAlertType = &objects.InventoryAlertTypeLowQuantity{
			Threshold: input[itemVariationLocationOverridesInventoryAlertThreshold].(int),
		}
	default:
		return nil, errors.New("unknown inventory alert type found")
	}

	return result, nil
}

func itemVariationLocationOverridesObjectToSchema(input *objects.ItemVariationLocationOverrides) (map[string]interface{}, error) {
	result := map[string]interface{}{
		itemVariationLocationOverridesLocationID:     input.LocationID,
		itemVariationLocationOverridesTrackInventory: input.TrackInventory,
		itemVariationLocationOverridesPriceMoney:     schema.NewSet(schema.HashResource(moneySchema), []interface{}{moneyObjectToSchema(input.PriceMoney)}),
		itemVariationLocationOverridesPricingType:    catalogPricingTypeEnumToStr[input.PricingType],
	}

	switch t := input.InventoryAlertType.(type) {
	case *objects.InventoryAlertTypeNone:
		result[itemVariationLocationOverridesInventoryAlertType] = inventoryAlertTypeNone
	case *objects.InventoryAlertTypeLowQuantity:
		result[itemVariationLocationOverridesInventoryAlertType] = inventoryAlertTypeLowQuantity
		result[itemVariationLocationOverridesInventoryAlertThreshold] = t.Threshold
	default:
		return nil, errors.New("unknown inventory alert type found")
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
		},
		catalogItemVariationUPC: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemVariationOrdinal: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
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
		},
		catalogItemVariationInventoryAlertType: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  inventoryAlertTypeNone,
		},
		catalogItemVariationInventoryAlertThreshold: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		catalogItemVariationUserData: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemVariationServiceDuration: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		catalogItemVariationAvailableForBooking: &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		catalogItemVariationItemOptionValues: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     catalogItemOptionValueForItemVariationSchema,
		},
		catalogItemVariationMeasurementUnitID: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		catalogItemVariationTeamMemberIDs: &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
		PriceMoney:          moneySchemaToObject(input[catalogItemVariationPriceMoney].(*schema.Set).List()[0].(map[string]interface{})),
		TrackInventory:      input[catalogItemVariationTrackInventory].(bool),
		UserData:            input[catalogItemVariationUserData].(string),
		ServiceDuration:     input[catalogItemVariationServiceDuration].(int),
		AvailableForBooking: input[catalogItemVariationAvailableForBooking].(bool),
		MeasurementUnitID:   input[catalogItemVariationMeasurementUnitID].(string),
		PricingType:         catalogPricingTypeStrToEnum[input[catalogItemVariationPricingType].(string)],
	}

	if ids := input[catalogItemVariationTeamMemberIDs].([]interface{}); len(ids) > 0 {
		result.TeamMemberIDs = make([]string, len(ids))
		for i, id := range ids {
			result.TeamMemberIDs[i] = id.(string)
		}
	}

	if overrides := input[catalogItemVariationLocationOverrides].([]interface{}); len(overrides) > 0 {
		result.LocationOverrides = make([]*objects.ItemVariationLocationOverrides, len(overrides))
		for i, override := range overrides {
			var err error
			result.LocationOverrides[i], err = itemVariationLocationOverridesSchemaToObject(override.(map[string]interface{}))
			if err != nil {
				return nil, fmt.Errorf("error parsing location override: %w", err)
			}
		}
	}

	switch input[itemVariationLocationOverridesInventoryAlertType].(string) {
	case inventoryAlertTypeNone:
		result.InventoryAlertType = &objects.InventoryAlertTypeNone{}
	case inventoryAlertTypeLowQuantity:
		result.InventoryAlertType = &objects.InventoryAlertTypeLowQuantity{
			Threshold: input[itemVariationLocationOverridesInventoryAlertThreshold].(int),
		}
	default:
		return nil, errors.New("unknown inventory alert type found")
	}

	if values := input[catalogItemVariationItemOptionValues].([]interface{}); len(values) > 0 {
		result.ItemOptionValues = make([]*objects.CatalogItemOptionValueForItemVariation, len(values))
		for i, value := range values {
			result.ItemOptionValues[i] = catalogItemOptionValueForItemVariationSchemaToObject(value.(map[string]interface{}))
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
		catalogItemVariationPriceMoney:          schema.NewSet(schema.HashResource(moneySchema), []interface{}{moneyObjectToSchema(input.PriceMoney)}),
		catalogItemVariationTrackInventory:      input.TrackInventory,
		catalogItemVariationUserData:            input.UserData,
		catalogItemVariationServiceDuration:     input.ServiceDuration,
		catalogItemVariationAvailableForBooking: input.AvailableForBooking,
		catalogItemVariationMeasurementUnitID:   input.MeasurementUnitID,
		catalogItemVariationPricingType:         catalogPricingTypeEnumToStr[input.PricingType],
	}

	if input.TeamMemberIDs != nil {
		result[catalogItemVariationTeamMemberIDs] = input.TeamMemberIDs
	}

	if input.LocationOverrides != nil {
		overrides := make([]interface{}, len(input.LocationOverrides))
		for i, override := range input.LocationOverrides {
			var err error
			overrides[i], err = itemVariationLocationOverridesObjectToSchema(override)
			if err != nil {
				return nil, fmt.Errorf("error calculating item variation override: %w", err)
			}
		}

		result[catalogItemVariationLocationOverrides] = schema.NewSet(schema.HashResource(itemVariationLocationOverridesSchema), overrides)
	}

	switch t := input.InventoryAlertType.(type) {
	case *objects.InventoryAlertTypeNone:
		result[itemVariationLocationOverridesInventoryAlertType] = inventoryAlertTypeNone
	case *objects.InventoryAlertTypeLowQuantity:
		result[itemVariationLocationOverridesInventoryAlertType] = inventoryAlertTypeLowQuantity
		result[itemVariationLocationOverridesInventoryAlertThreshold] = t.Threshold
	default:
		return nil, errors.New("unknown inventory alert type found")
	}

	if input.ItemOptionValues != nil {
		values := make([]interface{}, len(input.ItemOptionValues))
		for i, value := range input.ItemOptionValues {
			values[i] = catalogItemOptionValueForItemVariationObjectToSchema(value)
		}

		result[catalogItemVariationItemOptionValues] = schema.NewSet(schema.HashResource(catalogItemOptionValueForItemVariationSchema), values)
	}

	return result, nil
}
