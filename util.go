package main

import (
	"context"
	"fmt"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/catalog"
	"github.com/Houndie/square-go/objects"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceToObject func(*schema.ResourceData) (*objects.CatalogObject, error)

type ObjectToResource func(*objects.CatalogObject, *schema.ResourceData) error

func resourceCatalogUpsert(resourceToObject ResourceToObject, objectToResource ObjectToResource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, ok := m.(*square.Client)
		if !ok {
			return diag.Errorf("unable to create client from interface")
		}

		idempotencyKey, err := uuid.NewV4()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating idempotency key: %w", err))
		}

		object, err := resourceToObject(d)
		if err != nil {
			return diag.FromErr(err)
		}

		res, err := client.Catalog.UpsertObject(ctx, &catalog.UpsertObjectRequest{
			IdempotencyKey: idempotencyKey.String(),
			Object:         object,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error making network call to upsert object: %w", err))
		}

		if err := objectToResource(res.CatalogObject, d); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func resourceCatalogRead(objectToResource ObjectToResource) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, ok := m.(*square.Client)
		if !ok {
			return diag.Errorf("unable to create client from interface")
		}

		res, err := client.Catalog.RetrieveObject(ctx, &catalog.RetrieveObjectRequest{
			ObjectID: d.Id(),
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error making network call to upsert object: %w", err))
		}

		if err := objectToResource(res.Object, d); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func resourceCatalogDelete() func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, ok := m.(*square.Client)
		if !ok {
			return diag.Errorf("unable to create client from interface")
		}

		_, err := client.Catalog.DeleteObject(ctx, &catalog.DeleteObjectRequest{
			ObjectID: d.Id(),
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error making network call to upsert object: %w", err))
		}

		d.SetId("")

		return nil
	}
}
