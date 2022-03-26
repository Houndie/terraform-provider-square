package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/objects"
	"github.com/Houndie/square-go/options"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ProviderAccessToken  = "access_token"
	ProviderEnvironment  = "environment"
	ProviderTimeout      = "timeout"
	ProviderMaxRetryTime = "max_retry_time_seconds"
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
			ProviderMaxRetryTime: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"square_catalog_item":     resourceCatalogItem(),
			"square_catalog_discount": resourceCatalogDiscount(),
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

			o := []options.ClientOption{
				options.WithHTTPClient(&http.Client{
					Timeout: time.Duration(d.Get(ProviderTimeout).(int)) * time.Second, //nolint:durationcheck
				}),
			}

			if t := d.Get(ProviderMaxRetryTime).(int); t != -1 {
				o = append(o, options.WithRateLimit(time.Duration(t)*time.Second))
			}

			client, err := square.NewClient(d.Get(ProviderAccessToken).(string), environment, o...)
			if err != nil {
				return nil, diag.FromErr(fmt.Errorf("error creating square client: %w", err))
			}

			return client, nil
		},
	}
}
