# SessionData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | session id | 
**Type** | **string** |  | 
**Attributes** | [**SessionAttributes**](SessionAttributes.md) |  | 

## Methods

### NewSessionData

`func NewSessionData(id string, type_ string, attributes SessionAttributes, ) *SessionData`

NewSessionData instantiates a new SessionData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSessionDataWithDefaults

`func NewSessionDataWithDefaults() *SessionData`

NewSessionDataWithDefaults instantiates a new SessionData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SessionData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SessionData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SessionData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *SessionData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *SessionData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *SessionData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *SessionData) GetAttributes() SessionAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *SessionData) GetAttributesOk() (*SessionAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *SessionData) SetAttributes(v SessionAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


