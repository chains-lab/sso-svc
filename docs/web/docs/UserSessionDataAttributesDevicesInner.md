# UserSessionDataAttributesDevicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | session id | 
**DeviceName** | **string** | device name | 
**Client** | **string** | client name and version | 
**Location** | **string** | location | 
**LastUsed** | **time.Time** | last used date | 

## Methods

### NewUserSessionDataAttributesDevicesInner

`func NewUserSessionDataAttributesDevicesInner(id string, deviceName string, client string, location string, lastUsed time.Time, ) *UserSessionDataAttributesDevicesInner`

NewUserSessionDataAttributesDevicesInner instantiates a new UserSessionDataAttributesDevicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserSessionDataAttributesDevicesInnerWithDefaults

`func NewUserSessionDataAttributesDevicesInnerWithDefaults() *UserSessionDataAttributesDevicesInner`

NewUserSessionDataAttributesDevicesInnerWithDefaults instantiates a new UserSessionDataAttributesDevicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UserSessionDataAttributesDevicesInner) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UserSessionDataAttributesDevicesInner) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UserSessionDataAttributesDevicesInner) SetId(v string)`

SetId sets Id field to given value.


### GetDeviceName

`func (o *UserSessionDataAttributesDevicesInner) GetDeviceName() string`

GetDeviceName returns the DeviceName field if non-nil, zero value otherwise.

### GetDeviceNameOk

`func (o *UserSessionDataAttributesDevicesInner) GetDeviceNameOk() (*string, bool)`

GetDeviceNameOk returns a tuple with the DeviceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceName

`func (o *UserSessionDataAttributesDevicesInner) SetDeviceName(v string)`

SetDeviceName sets DeviceName field to given value.


### GetClient

`func (o *UserSessionDataAttributesDevicesInner) GetClient() string`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *UserSessionDataAttributesDevicesInner) GetClientOk() (*string, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *UserSessionDataAttributesDevicesInner) SetClient(v string)`

SetClient sets Client field to given value.


### GetLocation

`func (o *UserSessionDataAttributesDevicesInner) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *UserSessionDataAttributesDevicesInner) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *UserSessionDataAttributesDevicesInner) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetLastUsed

`func (o *UserSessionDataAttributesDevicesInner) GetLastUsed() time.Time`

GetLastUsed returns the LastUsed field if non-nil, zero value otherwise.

### GetLastUsedOk

`func (o *UserSessionDataAttributesDevicesInner) GetLastUsedOk() (*time.Time, bool)`

GetLastUsedOk returns a tuple with the LastUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUsed

`func (o *UserSessionDataAttributesDevicesInner) SetLastUsed(v time.Time)`

SetLastUsed sets LastUsed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


