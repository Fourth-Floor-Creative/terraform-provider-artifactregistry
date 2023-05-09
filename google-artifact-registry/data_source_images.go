package google_artifact_registry

import (
	"context"
	"fmt"
	artifact_registry_docker_images_client "github.com/Fourth-Floor-Creative/terraform-provider-artifact-registry/artifact-registry-docker-images-client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"image_size_bytes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upload_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"media_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*artifact_registry_docker_images_client.Client)

	projectID := d.Get("project").(string)
	location := d.Get("location").(string)
	repository := d.Get("repository").(string)

	images, err := client.ListImages(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("images", images); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", projectID, location, repository))

	return nil
}
