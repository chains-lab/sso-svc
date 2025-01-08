# \SessionsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UserChangeDeleteSessionDelete**](SessionsAPI.md#UserChangeDeleteSessionDelete) | **Delete** /user/change/delete-session | Delete the sessions
[**UserChangeSessionsGet**](SessionsAPI.md#UserChangeSessionsGet) | **Get** /user/change/sessions | Get user&#39;s session
[**UserChangeTerminateSessionsDelete**](SessionsAPI.md#UserChangeTerminateSessionsDelete) | **Delete** /user/change/terminate-sessions | Terminate all sessions



## UserChangeDeleteSessionDelete

> UserSessions UserChangeDeleteSessionDelete(ctx).DeleteSession(deleteSession).Execute()

Delete the sessions



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
	deleteSession := *openapiclient.NewDeleteSession(*openapiclient.NewDeleteSessionData("Type_example", *openapiclient.NewDeleteSessionDataAttributes("SessionId_example"))) // DeleteSession | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SessionsAPI.UserChangeDeleteSessionDelete(context.Background()).DeleteSession(deleteSession).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.UserChangeDeleteSessionDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UserChangeDeleteSessionDelete`: UserSessions
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.UserChangeDeleteSessionDelete`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserChangeDeleteSessionDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deleteSession** | [**DeleteSession**](DeleteSession.md) |  | 

### Return type

[**UserSessions**](UserSessions.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserChangeSessionsGet

> UserSessions UserChangeSessionsGet(ctx).Execute()

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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SessionsAPI.UserChangeSessionsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.UserChangeSessionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UserChangeSessionsGet`: UserSessions
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.UserChangeSessionsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiUserChangeSessionsGetRequest struct via the builder pattern


### Return type

[**UserSessions**](UserSessions.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserChangeTerminateSessionsDelete

> UserChangeTerminateSessionsDelete200Response UserChangeTerminateSessionsDelete(ctx).Execute()

Terminate all sessions



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
	resp, r, err := apiClient.SessionsAPI.UserChangeTerminateSessionsDelete(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SessionsAPI.UserChangeTerminateSessionsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UserChangeTerminateSessionsDelete`: UserChangeTerminateSessionsDelete200Response
	fmt.Fprintf(os.Stdout, "Response from `SessionsAPI.UserChangeTerminateSessionsDelete`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiUserChangeTerminateSessionsDeleteRequest struct via the builder pattern


### Return type

[**UserChangeTerminateSessionsDelete200Response**](UserChangeTerminateSessionsDelete200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

