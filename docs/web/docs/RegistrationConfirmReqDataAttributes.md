# RegistrationConfirmReqDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FirstPassword** | **string** | First password | 
**SecondPassword** | **string** | Second password | 
**Email** | **string** | User email | 
**Username** | Pointer to **string** | User username | [optional] 

## Methods

### NewRegistrationConfirmReqDataAttributes

`func NewRegistrationConfirmReqDataAttributes(firstPassword string, secondPassword string, email string, ) *RegistrationConfirmReqDataAttributes`

NewRegistrationConfirmReqDataAttributes instantiates a new RegistrationConfirmReqDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationConfirmReqDataAttributesWithDefaults

`func NewRegistrationConfirmReqDataAttributesWithDefaults() *RegistrationConfirmReqDataAttributes`

NewRegistrationConfirmReqDataAttributesWithDefaults instantiates a new RegistrationConfirmReqDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFirstPassword

`func (o *RegistrationConfirmReqDataAttributes) GetFirstPassword() string`

GetFirstPassword returns the FirstPassword field if non-nil, zero value otherwise.

### GetFirstPasswordOk

`func (o *RegistrationConfirmReqDataAttributes) GetFirstPasswordOk() (*string, bool)`

GetFirstPasswordOk returns a tuple with the FirstPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstPassword

`func (o *RegistrationConfirmReqDataAttributes) SetFirstPassword(v string)`

SetFirstPassword sets FirstPassword field to given value.


### GetSecondPassword

`func (o *RegistrationConfirmReqDataAttributes) GetSecondPassword() string`

GetSecondPassword returns the SecondPassword field if non-nil, zero value otherwise.

### GetSecondPasswordOk

`func (o *RegistrationConfirmReqDataAttributes) GetSecondPasswordOk() (*string, bool)`

GetSecondPasswordOk returns a tuple with the SecondPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecondPassword

`func (o *RegistrationConfirmReqDataAttributes) SetSecondPassword(v string)`

SetSecondPassword sets SecondPassword field to given value.


### GetEmail

`func (o *RegistrationConfirmReqDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegistrationConfirmReqDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegistrationConfirmReqDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetUsername

`func (o *RegistrationConfirmReqDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *RegistrationConfirmReqDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *RegistrationConfirmReqDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *RegistrationConfirmReqDataAttributes) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


