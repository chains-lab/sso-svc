# \LoginAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsSsoV1PublicGoogleLoginPost**](LoginAPI.md#ReNewsSsoV1PublicGoogleLoginPost) | **Post** /re-news/sso/v1/public/google/login | 



## ReNewsSsoV1PublicGoogleLoginPost

> TokensPair ReNewsSsoV1PublicGoogleLoginPost(ctx).Execute()





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
	resp, r, err := apiClient.LoginAPI.ReNewsSsoV1PublicGoogleLoginPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoginAPI.ReNewsSsoV1PublicGoogleLoginPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsSsoV1PublicGoogleLoginPost`: TokensPair
	fmt.Fprintf(os.Stdout, "Response from `LoginAPI.ReNewsSsoV1PublicGoogleLoginPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsSsoV1PublicGoogleLoginPostRequest struct via the builder pattern


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

