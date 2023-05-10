package provider

import (
	"context"
	"fmt"
	artifactregistrydockerimagesclient "github.com/Fourth-Floor-Creative/terraform-provider-artifact-registry/artifact-registry-docker-images-client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	Images CustomImageValue `tfsdk:"images"`
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

var _ attr.Type = ImageListType{}

type ImageListValue struct {
	types.List
}

var _ attr.Value = ImageListValue{}

func (ilt ImageListType) ElementType() attr.Type {
	return CustomImageValueType{}
}

func (ilt ImageListType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() || in.IsNull() {
		return ImageListValue{}, nil
	}
	val, err := ilt.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return ImageListValue{}, nil
	}
	return ImageListValue{
		List: val.(types.List),
	}, nil
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

func (v CustomImageValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	tags := make([]tftypes.Value, 0, len(v.Tags))
	for _, tag := range v.Tags {
		tags = append(tags, tftypes.NewValue(tftypes.String, tag))
	}

	result := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":             tftypes.String,
			"uri":              tftypes.String,
			"tags":             tftypes.List{ElementType: tftypes.DynamicPseudoType},
			"image_size_bytes": tftypes.String,
			"upload_time":      tftypes.String,
			"media_type":       tftypes.String,
			"build_time":       tftypes.String,
			"update_time":      tftypes.String,
		},
	}, map[string]tftypes.Value{
		"name":             tftypes.NewValue(tftypes.String, v.Name),
		"uri":              tftypes.NewValue(tftypes.String, v.URI),
		"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.DynamicPseudoType}, tags),
		"image_size_bytes": tftypes.NewValue(tftypes.String, v.ImageSizeBytes),
		"upload_time":      tftypes.NewValue(tftypes.String, v.UploadTime),
		"media_type":       tftypes.NewValue(tftypes.String, v.MediaType),
		"build_time":       tftypes.NewValue(tftypes.String, v.BuildTime),
		"update_time":      tftypes.NewValue(tftypes.String, v.UpdateTime),
	})

	return result, nil
}

func (v CustomImageValue) Type(ctx context.Context) attr.Type {
	return CustomImageValueType{}
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
	client := a.client
	images, err := client.ListImages(ctx)
	if err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to list images", err.Error()))
		return
	}

	// Convert this data to a list of CustomImageValue
	var imagesList []attr.Value
	for _, image := range images {
		// Create a CustomImageValue for each image
		imageValue := CustomImageValue{
			Name:           image.Name,
			URI:            image.Uri,
			Tags:           image.Tags,
			ImageSizeBytes: image.ImageSizeBytes,
			UploadTime:     image.UploadTime,
			MediaType:      image.MediaType,
			BuildTime:      image.BuildTime,
			UpdateTime:     image.UpdateTime,
		}

		imagesList = append(imagesList, imageValue)
	}
	response.State.SetAttribute(ctx, path.Root("images"), imagesList)
}
