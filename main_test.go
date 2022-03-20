package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/catalog"
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("sweep_catalog_objects", &resource.Sweeper{
		Name: "sweep_catalog_objects",
		F: func(r string) error {
			token := os.Getenv("TEST_TOKEN")
			if token == "" {
				return fmt.Errorf("Cannot sweep, test token not set")
			}

			client, err := square.NewClient(token, objects.Sandbox, &http.Client{})
			if err != nil {
				return fmt.Errorf("error creating square client in sweeper: %w", err)
			}

			listRes, err := client.Catalog.List(context.Background(), &catalog.ListRequest{})
			if err != nil {
				return fmt.Errorf("error listing catalog items: %w", err)
			}

			deleteIDs := []string{}
			for listRes.Objects.Next() {
				deleteIDs = append(deleteIDs, listRes.Objects.Value().Object.ID)
			}
			if err := listRes.Objects.Error(); err != nil {
				return fmt.Errorf("error iterating through list objects: %w", err)
			}

			_, err = client.Catalog.BatchDelete(context.Background(), &catalog.BatchDeleteRequest{
				ObjectIDs: deleteIDs,
			})
			if err != nil {
				return fmt.Errorf("error batch deleting objects: %w", err)
			}

			return nil
		},
	})
}
