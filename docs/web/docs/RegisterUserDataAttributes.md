# RegisterUserDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | The user&#39;s email address. | 
**Password** | **string** | The user&#39;s password. | 
**ConfirmPassword** | **string** | Confirmation of the user&#39;s password. Must match the password field. | 

## Methods

### NewRegisterUserDataAttributes

`func NewRegisterUserDataAttributes(email string, password string, confirmPassword string, ) *RegisterUserDataAttributes`

NewRegisterUserDataAttributes instantiates a new RegisterUserDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterUserDataAttributesWithDefaults

`func NewRegisterUserDataAttributesWithDefaults() *RegisterUserDataAttributes`

NewRegisterUserDataAttributesWithDefaults instantiates a new RegisterUserDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *RegisterUserDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegisterUserDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegisterUserDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *RegisterUserDataAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *RegisterUserDataAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *RegisterUserDataAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetConfirmPassword

`func (o *RegisterUserDataAttributes) GetConfirmPassword() string`

GetConfirmPassword returns the ConfirmPassword field if non-nil, zero value otherwise.

### GetConfirmPasswordOk

`func (o *RegisterUserDataAttributes) GetConfirmPasswordOk() (*string, bool)`

GetConfirmPasswordOk returns a tuple with the ConfirmPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfirmPassword

`func (o *RegisterUserDataAttributes) SetConfirmPassword(v string)`

SetConfirmPassword sets ConfirmPassword field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


