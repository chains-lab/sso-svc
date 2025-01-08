# ResetPasswordDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | Pointer to **string** | User email (required if username is not provided). | [optional] 
**Username** | Pointer to **string** | User username (required if email is not provided). | [optional] 
**FirstPassword** | **string** | User password. | 
**SecondPassword** | **string** | User password. | 

## Methods

### NewResetPasswordDataAttributes

`func NewResetPasswordDataAttributes(firstPassword string, secondPassword string, ) *ResetPasswordDataAttributes`

NewResetPasswordDataAttributes instantiates a new ResetPasswordDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResetPasswordDataAttributesWithDefaults

`func NewResetPasswordDataAttributesWithDefaults() *ResetPasswordDataAttributes`

NewResetPasswordDataAttributesWithDefaults instantiates a new ResetPasswordDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *ResetPasswordDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *ResetPasswordDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *ResetPasswordDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *ResetPasswordDataAttributes) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetUsername

`func (o *ResetPasswordDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *ResetPasswordDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *ResetPasswordDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *ResetPasswordDataAttributes) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetFirstPassword

`func (o *ResetPasswordDataAttributes) GetFirstPassword() string`

GetFirstPassword returns the FirstPassword field if non-nil, zero value otherwise.

### GetFirstPasswordOk

`func (o *ResetPasswordDataAttributes) GetFirstPasswordOk() (*string, bool)`

GetFirstPasswordOk returns a tuple with the FirstPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstPassword

`func (o *ResetPasswordDataAttributes) SetFirstPassword(v string)`

SetFirstPassword sets FirstPassword field to given value.


### GetSecondPassword

`func (o *ResetPasswordDataAttributes) GetSecondPassword() string`

GetSecondPassword returns the SecondPassword field if non-nil, zero value otherwise.

### GetSecondPasswordOk

`func (o *ResetPasswordDataAttributes) GetSecondPasswordOk() (*string, bool)`

GetSecondPasswordOk returns a tuple with the SecondPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecondPassword

`func (o *ResetPasswordDataAttributes) SetSecondPassword(v string)`

SetSecondPassword sets SecondPassword field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


