# RegistrationUserDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | The user&#39;s email address. | 
**Password** | **string** | The user&#39;s password. | 
**ConfirmPassword** | **string** | Confirmation of the user&#39;s password. Must match the password field. | 

## Methods

### NewRegistrationUserDataAttributes

`func NewRegistrationUserDataAttributes(email string, password string, confirmPassword string, ) *RegistrationUserDataAttributes`

NewRegistrationUserDataAttributes instantiates a new RegistrationUserDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationUserDataAttributesWithDefaults

`func NewRegistrationUserDataAttributesWithDefaults() *RegistrationUserDataAttributes`

NewRegistrationUserDataAttributesWithDefaults instantiates a new RegistrationUserDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *RegistrationUserDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegistrationUserDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegistrationUserDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *RegistrationUserDataAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *RegistrationUserDataAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *RegistrationUserDataAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetConfirmPassword

`func (o *RegistrationUserDataAttributes) GetConfirmPassword() string`

GetConfirmPassword returns the ConfirmPassword field if non-nil, zero value otherwise.

### GetConfirmPasswordOk

`func (o *RegistrationUserDataAttributes) GetConfirmPasswordOk() (*string, bool)`

GetConfirmPasswordOk returns a tuple with the ConfirmPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfirmPassword

`func (o *RegistrationUserDataAttributes) SetConfirmPassword(v string)`

SetConfirmPassword sets ConfirmPassword field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


