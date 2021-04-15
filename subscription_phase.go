package main

import (
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	subscriptionPhaseCadence             = "cadence"
	subscriptionPhaseRecurringPriceMoney = "recurring_price_money"
	subscriptionPhaseOrdinal             = "ordinal"
	subscriptionPhasePeriods             = "periods"
	subscriptionPhaseUID                 = "uid"

	subscriptionCadenceDaily           = "DAILY"
	subscriptionCadenceWeekly          = "WEEKLY"
	subscriptionCadenceEveryTwoWeeks   = "EVERY_TWO_WEEKS"
	subscriptionCadenceThirtyDays      = "THIRTY_DAYS"
	subscriptionCadenceSixtyDays       = "SIXTY_DAYS"
	subscriptionCadenceNinetyDays      = "NINETY_DAYS"
	subscriptionCadenceMonthly         = "MONTHLY"
	subscriptionCadenceEveryTwoMonths  = "EVERY_TWO_MONTHS"
	subscriptionCadenceQuarterly       = "QUARTERLY"
	subscriptionCadenceEveryFourMonths = "EVERY_FOUR_MONTHS"
	subscriptionCadenceEverySixMonths  = "EVERY_SIX_MONTHS"
	subscriptionCadenceAnnual          = "ANNUAL"
	subscriptionCadenceEveryTwoYears   = "EVERY_TWO_YEARS"
)

var (
	subscriptionCadenceStrToEnum = map[string]objects.SubscriptionCadence{
		subscriptionCadenceDaily:           objects.SubscriptionCadenceDaily,
		subscriptionCadenceWeekly:          objects.SubscriptionCadenceWeekly,
		subscriptionCadenceEveryTwoWeeks:   objects.SubscriptionCadenceEveryTwoWeeks,
		subscriptionCadenceThirtyDays:      objects.SubscriptionCadenceThirtyDays,
		subscriptionCadenceSixtyDays:       objects.SubscriptionCadenceSixtyDays,
		subscriptionCadenceNinetyDays:      objects.SubscriptionCadenceNinetyDays,
		subscriptionCadenceMonthly:         objects.SubscriptionCadenceMonthly,
		subscriptionCadenceEveryTwoMonths:  objects.SubscriptionCadenceEveryTwoMonths,
		subscriptionCadenceQuarterly:       objects.SubscriptionCadenceQuarterly,
		subscriptionCadenceEveryFourMonths: objects.SubscriptionCadenceEveryFourMonths,
		subscriptionCadenceEverySixMonths:  objects.SubscriptionCadenceEverySixMonths,
		subscriptionCadenceAnnual:          objects.SubscriptionCadenceAnnual,
		subscriptionCadenceEveryTwoYears:   objects.SubscriptionCadenceEveryTwoYears,
	}

	subscriptionCadenceEnumToStr = map[objects.SubscriptionCadence]string{
		objects.SubscriptionCadenceDaily:           subscriptionCadenceDaily,
		objects.SubscriptionCadenceWeekly:          subscriptionCadenceWeekly,
		objects.SubscriptionCadenceEveryTwoWeeks:   subscriptionCadenceEveryTwoWeeks,
		objects.SubscriptionCadenceThirtyDays:      subscriptionCadenceThirtyDays,
		objects.SubscriptionCadenceSixtyDays:       subscriptionCadenceSixtyDays,
		objects.SubscriptionCadenceNinetyDays:      subscriptionCadenceNinetyDays,
		objects.SubscriptionCadenceMonthly:         subscriptionCadenceMonthly,
		objects.SubscriptionCadenceEveryTwoMonths:  subscriptionCadenceEveryTwoMonths,
		objects.SubscriptionCadenceQuarterly:       subscriptionCadenceQuarterly,
		objects.SubscriptionCadenceEveryFourMonths: subscriptionCadenceEveryFourMonths,
		objects.SubscriptionCadenceEverySixMonths:  subscriptionCadenceEverySixMonths,
		objects.SubscriptionCadenceAnnual:          subscriptionCadenceAnnual,
		objects.SubscriptionCadenceEveryTwoYears:   subscriptionCadenceEveryTwoYears,
	}

	subscriptionCadenceValidate = stringInSlice([]string{
		subscriptionCadenceDaily,
		subscriptionCadenceWeekly,
		subscriptionCadenceEveryTwoWeeks,
		subscriptionCadenceThirtyDays,
		subscriptionCadenceSixtyDays,
		subscriptionCadenceNinetyDays,
		subscriptionCadenceMonthly,
		subscriptionCadenceEveryTwoMonths,
		subscriptionCadenceQuarterly,
		subscriptionCadenceEveryFourMonths,
		subscriptionCadenceEverySixMonths,
		subscriptionCadenceAnnual,
		subscriptionCadenceEveryTwoYears,
	}, false)
)

var subscriptionPhaseSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		subscriptionPhaseCadence: &schema.Schema{
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: subscriptionCadenceValidate,
		},
		subscriptionPhaseRecurringPriceMoney: &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem:     moneySchema,
		},
		subscriptionPhaseOrdinal: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		subscriptionPhasePeriods: &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		subscriptionPhaseUID: &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

func subscriptionPhaseSchemaToObject(input map[string]interface{}) *objects.SubscriptionPhase {
	return &objects.SubscriptionPhase{
		Cadence:             subscriptionCadenceStrToEnum[input[subscriptionPhaseCadence].(string)],
		RecurringPriceMoney: moneySchemaToObject(input[subscriptionPhaseRecurringPriceMoney].([]map[string]interface{})[0]),
		Ordinal:             input[subscriptionPhaseOrdinal].(int),
		Periods:             input[subscriptionPhasePeriods].(int),
		UID:                 input[subscriptionPhaseUID].(string),
	}
}

func subscriptionPhaseObjectToSchema(input *objects.SubscriptionPhase) map[string]interface{} {
	return map[string]interface{}{
		subscriptionPhaseCadence:             subscriptionCadenceEnumToStr[input.Cadence],
		subscriptionPhaseRecurringPriceMoney: []map[string]interface{}{moneyObjectToSchema(input.RecurringPriceMoney)},
		subscriptionPhaseOrdinal:             input.Ordinal,
		subscriptionPhasePeriods:             input.Periods,
		subscriptionPhaseUID:                 input.UID,
	}
}
