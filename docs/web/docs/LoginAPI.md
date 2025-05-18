# \LoginAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1OwnGoogleLoginPost**](LoginAPI.md#ReNewsChainsAuthV1OwnGoogleLoginPost) | **Post** /re-news//chains/auth/v1/own/google/login | 
[**ReNewsChainsAuthV1OwnLoginPost**](LoginAPI.md#ReNewsChainsAuthV1OwnLoginPost) | **Post** /re-news/chains/auth/v1/own/login | 
[**ReNewsChainsAuthV1OwnRefreshPost**](LoginAPI.md#ReNewsChainsAuthV1OwnRefreshPost) | **Post** /re-news/chains/auth/v1/own/refresh | Refresh Access Token
[**ReNewsChainsAuthV1PublicRefreshPost**](LoginAPI.md#ReNewsChainsAuthV1PublicRefreshPost) | **Post** /re-news//chains/auth/v1/public/refresh | Refresh Access Token



## ReNewsChainsAuthV1OwnGoogleLoginPost

> TokensPair ReNewsChainsAuthV1OwnGoogleLoginPost(ctx).Execute()





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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoginAPI.ReNewsChainsAuthV1OwnGoogleLoginPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.ReNewsChainsAuthV1OwnGoogleLoginPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnGoogleLoginPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `LoginAPI.ReNewsChainsAuthV1OwnGoogleLoginPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnGoogleLoginPostRequest struct via the builder pattern


### Return type

[**TokensPair**](TokensPair.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1OwnLoginPost

> TokensPair ReNewsChainsAuthV1OwnLoginPost(ctx).Execute()





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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoginAPI.ReNewsChainsAuthV1OwnLoginPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.ReNewsChainsAuthV1OwnLoginPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnLoginPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `LoginAPI.ReNewsChainsAuthV1OwnLoginPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnLoginPostRequest struct via the builder pattern


### Return type

[**TokensPair**](TokensPair.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


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
	resp, r, err := apiClient.LoginAPI.ReNewsChainsAuthV1OwnRefreshPost(context.Background()).RefreshToken(refreshToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.ReNewsChainsAuthV1OwnRefreshPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnRefreshPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `LoginAPI.ReNewsChainsAuthV1OwnRefreshPost`: %v\n", resp)
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


## ReNewsChainsAuthV1PublicRefreshPost

> TokensPair ReNewsChainsAuthV1PublicRefreshPost(ctx).RefreshToken(refreshToken).Execute()

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
	resp, r, err := apiClient.LoginAPI.ReNewsChainsAuthV1PublicRefreshPost(context.Background()).RefreshToken(refreshToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.ReNewsChainsAuthV1PublicRefreshPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1PublicRefreshPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `LoginAPI.ReNewsChainsAuthV1PublicRefreshPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1PublicRefreshPostRequest struct via the builder pattern


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

