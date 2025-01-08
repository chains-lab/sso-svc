# RegistrationReqData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique identifier username | 
**Type** | **string** |  | 
**Attributes** | [**RegistrationReqDataAttributes**](RegistrationReqDataAttributes.md) |  | 

## Methods

### NewRegistrationReqData

`func NewRegistrationReqData(id string, type_ string, attributes RegistrationReqDataAttributes, ) *RegistrationReqData`

NewRegistrationReqData instantiates a new RegistrationReqData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationReqDataWithDefaults

`func NewRegistrationReqDataWithDefaults() *RegistrationReqData`

NewRegistrationReqDataWithDefaults instantiates a new RegistrationReqData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RegistrationReqData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RegistrationReqData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RegistrationReqData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *RegistrationReqData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *RegistrationReqData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *RegistrationReqData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *RegistrationReqData) GetAttributes() RegistrationReqDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *RegistrationReqData) GetAttributesOk() (*RegistrationReqDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *RegistrationReqData) SetAttributes(v RegistrationReqDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


