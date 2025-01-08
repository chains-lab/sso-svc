# RegistrationConfirmDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FirstPassword** | **string** | First password | 
**SecondPassword** | **string** | Second password | 
**Email** | **string** | User email | 
**Username** | Pointer to **string** | User username | [optional] 

## Methods

### NewRegistrationConfirmDataAttributes

`func NewRegistrationConfirmDataAttributes(firstPassword string, secondPassword string, email string, ) *RegistrationConfirmDataAttributes`

NewRegistrationConfirmDataAttributes instantiates a new RegistrationConfirmDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationConfirmDataAttributesWithDefaults

`func NewRegistrationConfirmDataAttributesWithDefaults() *RegistrationConfirmDataAttributes`

NewRegistrationConfirmDataAttributesWithDefaults instantiates a new RegistrationConfirmDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFirstPassword

`func (o *RegistrationConfirmDataAttributes) GetFirstPassword() string`

GetFirstPassword returns the FirstPassword field if non-nil, zero value otherwise.

### GetFirstPasswordOk

`func (o *RegistrationConfirmDataAttributes) GetFirstPasswordOk() (*string, bool)`

GetFirstPasswordOk returns a tuple with the FirstPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstPassword

`func (o *RegistrationConfirmDataAttributes) SetFirstPassword(v string)`

SetFirstPassword sets FirstPassword field to given value.


### GetSecondPassword

`func (o *RegistrationConfirmDataAttributes) GetSecondPassword() string`

GetSecondPassword returns the SecondPassword field if non-nil, zero value otherwise.

### GetSecondPasswordOk

`func (o *RegistrationConfirmDataAttributes) GetSecondPasswordOk() (*string, bool)`

GetSecondPasswordOk returns a tuple with the SecondPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecondPassword

`func (o *RegistrationConfirmDataAttributes) SetSecondPassword(v string)`

SetSecondPassword sets SecondPassword field to given value.


### GetEmail

`func (o *RegistrationConfirmDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegistrationConfirmDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegistrationConfirmDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetUsername

`func (o *RegistrationConfirmDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *RegistrationConfirmDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *RegistrationConfirmDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *RegistrationConfirmDataAttributes) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


