# UserSessionsDataAttributesDevicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | session id | 
**DeviceName** | **string** | device name | 
**Client** | **string** | client name and version | 
**Location** | **string** | location | 
**LastUsed** | **time.Time** | last used date | 

## Methods

### NewUserSessionsDataAttributesDevicesInner

`func NewUserSessionsDataAttributesDevicesInner(id string, deviceName string, client string, location string, lastUsed time.Time, ) *UserSessionsDataAttributesDevicesInner`

NewUserSessionsDataAttributesDevicesInner instantiates a new UserSessionsDataAttributesDevicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserSessionsDataAttributesDevicesInnerWithDefaults

`func NewUserSessionsDataAttributesDevicesInnerWithDefaults() *UserSessionsDataAttributesDevicesInner`

NewUserSessionsDataAttributesDevicesInnerWithDefaults instantiates a new UserSessionsDataAttributesDevicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UserSessionsDataAttributesDevicesInner) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UserSessionsDataAttributesDevicesInner) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UserSessionsDataAttributesDevicesInner) SetId(v string)`

SetId sets Id field to given value.


### GetDeviceName

`func (o *UserSessionsDataAttributesDevicesInner) GetDeviceName() string`

GetDeviceName returns the DeviceName field if non-nil, zero value otherwise.

### GetDeviceNameOk

`func (o *UserSessionsDataAttributesDevicesInner) GetDeviceNameOk() (*string, bool)`

GetDeviceNameOk returns a tuple with the DeviceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceName

`func (o *UserSessionsDataAttributesDevicesInner) SetDeviceName(v string)`

SetDeviceName sets DeviceName field to given value.


### GetClient

`func (o *UserSessionsDataAttributesDevicesInner) GetClient() string`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *UserSessionsDataAttributesDevicesInner) GetClientOk() (*string, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *UserSessionsDataAttributesDevicesInner) SetClient(v string)`

SetClient sets Client field to given value.


### GetLocation

`func (o *UserSessionsDataAttributesDevicesInner) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *UserSessionsDataAttributesDevicesInner) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *UserSessionsDataAttributesDevicesInner) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetLastUsed

`func (o *UserSessionsDataAttributesDevicesInner) GetLastUsed() time.Time`

GetLastUsed returns the LastUsed field if non-nil, zero value otherwise.

### GetLastUsedOk

`func (o *UserSessionsDataAttributesDevicesInner) GetLastUsedOk() (*time.Time, bool)`

GetLastUsedOk returns a tuple with the LastUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUsed

`func (o *UserSessionsDataAttributesDevicesInner) SetLastUsed(v time.Time)`

SetLastUsed sets LastUsed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


