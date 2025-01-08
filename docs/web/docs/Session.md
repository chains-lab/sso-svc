# Session

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | session id | 
**Client** | **string** | client name and version | 
**IpFirst** | **string** | IP address | 
**IpLast** | **string** | IP address | 
**LastUsed** | **time.Time** | last used date | 
**CreatedAt** | **time.Time** | session creation date | 

## Methods

### NewSession

`func NewSession(id string, client string, ipFirst string, ipLast string, lastUsed time.Time, createdAt time.Time, ) *Session`

NewSession instantiates a new Session object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSessionWithDefaults

`func NewSessionWithDefaults() *Session`

NewSessionWithDefaults instantiates a new Session object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Session) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Session) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Session) SetId(v string)`

SetId sets Id field to given value.


### GetClient

`func (o *Session) GetClient() string`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *Session) GetClientOk() (*string, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *Session) SetClient(v string)`

SetClient sets Client field to given value.


### GetIpFirst

`func (o *Session) GetIpFirst() string`

GetIpFirst returns the IpFirst field if non-nil, zero value otherwise.

### GetIpFirstOk

`func (o *Session) GetIpFirstOk() (*string, bool)`

GetIpFirstOk returns a tuple with the IpFirst field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpFirst

`func (o *Session) SetIpFirst(v string)`

SetIpFirst sets IpFirst field to given value.


### GetIpLast

`func (o *Session) GetIpLast() string`

GetIpLast returns the IpLast field if non-nil, zero value otherwise.

### GetIpLastOk

`func (o *Session) GetIpLastOk() (*string, bool)`

GetIpLastOk returns a tuple with the IpLast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpLast

`func (o *Session) SetIpLast(v string)`

SetIpLast sets IpLast field to given value.


### GetLastUsed

`func (o *Session) GetLastUsed() time.Time`

GetLastUsed returns the LastUsed field if non-nil, zero value otherwise.

### GetLastUsedOk

`func (o *Session) GetLastUsedOk() (*time.Time, bool)`

GetLastUsedOk returns a tuple with the LastUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUsed

`func (o *Session) SetLastUsed(v time.Time)`

SetLastUsed sets LastUsed field to given value.


### GetCreatedAt

`func (o *Session) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Session) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Session) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


