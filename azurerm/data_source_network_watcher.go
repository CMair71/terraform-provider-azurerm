package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNetworkWatcher() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetworkWatcherRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"location":            locationForDataSourceSchema(),
			"tags":                tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmNetworkWatcherRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).watcherClient

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Network Watcher %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Network Watcher %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}
