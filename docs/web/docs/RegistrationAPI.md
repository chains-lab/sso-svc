# \RegistrationAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RegistertionInitiatePost**](RegistrationAPI.md#RegistertionInitiatePost) | **Post** /registertion-initiate | Send a request to register a new user
[**RegistrationCompletePost**](RegistrationAPI.md#RegistrationCompletePost) | **Post** /registration-complete | Confirm register a new user



## RegistertionInitiatePost

> RegistertionInitiatePost(ctx).RegistrationInitiate(registrationInitiate).Execute()

Send a request to register a new user



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
	registrationInitiate := *openapiclient.NewRegistrationInitiate(*openapiclient.NewRegistrationInitiateData("Type_example", *openapiclient.NewRegistrationInitiateDataAttributes("Email_example"))) // RegistrationInitiate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RegistrationAPI.RegistertionInitiatePost(context.Background()).RegistrationInitiate(registrationInitiate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RegistrationAPI.RegistertionInitiatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRegistertionInitiatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registrationInitiate** | [**RegistrationInitiate**](RegistrationInitiate.md) |  | 

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


## RegistrationCompletePost

> RegistrationCompletePost201Response RegistrationCompletePost(ctx).RegistrationComplete(registrationComplete).Execute()

Confirm register a new user



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
	registrationComplete := *openapiclient.NewRegistrationComplete(*openapiclient.NewRegistrationCompleteData("Type_example", *openapiclient.NewRegistrationCompleteDataAttributes("Password_example", "Email_example"))) // RegistrationComplete | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RegistrationAPI.RegistrationCompletePost(context.Background()).RegistrationComplete(registrationComplete).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RegistrationAPI.RegistrationCompletePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RegistrationCompletePost`: RegistrationCompletePost201Response
	fmt.Fprintf(os.Stdout, "Response from `RegistrationAPI.RegistrationCompletePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRegistrationCompletePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registrationComplete** | [**RegistrationComplete**](RegistrationComplete.md) |  | 

### Return type

[**RegistrationCompletePost201Response**](RegistrationCompletePost201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/vnd.api+json
- **Accept**: application/json, application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

