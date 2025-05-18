# \LogoutAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1OwnLogoutPost**](LogoutAPI.md#ReNewsChainsAuthV1OwnLogoutPost) | **Post** /re-news/chains/auth/v1/own/logout | 



## ReNewsChainsAuthV1OwnLogoutPost

> TokensPair ReNewsChainsAuthV1OwnLogoutPost(ctx).Execute()





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
	resp, r, err := apiClient.LogoutAPI.ReNewsChainsAuthV1OwnLogoutPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LogoutAPI.ReNewsChainsAuthV1OwnLogoutPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnLogoutPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `LogoutAPI.ReNewsChainsAuthV1OwnLogoutPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnLogoutPostRequest struct via the builder pattern


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

