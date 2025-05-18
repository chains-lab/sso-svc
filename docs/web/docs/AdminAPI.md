# \AdminAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReNewsChainsAuthV1AdminAccountIdGet**](AdminAPI.md#ReNewsChainsAuthV1AdminAccountIdGet) | **Get** /re-news/chains/auth/v1/admin/{account_id} | admin get user
[**ReNewsChainsAuthV1AdminAccountIdSessionsDelete**](AdminAPI.md#ReNewsChainsAuthV1AdminAccountIdSessionsDelete) | **Delete** /re-news/chains/auth/v1/admin/{account_id}/sessions | admin delete user
[**ReNewsChainsAuthV1AdminAccountIdSessionsGet**](AdminAPI.md#ReNewsChainsAuthV1AdminAccountIdSessionsGet) | **Get** /re-news/chains/auth/v1/admin/{account_id}/sessions | admin get sessions
[**ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet**](AdminAPI.md#ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet) | **Get** /re-news/chains/auth/v1/admin/{account_id}/sessions/{session_id} | admin get session
[**ReNewsChainsAuthV1PrivateAccountAccountIdGet**](AdminAPI.md#ReNewsChainsAuthV1PrivateAccountAccountIdGet) | **Get** /re-news//chains/auth/v1/private/account/{account_id} | admin get user
[**ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsDelete**](AdminAPI.md#ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsDelete) | **Delete** /re-news//chains/auth/v1/private/accounts/{account_id}/sessions | admin delete user
[**ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet**](AdminAPI.md#ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet) | **Get** /re-news//chains/auth/v1/private/accounts/{account_id}/sessions | admin get sessions
[**ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet**](AdminAPI.md#ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet) | **Get** /re-news//chains/auth/v1/private/accounts/{account_id}/sessions/{session_id} | admin get session



## ReNewsChainsAuthV1AdminAccountIdGet

> Account ReNewsChainsAuthV1AdminAccountIdGet(ctx, accountId).Execute()

admin get user



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ReNewsChainsAuthV1AdminAccountIdGet(context.Background(), accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1AdminAccountIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1AdminAccountIdGet`: Account
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ReNewsChainsAuthV1AdminAccountIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1AdminAccountIdGetRequest struct via the builder pattern


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


## ReNewsChainsAuthV1AdminAccountIdSessionsDelete

> ReNewsChainsAuthV1AdminAccountIdSessionsDelete(ctx, accountId).Execute()

admin delete user



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsDelete(context.Background(), accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1AdminAccountIdSessionsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1AdminAccountIdSessionsGet

> SessionsCollection ReNewsChainsAuthV1AdminAccountIdSessionsGet(ctx, accountId).Execute()

admin get sessions



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsGet(context.Background(), accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1AdminAccountIdSessionsGet`: SessionsCollection
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1AdminAccountIdSessionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SessionsCollection**](SessionsCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet

> Session ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet(ctx, accountId, sessionId).Execute()

admin get session



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
	sessionId := "sessionId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet(context.Background(), accountId, sessionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet`: Session
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 
**sessionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1AdminAccountIdSessionsSessionIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Session**](Session.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1PrivateAccountAccountIdGet

> Account ReNewsChainsAuthV1PrivateAccountAccountIdGet(ctx, accountId).Execute()

admin get user



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ReNewsChainsAuthV1PrivateAccountAccountIdGet(context.Background(), accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1PrivateAccountAccountIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1PrivateAccountAccountIdGet`: Account
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ReNewsChainsAuthV1PrivateAccountAccountIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1PrivateAccountAccountIdGetRequest struct via the builder pattern


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


## ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsDelete

> ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsDelete(ctx, accountId).Execute()

admin delete user



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsDelete(context.Background(), accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1PrivateAccountsAccountIdSessionsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet

> SessionsCollection ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet(ctx, accountId).Execute()

admin get sessions



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet(context.Background(), accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet`: SessionsCollection
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1PrivateAccountsAccountIdSessionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SessionsCollection**](SessionsCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet

> Session ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet(ctx, accountId, sessionId).Execute()

admin get session



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
	sessionId := "sessionId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet(context.Background(), accountId, sessionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet`: Session
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.ReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **string** |  | 
**sessionId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReNewsChainsAuthV1PrivateAccountsAccountIdSessionsSessionIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Session**](Session.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/vnd.api+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

