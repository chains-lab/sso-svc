# ResetPasswordCompleteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | User email (required if username is not provided) | 
**Password** | **string** | New user password | 

## Methods

### NewResetPasswordCompleteDataAttributes

`func NewResetPasswordCompleteDataAttributes(email string, password string, ) *ResetPasswordCompleteDataAttributes`

NewResetPasswordCompleteDataAttributes instantiates a new ResetPasswordCompleteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResetPasswordCompleteDataAttributesWithDefaults

`func NewResetPasswordCompleteDataAttributesWithDefaults() *ResetPasswordCompleteDataAttributes`

NewResetPasswordCompleteDataAttributesWithDefaults instantiates a new ResetPasswordCompleteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *ResetPasswordCompleteDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *ResetPasswordCompleteDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *ResetPasswordCompleteDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *ResetPasswordCompleteDataAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *ResetPasswordCompleteDataAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *ResetPasswordCompleteDataAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


