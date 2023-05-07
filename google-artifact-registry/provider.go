package google_artifact_registry

import (
	"context"
	artifactregistrydockerimagesclient "github.com/Fourth-Floor-Creative/terraform-provider-artifact-registry/artifact-registry-docker-images-client"
	"golang.org/x/oauth2/google"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"artifact_registry_images": dataSourceImages(),
		},
		ConfigureContextFunc: configureProvider,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project ID where the Artifact Registry repository is located.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location of the Artifact Registry repository.",
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Artifact Registry repository.",
			},
		},
	}
}

var requiredConfig = []string{
	"project",
	"location",
	"repository",
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	for _, key := range requiredConfig {
		if _, ok := d.GetOk(key); !ok {
			return nil, diag.Errorf("missing required configuration '%s'", key)
		}
	}

	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/cloud-platform.read-only")
	if err != nil {
		return nil, nil
	}

	client, err := artifactregistrydockerimagesclient.NewClient(nil, &artifactregistrydockerimagesclient.Options{
		Credentials: credentials,
		ProjectID:   d.Get("project").(string),
		Location:    d.Get("location").(string),
		Repository:  d.Get("repository").(string),
	})
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}
