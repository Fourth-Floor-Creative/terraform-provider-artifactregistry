package provider

import (
	"context"
	"fmt"
	artifactregistrydockerimagesclient "github.com/Fourth-Floor-Creative/terraform-provider-artifact-registry/artifact-registry-docker-images-client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ArtifactRegistryImagesDataSource{}

func NewArtifactRegistryImagesData() datasource.DataSource {
	return &ArtifactRegistryImagesDataSource{}
}

// ArtifactRegistryImagesDataSource defines the data source implementation.
type ArtifactRegistryImagesDataSource struct {
	client *artifactregistrydockerimagesclient.Client
}

// ArtifactRegistryImagesDataSourceModel defines the data source model.
type ArtifactRegistryImagesDataSourceModel struct {
	Images []CustomImageValue `tfsdk:"images"`
}

func (a *ArtifactRegistryImagesDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_artifact_registry_images"
}

type CustomImageValueType struct {
	types.ObjectType
}

func (civt CustomImageValueType) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":             types.StringType,
		"uri":              types.StringType,
		"tags":             types.ListType{ElemType: types.StringType},
		"image_size_bytes": types.StringType,
		"upload_time":      types.StringType,
		"media_type":       types.StringType,
		"build_time":       types.StringType,
		"update_time":      types.StringType,
	}
}

type ImageListType struct {
	types.ListType
}

type ImageListValue struct {
	types.List
}

func (ilt ImageListType) ElementType() attr.Type {
	return CustomImageValueType{}
}

func (ilt ImageListType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := ilt.ListType.ValueFromTerraform(ctx, in)

	return ImageListValue{
		List: val.(types.List),
	}, err
}

type CustomImageValue struct {
	types.List
	Name           string   `tfsdk:"name"`
	URI            string   `tfsdk:"uri"`
	Tags           []string `tfsdk:"tags"`
	ImageSizeBytes string   `tfsdk:"image_size_bytes"`
	UploadTime     string   `tfsdk:"upload_time"`
	MediaType      string   `tfsdk:"media_type"`
	BuildTime      string   `tfsdk:"build_time"`
	UpdateTime     string   `tfsdk:"update_time"`
}

func (a *ArtifactRegistryImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*artifactregistrydockerimagesclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *artifactregistrydockerimagesclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	a.client = client
}

func (a *ArtifactRegistryImagesDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This data source provides a list of images in a repository.",
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of images in the repository.",
				CustomType: ImageListType{
					types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"name": types.StringType,
								"uri":  types.StringType,
								"tags": types.ListType{
									ElemType: types.StringType,
								},
								"image_size_bytes": types.StringType,
								"upload_time":      types.StringType,
								"media_type":       types.StringType,
								"build_time":       types.StringType,
								"update_time":      types.StringType,
							},
						},
					},
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"uri": schema.StringAttribute{
							Computed: true,
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"image_size_bytes": schema.StringAttribute{
							Computed: true,
						},
						"upload_time": schema.StringAttribute{
							Computed: true,
						},
						"media_type": schema.StringAttribute{
							Computed: true,
						},
						"build_time": schema.StringAttribute{
							Computed: true,
						},
						"update_time": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (a *ArtifactRegistryImagesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data ArtifactRegistryImagesDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	client := a.client
	images, err := client.ListImages(ctx)
	if err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to list images", err.Error()))
		return
	}

	data.Images = make([]CustomImageValue, len(images))
	for i, image := range images {
		data.Images[i] = CustomImageValue{
			Name:           image.Name,
			URI:            image.Uri,
			Tags:           image.Tags,
			ImageSizeBytes: image.ImageSizeBytes,
			UploadTime:     image.UploadTime,
			MediaType:      image.MediaType,
			BuildTime:      image.BuildTime,
			UpdateTime:     image.UpdateTime,
		}
	}

	// Add the images to the response
	diagnostic := response.State.Set(ctx, &data)
	if diagnostic.HasError() {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to set images", "failed to set images"))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
