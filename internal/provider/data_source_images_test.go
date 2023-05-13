package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	config := `
provider "artifactregistry" {
	project = "devops-339608"
	location = "europe"
	repository = "services"
}
data "artifactregistry_artifact_registry_images" "test" {}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.artifactregistry_artifact_registry_images.test", "images.0.%", "9"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.artifactregistry_artifact_registry_images.test", "latest_images.campaign-service.%", "9"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.artifactregistry_artifact_registry_images.test", "latest_images.campaign-service.development_tagged_uri", "europe-docker.pkg.dev/devops-339608/services/campaign-service:development-9681cde"),
				),
			},
		},
	})
}
