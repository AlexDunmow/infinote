// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	db "boilerplate/db"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// AuthProvider is an autogenerated mock type for the AuthProvider type
type AuthProvider struct {
	mock.Mock
}

// GenerateJWT provides a mock function with given fields: ctx, user, userAgent
func (_m *AuthProvider) GenerateJWT(ctx context.Context, user *db.User, userAgent string) (string, error) {
	ret := _m.Called(ctx, user, userAgent)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *db.User, string) string); ok {
		r0 = rf(ctx, user, userAgent)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *db.User, string) error); ok {
		r1 = rf(ctx, user, userAgent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}