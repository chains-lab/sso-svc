# LoginReqData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Attributes** | [**LoginReqDataAttributes**](LoginReqDataAttributes.md) |  | 

## Methods

### NewLoginReqData

`func NewLoginReqData(type_ string, attributes LoginReqDataAttributes, ) *LoginReqData`

NewLoginReqData instantiates a new LoginReqData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginReqDataWithDefaults

`func NewLoginReqDataWithDefaults() *LoginReqData`

NewLoginReqDataWithDefaults instantiates a new LoginReqData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *LoginReqData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *LoginReqData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *LoginReqData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *LoginReqData) GetAttributes() LoginReqDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *LoginReqData) GetAttributesOk() (*LoginReqDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *LoginReqData) SetAttributes(v LoginReqDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


