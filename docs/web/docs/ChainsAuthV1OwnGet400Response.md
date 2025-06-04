# ChainsAuthV1OwnGet400Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Errors** | [**[]ChainsAuthV1OwnGet400ResponseErrorsInner**](ChainsAuthV1OwnGet400ResponseErrorsInner.md) | Non empty array of errors occurred during request processing | 

## Methods

### NewChainsAuthV1OwnGet400Response

`func NewChainsAuthV1OwnGet400Response(errors []ChainsAuthV1OwnGet400ResponseErrorsInner, ) *ChainsAuthV1OwnGet400Response`

NewChainsAuthV1OwnGet400Response instantiates a new ChainsAuthV1OwnGet400Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChainsAuthV1OwnGet400ResponseWithDefaults

`func NewChainsAuthV1OwnGet400ResponseWithDefaults() *ChainsAuthV1OwnGet400Response`

NewChainsAuthV1OwnGet400ResponseWithDefaults instantiates a new ChainsAuthV1OwnGet400Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrors

`func (o *ChainsAuthV1OwnGet400Response) GetErrors() []ChainsAuthV1OwnGet400ResponseErrorsInner`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *ChainsAuthV1OwnGet400Response) GetErrorsOk() (*[]ChainsAuthV1OwnGet400ResponseErrorsInner, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *ChainsAuthV1OwnGet400Response) SetErrors(v []ChainsAuthV1OwnGet400ResponseErrorsInner)`

SetErrors sets Errors field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


