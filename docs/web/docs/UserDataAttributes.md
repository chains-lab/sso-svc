# UserDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | Email | 
**Role** | **string** | Role | 
**EmailVerified** | **bool** | Email Verified | 
**CreatedAt** | **time.Time** | Created At | 
**CityId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | City ID | [optional] 
**CityRole** | Pointer to **string** | Admin | [optional] 
**CompanyId** | Pointer to [**uuid.UUID**](uuid.UUID.md) | Company ID | [optional] 
**CompanyRole** | Pointer to **string** | Moder | [optional] 
**UpdatedAt** | **time.Time** | Updated At | 

## Methods

### NewUserDataAttributes

`func NewUserDataAttributes(email string, role string, emailVerified bool, createdAt time.Time, updatedAt time.Time, ) *UserDataAttributes`

NewUserDataAttributes instantiates a new UserDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserDataAttributesWithDefaults

`func NewUserDataAttributesWithDefaults() *UserDataAttributes`

NewUserDataAttributesWithDefaults instantiates a new UserDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *UserDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *UserDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *UserDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetRole

`func (o *UserDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *UserDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *UserDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.


### GetEmailVerified

`func (o *UserDataAttributes) GetEmailVerified() bool`

GetEmailVerified returns the EmailVerified field if non-nil, zero value otherwise.

### GetEmailVerifiedOk

`func (o *UserDataAttributes) GetEmailVerifiedOk() (*bool, bool)`

GetEmailVerifiedOk returns a tuple with the EmailVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailVerified

`func (o *UserDataAttributes) SetEmailVerified(v bool)`

SetEmailVerified sets EmailVerified field to given value.


### GetCreatedAt

`func (o *UserDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *UserDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *UserDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetCityId

`func (o *UserDataAttributes) GetCityId() uuid.UUID`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *UserDataAttributes) GetCityIdOk() (*uuid.UUID, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *UserDataAttributes) SetCityId(v uuid.UUID)`

SetCityId sets CityId field to given value.

### HasCityId

`func (o *UserDataAttributes) HasCityId() bool`

HasCityId returns a boolean if a field has been set.

### GetCityRole

`func (o *UserDataAttributes) GetCityRole() string`

GetCityRole returns the CityRole field if non-nil, zero value otherwise.

### GetCityRoleOk

`func (o *UserDataAttributes) GetCityRoleOk() (*string, bool)`

GetCityRoleOk returns a tuple with the CityRole field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityRole

`func (o *UserDataAttributes) SetCityRole(v string)`

SetCityRole sets CityRole field to given value.

### HasCityRole

`func (o *UserDataAttributes) HasCityRole() bool`

HasCityRole returns a boolean if a field has been set.

### GetCompanyId

`func (o *UserDataAttributes) GetCompanyId() uuid.UUID`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *UserDataAttributes) GetCompanyIdOk() (*uuid.UUID, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *UserDataAttributes) SetCompanyId(v uuid.UUID)`

SetCompanyId sets CompanyId field to given value.

### HasCompanyId

`func (o *UserDataAttributes) HasCompanyId() bool`

HasCompanyId returns a boolean if a field has been set.

### GetCompanyRole

`func (o *UserDataAttributes) GetCompanyRole() string`

GetCompanyRole returns the CompanyRole field if non-nil, zero value otherwise.

### GetCompanyRoleOk

`func (o *UserDataAttributes) GetCompanyRoleOk() (*string, bool)`

GetCompanyRoleOk returns a tuple with the CompanyRole field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyRole

`func (o *UserDataAttributes) SetCompanyRole(v string)`

SetCompanyRole sets CompanyRole field to given value.

### HasCompanyRole

`func (o *UserDataAttributes) HasCompanyRole() bool`

HasCompanyRole returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *UserDataAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *UserDataAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *UserDataAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


