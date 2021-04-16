package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogSubscriptionPlanName   = "name"
	catalogSubscriptionPlanPhases = "phases"
)

var catalogSubscriptionPlanSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		catalogSubscriptionPlanName: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		catalogSubscriptionPlanPhases: &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     subscriptionPhaseSchema,
		},
	},
}

func catalogSubscriptionPlanSchemaToObject(input map[string]interface{}) *objects.CatalogSubscriptionPlan {
	result := &objects.CatalogSubscriptionPlan{
		Name: input[catalogSubscriptionPlanName].(string),
	}

	phases := input[catalogSubscriptionPlanPhases].([]interface{})
	result.Phases = make([]*objects.SubscriptionPhase, len(phases))

	for i, phase := range phases {
		result.Phases[i] = subscriptionPhaseSchemaToObject(phase.(map[string]interface{}))
	}

	return result
}

func catalogSubscriptionPlanObjectToSchema(input *objects.CatalogSubscriptionPlan) map[string]interface{} {
	result := map[string]interface{}{
		catalogSubscriptionPlanName: input.Name,
	}

	phases := make([]interface{}, len(input.Phases))
	for i, phase := range input.Phases {
		phases[i] = subscriptionPhaseObjectToSchema(phase)
	}

	result[catalogSubscriptionPlanPhases] = phases

	return result
}
