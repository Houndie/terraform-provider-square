package main

import (
	"errors"
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	measurementUnitCustomName         = "name"
	measurementUnitCustomAbbreviation = "abbreviation"
)

var measurementUnitCustomSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		measurementUnitCustomName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		measurementUnitCustomAbbreviation: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func measurementUnitCustomSchemaToObject(input map[string]interface{}) *objects.MeasurementUnitCustom {
	return &objects.MeasurementUnitCustom{
		Name:         input[measurementUnitCustomName].(string),
		Abbreviation: input[measurementUnitCustomAbbreviation].(string),
	}
}

func measurementUnitCustomObjectToSchema(input *objects.MeasurementUnitCustom) map[string]interface{} {
	return map[string]interface{}{
		catalogProductSetName:             input.Name,
		measurementUnitCustomAbbreviation: input.Abbreviation,
	}
}

const (
	measurementUnitCustomUnit  = "custom_unit"
	measurementUnitAreaUnit    = "area_unit"
	measurementUnitLengthUnit  = "length_unit"
	measurementUnitVolumeUnit  = "volume_unit"
	measurementUnitWeightUnit  = "weight_unit"
	measurementUnitGenericUnit = "generic_unit"
	measurementUnitTimeUnit    = "time_unit"
	measurementUnitType        = "type"

	measurementUnitTypeCustom  = "TYPE_CUSTOM"
	measurementUnitTypeArea    = "TYPE_AREA"
	measurementUnitTypeLength  = "TYPE_LENGTH"
	measurementUnitTypeVolume  = "TYPE_VOLUME"
	measurementUnitTypeWeight  = "TYPE_WEIGHT"
	measurementUnitTypeGeneric = "TYPE_GENERIC"
	measurementUnitTypeTime    = "TYPE_TIME"
)

var measurementUnitTypeValidate = stringInSlice([]string{
	measurementUnitTypeCustom,
	measurementUnitTypeArea,
	measurementUnitTypeLength,
	measurementUnitTypeVolume,
	measurementUnitTypeWeight,
	measurementUnitTypeGeneric,
	measurementUnitTypeTime,
}, false)

var measurementUnitSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		measurementUnitCustomUnit: &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem:     measurementUnitCustomSchema,
		},
		measurementUnitAreaUnit: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		measurementUnitLengthUnit: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		measurementUnitVolumeUnit: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		measurementUnitWeightUnit: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		measurementUnitGenericUnit: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		measurementUnitTimeUnit: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		measurementUnitType: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: measurementUnitTypeValidate,
		},
	},
}

func measurementUnitSchemaToObject(input map[string]interface{}) (*objects.MeasurementUnit, error) {
	result := &objects.MeasurementUnit{}

	switch input[measurementUnitType].(string) {
	case measurementUnitTypeCustom:
		result.Type = measurementUnitCustomSchemaToObject(input[measurementUnitCustomUnit].([]map[string]interface{})[0])
	case measurementUnitTypeArea:
		unit := input[measurementUnitAreaUnit].(string)
		if unit == "" {
			return nil, errors.New("unit type set to area, but no area unit provided")
		}

		result.Type = objects.MeasurementUnitArea(unit)
	case measurementUnitTypeLength:
		unit := input[measurementUnitLengthUnit].(string)
		if unit == "" {
			return nil, errors.New("unit type set to length, but no length unit provided")
		}

		result.Type = objects.MeasurementUnitLength(unit)
	case measurementUnitTypeVolume:
		unit := input[measurementUnitVolumeUnit].(string)
		if unit == "" {
			return nil, errors.New("unit type set to volume, but no volume unit provided")
		}

		result.Type = objects.MeasurementUnitVolume(unit)
	case measurementUnitTypeWeight:
		unit := input[measurementUnitWeightUnit].(string)
		if unit == "" {
			return nil, errors.New("unit type set to weight, but no weight unit provided")
		}

		result.Type = objects.MeasurementUnitWeight(unit)
	case measurementUnitTypeGeneric:
		unit := input[measurementUnitGenericUnit].(string)
		if unit == "" {
			return nil, errors.New("unit type set to generic, but no generic unit provided")
		}

		result.Type = objects.MeasurementUnitGeneric(unit)
	case measurementUnitTypeTime:
		unit := input[measurementUnitTimeUnit].(string)
		if unit == "" {
			return nil, errors.New("unit type set to time, but no time unit provided")
		}

		result.Type = objects.MeasurementUnitTime(unit)
	default:
		return nil, fmt.Errorf("unknown type found: %s", input[measurementUnitType].(string))
	}

	return result, nil
}

func measurementUnitObjectToSchema(input *objects.MeasurementUnit) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	switch t := input.Type.(type) {
	case *objects.MeasurementUnitCustom:
		result[measurementUnitType] = measurementUnitTypeCustom
		result[measurementUnitCustomUnit] = []map[string]interface{}{measurementUnitCustomObjectToSchema(t)}
	case objects.MeasurementUnitArea:
		result[measurementUnitType] = measurementUnitTypeArea
		result[measurementUnitAreaUnit] = string(t)
	case objects.MeasurementUnitLength:
		result[measurementUnitType] = measurementUnitTypeLength
		result[measurementUnitLengthUnit] = string(t)
	case objects.MeasurementUnitVolume:
		result[measurementUnitType] = measurementUnitTypeVolume
		result[measurementUnitVolumeUnit] = string(t)
	case objects.MeasurementUnitWeight:
		result[measurementUnitType] = measurementUnitTypeWeight
		result[measurementUnitWeightUnit] = string(t)
	case objects.MeasurementUnitGeneric:
		result[measurementUnitType] = measurementUnitTypeGeneric
		result[measurementUnitGenericUnit] = string(t)
	case objects.MeasurementUnitTime:
		result[measurementUnitType] = measurementUnitTypeTime
		result[measurementUnitTimeUnit] = string(t)
	default:
		return nil, errors.New("unknown unit type found")
	}

	return result, nil
}
