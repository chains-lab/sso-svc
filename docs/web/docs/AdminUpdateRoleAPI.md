# \AdminUpdateRoleAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1AdminAccountIdRolePatch**](AdminUpdateRoleAPI.md#ReNewsChainsAuthV1AdminAccountIdRolePatch) | **Patch** /re-news/chains/auth/v1/admin/{account_id}/{role} | admin role update
[**ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch**](AdminUpdateRoleAPI.md#ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch) | **Patch** /re-news//chains/auth/v1/private/accounts/{account_id}/{role} | admin role update



## ReNewsChainsAuthV1AdminAccountIdRolePatch

> Account ReNewsChainsAuthV1AdminAccountIdRolePatch(ctx, accountId, role).Execute()

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
	resp, r, err := apiClient.AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminAccountIdRolePatch(context.Background(), accountId, role).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminAccountIdRolePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1AdminAccountIdRolePatch`: Account
	fmt.Fprintf(os.Stdout, "Response from `AdminUpdateRoleAPI.ReNewsChainsAuthV1AdminAccountIdRolePatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 
**role** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1AdminAccountIdRolePatchRequest struct via the builder pattern


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


## ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch

> Account ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch(ctx, accountId, role).Execute()

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
	resp, r, err := apiClient.AdminUpdateRoleAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch(context.Background(), accountId, role).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminUpdateRoleAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch`: Account
	fmt.Fprintf(os.Stdout, "Response from `AdminUpdateRoleAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdRolePatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 
**role** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1PrivateAccountsAccountIdRolePatchRequest struct via the builder pattern


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

