package artifact_registry_docker_images_client

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"golang.org/x/oauth2/google"
)

const (
	apiBaseUrl = "https://artifactregistry.googleapis.com/v1/"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

// Error implements go error interface.
func (msg *ErrorMessage) Error() string {
	return fmt.Sprintf("API Error: %s", msg.Message)
}

// ClientAccessTokenOauthResponse is the response from the TikTok OAuth endpoint
type ClientAccessTokenOauthResponse struct {
	AccessToken string `json:"access_token"`
	// Tokens are valid for 2 hours
	ExpiresIn int64  `json:"expires_in"`
	TokenType string `json:"token_type"`
}

type Client struct {
	*req.Client
	ProjectID  string
	Location   string
	Repository string
}

type Options struct {
	Credentials *google.Credentials
	ProjectID   string
	Location    string
	Repository  string
}

// NewClient creates a new TikTok client and authenticates it
func NewClient(reqClient *req.Client, options *Options) (*Client, error) {
	token, err := options.Credentials.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	if reqClient == nil {
		reqClient = req.NewClient()
	}
	reqClient.
		SetBaseURL(apiBaseUrl).
		SetCommonErrorResult(&ErrorMessage{}).
		EnableDumpEachRequest().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
				return nil
			}
			if errMsg, ok := resp.ErrorResult().(*ErrorMessage); ok {
				resp.Err = errMsg // Convert api error into go error
				return nil
			}
			if !resp.IsSuccessState() {
				// Neither a success response nor an error response, record details to help troubleshooting
				resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
			}
			return nil
		}).
		SetCommonBearerAuthToken(token.AccessToken)

	newClient := &Client{
		Client:     reqClient,
		ProjectID:  options.ProjectID,
		Location:   options.Location,
		Repository: options.Repository,
	}
	return newClient, nil
}

type ListImagesResponse struct {
	DockerImages  []DockerImage `json:"dockerImages"`
	NextPageToken string        `json:"nextPageToken"`
}

type DockerImage struct {
	Name           string   `json:"name"`
	Uri            string   `json:"uri"`
	Tags           []string `json:"tags"`
	ImageSizeBytes string   `json:"imageSizeBytes"`
	UploadTime     string   `json:"uploadTime"`
	MediaType      string   `json:"mediaType"`
	BuildTime      string   `json:"buildTime"`
	UpdateTime     string   `json:"updateTime"`
}

// ListImages hits https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/get
// to list the images in the registry.
func (c *Client) ListImages(ctx context.Context) ([]DockerImage, error) {
	var dockerImages []DockerImage
	hasNextPage := true
	var nextPageToken string
	for hasNextPage {
		var listImagesResponse ListImagesResponse
		request := c.R().SetURL(fmt.Sprintf("projects/%s/locations/%s/repositories/%s/dockerImages/", c.ProjectID, c.Location, c.Repository)).
			SetSuccessResult(&listImagesResponse).
			SetQueryParam("pageSize", "200")
		if nextPageToken != "" {
			request.SetQueryParam("pageToken", nextPageToken)
		}
		res := request.Do(ctx)
		if res.IsErrorState() {
			return nil, res.Err
		}
		dockerImages = append(dockerImages, listImagesResponse.DockerImages...)
		hasNextPage = listImagesResponse.NextPageToken != ""
		nextPageToken = listImagesResponse.NextPageToken
	}
	return dockerImages, nil
}
