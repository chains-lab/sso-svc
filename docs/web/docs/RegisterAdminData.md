# RegisterAdminData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Attributes** | [**RegisterUserDataAttributes**](RegisterUserDataAttributes.md) |  | 

## Methods

### NewRegisterAdminData

`func NewRegisterAdminData(type_ string, attributes RegisterUserDataAttributes, ) *RegisterAdminData`

NewRegisterAdminData instantiates a new RegisterAdminData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterAdminDataWithDefaults

`func NewRegisterAdminDataWithDefaults() *RegisterAdminData`

NewRegisterAdminDataWithDefaults instantiates a new RegisterAdminData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *RegisterAdminData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *RegisterAdminData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *RegisterAdminData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *RegisterAdminData) GetAttributes() RegisterUserDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *RegisterAdminData) GetAttributesOk() (*RegisterUserDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *RegisterAdminData) SetAttributes(v RegisterUserDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


