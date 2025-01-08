# RegistrationConfirmData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique identifier username | 
**Type** | **string** |  | 
**Attributes** | [**RegistrationConfirmDataAttributes**](RegistrationConfirmDataAttributes.md) |  | 

## Methods

### NewRegistrationConfirmData

`func NewRegistrationConfirmData(id string, type_ string, attributes RegistrationConfirmDataAttributes, ) *RegistrationConfirmData`

NewRegistrationConfirmData instantiates a new RegistrationConfirmData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationConfirmDataWithDefaults

`func NewRegistrationConfirmDataWithDefaults() *RegistrationConfirmData`

NewRegistrationConfirmDataWithDefaults instantiates a new RegistrationConfirmData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RegistrationConfirmData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RegistrationConfirmData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RegistrationConfirmData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *RegistrationConfirmData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *RegistrationConfirmData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *RegistrationConfirmData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *RegistrationConfirmData) GetAttributes() RegistrationConfirmDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *RegistrationConfirmData) GetAttributesOk() (*RegistrationConfirmDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *RegistrationConfirmData) SetAttributes(v RegistrationConfirmDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


