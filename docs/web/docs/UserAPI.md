# \UserAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ResetPasswordCompletePost**](UserAPI.md#ResetPasswordCompletePost) | **Post** /reset-password-complete | Reset password
[**ResetPasswordInitiatePost**](UserAPI.md#ResetPasswordInitiatePost) | **Post** /reset-password-initiate | Send a request to change user password



## ResetPasswordCompletePost

> ResetPasswordCompletePost(ctx).ResetPasswordComplete(resetPasswordComplete).Execute()

Reset password



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
	resetPasswordComplete := *openapiclient.NewResetPasswordComplete(*openapiclient.NewResetPasswordCompleteData("Type_example", *openapiclient.NewResetPasswordCompleteDataAttributes("Email_example", "Password_example"))) // ResetPasswordComplete | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.UserAPI.ResetPasswordCompletePost(context.Background()).ResetPasswordComplete(resetPasswordComplete).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserAPI.ResetPasswordCompletePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiResetPasswordCompletePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resetPasswordComplete** | [**ResetPasswordComplete**](ResetPasswordComplete.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetPasswordInitiatePost

> ResetPasswordInitiatePost(ctx).ResetPasswordInitiate(resetPasswordInitiate).Execute()

Send a request to change user password



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
	resetPasswordInitiate := *openapiclient.NewResetPasswordInitiate(*openapiclient.NewResetPasswordInitiateData("Type_example", *openapiclient.NewResetPasswordInitiateDataAttributes("Email_example"))) // ResetPasswordInitiate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.UserAPI.ResetPasswordInitiatePost(context.Background()).ResetPasswordInitiate(resetPasswordInitiate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserAPI.ResetPasswordInitiatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiResetPasswordInitiatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resetPasswordInitiate** | [**ResetPasswordInitiate**](ResetPasswordInitiate.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

