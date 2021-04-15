package main

import (
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogMeasurementUnitMeasurementUnit = "measurement_unit"
	catalogMeasurementUnitPrecision       = "precision"
)

var catalogMeasurementUnitSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogMeasurementUnitMeasurementUnit: &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem:     measurementUnitSchema,
		},
		catalogMeasurementUnitPrecision: &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
	},
}

func catalogMeasurementUnitSchemaToObject(input map[string]interface{}) (*objects.CatalogMeasurementUnit, error) {
	measurementUnit, err := measurementUnitSchemaToObject(input[catalogMeasurementUnitMeasurementUnit].([]map[string]interface{})[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing mesurement unit: %w", err)
	}

	return &objects.CatalogMeasurementUnit{
		MeasurementUnit: measurementUnit,
		Precision:       input[catalogMeasurementUnitPrecision].(int),
	}, nil
}

func catalogMeasurementUnitObjectToSchema(input *objects.CatalogMeasurementUnit) (map[string]interface{}, error) {
	measurementUnit, err := measurementUnitObjectToSchema(input.MeasurementUnit)
	if err != nil {
		return nil, fmt.Errorf("error parsing measurement unit: %w", err)
	}

	return map[string]interface{}{
		catalogMeasurementUnitMeasurementUnit: []map[string]interface{}{measurementUnit},
		catalogMeasurementUnitPrecision:       input.Precision,
	}, nil
}
