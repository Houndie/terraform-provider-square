package main

/*func upsertConfig(token string) string {
	return providerBlock(token) + `

resource "square_catalog_object" "test_object" {
	type = "ITEM"
	item_data {
		name = "my-item"
	}
}`
}

func addVariation(token string) string {
	return upsertConfig(token) + `

resource "square_catalog_object" "test_variation" {
	type = "ITEM_VARIATION"

	item_variation_data {
		item_id = square_catalog_object.test_object.id
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

func TestAccCatalogItemVariation(t *testing.T) {
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
				Config: upsertConfig(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectExists("square_catalog_object.test_object"),
					resource.TestCheckResourceAttr("square_catalog_object.test_object", "type", "ITEM"),
					resource.TestCheckResourceAttr("square_catalog_object.test_object", "item_data.0.name", "my-item"),
					checkCatalogObjectRemote("square_catalog_object.test_object", token),
				),
			},
			{
				Config: addVariation(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectExists("square_catalog_object.test_object"),
					checkCatalogObjectExists("square_catalog_object.test_variation"),
					resource.TestCheckResourceAttr("square_catalog_object.test_variation", "type", "ITEM_VARIATION"),
					resource.TestCheckResourceAttr("square_catalog_object.test_variation", "item_variation_data.0.name", "variation1"),
					resource.TestCheckResourceAttr("square_catalog_object.test_variation", "item_variation_data.0.pricing_type", "FIXED_PRICING"),
					resource.TestCheckResourceAttr("square_catalog_object.test_variation", "item_variation_data.0.price_money.0.amount", "5"),
					resource.TestCheckResourceAttr("square_catalog_object.test_variation", "item_variation_data.0.price_money.0.currency", "USD"),
					checkCatalogObjectRemote("square_catalog_object.test_object", token),
					checkCatalogObjectRemote("square_catalog_object.test_variation", token),
				),
			},
			{
				Config: removedConfig(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectDoesntExist("square_catalog_object.test_object"),
					checkCatalogObjectDoesntExist("square_catalog_object.test_variation"),
				),
			},
		},
	})
}

func discountConfig(token string) string {
	return providerBlock(token) + `

resource "square_catalog_object" "test_discount" {
	type = "DISCOUNT"

	discount_data {
		name = "discount1"
		discount_type = "FIXED_AMOUNT"
		amount_money {
			amount = 5
			currency = "USD"
		}
	}
}`
}

func TestAccCatalogItemDiscount(t *testing.T) {
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
				Config: discountConfig(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectExists("square_catalog_object.test_discount"),
					resource.TestCheckResourceAttr("square_catalog_object.test_discount", "type", "DISCOUNT"),
					resource.TestCheckResourceAttr("square_catalog_object.test_discount", "discount_data.0.name", "discount1"),
					resource.TestCheckResourceAttr("square_catalog_object.test_discount", "discount_data.0.discount_type", "FIXED_AMOUNT"),
					resource.TestCheckResourceAttr("square_catalog_object.test_discount", "discount_data.0.amount_money.0.amount", "5"),
					resource.TestCheckResourceAttr("square_catalog_object.test_discount", "discount_data.0.amount_money.0.currency", "USD"),
					checkCatalogObjectRemote("square_catalog_object.test_discount", token),
				),
			},
			{
				Config: removedConfig(token),
				Check: resource.ComposeTestCheckFunc(
					checkCatalogObjectDoesntExist("square_catalog_object.test_discount"),
				),
			},
		},
	})
}*/
