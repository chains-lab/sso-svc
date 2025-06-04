# \SessionAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1OwnRefreshPost**](SessionAPI.md#ReNewsChainsAuthV1OwnRefreshPost) | **Post** /re-news/chains/auth/v1/own/refresh | Refresh Access Token



## ReNewsChainsAuthV1OwnRefreshPost

> TokensPair ReNewsChainsAuthV1OwnRefreshPost(ctx).RefreshToken(refreshToken).Execute()

Refresh Access Token



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	refreshToken := *openapiclient.NewRefreshToken(*openapiclient.NewRefreshTokenData("Type_example", *openapiclient.NewRefreshTokenDataAttributes("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."))) // RefreshToken | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SessionAPI.ReNewsChainsAuthV1OwnRefreshPost(context.Background()).RefreshToken(refreshToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionAPI.ReNewsChainsAuthV1OwnRefreshPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnRefreshPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `SessionAPI.ReNewsChainsAuthV1OwnRefreshPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnRefreshPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **refreshToken** | [**RefreshToken**](RefreshToken.md) |  | 

### Return type

[**TokensPair**](TokensPair.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

