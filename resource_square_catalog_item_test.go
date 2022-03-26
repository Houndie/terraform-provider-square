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

func providerBlock(token string) string {
	return fmt.Sprintf(`
provider "square" {
	access_token = "%s"
	environment = "sandbox"
}

`, token)
}

const catalogItemBlock = `resource "square_catalog_item" "test_item" {
	name = "my-item"

	variation {
		name = "variation1"
		pricing_type = "FIXED_PRICING"
		amount = 500
	}

	variation {
		name = "variation2"
		pricing_type = "VARIABLE_PRICING"
	}
}

`

func TestCatalogItem(t *testing.T) {
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
				Config: providerBlock(token) + catalogItemBlock,
				Check: resource.ComposeTestCheckFunc(
					checkCatalogItem("square_catalog_item.test_item", &objects.CatalogObject{
						Type: &objects.CatalogItem{
							Name: "my-item",
							Variations: []*objects.CatalogObject{
								{
									Type: &objects.CatalogItemVariation{
										Name:        "variation1",
										PricingType: objects.CatalogPricingTypeFixed,
										PriceMoney: &objects.Money{
											Amount: 500,
										},
									},
								},
								{
									Type: &objects.CatalogItemVariation{
										Name:        "variation2",
										PricingType: objects.CatalogPricingTypeVariable,
									},
								},
							},
						},
					}),
					checkCatalogItemRemote("square_catalog_item.test_item", token),
				),
			},
			{
				Config: providerBlock(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogItemDoesntExist("square_catalog_item.test_item"),
					checkCatalogItemDoesntExistRemote("my-item", token),
				),
			},
		},
	})
}

func checkCatalogItem(resourceName string, expected *objects.CatalogObject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}

		return compareCatalogItemToResource(rs.Primary, expected)
	}
}

func compareCatalogItemToResource(d *terraform.InstanceState, s *objects.CatalogObject) error {
	if s.ID != "" && d.ID != s.ID {
		return fmt.Errorf("unexpected id")
	}

	if strings.HasPrefix(d.ID, "#") {
		return fmt.Errorf("no id assigned from server")
	}

	item, ok := s.Type.(*objects.CatalogItem)
	if !ok {
		return fmt.Errorf("provided catalog object is not an item")
	}

	if d.Attributes["name"] != item.Name {
		return fmt.Errorf("unexpected name")
	}

	for _, variation := range item.Variations {
		v, ok := variation.Type.(*objects.CatalogItemVariation)
		if !ok {
			return fmt.Errorf("provided catalog object is not a variation")
		}

		i := 0
		found := false

		for ; i < len(item.Variations); i++ {
			if d.Attributes[fmt.Sprintf("variation.%d.name", i)] == v.Name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unable to find terraform variation matching square variation name")
		}

		vID := d.Attributes[fmt.Sprintf("variation.%d.id", i)]
		if strings.HasPrefix(vID, "#") {
			return fmt.Errorf("no variation id assigned from the server")
		}

		if vID == d.ID {
			return fmt.Errorf("variation id matches item id")
		}

		if variation.ID != "" && vID != variation.ID {
			return fmt.Errorf("unexpected variation id")
		}

		switch d.Attributes[fmt.Sprintf("variation.%d.pricing_type", i)] {
		case "FIXED_PRICING":
			if v.PricingType != objects.CatalogPricingTypeFixed {
				return fmt.Errorf("unexpected pricing type")
			}
		case "VARIABLE_PRICING":
			if v.PricingType != objects.CatalogPricingTypeVariable {
				return fmt.Errorf("unexpected pricing type")
			}
		default:
			return fmt.Errorf("unexpected pricing type")
		}

		if v.PricingType == objects.CatalogPricingTypeFixed {
			amount, err := strconv.Atoi(d.Attributes[fmt.Sprintf("variation.%d.amount", i)])
			if err != nil {
				return fmt.Errorf("terraform amount is not a string")
			}

			if amount != v.PriceMoney.Amount {
				return fmt.Errorf("unexpected amount")
			}
		}
	}

	return nil
}

func checkCatalogItemDoesntExist(resourceName string) resource.TestCheckFunc {
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
func checkCatalogItemDoesntExistRemote(itemName, apiKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := square.NewClient(apiKey, objects.Sandbox, options.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}))
		if err != nil {
			return fmt.Errorf("error creating square client: %w", err)
		}

		res, err := client.Catalog.List(context.Background(), &catalog.ListRequest{
			Types: []objects.CatalogObjectEnumType{objects.CatalogObjectEnumTypeItem},
		})
		if err != nil {
			return fmt.Errorf("error listing all catalog items: %w", err)
		}

		for res.Objects.Next() {
			v := res.Objects.Value()

			item, ok := v.Object.Type.(*objects.CatalogItem)
			if !ok {
				return fmt.Errorf("object is not a catalog item")
			}

			if item.Name == itemName {
				return fmt.Errorf("found catalog item that should be deleted")
			}
		}

		return nil
	}
}

func checkCatalogItemRemote(resourceName, apiKey string) resource.TestCheckFunc {
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

		return compareCatalogItemToResource(rs.Primary, res.Object)
	}
}
