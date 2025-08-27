# UserSessionAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserId** | **string** | user id | 
**Client** | **string** | client name and version | 
**Ip** | **string** | IP address | 
**CreatedAt** | **time.Time** | session creation date | 
**LastUsed** | **time.Time** | last used date | 

## Methods

### NewUserSessionAttributes

`func NewUserSessionAttributes(userId string, client string, ip string, createdAt time.Time, lastUsed time.Time, ) *UserSessionAttributes`

NewUserSessionAttributes instantiates a new UserSessionAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserSessionAttributesWithDefaults

`func NewUserSessionAttributesWithDefaults() *UserSessionAttributes`

NewUserSessionAttributesWithDefaults instantiates a new UserSessionAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserId

`func (o *UserSessionAttributes) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *UserSessionAttributes) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *UserSessionAttributes) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetClient

`func (o *UserSessionAttributes) GetClient() string`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *UserSessionAttributes) GetClientOk() (*string, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *UserSessionAttributes) SetClient(v string)`

SetClient sets Client field to given value.


### GetIp

`func (o *UserSessionAttributes) GetIp() string`

GetIp returns the Ip field if non-nil, zero value otherwise.

### GetIpOk

`func (o *UserSessionAttributes) GetIpOk() (*string, bool)`

GetIpOk returns a tuple with the Ip field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIp

`func (o *UserSessionAttributes) SetIp(v string)`

SetIp sets Ip field to given value.


### GetCreatedAt

`func (o *UserSessionAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *UserSessionAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *UserSessionAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetLastUsed

`func (o *UserSessionAttributes) GetLastUsed() time.Time`

GetLastUsed returns the LastUsed field if non-nil, zero value otherwise.

### GetLastUsedOk

`func (o *UserSessionAttributes) GetLastUsedOk() (*time.Time, bool)`

GetLastUsedOk returns a tuple with the LastUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUsed

`func (o *UserSessionAttributes) SetLastUsed(v time.Time)`

SetLastUsed sets LastUsed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


