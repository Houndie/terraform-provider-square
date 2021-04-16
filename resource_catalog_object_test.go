package main

import (
	"fmt"
	"os"
	"testing"

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
				Check:  checkCatalogObjectExists("catalog_object.test_object"),
			},
			{
				Config: removedConfig(token),
				Check:  checkCatalogObjectDoesntExist("catalog_object.test_object"),
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
