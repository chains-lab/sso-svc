# Login

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**LoginData**](LoginData.md) |  | 

## Methods

### NewLogin

`func NewLogin(data LoginData, ) *Login`

NewLogin instantiates a new Login object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginWithDefaults

`func NewLoginWithDefaults() *Login`

NewLoginWithDefaults instantiates a new Login object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *Login) GetData() LoginData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Login) GetDataOk() (*LoginData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Login) SetData(v LoginData)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


