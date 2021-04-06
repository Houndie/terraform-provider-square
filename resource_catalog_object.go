package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/catalog"
	"github.com/Houndie/square-go/objects"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

const (
	catalogObjectType                          = "type"
	catalogObjectAbsentAtLocationIDs           = "absent_at_location_ids"
	catalogObjectCatalogV1IDs                  = "catalog_v1_ids"
	catalogObjectCategoryData                  = "category_data"
	catalogObjectCustomAttributeDefinitionData = "custom_attribute_definition_data"
	catalogObjectCustomAttributeValues         = "custom_attribute_values"
	catalogObjectDiscountData                  = "discount_data"
	catalogObjectImageData                     = "image_data"
	catalogObjectImageID                       = "image_id"
	catalogObjectItemData                      = "item_data"
	catalogObjectItemOptionData                = "item_option_data"
	catalogObjectItemOptionValueData           = "item_option_value_data"
	catalogObjectItemVariationData             = "item_variation_data"
	catalogObjectMeasurementUnitData           = "measurement_unit_data"
	catalogObjectModifierData                  = "modifier_data"
	catalogObjectModifierListData              = "modifier_list_data"
	catalogObjectPresentAtAllLocations         = "present_at_all_locations"
	catalogObjectPresentAtLocationIDs          = "present_at_location_ids"
	catalogObjectProductSetData                = "product_set_data"
	catalogObjectPricingRuleData               = "pricing_rule_data"
	catalogObjectQuickAmountsSettingsData      = "quick_amounts_settings_data"
	catalogObjectSubscriptionPlanData          = "subscription_plan_data"
	catalogObjectTaxData                       = "tax_data"
	catalogObjectTimePeriodData                = "time_period_data"
	catalogObjectUpdatedAt                     = "updated_at"
	catalogObjectVersion                       = "version"

	catalogObjectTypeItem                      = "ITEM"
	catalogObjectTypeImage                     = "IMAGE"
	catalogObjectTypeCategory                  = "CATEGORY"
	catalogObjectTypeItemVariation             = "ITEM_VARIATION"
	catalogObjectTypeTax                       = "TAX"
	catalogObjectTypeDiscount                  = "DISCOUNT"
	catalogObjectTypeModifierList              = "MODIFIER_LIST"
	catalogObjectTypeModifier                  = "MODIFIER"
	catalogObjectTypePricingRule               = "PRICING_RULE"
	catalogObjectTypeProductSet                = "PRODUCT_SET"
	catalogObjectTypeTimePeriod                = "TIME_PERIOD"
	catalogObjectTypeMeasurementUnit           = "MEASUREMENT_UNIT"
	catalogObjectTypeSubscriptionPlan          = "SUBSCRIPTION_PLAN"
	catalogObjectTypeItemOption                = "ITEM_OPTION"
	catalogObjectTypeItemOptionVal             = "ITEM_OPTION_VAL"
	catalogObjectTypeCustomAttributeDefinition = "CUSTOM_ATTRIBUTE_DEFINITION"
	catalogObjectTypeQuickAmountsSettings      = "QUICK_AMOUNTS_SETTINGS"
)

// We don't store the variations struct to prevent drift with the server.
type TerraformCatalogItem struct {
	*objects.CatalogItem
	Variations []*objects.CatalogObject `json:"-"`
}

func resourceCatalogObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCatalogObjectCreate,
		ReadContext:   resourceCatalogObjectRead,
		UpdateContext: resourceCatalogObjectUpdate,
		DeleteContext: resourceCatalogObjectDelete,
		Schema: map[string]*schema.Schema{
			catalogObjectType: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			catalogObjectAbsentAtLocationIDs: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     schema.TypeString,
			},
			catalogObjectCatalogV1IDs: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     catalogV1IDSchema,
			},
			catalogObjectCategoryData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectCustomAttributeDefinitionData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectCustomAttributeValues: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectDiscountData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectImageData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectImageID: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectItemData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectItemOptionData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectItemOptionValueData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectItemVariationData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectMeasurementUnitData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectModifierData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectModifierListData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectPresentAtAllLocations: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			catalogObjectPresentAtLocationIDs: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Default:  true,
			},
			catalogObjectProductSetData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectPricingRuleData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectQuickAmountsSettingsData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectSubscriptionPlanData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectTaxData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectTimePeriodData: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			catalogObjectUpdatedAt: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			catalogObjectVersion: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceCatalogObjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(square.Client)
	if !ok {
		return diag.Errorf("unable to create client from interface")
	}

	idempotencyKey, err := uuid.NewV4()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating idempotency key: %w", err))
	}

	object, err := resourceToCatalogObject(d)
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

	if err := catalogObjectToResource(res.CatalogObject, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCatalogObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(square.Client)
	if !ok {
		return diag.Errorf("unable to create client from interface")
	}

	res, err := client.Catalog.RetrieveObject(ctx, &catalog.RetrieveObjectRequest{
		ObjectID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error making network call to retrieve object: %w", err))
	}

	if err := catalogObjectToResource(res.Object, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCatalogObjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(square.Client)
	if !ok {
		return diag.Errorf("unable to create client from interface")
	}

	idempotencyKey, err := uuid.NewV4()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating idempotency key: %w", err))
	}

	object, err := resourceToCatalogObject(d)
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

	if err := catalogObjectToResource(res.CatalogObject, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCatalogObjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(square.Client)
	if !ok {
		return diag.Errorf("unable to create client from interface")
	}

	_, err := client.Catalog.DeleteObject(ctx, &catalog.DeleteObjectRequest{
		ObjectID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error making network call to delete object: %w", err))
	}

	d.SetId("")

	return nil
}

func resourceToCatalogObject(d *schema.ResourceData) (*objects.CatalogObject, error) {
	id := d.Id()
	if id == "" {
		id = "#id"
	}

	object := &objects.CatalogObject{
		ID: id,
	}

	if version := d.Get(catalogObjectVersion); version != nil {
		object.Version = version.(int)
	}

	if err := decode(d.Get(catalogObjectCustomAttributeValues), &object.CustomAttributeValues); err != nil {
		return nil, fmt.Errorf("error parsing custom attribute values: %w", err)
	}

	if err := decode(d.Get(catalogObjectCatalogV1IDs), &object.CatalogV1IDs); err != nil {
		return nil, fmt.Errorf("error parsing catalog v1 IDs: %w", err)
	}

	if presentAtAllLocations := d.Get(catalogObjectPresentAtAllLocations); presentAtAllLocations != nil {
		object.PresentAtAllLocations = presentAtAllLocations.(bool)
	}

	if presentAtLocationIDs := d.Get(catalogObjectPresentAtLocationIDs); presentAtLocationIDs != nil {
		list := presentAtLocationIDs.([]interface{})

		object.PresentAtLocationIDs = make([]string, len(list))

		for i, location := range list {
			var ok bool
			if object.PresentAtLocationIDs[i], ok = location.(string); !ok {
				return nil, fmt.Errorf("cannot convert present at location id to string")
			}
		}
	}

	if absentAtLocationIDs := d.Get(catalogObjectAbsentAtLocationIDs); absentAtLocationIDs != nil {
		list := absentAtLocationIDs.([]interface{})

		object.AbsentAtLocationIDs = make([]string, len(list))

		for i, location := range list {
			var ok bool
			if object.AbsentAtLocationIDs[i], ok = location.(string); !ok {
				return nil, fmt.Errorf("cannot convert absent at location id to string")
			}
		}
	}

	if imageID := d.Get(catalogObjectImageID); imageID != nil {
		object.ImageID = imageID.(string)
	}

	switch d.Get(catalogObjectType) {
	case catalogObjectTypeItem:
		t := &TerraformCatalogItem{}
		if err := parseItemData(d, catalogObjectTypeItem, catalogObjectItemData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeCategory:
		t := &objects.CatalogCategory{}
		if err := parseItemData(d, catalogObjectTypeCategory, catalogObjectCategoryData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeItemVariation:
		t := &objects.CatalogItemVariation{}
		if err := parseItemData(d, catalogObjectTypeItemVariation, catalogObjectItemVariationData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeTax:
		t := &objects.CatalogTax{}
		if err := parseItemData(d, catalogObjectTypeTax, catalogObjectTaxData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeDiscount:
		t := &objects.CatalogDiscount{}
		if err := parseItemData(d, catalogObjectTypeDiscount, catalogObjectDiscountData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeModifierList:
		t := &objects.CatalogModifierList{}
		if err := parseItemData(d, catalogObjectTypeModifierList, catalogObjectModifierListData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeModifier:
		t := &objects.CatalogModifier{}
		if err := parseItemData(d, catalogObjectTypeModifier, catalogObjectModifierData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypePricingRule:
		t := &objects.CatalogPricingRule{}
		if err := parseItemData(d, catalogObjectTypePricingRule, catalogObjectPricingRuleData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeProductSet:
		t := &objects.CatalogProductSet{}
		if err := parseItemData(d, catalogObjectTypeProductSet, catalogObjectProductSetData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeTimePeriod:
		t := &objects.CatalogTimePeriod{}
		if err := parseItemData(d, catalogObjectTypeTimePeriod, catalogObjectTimePeriodData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeMeasurementUnit:
		t := &objects.CatalogMeasurementUnit{}
		if err := parseItemData(d, catalogObjectTypeMeasurementUnit, catalogObjectMeasurementUnitData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeSubscriptionPlan:
		t := &objects.CatalogSubscriptionPlan{}
		if err := parseItemData(d, catalogObjectTypeSubscriptionPlan, catalogObjectSubscriptionPlanData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeItemOption:
		t := &objects.CatalogItemOption{}
		if err := parseItemData(d, catalogObjectTypeItemOption, catalogObjectItemOptionData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeItemOptionVal:
		t := &objects.CatalogItemOptionValue{}
		if err := parseItemData(d, catalogObjectTypeItemOptionVal, catalogObjectItemOptionValueData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeCustomAttributeDefinition:
		t := &objects.CatalogCustomAttributeDefinition{}
		if err := parseItemData(d, catalogObjectTypeCustomAttributeDefinition, catalogObjectCustomAttributeDefinitionData, t); err != nil {
			return nil, err
		}

		object.Type = t
	case catalogObjectTypeQuickAmountsSettings:
		t := &objects.CatalogQuickAmountsSettings{}
		if err := parseItemData(d, catalogObjectTypeQuickAmountsSettings, catalogObjectQuickAmountsSettingsData, t); err != nil {
			return nil, err
		}

		object.Type = t
	}

	return object, nil
}

func setError(key string, err error) error {
	return fmt.Errorf("error setting %s: %w", key, err)
}

func catalogObjectToResource(object *objects.CatalogObject, d *schema.ResourceData) error {
	d.SetId(object.ID)

	if err := d.Set(catalogObjectUpdatedAt, object.UpdatedAt.Format(time.RFC3339)); err != nil {
		return setError(catalogObjectUpdatedAt, err)
	}

	if object.Version != 0 {
		if err := d.Set(catalogObjectVersion, object.Version); err != nil {
			return setError(catalogObjectVersion, err)
		}
	} else {
		if err := d.Set(catalogObjectVersion, nil); err != nil {
			return setError(catalogObjectVersion, err)
		}
	}

	var catalogCustomAttributeValues map[string]interface{}
	if object.CustomAttributeValues != nil {
		catalogCustomAttributeValues = map[string]interface{}{}
		for key, value := range object.CustomAttributeValues {
			catalogCustomAttributeValues[key] = value
		}
	}

	if err := d.Set(catalogObjectCustomAttributeValues, catalogCustomAttributeValues); err != nil {
		return setError(catalogObjectCustomAttributeValues, err)
	}

	var catalogV1IDs []interface{}
	if object.CatalogV1IDs != nil {
		catalogV1IDs := []interface{}{}
		if err := decode(object.CatalogV1IDs, &catalogV1IDs); err != nil {
			return err
		}
	}

	if err := d.Set(catalogObjectCatalogV1IDs, catalogV1IDs); err != nil {
		return setError(catalogObjectCatalogV1IDs, err)
	}

	if err := d.Set(catalogObjectPresentAtAllLocations, object.PresentAtAllLocations); err != nil {
		return setError(catalogObjectPresentAtAllLocations, err)
	}

	var presentAtLocationIDs []interface{}
	if object.PresentAtLocationIDs != nil {
		presentAtLocationIDs = make([]interface{}, len(object.PresentAtLocationIDs))
		for i, locationID := range object.PresentAtLocationIDs {
			presentAtLocationIDs[i] = locationID
		}
	}

	if err := d.Set(catalogObjectPresentAtLocationIDs, presentAtLocationIDs); err != nil {
		return setError(catalogObjectPresentAtLocationIDs, err)
	}

	var absentAtLocationIDs []interface{}
	if object.AbsentAtLocationIDs != nil {
		absentAtLocationIDs = make([]interface{}, len(object.AbsentAtLocationIDs))
		for i, locationID := range object.AbsentAtLocationIDs {
			absentAtLocationIDs[i] = locationID
		}
	}

	if err := d.Set(catalogObjectAbsentAtLocationIDs, absentAtLocationIDs); err != nil {
		return setError(catalogObjectAbsentAtLocationIDs, err)
	}

	var imageID interface{}
	if object.ImageID != "" {
		imageID = object.ImageID
	}

	if err := d.Set(catalogObjectImageID, imageID); err != nil {
		return setError(catalogObjectImageID, err)
	}

	switch t := object.Type.(type) {
	case *objects.CatalogItem:
		myCatalogItem := TerraformCatalogItem{CatalogItem: t}
		if err := setCatalogObjectType(d, myCatalogItem, catalogObjectTypeItem, catalogObjectItemData); err != nil {
			return err
		}
	case *objects.CatalogImage:
		if err := setCatalogObjectType(d, t, catalogObjectTypeImage, catalogObjectImageData); err != nil {
			return err
		}
	case *objects.CatalogCategory:
		if err := setCatalogObjectType(d, t, catalogObjectTypeCategory, catalogObjectCategoryData); err != nil {
			return err
		}
	case *objects.CatalogItemVariation:
		if err := setCatalogObjectType(d, t, catalogObjectTypeItemVariation, catalogObjectItemVariationData); err != nil {
			return err
		}
	case *objects.CatalogTax:
		if err := setCatalogObjectType(d, t, catalogObjectTypeTax, catalogObjectTaxData); err != nil {
			return err
		}
	case *objects.CatalogDiscount:
		if err := setCatalogObjectType(d, t, catalogObjectTypeDiscount, catalogObjectDiscountData); err != nil {
			return err
		}
	case *objects.CatalogModifierList:
		if err := setCatalogObjectType(d, t, catalogObjectTypeModifierList, catalogObjectModifierListData); err != nil {
			return err
		}
	case *objects.CatalogModifier:
		if err := setCatalogObjectType(d, t, catalogObjectTypeModifier, catalogObjectModifierData); err != nil {
			return err
		}
	case *objects.CatalogPricingRule:
		if err := setCatalogObjectType(d, t, catalogObjectTypePricingRule, catalogObjectPricingRuleData); err != nil {
			return err
		}
	case *objects.CatalogProductSet:
		if err := setCatalogObjectType(d, t, catalogObjectTypeProductSet, catalogObjectProductSetData); err != nil {
			return err
		}
	case *objects.CatalogTimePeriod:
		if err := setCatalogObjectType(d, t, catalogObjectTypeTimePeriod, catalogObjectTimePeriodData); err != nil {
			return err
		}
	case *objects.CatalogMeasurementUnit:
		if err := setCatalogObjectType(d, t, catalogObjectTypeMeasurementUnit, catalogObjectMeasurementUnitData); err != nil {
			return err
		}
	case *objects.CatalogSubscriptionPlan:
		if err := setCatalogObjectType(d, t, catalogObjectTypeSubscriptionPlan, catalogObjectSubscriptionPlanData); err != nil {
			return err
		}
	case *objects.CatalogItemOption:
		if err := setCatalogObjectType(d, t, catalogObjectTypeItemOption, catalogObjectItemOptionData); err != nil {
			return err
		}
	case *objects.CatalogItemOptionValue:
		if err := setCatalogObjectType(d, t, catalogObjectTypeItemOptionVal, catalogObjectItemOptionValueData); err != nil {
			return err
		}
	case *objects.CatalogCustomAttributeDefinition:
		if err := setCatalogObjectType(d, t, catalogObjectTypeCustomAttributeDefinition, catalogObjectCustomAttributeDefinitionData); err != nil {
			return err
		}
	case *objects.CatalogQuickAmountsSettings:
		if err := setCatalogObjectType(d, t, catalogObjectTypeQuickAmountsSettings, catalogObjectQuickAmountsSettingsData); err != nil {
			return err
		}
	}

	return nil
}

func setCatalogObjectType(d *schema.ResourceData, t interface{}, typeName, field string) error {
	interfaceType := map[string]interface{}{}
	if err := decode(t, &interfaceType); err != nil {
		return err
	}

	if err := d.Set(catalogObjectType, typeName); err != nil {
		return setError(catalogObjectType, err)
	}

	for _, dataField := range []string{
		catalogObjectTypeItem,
		catalogObjectTypeImage,
		catalogObjectTypeCategory,
		catalogObjectTypeItemVariation,
		catalogObjectTypeTax,
		catalogObjectTypeDiscount,
		catalogObjectTypeModifierList,
		catalogObjectTypeModifier,
		catalogObjectTypePricingRule,
		catalogObjectTypeProductSet,
		catalogObjectTypeTimePeriod,
		catalogObjectTypeMeasurementUnit,
		catalogObjectTypeSubscriptionPlan,
		catalogObjectTypeItemOption,
		catalogObjectTypeItemOptionVal,
		catalogObjectTypeCustomAttributeDefinition,
		catalogObjectTypeQuickAmountsSettings,
	} {
		if dataField == field {
			if err := d.Set(catalogObjectItemData, interfaceType); err != nil {
				return setError(catalogObjectItemData, err)
			}

			continue
		}

		if err := d.Set(catalogObjectItemData, nil); err != nil {
			return setError(catalogObjectItemData, err)
		}
	}

	return nil
}

func parseItemData(d *schema.ResourceData, key, field string, output interface{}) error {
	itemData := d.Get(field)
	if itemData == nil {
		return fmt.Errorf("catalog object type set to %s but no %s found", key, field)
	}

	if err := decode(itemData, output); err != nil {
		return fmt.Errorf("error parsing catalog object %s: %w", field, err)
	}

	return nil
}

func decode(input, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  output,
	})
	if err != nil {
		return fmt.Errorf("error creating new decoder: %w", err)
	}

	if err := decoder.Decode(input); err != nil {
		return fmt.Errorf("error decoding input: %w", err)
	}

	return nil
}
