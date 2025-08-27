# RegistrationDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | The user&#39;s email address. | 
**Password** | **string** | The user&#39;s password. | 
**ConfirmPassword** | **string** | Confirmation of the user&#39;s password. Must match the password field. | 

## Methods

### NewRegistrationDataAttributes

`func NewRegistrationDataAttributes(email string, password string, confirmPassword string, ) *RegistrationDataAttributes`

NewRegistrationDataAttributes instantiates a new RegistrationDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationDataAttributesWithDefaults

`func NewRegistrationDataAttributesWithDefaults() *RegistrationDataAttributes`

NewRegistrationDataAttributesWithDefaults instantiates a new RegistrationDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *RegistrationDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegistrationDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegistrationDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *RegistrationDataAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *RegistrationDataAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *RegistrationDataAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetConfirmPassword

`func (o *RegistrationDataAttributes) GetConfirmPassword() string`

GetConfirmPassword returns the ConfirmPassword field if non-nil, zero value otherwise.

### GetConfirmPasswordOk

`func (o *RegistrationDataAttributes) GetConfirmPasswordOk() (*string, bool)`

GetConfirmPasswordOk returns a tuple with the ConfirmPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfirmPassword

`func (o *RegistrationDataAttributes) SetConfirmPassword(v string)`

SetConfirmPassword sets ConfirmPassword field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


