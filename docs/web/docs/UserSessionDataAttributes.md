# UserSessionDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Client** | **string** | client name and version | 
**IpFirst** | **string** | IP address | 
**IpLast** | **string** | IP address | 
**LastUsed** | **time.Time** | last used date | 
**CreatedAt** | **time.Time** | session creation date | 

## Methods

### NewUserSessionDataAttributes

`func NewUserSessionDataAttributes(client string, ipFirst string, ipLast string, lastUsed time.Time, createdAt time.Time, ) *UserSessionDataAttributes`

NewUserSessionDataAttributes instantiates a new UserSessionDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserSessionDataAttributesWithDefaults

`func NewUserSessionDataAttributesWithDefaults() *UserSessionDataAttributes`

NewUserSessionDataAttributesWithDefaults instantiates a new UserSessionDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClient

`func (o *UserSessionDataAttributes) GetClient() string`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *UserSessionDataAttributes) GetClientOk() (*string, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *UserSessionDataAttributes) SetClient(v string)`

SetClient sets Client field to given value.


### GetIpFirst

`func (o *UserSessionDataAttributes) GetIpFirst() string`

GetIpFirst returns the IpFirst field if non-nil, zero value otherwise.

### GetIpFirstOk

`func (o *UserSessionDataAttributes) GetIpFirstOk() (*string, bool)`

GetIpFirstOk returns a tuple with the IpFirst field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpFirst

`func (o *UserSessionDataAttributes) SetIpFirst(v string)`

SetIpFirst sets IpFirst field to given value.


### GetIpLast

`func (o *UserSessionDataAttributes) GetIpLast() string`

GetIpLast returns the IpLast field if non-nil, zero value otherwise.

### GetIpLastOk

`func (o *UserSessionDataAttributes) GetIpLastOk() (*string, bool)`

GetIpLastOk returns a tuple with the IpLast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpLast

`func (o *UserSessionDataAttributes) SetIpLast(v string)`

SetIpLast sets IpLast field to given value.


### GetLastUsed

`func (o *UserSessionDataAttributes) GetLastUsed() time.Time`

GetLastUsed returns the LastUsed field if non-nil, zero value otherwise.

### GetLastUsedOk

`func (o *UserSessionDataAttributes) GetLastUsedOk() (*time.Time, bool)`

GetLastUsedOk returns a tuple with the LastUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUsed

`func (o *UserSessionDataAttributes) SetLastUsed(v time.Time)`

SetLastUsed sets LastUsed field to given value.


### GetCreatedAt

`func (o *UserSessionDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *UserSessionDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *UserSessionDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


