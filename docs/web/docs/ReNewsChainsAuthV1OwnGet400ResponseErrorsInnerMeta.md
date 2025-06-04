# ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ErrorId** | **string** | Error ID is a unique identifier for the error, used for debugging and tracing | 
**RequestId** | **string** | Request ID is a unique identifier for the request, used for debugging and tracing | 
**Parameter** | Pointer to **string** | Parameter is the name of the request parameter that caused the error, if applicable | [optional] 
**Pointer** | Pointer to **string** | Pointer is a JSON Pointer that identifies the part of the request document that caused the error, if applicable | [optional] 
**Timestamp** | **time.Time** | Timestamp is the time when the error occurred, in ISO 8601 format | 

## Methods

### NewReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta

`func NewReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta(errorId string, requestId string, timestamp time.Time, ) *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta`

NewReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta instantiates a new ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMetaWithDefaults

`func NewReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMetaWithDefaults() *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta`

NewReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMetaWithDefaults instantiates a new ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrorId

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetErrorId() string`

GetErrorId returns the ErrorId field if non-nil, zero value otherwise.

### GetErrorIdOk

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetErrorIdOk() (*string, bool)`

GetErrorIdOk returns a tuple with the ErrorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorId

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) SetErrorId(v string)`

SetErrorId sets ErrorId field to given value.


### GetRequestId

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetRequestId() string`

GetRequestId returns the RequestId field if non-nil, zero value otherwise.

### GetRequestIdOk

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetRequestIdOk() (*string, bool)`

GetRequestIdOk returns a tuple with the RequestId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestId

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) SetRequestId(v string)`

SetRequestId sets RequestId field to given value.


### GetParameter

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetParameter() string`

GetParameter returns the Parameter field if non-nil, zero value otherwise.

### GetParameterOk

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetParameterOk() (*string, bool)`

GetParameterOk returns a tuple with the Parameter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParameter

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) SetParameter(v string)`

SetParameter sets Parameter field to given value.

### HasParameter

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) HasParameter() bool`

HasParameter returns a boolean if a field has been set.

### GetPointer

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetPointer() string`

GetPointer returns the Pointer field if non-nil, zero value otherwise.

### GetPointerOk

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetPointerOk() (*string, bool)`

GetPointerOk returns a tuple with the Pointer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPointer

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) SetPointer(v string)`

SetPointer sets Pointer field to given value.

### HasPointer

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) HasPointer() bool`

HasPointer returns a boolean if a field has been set.

### GetTimestamp

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *ReNewsChainsAuthV1OwnGet400ResponseErrorsInnerMeta) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


