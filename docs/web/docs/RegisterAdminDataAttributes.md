# RegisterAdminDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** | The user&#39;s email address. | 
**Password** | **string** | The user&#39;s password. | 
**Role** | **string** | The role assigned to the registering user. | 

## Methods

### NewRegisterAdminDataAttributes

`func NewRegisterAdminDataAttributes(email string, password string, role string, ) *RegisterAdminDataAttributes`

NewRegisterAdminDataAttributes instantiates a new RegisterAdminDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterAdminDataAttributesWithDefaults

`func NewRegisterAdminDataAttributesWithDefaults() *RegisterAdminDataAttributes`

NewRegisterAdminDataAttributesWithDefaults instantiates a new RegisterAdminDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *RegisterAdminDataAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *RegisterAdminDataAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *RegisterAdminDataAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *RegisterAdminDataAttributes) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *RegisterAdminDataAttributes) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *RegisterAdminDataAttributes) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetRole

`func (o *RegisterAdminDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *RegisterAdminDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *RegisterAdminDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


