package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/catalog"
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func providerBlock(token string) string {
	return fmt.Sprintf(`
provider "test-provider" {
	access_token = "%s"
	environment = "sandbox"
}`, token)
}

func upsertConfig(token string) string {
	return providerBlock(token) + `

resource "catalog_object" "test_object" {
	provider = "test-provider"
	type = "ITEM"
	item_data {
		name = "my-item"
	}
}`
}

func addVariation(token string) string {
	return upsertConfig(token) + `

resource "catalog_object" "test_variation" {
	provider = "test-provider"
	type = "ITEM_VARIATION"

	item_variation_data {
		item_id = catalog_object.test_object.id
		name = "variation1"
		pricing_type = "FIXED_PRICING"

		price_money {
			amount = 5
			currency = "USD"
		}
	}
}`
}

func removedConfig(token string) string {
	return providerBlock(token)
}

func TestAccCatalogObject(t *testing.T) {
	t.Parallel()

	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Log("Test skipped as TEST_TOKEN not set")
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"test-provider": func() (*schema.Provider, error) { return Provider(), nil }, //nolint:unparam
		},
		Steps: []resource.TestStep{
			{
				Config: upsertConfig(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectExists("catalog_object.test_object"),
					resource.TestCheckResourceAttr("catalog_object.test_object", "type", "ITEM"),
					resource.TestCheckResourceAttr("catalog_object.test_object", "item_data.0.name", "my-item"),
					checkCatalogObjectRemote("catalog_object.test_object", token),
				),
			},
			{
				Config: addVariation(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectExists("catalog_object.test_object"),
					checkCatalogObjectExists("catalog_object.test_variation"),
					resource.TestCheckResourceAttr("catalog_object.test_variation", "type", "ITEM_VARIATION"),
					resource.TestCheckResourceAttr("catalog_object.test_variation", "item_variation_data.0.name", "variation1"),
					resource.TestCheckResourceAttr("catalog_object.test_variation", "item_variation_data.0.pricing_type", "FIXED_PRICING"),
					resource.TestCheckResourceAttr("catalog_object.test_variation", "item_variation_data.0.price_money.0.amount", "5"),
					resource.TestCheckResourceAttr("catalog_object.test_variation", "item_variation_data.0.price_money.0.currency", "USD"),
					checkCatalogObjectRemote("catalog_object.test_object", token),
					checkCatalogObjectRemote("catalog_object.test_variation", token),
				),
			},
			{
				Config: removedConfig(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectDoesntExist("catalog_object.test_object"),
					checkCatalogObjectDoesntExist("catalog_object.test_variation"),
				),
			},
		},
	})
}

func checkCatalogObjectExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}

		return nil
	}
}

func checkCatalogObjectDoesntExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("Found: %s", resourceName)
		}

		return nil
	}
}

func checkCatalogObjectRemote(resourceName, apiKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client, err := square.NewClient(apiKey, objects.Sandbox, &http.Client{
			Timeout: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error creating square client: %w", err)
		}

		_, err = client.Catalog.RetrieveObject(context.Background(), &catalog.RetrieveObjectRequest{
			ObjectID: rs.Primary.ID,
		})

		if err != nil {
			return fmt.Errorf("error retrieving remote object: %w", err)
		}

		return nil
	}
}
