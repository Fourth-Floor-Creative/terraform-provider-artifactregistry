package main

import (
	googleartifactregistry "github.com/Fourth-Floor-Creative/terraform-provider-artifact-registry/google-artifact-registry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: googleartifactregistry.New})
}
