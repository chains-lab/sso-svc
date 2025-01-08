# ChangeUsernameDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NewUsername** | Pointer to **string** | New username | [optional] 
**Password** | **string** | User password | 

## Methods

### NewChangeUsernameDataAttributes

`func NewChangeUsernameDataAttributes(password string, ) *ChangeUsernameDataAttributes`

NewChangeUsernameDataAttributes instantiates a new ChangeUsernameDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChangeUsernameDataAttributesWithDefaults

`func NewChangeUsernameDataAttributesWithDefaults() *ChangeUsernameDataAttributes`

NewChangeUsernameDataAttributesWithDefaults instantiates a new ChangeUsernameDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNewUsername

`func (o *ChangeUsernameDataAttributes) GetNewUsername() string`

GetNewUsername returns the NewUsername field if non-nil, zero value otherwise.

### GetNewUsernameOk

`func (o *ChangeUsernameDataAttributes) GetNewUsernameOk() (*string, bool)`

GetNewUsernameOk returns a tuple with the NewUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewUsername

`func (o *ChangeUsernameDataAttributes) SetNewUsername(v string)`

SetNewUsername sets NewUsername field to given value.

### HasNewUsername

`func (o *ChangeUsernameDataAttributes) HasNewUsername() bool`

HasNewUsername returns a boolean if a field has been set.

### GetPassword

`func (o *ChangeUsernameDataAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *ChangeUsernameDataAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *ChangeUsernameDataAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


