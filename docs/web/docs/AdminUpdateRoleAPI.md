# \AdminUpdateRoleAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1AdminUserIdRolePost**](AdminUpdateRoleAPI.md#ReNewsChainsAuthV1AdminUserIdRolePost) | **Post** /re-news/chains/auth/v1/admin/{user_id}/{role} | admin role update



## ReNewsChainsAuthV1AdminUserIdRolePost

> User ReNewsChainsAuthV1AdminUserIdRolePost(ctx, userId, role).Execute()

admin role update



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
	userId := "userId_example" // string | 
	role := "role_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminUserIdRolePost(context.Background(), userId, role).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminUserIdRolePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1AdminUserIdRolePost`: User
	fmt.Fprintf(os.Stdout, "Response from `AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminUserIdRolePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** |  | 
**role** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1AdminUserIdRolePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

