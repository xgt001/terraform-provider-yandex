package yandex

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/resourcemanager/v1"
	"github.com/yandex-cloud/go-sdk/sdkresolvers"
)

func dataSourceYandexResourceManagerCloud() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceYandexResourceManagerCloudRead,
		Schema: map[string]*schema.Schema{
			"cloud_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceYandexResourceManagerCloudRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := config.Context()

	err := checkOneOf(d, "cloud_id", "name")
	if err != nil {
		return err
	}

	cloudID := d.Get("cloud_id").(string)
	cloudName, cloudNameOk := d.GetOk("name")

	if cloudNameOk {
		cloudID, err = resolveObjectID(ctx, config, cloudName.(string), sdkresolvers.CloudResolver)
		if err != nil {
			return fmt.Errorf("failed to resolve data source cloud by name: %v", err)
		}
	}

	cloud, err := config.sdk.ResourceManager().Cloud().Get(ctx, &resourcemanager.GetCloudRequest{
		CloudId: cloudID,
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("cloud with ID %q", cloudID))
	}

	createdAt, err := getTimestamp(cloud.CreatedAt)
	if err != nil {
		return err
	}

	d.Set("cloud_id", cloud.Id)
	d.Set("name", cloud.Name)
	d.Set("description", cloud.Description)
	d.Set("created_at", createdAt)
	d.SetId(cloud.Id)

	return nil
}
