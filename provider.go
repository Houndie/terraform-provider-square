package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ProviderAccessToken = "access_token"
	ProviderEnvironment = "environment"
	ProviderTimeout     = "timeout"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			ProviderAccessToken: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			ProviderEnvironment: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "sandbox",
			},
			ProviderTimeout: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30, //nolint:gomnd
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"catalog_object": resourceCatalogObject(),
		},
		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			var environment objects.Environment
			switch d.Get(ProviderEnvironment).(string) {
			case "production":
				environment = objects.Production
			case "sandbox":
				environment = objects.Sandbox
			default:
				return nil, diag.Errorf("unknown provider environment: %s", d.Get(ProviderEnvironment).(string))
			}

			httpClient := &http.Client{
				Timeout: time.Duration(d.Get(ProviderTimeout).(int)) * time.Second, //nolint:durationcheck
			}

			client, err := square.NewClient(d.Get(ProviderAccessToken).(string), environment, httpClient)
			if err != nil {
				return nil, diag.FromErr(fmt.Errorf("error creating square client: %w", err))
			}

			return client, nil
		},
	}
}
