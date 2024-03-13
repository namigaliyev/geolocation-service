// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "geolocation-service/internal/model"
)

// NoSQLDb is an autogenerated mock type for the NoSQLDb type
type NoSQLDb struct {
	mock.Mock
}

// CheckHealth provides a mock function with given fields: ctx
func (_m *NoSQLDb) CheckHealth(ctx context.Context) bool {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CheckHealth")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *NoSQLDb) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllGeolocations provides a mock function with given fields: ctx
func (_m *NoSQLDb) DeleteAllGeolocations(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAllGeolocations")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMetricByProcessID provides a mock function with given fields: ctx, processID
func (_m *NoSQLDb) DeleteMetricByProcessID(ctx context.Context, processID string) error {
	ret := _m.Called(ctx, processID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMetricByProcessID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, processID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetGeolocationByIpAddress provides a mock function with given fields: ctx, ipAddress
func (_m *NoSQLDb) GetGeolocationByIpAddress(ctx context.Context, ipAddress string) (*model.Geolocation, error) {
	ret := _m.Called(ctx, ipAddress)

	if len(ret) == 0 {
		panic("no return value specified for GetGeolocationByIpAddress")
	}

	var r0 *model.Geolocation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Geolocation, error)); ok {
		return rf(ctx, ipAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Geolocation); ok {
		r0 = rf(ctx, ipAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Geolocation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ipAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGeolocations provides a mock function with given fields: ctx
func (_m *NoSQLDb) GetGeolocations(ctx context.Context) ([]*model.Geolocation, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetGeolocations")
	}

	var r0 []*model.Geolocation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.Geolocation, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Geolocation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Geolocation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetricByProcessID provides a mock function with given fields: ctx, processID
func (_m *NoSQLDb) GetMetricByProcessID(ctx context.Context, processID string) (*model.Metric, error) {
	ret := _m.Called(ctx, processID)

	if len(ret) == 0 {
		panic("no return value specified for GetMetricByProcessID")
	}

	var r0 *model.Metric
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Metric, error)); ok {
		return rf(ctx, processID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Metric); ok {
		r0 = rf(ctx, processID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Metric)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, processID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPageableGeolocations provides a mock function with given fields: ctx, offset, limit
func (_m *NoSQLDb) GetPageableGeolocations(ctx context.Context, offset int, limit int) (*model.PageResponse, error) {
	ret := _m.Called(ctx, offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetPageableGeolocations")
	}

	var r0 *model.PageResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) (*model.PageResponse, error)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) *model.PageResponse); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PageResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveGeolocations provides a mock function with given fields: ctx, locations
func (_m *NoSQLDb) SaveGeolocations(ctx context.Context, locations []*model.Geolocation) error {
	ret := _m.Called(ctx, locations)

	if len(ret) == 0 {
		panic("no return value specified for SaveGeolocations")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Geolocation) error); ok {
		r0 = rf(ctx, locations)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveMetric provides a mock function with given fields: ctx, metrics
func (_m *NoSQLDb) SaveMetric(ctx context.Context, metrics *model.Metric) error {
	ret := _m.Called(ctx, metrics)

	if len(ret) == 0 {
		panic("no return value specified for SaveMetric")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Metric) error); ok {
		r0 = rf(ctx, metrics)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMetric provides a mock function with given fields: ctx, metric
func (_m *NoSQLDb) UpdateMetric(ctx context.Context, metric *model.Metric) error {
	ret := _m.Called(ctx, metric)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMetric")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Metric) error); ok {
		r0 = rf(ctx, metric)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewNoSQLDb creates a new instance of NoSQLDb. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNoSQLDb(t interface {
	mock.TestingT
	Cleanup(func())
}) *NoSQLDb {
	mock := &NoSQLDb{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
