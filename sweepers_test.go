package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"net/http"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/catalog"
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("delete all objects", &resource.Sweeper{
		Name: "delete all objects",
		F: func(region string) error {
			token := os.Getenv("TEST_TOKEN")
			if token == "" {
				return nil
			}

			client, err := square.NewClient(token, objects.Sandbox, &http.Client{
				Timeout: 10 * time.Second,
			})

			if err != nil {
				return fmt.Errorf("error creating client for sweeper")
			}

			res, err := client.Catalog.List(context.Background(), &catalog.ListRequest{})
			if err != nil {
				return fmt.Errorf("error listing catalog objects: %w", err)
			}

			for res.Objects.Next() {
				_, err := client.Catalog.DeleteObject(context.Background(), &catalog.DeleteObjectRequest{
					ObjectID: res.Objects.Value().Object.ID,
				})
				if err != nil {
					return fmt.Errorf("error deleting object with ID %s: %w", res.Objects.Value().Object.ID, err)
				}
			}

			if err := res.Objects.Error(); err != nil {
				return fmt.Errorf("error iterating through listed catalog objects: %w", err)
			}

			return nil
		},
	})
}
