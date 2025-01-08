# RegistrationReqDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | User email | 
**Username** | Pointer to **string** | User username | [optional] 

## Methods

### NewRegistrationReqDataAttributes

`func NewRegistrationReqDataAttributes(email string, ) *RegistrationReqDataAttributes`

NewRegistrationReqDataAttributes instantiates a new RegistrationReqDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegistrationReqDataAttributesWithDefaults

`func NewRegistrationReqDataAttributesWithDefaults() *RegistrationReqDataAttributes`

NewRegistrationReqDataAttributesWithDefaults instantiates a new RegistrationReqDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *RegistrationReqDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegistrationReqDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegistrationReqDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetUsername

`func (o *RegistrationReqDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *RegistrationReqDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *RegistrationReqDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *RegistrationReqDataAttributes) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


