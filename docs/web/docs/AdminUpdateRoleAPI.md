# \AdminUpdateRoleAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1AdminAccountIdRolePost**](AdminUpdateRoleAPI.md#ReNewsChainsAuthV1AdminAccountIdRolePost) | **Post** /re-news/chains/auth/v1/admin/{account_id}/{role} | admin role update



## ReNewsChainsAuthV1AdminAccountIdRolePost

> Account ReNewsChainsAuthV1AdminAccountIdRolePost(ctx, accountId, role).Execute()

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
	accountId := "accountId_example" // string | 
	role := "role_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminAccountIdRolePost(context.Background(), accountId, role).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminAccountIdRolePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1AdminAccountIdRolePost`: Account
	fmt.Fprintf(os.Stdout, "Response from `AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminAccountIdRolePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 
**role** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1AdminAccountIdRolePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Account**](Account.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

