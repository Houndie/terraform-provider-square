package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/catalog"
	"github.com/Houndie/square-go/objects"
	"github.com/Houndie/square-go/options"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const catalogDiscountBlock1 = `resource "square_catalog_discount" "test_discount_1" {
	name = "my-discount1"
	type = "FIXED_AMOUNT"
	amount = 500
}

`

const catalogDiscountBlock2 = `resource "square_catalog_discount" "test_discount_2" {
	name = "my-discount2"
	type = "FIXED_PERCENTAGE"
	percentage = "100.0"
}

`

func TestCatalogDiscount(t *testing.T) {
	t.Parallel()

	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Log("Test skipped as TEST_TOKEN not set")
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"square": func() (*schema.Provider, error) { return Provider(), nil }, //nolint:unparam
		},
		Steps: []resource.TestStep{
			{
				Config: providerBlock(token) + catalogDiscountBlock1,
				Check: resource.ComposeTestCheckFunc(
					checkCatalogDiscount("square_catalog_discount.test_discount_1", &objects.CatalogObject{
						Type: &objects.CatalogDiscount{
							Name: "my-discount1",
							DiscountType: &objects.CatalogDiscountFixedAmount{
								AmountMoney: &objects.Money{
									Amount: 500,
								},
							},
						},
					}),
					checkCatalogDiscountRemote("square_catalog_discount.test_discount_1", token),
				),
			},
			{
				Config: providerBlock(token) + catalogDiscountBlock2,
				Check: resource.ComposeTestCheckFunc(
					checkCatalogDiscountDoesntExist("square_catalog_discount.test_discount_1"),
					checkCatalogDiscountDoesntExistRemote("my-discount1", token),
					checkCatalogDiscount("square_catalog_discount.test_discount_2", &objects.CatalogObject{
						Type: &objects.CatalogDiscount{
							Name: "my-discount2",
							DiscountType: &objects.CatalogDiscountFixedPercentage{
								Percentage: "100.0",
							},
						},
					}),
					checkCatalogDiscountRemote("square_catalog_discount.test_discount_2", token),
				),
			},
			{
				Config: providerBlock(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogDiscountDoesntExist("square_catalog_discount.test_discount_2"),
					checkCatalogDiscountDoesntExistRemote("my-discount2", token),
				),
			},
		},
	})
}

func checkCatalogDiscount(resourceName string, expected *objects.CatalogObject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}

		return compareCatalogDiscountToResource(rs.Primary, expected)
	}
}

func compareCatalogDiscountToResource(d *terraform.InstanceState, s *objects.CatalogObject) error {
	if s.ID != "" && d.ID != s.ID {
		return fmt.Errorf("unexpected id")
	}

	if strings.HasPrefix(d.ID, "#") {
		return fmt.Errorf("no id assigned from server")
	}

	discount, ok := s.Type.(*objects.CatalogDiscount)
	if !ok {
		return fmt.Errorf("provided catalog object is not a discount")
	}

	if d.Attributes["name"] != discount.Name {
		return fmt.Errorf("unexpected name")
	}

	switch t := discount.DiscountType.(type) {
	case *objects.CatalogDiscountFixedPercentage:
		if d.Attributes["type"] != catalogDiscountFixedPercentage {
			return fmt.Errorf("incorrect discount type")
		}

		if d.Attributes["percentage"] != t.Percentage {
			return fmt.Errorf("incorrect percentage")
		}
	case *objects.CatalogDiscountVariablePercentage:
		if d.Attributes["type"] != catalogDiscountVariablePercentage {
			return fmt.Errorf("incorrect discount type")
		}

		if d.Attributes["percentage"] != t.Percentage {
			return fmt.Errorf("incorrect percentage")
		}
	case *objects.CatalogDiscountFixedAmount:
		if d.Attributes["type"] != catalogDiscountFixedAmount {
			return fmt.Errorf("incorrect discount type")
		}

		amount, err := strconv.Atoi(d.Attributes["amount"])
		if err != nil {
			return fmt.Errorf("error converting amount to int: %w", err)
		}

		if amount != t.AmountMoney.Amount {
			return fmt.Errorf("incorrect percentage")
		}
	case *objects.CatalogDiscountVariableAmount:
		if d.Attributes["type"] != catalogDiscountVariableAmount {
			return fmt.Errorf("incorrect discount type")
		}

		amount, err := strconv.Atoi(d.Attributes["amount"])
		if err != nil {
			return fmt.Errorf("error converting amount to int: %w", err)
		}

		if amount != t.AmountMoney.Amount {
			return fmt.Errorf("incorrect percentage")
		}
	}

	return nil
}

func checkCatalogDiscountDoesntExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("Found: %s", resourceName)
		}

		return nil
	}
}

//nolint:dupl
func checkCatalogDiscountDoesntExistRemote(itemName, apiKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := square.NewClient(apiKey, objects.Sandbox, options.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}))
		if err != nil {
			return fmt.Errorf("error creating square client: %w", err)
		}

		res, err := client.Catalog.List(context.Background(), &catalog.ListRequest{
			Types: []objects.CatalogObjectEnumType{objects.CatalogObjectEnumTypeDiscount},
		})
		if err != nil {
			return fmt.Errorf("error listing all catalog discounts: %w", err)
		}

		for res.Objects.Next() {
			v := res.Objects.Value()

			item, ok := v.Object.Type.(*objects.CatalogDiscount)
			if !ok {
				return fmt.Errorf("object is not a catalog discount")
			}

			if item.Name == itemName {
				return fmt.Errorf("found catalog discount that should be deleted")
			}
		}

		return nil
	}
}

func checkCatalogDiscountRemote(resourceName, apiKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client, err := square.NewClient(apiKey, objects.Sandbox, options.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}))
		if err != nil {
			return fmt.Errorf("error creating square client: %w", err)
		}

		res, err := client.Catalog.RetrieveObject(context.Background(), &catalog.RetrieveObjectRequest{
			ObjectID: rs.Primary.ID,
		})

		if err != nil {
			return fmt.Errorf("error retrieving remote object: %w", err)
		}

		return compareCatalogDiscountToResource(rs.Primary, res.Object)
	}
}
