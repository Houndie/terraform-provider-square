package main

import (
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var variationSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"item_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"pricing_type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"amount": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	},
}

func resourceCatalogItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"variation": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     variationSchema,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		CreateContext: resourceCatalogUpsert(catalogItemResourceToObject, catalogItemObjectToResource),
		ReadContext:   resourceCatalogRead(catalogItemObjectToResource),
		UpdateContext: resourceCatalogUpsert(catalogItemResourceToObject, catalogItemObjectToResource),
		DeleteContext: resourceCatalogDelete(),
	}
}

func catalogItemResourceToObject(d *schema.ResourceData) (*objects.CatalogObject, error) {
	id := d.Id()
	if id == "" {
		id = "#id"
	}

	dVariations := d.Get("variation").(*schema.Set)
	variations := make([]*objects.CatalogObject, dVariations.Len())

	for i, v := range dVariations.List() {
		mv := v.(map[string]interface{})

		vid := mv["id"].(string)
		if vid == "" {
			uid, err := uuid.NewV4()
			if err != nil {
				return nil, fmt.Errorf("error generating uuid for new variation: %w", err)
			}

			vid = "#" + uid.String()
		}

		var (
			pricingType objects.CatalogPricingType
			money       *objects.Money
		)

		switch mv["pricing_type"].(string) {
		case "FIXED_PRICING":
			pricingType = objects.CatalogPricingTypeFixed
			money = &objects.Money{
				Amount:   mv["amount"].(int),
				Currency: "USD",
			}
		case "VARIABLE_PRICING":
			pricingType = objects.CatalogPricingTypeVariable
		}

		variations[i] = &objects.CatalogObject{
			ID: vid,
			Type: &objects.CatalogItemVariation{
				ItemID:      id,
				Name:        mv["name"].(string),
				PricingType: pricingType,
				PriceMoney:  money,
			},
		}
	}

	return &objects.CatalogObject{
		ID: id,
		Type: &objects.CatalogItem{
			Name:       d.Get("name").(string),
			Variations: variations,
		},
		Version: d.Get("version").(int),
	}, nil
}

func catalogItemObjectToResource(o *objects.CatalogObject, d *schema.ResourceData) error {
	d.SetId(o.ID)

	item, ok := o.Type.(*objects.CatalogItem)
	if !ok {
		return fmt.Errorf("catalog object is not a catalog item")
	}

	if err := d.Set("name", item.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}

	if len(item.Variations) < 1 {
		return fmt.Errorf("expected at least one item variation")
	}

	variations := make([]interface{}, len(item.Variations))

	for i, vo := range item.Variations {
		v, ok := vo.Type.(*objects.CatalogItemVariation)
		if !ok {
			return fmt.Errorf("catalog object is not a catalog item variation")
		}

		var amount int
		if v.PricingType == objects.CatalogPricingTypeFixed {
			amount = v.PriceMoney.Amount
		}

		variations[i] = map[string]interface{}{
			"id":           vo.ID,
			"name":         v.Name,
			"pricing_type": string(v.PricingType),
			"amount":       amount,
		}
	}

	if err := d.Set("variation", schema.NewSet(schema.HashResource(variationSchema), variations)); err != nil {
		return fmt.Errorf("error setting variations: %w", err)
	}

	if err := d.Set("version", o.Version); err != nil {
		return fmt.Errorf("error setting version: %w", err)
	}

	return nil
}
