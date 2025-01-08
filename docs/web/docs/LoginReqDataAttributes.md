# LoginReqDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | Pointer to **string** | User email (required if username is not provided). | [optional] 
**Username** | Pointer to **string** | User username (required if email is not provided). | [optional] 
**Password** | **string** | User password. | 
**FactoryId** | **string** | Unique identifier for the user&#39;s device. | 
**DeviceName** | **string** | Human-readable name for the user&#39;s device (e.g., &#39;iPhone 14&#39;). | 
**OsVersion** | **string** | Operating system version of the user&#39;s device. | 

## Methods

### NewLoginReqDataAttributes

`func NewLoginReqDataAttributes(password string, factoryId string, deviceName string, osVersion string, ) *LoginReqDataAttributes`

NewLoginReqDataAttributes instantiates a new LoginReqDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginReqDataAttributesWithDefaults

`func NewLoginReqDataAttributesWithDefaults() *LoginReqDataAttributes`

NewLoginReqDataAttributesWithDefaults instantiates a new LoginReqDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *LoginReqDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *LoginReqDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *LoginReqDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *LoginReqDataAttributes) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetUsername

`func (o *LoginReqDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *LoginReqDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *LoginReqDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *LoginReqDataAttributes) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetPassword

`func (o *LoginReqDataAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *LoginReqDataAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *LoginReqDataAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetFactoryId

`func (o *LoginReqDataAttributes) GetFactoryId() string`

GetFactoryId returns the FactoryId field if non-nil, zero value otherwise.

### GetFactoryIdOk

`func (o *LoginReqDataAttributes) GetFactoryIdOk() (*string, bool)`

GetFactoryIdOk returns a tuple with the FactoryId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFactoryId

`func (o *LoginReqDataAttributes) SetFactoryId(v string)`

SetFactoryId sets FactoryId field to given value.


### GetDeviceName

`func (o *LoginReqDataAttributes) GetDeviceName() string`

GetDeviceName returns the DeviceName field if non-nil, zero value otherwise.

### GetDeviceNameOk

`func (o *LoginReqDataAttributes) GetDeviceNameOk() (*string, bool)`

GetDeviceNameOk returns a tuple with the DeviceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceName

`func (o *LoginReqDataAttributes) SetDeviceName(v string)`

SetDeviceName sets DeviceName field to given value.


### GetOsVersion

`func (o *LoginReqDataAttributes) GetOsVersion() string`

GetOsVersion returns the OsVersion field if non-nil, zero value otherwise.

### GetOsVersionOk

`func (o *LoginReqDataAttributes) GetOsVersionOk() (*string, bool)`

GetOsVersionOk returns a tuple with the OsVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOsVersion

`func (o *LoginReqDataAttributes) SetOsVersion(v string)`

SetOsVersion sets OsVersion field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


