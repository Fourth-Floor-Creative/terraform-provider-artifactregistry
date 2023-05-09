package provider

import (
	"context"
	artifactregistrydockerimagesclient "github.com/Fourth-Floor-Creative/terraform-provider-artifact-registry/artifact-registry-docker-images-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2/google"
)

// Ensure ArtifactRegistryProvider satisfies various provider interfaces.
var _ provider.Provider = &ArtifactRegistryProvider{}

// ArtifactRegistryProvider defines the provider implementation.
type ArtifactRegistryProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *ArtifactRegistryProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewArtifactRegistryImagesData,
	}
}

func (p *ArtifactRegistryProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// ArtifactRegistryProviderModel defines the provider data model.
type ArtifactRegistryProviderModel struct {
	Project    types.String `tfsdk:"project"`
	Location   types.String `tfsdk:"location"`
	Repository types.String `tfsdk:"repository"`
}

func (p *ArtifactRegistryProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "artifactregistry"
	resp.Version = p.version
}

func (p *ArtifactRegistryProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Required:    true,
				Description: "The project ID where the Artifact Registry repository is located.",
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location of the Artifact Registry repository.",
			},
			"repository": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Artifact Registry repository.",
			},
		},
	}
}

func (p *ArtifactRegistryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ArtifactRegistryProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/cloud-platform.read-only")
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to find default credentials", err.Error()))
	}

	registryAPIClient, err := artifactregistrydockerimagesclient.NewClient(nil, &artifactregistrydockerimagesclient.Options{
		Credentials: credentials,
		ProjectID:   data.Project.ValueString(),
		Location:    data.Location.ValueString(),
		Repository:  data.Repository.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to create Artifact Registry client", err.Error()))
	}
	resp.DataSourceData = registryAPIClient
	resp.ResourceData = registryAPIClient
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ArtifactRegistryProvider{
			version: version,
		}
	}
}
