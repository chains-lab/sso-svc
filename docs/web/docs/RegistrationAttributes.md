# RegistrationAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | Pointer to **string** | Unique username | [optional] 
**Email** | Pointer to **string** | User email | [optional] 
**Password** | Pointer to **string** | User password | [optional] 

## Methods

### NewRegistrationAttributes

`func NewRegistrationAttributes() *RegistrationAttributes`

NewRegistrationAttributes instantiates a new RegistrationAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationAttributesWithDefaults

`func NewRegistrationAttributesWithDefaults() *RegistrationAttributes`

NewRegistrationAttributesWithDefaults instantiates a new RegistrationAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *RegistrationAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *RegistrationAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *RegistrationAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *RegistrationAttributes) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetEmail

`func (o *RegistrationAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegistrationAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegistrationAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *RegistrationAttributes) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetPassword

`func (o *RegistrationAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *RegistrationAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *RegistrationAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *RegistrationAttributes) HasPassword() bool`

HasPassword returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


