package main

import (
	"fmt"

	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCatalogDiscount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"percentage": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"amount": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		CreateContext: resourceCatalogUpsert(catalogDiscountResourceToObject, catalogDiscountObjectToResource),
		ReadContext:   resourceCatalogRead(catalogDiscountObjectToResource),
		UpdateContext: resourceCatalogUpsert(catalogDiscountResourceToObject, catalogDiscountObjectToResource),
		DeleteContext: resourceCatalogDelete(),
	}
}

const (
	catalogDiscountFixedAmount        = "FIXED_AMOUNT"
	catalogDiscountVariableAmount     = "VARIABLE_AMOUNT"
	catalogDiscountFixedPercentage    = "FIXED_PERCENTAGE"
	catalogDiscountVariablePercentage = "VARIABLE_PERCENTAGE"
)

func catalogDiscountResourceToObject(d *schema.ResourceData) (*objects.CatalogObject, error) {
	id := d.Id()
	if id == "" {
		id = "#id"
	}

	var discountType objects.CatalogDiscountType

	switch d.Get("type") {
	case catalogDiscountFixedPercentage:
		discountType = &objects.CatalogDiscountFixedPercentage{
			Percentage: d.Get("percentage").(string),
		}
	case catalogDiscountVariablePercentage:
		discountType = &objects.CatalogDiscountVariablePercentage{
			Percentage: d.Get("percentage").(string),
		}
	case catalogDiscountFixedAmount:
		discountType = &objects.CatalogDiscountFixedAmount{
			AmountMoney: &objects.Money{
				Amount:   d.Get("amount").(int),
				Currency: "USD",
			},
		}
	case catalogDiscountVariableAmount:
		discountType = &objects.CatalogDiscountVariableAmount{
			AmountMoney: &objects.Money{
				Amount:   d.Get("amount").(int),
				Currency: "USD",
			},
		}
	}

	return &objects.CatalogObject{
		ID: id,
		Type: &objects.CatalogDiscount{
			Name:         d.Get("name").(string),
			DiscountType: discountType,
		},
	}, nil
}

func catalogDiscountObjectToResource(o *objects.CatalogObject, d *schema.ResourceData) error {
	d.SetId(o.ID)

	discount, ok := o.Type.(*objects.CatalogDiscount)
	if !ok {
		return fmt.Errorf("catalog object is not a catalog discount")
	}

	if err := d.Set("name", discount.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}

	switch t := discount.DiscountType.(type) {
	case *objects.CatalogDiscountFixedPercentage:
		if err := d.Set("type", catalogDiscountFixedPercentage); err != nil {
			return fmt.Errorf("error setting fixed percentage type")
		}

		if err := d.Set("percentage", t.Percentage); err != nil {
			return fmt.Errorf("error setting fixed percentage amount")
		}
	case *objects.CatalogDiscountVariablePercentage:
		if err := d.Set("type", catalogDiscountVariablePercentage); err != nil {
			return fmt.Errorf("error setting variable percentage type")
		}

		if err := d.Set("percentage", t.Percentage); err != nil {
			return fmt.Errorf("error setting variable percentage amount")
		}
	case *objects.CatalogDiscountFixedAmount:
		if err := d.Set("type", catalogDiscountFixedAmount); err != nil {
			return fmt.Errorf("error setting fixed amount type")
		}

		if err := d.Set("amount", t.AmountMoney.Amount); err != nil {
			return fmt.Errorf("error setting fixed amount amount")
		}
	case *objects.CatalogDiscountVariableAmount:
		if err := d.Set("type", catalogDiscountVariableAmount); err != nil {
			return fmt.Errorf("error setting variable amount type")
		}

		if err := d.Set("amount", t.AmountMoney.Amount); err != nil {
			return fmt.Errorf("error setting variable amount amount")
		}
	}

	return nil
}
