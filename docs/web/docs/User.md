# User

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**UserData**](UserData.md) |  | 

## Methods

### NewUser

`func NewUser(data UserData, ) *User`

NewUser instantiates a new User object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserWithDefaults

`func NewUserWithDefaults() *User`

NewUserWithDefaults instantiates a new User object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *User) GetData() UserData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *User) GetDataOk() (*UserData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *User) SetData(v UserData)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


