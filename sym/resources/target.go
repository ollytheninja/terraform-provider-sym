package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func Target() *schema.Resource {
	return &schema.Resource{
		Schema:        targetSchema(),
		CreateContext: createTarget,
		ReadContext:   readTarget,
		UpdateContext: updateTarget,
		DeleteContext: deleteTarget,
	}
}

func targetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":           required(schema.TypeString),
		"label":          required(schema.TypeString),
		"integration_id": required(schema.TypeString),
		"settings":       settingsMap(),
	}
}

func createTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	target := client.SymTarget{
		Label:         data.Get("label").(string),
		IntegrationId: data.Get("integration_id").(string),
		Type:          data.Get("type").(string),
		Settings:      getSettings(data),
	}

	id, err := c.Target.Create(target)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create sym target: " + err.Error(),
		})
	} else {
		data.SetId(id)
	}
	return diags
}

func readTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}

func updateTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}
func deleteTarget(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return notYetImplemented
}
