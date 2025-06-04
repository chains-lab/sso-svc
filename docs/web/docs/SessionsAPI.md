# \SessionsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1OwnRefreshPost**](SessionsAPI.md#ReNewsChainsAuthV1OwnRefreshPost) | **Post** /re-news/chains/auth/v1/own/refresh | Refresh Access Token
[**ReNewsChainsAuthV1OwnSessionsDelete**](SessionsAPI.md#ReNewsChainsAuthV1OwnSessionsDelete) | **Delete** /re-news/chains/auth/v1/own/sessions | Terminate user&#39;s sessions
[**ReNewsChainsAuthV1OwnSessionsGet**](SessionsAPI.md#ReNewsChainsAuthV1OwnSessionsGet) | **Get** /re-news/chains/auth/v1/own/sessions | Get user&#39;s sessions
[**ReNewsChainsAuthV1OwnSessionsSessionIdDelete**](SessionsAPI.md#ReNewsChainsAuthV1OwnSessionsSessionIdDelete) | **Delete** /re-news/chains/auth/v1/own/sessions/{session_id} | Terminate user&#39;s session
[**ReNewsChainsAuthV1OwnSessionsSessionIdGet**](SessionsAPI.md#ReNewsChainsAuthV1OwnSessionsSessionIdGet) | **Get** /re-news/chains/auth/v1/own/sessions/{session_id} | Get user&#39;s session



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
	resp, r, err := apiClient.SessionsAPI.ReNewsChainsAuthV1OwnRefreshPost(context.Background()).RefreshToken(refreshToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.ReNewsChainsAuthV1OwnRefreshPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnRefreshPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.ReNewsChainsAuthV1OwnRefreshPost`: %v\n", resp)
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


## ReNewsChainsAuthV1OwnSessionsDelete

> SessionsCollection ReNewsChainsAuthV1OwnSessionsDelete(ctx).Execute()

Terminate user's sessions



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
	resp, r, err := apiClient.SessionsAPI.ReNewsChainsAuthV1OwnSessionsDelete(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.ReNewsChainsAuthV1OwnSessionsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnSessionsDelete`: SessionsCollection
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.ReNewsChainsAuthV1OwnSessionsDelete`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnSessionsDeleteRequest struct via the builder pattern


### Return type

[**SessionsCollection**](SessionsCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1OwnSessionsGet

> SessionsCollection ReNewsChainsAuthV1OwnSessionsGet(ctx).Execute()

Get user's sessions



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
	resp, r, err := apiClient.SessionsAPI.ReNewsChainsAuthV1OwnSessionsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.ReNewsChainsAuthV1OwnSessionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnSessionsGet`: SessionsCollection
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.ReNewsChainsAuthV1OwnSessionsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnSessionsGetRequest struct via the builder pattern


### Return type

[**SessionsCollection**](SessionsCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1OwnSessionsSessionIdDelete

> SessionsCollection ReNewsChainsAuthV1OwnSessionsSessionIdDelete(ctx, sessionId).Execute()

Terminate user's session



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
	sessionId := "sessionId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SessionsAPI.ReNewsChainsAuthV1OwnSessionsSessionIdDelete(context.Background(), sessionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.ReNewsChainsAuthV1OwnSessionsSessionIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnSessionsSessionIdDelete`: SessionsCollection
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.ReNewsChainsAuthV1OwnSessionsSessionIdDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**sessionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnSessionsSessionIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SessionsCollection**](SessionsCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1OwnSessionsSessionIdGet

> Session ReNewsChainsAuthV1OwnSessionsSessionIdGet(ctx, sessionId).Execute()

Get user's session



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
	sessionId := "sessionId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SessionsAPI.ReNewsChainsAuthV1OwnSessionsSessionIdGet(context.Background(), sessionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.ReNewsChainsAuthV1OwnSessionsSessionIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnSessionsSessionIdGet`: Session
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.ReNewsChainsAuthV1OwnSessionsSessionIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**sessionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnSessionsSessionIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Session**](Session.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

