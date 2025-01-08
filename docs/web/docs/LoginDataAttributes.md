# LoginDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | User email | 
**DeviceName** | **string** | Device name | 

## Methods

### NewLoginDataAttributes

`func NewLoginDataAttributes(email string, deviceName string, ) *LoginDataAttributes`

NewLoginDataAttributes instantiates a new LoginDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginDataAttributesWithDefaults

`func NewLoginDataAttributesWithDefaults() *LoginDataAttributes`

NewLoginDataAttributesWithDefaults instantiates a new LoginDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *LoginDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *LoginDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *LoginDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetDeviceName

`func (o *LoginDataAttributes) GetDeviceName() string`

GetDeviceName returns the DeviceName field if non-nil, zero value otherwise.

### GetDeviceNameOk

`func (o *LoginDataAttributes) GetDeviceNameOk() (*string, bool)`

GetDeviceNameOk returns a tuple with the DeviceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceName

`func (o *LoginDataAttributes) SetDeviceName(v string)`

SetDeviceName sets DeviceName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


