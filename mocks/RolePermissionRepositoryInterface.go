// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	models "myapp/models"

	mock "github.com/stretchr/testify/mock"
)

// RolePermissionRepositoryInterface is an autogenerated mock type for the RolePermissionRepositoryInterface type
type RolePermissionRepositoryInterface struct {
	mock.Mock
}

// AssignRoleToUserRequest provides a mock function with given fields: req
func (_m *RolePermissionRepositoryInterface) AssignRoleToUserRequest(req models.AssignRoleToUserRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for AssignRoleToUserRequest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(models.AssignRoleToUserRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindUserRole provides a mock function with given fields: userID
func (_m *RolePermissionRepositoryInterface) FindUserRole(userID int64) (models.FindUserRoleResponse, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for FindUserRole")
	}

	var r0 models.FindUserRoleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (models.FindUserRoleResponse, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int64) models.FindUserRoleResponse); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(models.FindUserRoleResponse)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRolePermissionRepositoryInterface creates a new instance of RolePermissionRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRolePermissionRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *RolePermissionRepositoryInterface {
	mock := &RolePermissionRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
