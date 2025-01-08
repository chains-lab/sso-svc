# \LoginAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**LoginCompletePost**](LoginAPI.md#LoginCompletePost) | **Post** /login-complete | Confirm login
[**LoginInitiatePost**](LoginAPI.md#LoginInitiatePost) | **Post** /login-initiate | Send a request to login
[**RefreshPost**](LoginAPI.md#RefreshPost) | **Post** /refresh | Refresh Access Token
[**UserChangeLogoutPost**](LoginAPI.md#UserChangeLogoutPost) | **Post** /user/change/logout | Logout user



## LoginCompletePost

> LoginCompleteResp LoginCompletePost(ctx).LoginCompleteReq(loginCompleteReq).Execute()

Confirm login



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
	loginCompleteReq := *openapiclient.NewLoginCompleteReq(*openapiclient.NewLoginCompleteReqData("Type_example", *openapiclient.NewLoginCompleteReqDataAttributes("Email_example", "DeviceName_example"))) // LoginCompleteReq | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoginAPI.LoginCompletePost(context.Background()).LoginCompleteReq(loginCompleteReq).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.LoginCompletePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `LoginCompletePost`: LoginCompleteResp
	fmt.Fprintf(os.Stdout, "Response from `LoginAPI.LoginCompletePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLoginCompletePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginCompleteReq** | [**LoginCompleteReq**](LoginCompleteReq.md) |  | 

### Return type

[**LoginCompleteResp**](LoginCompleteResp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LoginInitiatePost

> LoginInitiatePost(ctx).LoginInitiate(loginInitiate).Execute()

Send a request to login



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
	loginInitiate := *openapiclient.NewLoginInitiate(*openapiclient.NewLoginInitiateData("Type_example", *openapiclient.NewLoginInitiateDataAttributes("Email_example", "Password_example"))) // LoginInitiate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.LoginAPI.LoginInitiatePost(context.Background()).LoginInitiate(loginInitiate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.LoginInitiatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLoginInitiatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginInitiate** | [**LoginInitiate**](LoginInitiate.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/vnd.api+json, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshPost

> RefreshResp RefreshPost(ctx).RefreshReq(refreshReq).Execute()

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
	refreshReq := *openapiclient.NewRefreshReq(*openapiclient.NewRefreshReqData("Type_example", *openapiclient.NewRefreshReqDataAttributes("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."))) // RefreshReq | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoginAPI.RefreshPost(context.Background()).RefreshReq(refreshReq).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.RefreshPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshPost`: RefreshResp
	fmt.Fprintf(os.Stdout, "Response from `LoginAPI.RefreshPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRefreshPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **refreshReq** | [**RefreshReq**](RefreshReq.md) |  | 

### Return type

[**RefreshResp**](RefreshResp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserChangeLogoutPost

> UserChangeLogoutPost(ctx).Execute()

Logout user



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
	r, err := apiClient.LoginAPI.UserChangeLogoutPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.UserChangeLogoutPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiUserChangeLogoutPostRequest struct via the builder pattern


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

