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
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
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
		percentage := d.Get("percentage").(string)
		if percentage == "" {
			return nil, fmt.Errorf("percentage required with a type of %s", catalogDiscountFixedPercentage)
		}

		discountType = &objects.CatalogDiscountFixedPercentage{
			Percentage: percentage,
		}
	case catalogDiscountVariablePercentage:
		percentage := d.Get("percentage").(string)
		if percentage == "" {
			return nil, fmt.Errorf("percentage required with a type of %s", catalogDiscountVariablePercentage)
		}

		discountType = &objects.CatalogDiscountVariablePercentage{
			Percentage: percentage,
		}
	case catalogDiscountFixedAmount:
		amount := d.Get("amount").(int)
		if amount == 0 {
			return nil, fmt.Errorf("amount required with a type of %s", catalogDiscountFixedAmount)
		}

		discountType = &objects.CatalogDiscountFixedAmount{
			AmountMoney: &objects.Money{
				Amount:   amount,
				Currency: "USD",
			},
		}
	case catalogDiscountVariableAmount:
		amount := d.Get("amount").(int)
		if amount == 0 {
			return nil, fmt.Errorf("amount required with a type of %s", catalogDiscountVariableAmount)
		}

		discountType = &objects.CatalogDiscountVariableAmount{
			AmountMoney: &objects.Money{
				Amount:   amount,
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
		Version: d.Get("version").(int),
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

	if err := d.Set("version", o.Version); err != nil {
		return fmt.Errorf("error setting version: %w", err)
	}

	return nil
}
