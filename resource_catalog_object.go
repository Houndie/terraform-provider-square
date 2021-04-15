package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Houndie/square-go"
	"github.com/Houndie/square-go/catalog"
	"github.com/Houndie/square-go/objects"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	catalogObjectType                          = "type"
	catalogObjectAbsentAtLocationIDs           = "absent_at_location_ids"
	catalogObjectCatalogV1IDs                  = "catalog_v1_ids"
	catalogObjectCategoryData                  = "category_data"
	catalogObjectCustomAttributeDefinitionData = "custom_attribute_definition_data"
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

var (
	catalogObjectTypeStrToEnum = map[string]objects.CatalogObjectEnumType{
		catalogObjectTypeItem:                      objects.CatalogObjectEnumTypeItem,
		catalogObjectTypeImage:                     objects.CatalogObjectEnumTypeImage,
		catalogObjectTypeCategory:                  objects.CatalogObjectEnumTypeCategory,
		catalogObjectTypeItemVariation:             objects.CatalogObjectEnumTypeItemVariation,
		catalogObjectTypeTax:                       objects.CatalogObjectEnumTypeTax,
		catalogObjectTypeDiscount:                  objects.CatalogObjectEnumTypeDiscount,
		catalogObjectTypeModifierList:              objects.CatalogObjectEnumTypeModifierList,
		catalogObjectTypeModifier:                  objects.CatalogObjectEnumTypeModifier,
		catalogObjectTypePricingRule:               objects.CatalogObjectEnumTypePricingRule,
		catalogObjectTypeProductSet:                objects.CatalogObjectEnumTypeProductSet,
		catalogObjectTypeTimePeriod:                objects.CatalogObjectEnumTypeTimePeriod,
		catalogObjectTypeMeasurementUnit:           objects.CatalogObjectEnumTypeMeasurementUnit,
		catalogObjectTypeSubscriptionPlan:          objects.CatalogObjectEnumTypeSubscriptionPlan,
		catalogObjectTypeItemOption:                objects.CatalogObjectEnumTypeItemOption,
		catalogObjectTypeItemOptionVal:             objects.CatalogObjectEnumTypeItemOptionVal,
		catalogObjectTypeCustomAttributeDefinition: objects.CatalogObjectEnumTypeCustomAttributeDefinition,
		catalogObjectTypeQuickAmountsSettings:      objects.CatalogObjectEnumTypeQuickAmountsSettings,
	}

	catalogObjectTypeEnumToStr = map[objects.CatalogObjectEnumType]string{
		objects.CatalogObjectEnumTypeItem:                      catalogObjectTypeItem,
		objects.CatalogObjectEnumTypeImage:                     catalogObjectTypeImage,
		objects.CatalogObjectEnumTypeCategory:                  catalogObjectTypeCategory,
		objects.CatalogObjectEnumTypeItemVariation:             catalogObjectTypeItemVariation,
		objects.CatalogObjectEnumTypeTax:                       catalogObjectTypeTax,
		objects.CatalogObjectEnumTypeDiscount:                  catalogObjectTypeDiscount,
		objects.CatalogObjectEnumTypeModifierList:              catalogObjectTypeModifierList,
		objects.CatalogObjectEnumTypeModifier:                  catalogObjectTypeModifier,
		objects.CatalogObjectEnumTypePricingRule:               catalogObjectTypePricingRule,
		objects.CatalogObjectEnumTypeProductSet:                catalogObjectTypeProductSet,
		objects.CatalogObjectEnumTypeTimePeriod:                catalogObjectTypeTimePeriod,
		objects.CatalogObjectEnumTypeMeasurementUnit:           catalogObjectTypeMeasurementUnit,
		objects.CatalogObjectEnumTypeSubscriptionPlan:          catalogObjectTypeSubscriptionPlan,
		objects.CatalogObjectEnumTypeItemOption:                catalogObjectTypeItemOption,
		objects.CatalogObjectEnumTypeItemOptionVal:             catalogObjectTypeItemOptionVal,
		objects.CatalogObjectEnumTypeCustomAttributeDefinition: catalogObjectTypeCustomAttributeDefinition,
		objects.CatalogObjectEnumTypeQuickAmountsSettings:      catalogObjectTypeQuickAmountsSettings,
	}

	catalogObjectTypeValidate = stringInSlice([]string{
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
	}, false)
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
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: catalogObjectTypeValidate,
			},
			catalogObjectAbsentAtLocationIDs: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			catalogObjectCatalogV1IDs: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     catalogV1IDSchema,
			},
			catalogObjectCategoryData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogCategorySchema,
			},
			catalogObjectCustomAttributeDefinitionData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogCustomAttributeDefinitionSchema,
			},
			catalogObjectDiscountData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogDiscountSchema,
			},
			catalogObjectImageData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogImageSchema,
			},
			catalogObjectImageID: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			catalogObjectItemData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogItemSchema,
			},
			catalogObjectItemOptionData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogItemOptionSchema,
			},
			catalogObjectItemOptionValueData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogItemOptionValueSchema,
			},
			catalogObjectItemVariationData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogItemVariationSchema,
			},
			catalogObjectMeasurementUnitData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogMeasurementUnitSchema,
			},
			catalogObjectModifierData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogModifierSchema,
			},
			catalogObjectModifierListData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogModifierListSchema,
			},
			catalogObjectPresentAtAllLocations: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			catalogObjectPresentAtLocationIDs: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			catalogObjectProductSetData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogProductSetSchema,
			},
			catalogObjectPricingRuleData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogPricingRuleSchema,
			},
			catalogObjectQuickAmountsSettingsData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogQuickAmountsSettingsSchema,
			},
			catalogObjectSubscriptionPlanData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogSubscriptionPlanSchema,
			},
			catalogObjectTaxData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogTaxSchema,
			},
			catalogObjectTimePeriodData: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     catalogTimePeriodSchema,
			},
			catalogObjectVersion: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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

	object, err := catalogObjectResourceToObject(d)
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

	if err := catalogObjectObjectToResource(res.CatalogObject, d); err != nil {
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

	if err := catalogObjectObjectToResource(res.Object, d); err != nil {
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

	object, err := catalogObjectResourceToObject(d)
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

	if err := catalogObjectObjectToResource(res.CatalogObject, d); err != nil {
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

func catalogObjectResourceToObject(d *schema.ResourceData) (*objects.CatalogObject, error) {
	result := &objects.CatalogObject{
		Version:               d.Get(catalogObjectVersion).(int),
		PresentAtAllLocations: d.Get(catalogObjectPresentAtAllLocations).(bool),
		ImageID:               d.Get(catalogObjectImageID).(string),
	}

	if id := d.Id(); id != "" {
		result.ID = id
	} else {
		result.ID = "#id"
	}

	if ids, ok := d.GetOk(catalogObjectCatalogV1IDs); ok {
		idList := ids.([]map[string]interface{})
		result.CatalogV1IDs = make([]*objects.CatalogV1ID, len(idList))
		for i, id := range idList {
			result.CatalogV1IDs[i] = catalogV1IDSchemaToObject(id)
		}
	}

	if ids, ok := d.GetOk(catalogObjectPresentAtLocationIDs); ok {
		result.PresentAtLocationIDs = ids.([]string)
	}

	if ids, ok := d.GetOk(catalogObjectAbsentAtLocationIDs); ok {
		result.AbsentAtLocationIDs = ids.([]string)
	}

	switch d.Get(catalogObjectType) {
	case catalogObjectTypeItem:
		data, ok := d.GetOk(catalogObjectItemData)
		if !ok {
			return nil, errors.New("object type item set, but item data not found")
		}

		result.Type = catalogItemSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeCategory:
		data, ok := d.GetOk(catalogObjectCategoryData)
		if !ok {
			return nil, errors.New("object type category set, but category data not found")
		}

		result.Type = catalogCategorySchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeItemVariation:
		data, ok := d.GetOk(catalogObjectItemVariationData)
		if !ok {
			return nil, errors.New("object type item variation set, but item variation data not found")
		}

		var err error
		if result.Type, err = catalogItemVariationSchemaToObject(data.([]map[string]interface{})[0]); err != nil {
			return nil, fmt.Errorf("error parsing item variation: %w", err)
		}
	case catalogObjectTypeTax:
		data, ok := d.GetOk(catalogObjectTaxData)
		if !ok {
			return nil, errors.New("object type tax set, but tax data not found")
		}

		result.Type = catalogTaxSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeDiscount:
		data, ok := d.GetOk(catalogObjectDiscountData)
		if !ok {
			return nil, errors.New("object type discount set, but discount data not found")
		}

		var err error
		if result.Type, err = catalogDiscountSchemaToObject(data.([]map[string]interface{})[0]); err != nil {
			return nil, fmt.Errorf("error parsing discount: %w", err)
		}
	case catalogObjectTypeModifierList:
		data, ok := d.GetOk(catalogObjectModifierListData)
		if !ok {
			return nil, errors.New("object type modifier list set, but modifier list data not found")
		}

		result.Type = catalogModifierListSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeModifier:
		data, ok := d.GetOk(catalogObjectModifierData)
		if !ok {
			return nil, errors.New("object type modifier set, but modifier data not found")
		}

		result.Type = catalogModifierSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeTimePeriod:
		data, ok := d.GetOk(catalogObjectTimePeriodData)
		if !ok {
			return nil, errors.New("object type time period set, but time period data not found")
		}

		result.Type = catalogTimePeriodSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeProductSet:
		data, ok := d.GetOk(catalogObjectProductSetData)
		if !ok {
			return nil, errors.New("object type product set set, but product set data not found")
		}

		var err error
		if result.Type, err = catalogProductSetSchemaToObject(data.([]map[string]interface{})[0]); err != nil {
			return nil, fmt.Errorf("error parsing product set: %w", err)
		}
	case catalogObjectTypePricingRule:
		data, ok := d.GetOk(catalogObjectPricingRuleData)
		if !ok {
			return nil, errors.New("object type pricing rule set, but pricing rule data not found")
		}

		var err error
		if result.Type, err = catalogPricingRuleSchemaToObject(data.([]map[string]interface{})[0]); err != nil {
			return nil, fmt.Errorf("error parsing pricing rule set: %w", err)
		}
	case catalogObjectTypeImage:
		data, ok := d.GetOk(catalogObjectImageData)
		if !ok {
			return nil, errors.New("object type image set, but image data not found")
		}

		result.Type = catalogImageSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeMeasurementUnit:
		data, ok := d.GetOk(catalogObjectMeasurementUnitData)
		if !ok {
			return nil, errors.New("object type measurement unit set, but measurement unit data not found")
		}

		var err error
		if result.Type, err = catalogMeasurementUnitSchemaToObject(data.([]map[string]interface{})[0]); err != nil {
			return nil, fmt.Errorf("error parsing pricing rule set: %w", err)
		}
	case catalogObjectTypeSubscriptionPlan:
		data, ok := d.GetOk(catalogObjectSubscriptionPlanData)
		if !ok {
			return nil, errors.New("object type subscription plan set, but subscription plan data not found")
		}

		result.Type = catalogSubscriptionPlanSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeItemOption:
		data, ok := d.GetOk(catalogObjectItemOptionData)
		if !ok {
			return nil, errors.New("object type item option set, but item option data not found")
		}

		result.Type = catalogItemOptionSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeItemOptionVal:
		data, ok := d.GetOk(catalogObjectItemOptionValueData)
		if !ok {
			return nil, errors.New("object type item option value set, but item option value data not found")
		}

		result.Type = catalogItemOptionValueSchemaToObject(data.([]map[string]interface{})[0])
	case catalogObjectTypeCustomAttributeDefinition:
		data, ok := d.GetOk(catalogObjectCustomAttributeDefinitionData)
		if !ok {
			return nil, errors.New("object type custom attribute definition set, but custom attribute definition data not found")
		}

		var err error
		if result.Type, err = catalogCustomAttributeDefinitionSchemaToObject(data.([]map[string]interface{})[0]); err != nil {
			return nil, fmt.Errorf("error parsing custom attribute definition: %w", err)
		}
	case catalogObjectTypeQuickAmountsSettings:
		data, ok := d.GetOk(catalogObjectQuickAmountsSettingsData)
		if !ok {
			return nil, errors.New("object type quick amounts settings set, but quick amounts settings data not found")
		}

		result.Type = catalogQuickAmountsSettingsSchemaToObject(data.([]map[string]interface{})[0])
	}

	return result, nil
}

func catalogObjectObjectToResource(input *objects.CatalogObject, d *schema.ResourceData) error {
	d.SetId(input.ID)

	if err := d.Set(catalogObjectVersion, input.Version); err != nil {
		return fmt.Errorf("error setting version: %w", err)
	}

	if err := d.Set(catalogObjectPresentAtAllLocations, input.PresentAtAllLocations); err != nil {
		return fmt.Errorf("error setting present at all locations: %w", err)
	}

	if err := d.Set(catalogObjectImageID, input.ImageID); err != nil {
		return fmt.Errorf("error setting image id: %w", err)
	}

	if input.CatalogV1IDs != nil {
		ids := make([]map[string]interface{}, len(input.CatalogV1IDs))
		for i, id := range input.CatalogV1IDs {
			ids[i] = catalogV1IDObjectToSchema(id)
		}

		if err := d.Set(catalogObjectCatalogV1IDs, ids); err != nil {
			return fmt.Errorf("error setting catalog v1 ids: %w", err)
		}
	}

	if input.PresentAtLocationIDs != nil {
		if err := d.Set(catalogObjectPresentAtLocationIDs, input.PresentAtLocationIDs); err != nil {
			return fmt.Errorf("error setting present at location ids: %w", err)
		}
	}

	if input.AbsentAtLocationIDs != nil {
		if err := d.Set(catalogObjectAbsentAtLocationIDs, input.AbsentAtLocationIDs); err != nil {
			return fmt.Errorf("error setting absent at location ids: %w", err)
		}
	}

	switch t := input.Type.(type) {
	case *objects.CatalogItem:
		if err := d.Set(catalogObjectType, catalogObjectTypeItem); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectItemData, []map[string]interface{}{catalogItemObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting item data: %w", err)
		}
	case *objects.CatalogCategory:
		if err := d.Set(catalogObjectType, catalogObjectTypeCategory); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectCategoryData, []map[string]interface{}{catalogCategoryObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting category data: %w", err)
		}
	case *objects.CatalogItemVariation:
		if err := d.Set(catalogObjectType, catalogObjectTypeItemVariation); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		s, err := catalogItemVariationObjectToSchema(t)
		if err != nil {
			return fmt.Errorf("error parsing item variation: %w", err)
		}

		if err := d.Set(catalogObjectItemVariationData, []map[string]interface{}{s}); err != nil {
			return fmt.Errorf("error setting category data: %w", err)
		}
	case *objects.CatalogTax:
		if err := d.Set(catalogObjectType, catalogObjectTypeTax); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectTaxData, []map[string]interface{}{catalogTaxObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting tax data: %w", err)
		}
	case *objects.CatalogDiscount:
		if err := d.Set(catalogObjectType, catalogObjectTypeDiscount); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		s, err := catalogDiscountObjectToSchema(t)
		if err != nil {
			return fmt.Errorf("error parsing discount: %w", err)
		}

		if err := d.Set(catalogObjectDiscountData, []map[string]interface{}{s}); err != nil {
			return fmt.Errorf("error setting discount data: %w", err)
		}
	case *objects.CatalogModifierList:
		if err := d.Set(catalogObjectType, catalogObjectTypeModifierList); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectModifierListData, []map[string]interface{}{catalogModifierListObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting modifier list data: %w", err)
		}
	case *objects.CatalogModifier:
		if err := d.Set(catalogObjectType, catalogObjectTypeModifier); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectModifierData, []map[string]interface{}{catalogModifierObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting modifier data: %w", err)
		}
	case *objects.CatalogTimePeriod:
		if err := d.Set(catalogObjectType, catalogObjectTypeTimePeriod); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectTimePeriodData, []map[string]interface{}{catalogTimePeriodObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting time period data: %w", err)
		}
	case *objects.CatalogProductSet:
		if err := d.Set(catalogObjectType, catalogObjectTypeProductSet); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		s, err := catalogProductSetObjectToSchema(t)
		if err != nil {
			return fmt.Errorf("error parsing product set: %w", err)
		}

		if err := d.Set(catalogObjectProductSetData, []map[string]interface{}{s}); err != nil {
			return fmt.Errorf("error setting product set data: %w", err)
		}
	case *objects.CatalogPricingRule:
		if err := d.Set(catalogObjectType, catalogObjectTypePricingRule); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectPricingRuleData, []map[string]interface{}{catalogPricingRuleObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting pricing rule data: %w", err)
		}
	case *objects.CatalogImage:
		if err := d.Set(catalogObjectType, catalogObjectTypeImage); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectImageData, []map[string]interface{}{catalogImageObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting image data: %w", err)
		}
	case *objects.CatalogMeasurementUnit:
		if err := d.Set(catalogObjectType, catalogObjectTypeMeasurementUnit); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		s, err := catalogMeasurementUnitObjectToSchema(t)
		if err != nil {
			return fmt.Errorf("error parsing measurement unit: %w", err)
		}

		if err := d.Set(catalogObjectMeasurementUnitData, []map[string]interface{}{s}); err != nil {
			return fmt.Errorf("error setting measurement unit data: %w", err)
		}
	case *objects.CatalogSubscriptionPlan:
		if err := d.Set(catalogObjectType, catalogObjectTypeSubscriptionPlan); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectSubscriptionPlanData, []map[string]interface{}{catalogSubscriptionPlanObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting subscription plan data: %w", err)
		}
	case *objects.CatalogItemOption:
		if err := d.Set(catalogObjectType, catalogObjectTypeItemOption); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectItemOptionData, []map[string]interface{}{catalogItemOptionObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting item option data: %w", err)
		}
	case *objects.CatalogItemOptionValue:
		if err := d.Set(catalogObjectType, catalogObjectTypeItemOptionVal); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectItemOptionValueData, []map[string]interface{}{catalogItemOptionValueObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting item option data: %w", err)
		}
	case *objects.CatalogCustomAttributeDefinition:
		if err := d.Set(catalogObjectType, catalogObjectTypeCustomAttributeDefinition); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		s, err := catalogCustomAttributeDefinitionObjectToSchema(t)
		if err != nil {
			return fmt.Errorf("error parsing custom attribute definition: %w", err)
		}

		if err := d.Set(catalogObjectCustomAttributeDefinitionData, []map[string]interface{}{s}); err != nil {
			return fmt.Errorf("error setting custom attribute definition data: %w", err)
		}

	case *objects.CatalogQuickAmountsSettings:
		if err := d.Set(catalogObjectType, catalogObjectTypeQuickAmountsSettings); err != nil {
			return fmt.Errorf("error setting catalog object type: %w", err)
		}

		if err := d.Set(catalogObjectQuickAmountsSettingsData, []map[string]interface{}{catalogQuickAmountsSettingsObjectToSchema(t)}); err != nil {
			return fmt.Errorf("error setting quick amounts settings data: %w", err)
		}
	}

	return nil
}
