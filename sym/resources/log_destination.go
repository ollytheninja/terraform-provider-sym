package resources

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/symopsio/terraform-provider-sym/sym/client"
	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func LogDestination() *schema.Resource {
	return &schema.Resource{
		Schema:        LogDestinationSchema(),
		CreateContext: createLogDestination,
		ReadContext:   readLogDestination,
		UpdateContext: updateLogDestination,
		DeleteContext: deleteLogDestination,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func LogDestinationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":           utils.Required(schema.TypeString),
		"integration_id": utils.Optional(schema.TypeString),
		"settings":       utils.SettingsMap(),
	}
}

func validateLogDestination(diags diag.Diagnostics, ld *client.LogDestination) diag.Diagnostics {
	if ld.IntegrationId == "" {
		if ld.Type == "http" {
			ld.IntegrationId = NullPlaceholder
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "LogDestination requires an Integration",
				Detail:   fmt.Sprintf("Please check the docs for %s LogDestinations and specify an `integration_id` in your config.", ld.Type),
			})
		}
	}

	return diags
}

func createLogDestination(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	destination := client.LogDestination{
		Type:          data.Get("type").(string),
		IntegrationId: data.Get("integration_id").(string),
		Settings:      getSettings(data),
	}

	if diags = validateLogDestination(diags, &destination); diags.HasError() {
		return diags
	}

	id, err := c.LogDestination.Create(destination)
	if err != nil {
		diags = utils.DiagsCheckError(diags, err, "Unable to create LogDestination")
	} else {
		data.SetId(id)
	}
	return diags
}

func readLogDestination(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	destination, err := c.LogDestination.Read(id)
	if err != nil {
		if isNotFoundError(err) {
			log.Println(notFoundWarning("LogDestination", id))
			data.SetId("")
			return nil
		}
		diags = append(diags, utils.DiagFromError(err, "Unable to read LogDestination"))
		return diags
	}

	diags = utils.DiagsCheckError(diags, data.Set("type", destination.Type), "Unable to read LogDestination type")
	diags = utils.DiagsCheckError(diags, data.Set("integration_id", destination.IntegrationId), "Unable to read LogDestination integration_id")
	diags = utils.DiagsCheckError(diags, data.Set("settings", destination.Settings), "Unable to read LogDestination settings")

	return diags
}

func updateLogDestination(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)

	destination := client.LogDestination{
		Id:            data.Id(),
		Type:          data.Get("type").(string),
		IntegrationId: data.Get("integration_id").(string),
		Settings:      getSettings(data),
	}

	if diags = validateLogDestination(diags, &destination); diags.HasError() {
		return diags
	}

	if _, err := c.LogDestination.Update(destination); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to update LogDestination"))
	}

	return diags
}

func deleteLogDestination(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.ApiClient)
	id := data.Id()

	if _, err := c.LogDestination.Delete(id); err != nil {
		diags = append(diags, utils.DiagFromError(err, "Unable to delete LogDestination"))
	}

	return diags
}
