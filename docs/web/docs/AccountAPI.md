# \AccountAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1OwnGet**](AccountAPI.md#ReNewsChainsAuthV1OwnGet) | **Get** /re-news/chains/auth/v1/own | Get account
[**ReNewsChainsAuthV1PublicAccountGet**](AccountAPI.md#ReNewsChainsAuthV1PublicAccountGet) | **Get** /re-news//chains/auth/v1/public/account | Get account



## ReNewsChainsAuthV1OwnGet

> Account ReNewsChainsAuthV1OwnGet(ctx).Execute()

Get account



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
	resp, r, err := apiClient.AccountAPI.ReNewsChainsAuthV1OwnGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.ReNewsChainsAuthV1OwnGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1OwnGet`: Account
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.ReNewsChainsAuthV1OwnGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1OwnGetRequest struct via the builder pattern


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


## ReNewsChainsAuthV1PublicAccountGet

> Account ReNewsChainsAuthV1PublicAccountGet(ctx).Execute()

Get account



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
	resp, r, err := apiClient.AccountAPI.ReNewsChainsAuthV1PublicAccountGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.ReNewsChainsAuthV1PublicAccountGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1PublicAccountGet`: Account
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.ReNewsChainsAuthV1PublicAccountGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1PublicAccountGetRequest struct via the builder pattern


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

