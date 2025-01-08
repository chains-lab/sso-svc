# ForgotPasswordReqDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | Pointer to **string** | User email (required if username is not provided). | [optional] 
**Username** | Pointer to **string** | User username (required if email is not provided). | [optional] 

## Methods

### NewForgotPasswordReqDataAttributes

`func NewForgotPasswordReqDataAttributes() *ForgotPasswordReqDataAttributes`

NewForgotPasswordReqDataAttributes instantiates a new ForgotPasswordReqDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewForgotPasswordReqDataAttributesWithDefaults

`func NewForgotPasswordReqDataAttributesWithDefaults() *ForgotPasswordReqDataAttributes`

NewForgotPasswordReqDataAttributesWithDefaults instantiates a new ForgotPasswordReqDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *ForgotPasswordReqDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *ForgotPasswordReqDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *ForgotPasswordReqDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *ForgotPasswordReqDataAttributes) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetUsername

`func (o *ForgotPasswordReqDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *ForgotPasswordReqDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *ForgotPasswordReqDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *ForgotPasswordReqDataAttributes) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


