package devicehealth

import (
	"testing"

	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testAppendOptions(t *testing.T, source []health.Option) {
	var (
		assert = assert.New(t)
		result = AppendOptions(source)

		expectedOptions = map[health.Option]bool{
			DeviceCount:                      true,
			TotalWRPRequestResponseProcessed: true,
			TotalPingMessagesReceived:        true,
			TotalPongMessagesReceived:        true,
			TotalConnectionEvents:            true,
			TotalDisconnectionEvents:         true,
		}
	)

	assert.NotEqual(source, result)

	for _, o := range result {
		delete(expectedOptions, o)
	}

	assert.Empty(expectedOptions)
}

func TestAppendOptions(t *testing.T) {
	t.Run("NilSlice", func(t *testing.T) { testAppendOptions(t, nil) })
	t.Run("EmptySlice", func(t *testing.T) { testAppendOptions(t, []health.Option{}) })
	t.Run("ExistingSlice", func(t *testing.T) { testAppendOptions(t, []health.Option{health.Stat("existing")}) })
}

func testListenerOnDeviceEventConnect(t *testing.T) {
	var (
		assert     = assert.New(t)
		dispatcher = new(mockDispatcher)
		listener   = &Listener{Dispatcher: dispatcher}

		expectedStats = health.Stats{
			DeviceCount:           1,
			TotalConnectionEvents: 1,
		}

		actualStats = health.Stats{}
	)

	dispatcher.On("SendEvent", mock.AnythingOfType("health.HealthFunc")).Once().
		Run(func(arguments mock.Arguments) {
			hf := arguments.Get(0).(health.HealthFunc)
			hf(actualStats)
		})

	listener.OnDeviceEvent(&device.Event{Type: device.Connect})
	assert.Equal(expectedStats, actualStats)

	dispatcher.AssertExpectations(t)
}

func testListenerOnDeviceEventDisconnect(t *testing.T) {
	var (
		assert     = assert.New(t)
		dispatcher = new(mockDispatcher)
		listener   = &Listener{Dispatcher: dispatcher}

		expectedStats = health.Stats{
			DeviceCount:              0,
			TotalConnectionEvents:    1,
			TotalDisconnectionEvents: 1,
		}

		actualStats = health.Stats{
			DeviceCount:           1,
			TotalConnectionEvents: 1,
		}
	)

	dispatcher.On("SendEvent", mock.AnythingOfType("health.HealthFunc")).Once().
		Run(func(arguments mock.Arguments) {
			hf := arguments.Get(0).(health.HealthFunc)
			hf(actualStats)
		})

	listener.OnDeviceEvent(&device.Event{Type: device.Disconnect})
	assert.Equal(expectedStats, actualStats)

	dispatcher.AssertExpectations(t)
}

func testListenerOnDeviceEventTransactionComplete(t *testing.T) {
	var (
		assert     = assert.New(t)
		dispatcher = new(mockDispatcher)
		listener   = &Listener{Dispatcher: dispatcher}

		expectedStats = health.Stats{
			TotalWRPRequestResponseProcessed: 1,
		}

		actualStats = health.Stats{}
	)

	dispatcher.On("SendEvent", mock.AnythingOfType("health.HealthFunc")).Once().
		Run(func(arguments mock.Arguments) {
			hf := arguments.Get(0).(health.HealthFunc)
			hf(actualStats)
		})

	listener.OnDeviceEvent(&device.Event{Type: device.TransactionComplete})
	assert.Equal(expectedStats, actualStats)

	dispatcher.AssertExpectations(t)
}

func testListenerOnDeviceEventPing(t *testing.T) {
	var (
		assert     = assert.New(t)
		dispatcher = new(mockDispatcher)
		listener   = &Listener{Dispatcher: dispatcher}

		expectedStats = health.Stats{
			TotalPingMessagesReceived: 1,
		}

		actualStats = health.Stats{}
	)

	dispatcher.On("SendEvent", mock.AnythingOfType("health.HealthFunc")).Once().
		Run(func(arguments mock.Arguments) {
			hf := arguments.Get(0).(health.HealthFunc)
			hf(actualStats)
		})

	listener.OnDeviceEvent(&device.Event{Type: device.Ping})
	assert.Equal(expectedStats, actualStats)

	dispatcher.AssertExpectations(t)
}

func testListenerOnDeviceEventPong(t *testing.T) {
	var (
		assert     = assert.New(t)
		dispatcher = new(mockDispatcher)
		listener   = &Listener{Dispatcher: dispatcher}

		expectedStats = health.Stats{
			TotalPongMessagesReceived: 1,
		}

		actualStats = health.Stats{}
	)

	dispatcher.On("SendEvent", mock.AnythingOfType("health.HealthFunc")).Once().
		Run(func(arguments mock.Arguments) {
			hf := arguments.Get(0).(health.HealthFunc)
			hf(actualStats)
		})

	listener.OnDeviceEvent(&device.Event{Type: device.Pong})
	assert.Equal(expectedStats, actualStats)

	dispatcher.AssertExpectations(t)
}

func TestListener(t *testing.T) {
	t.Run("OnDeviceEvent", func(t *testing.T) {
		t.Run("Connect", testListenerOnDeviceEventConnect)
		t.Run("Disconnect", testListenerOnDeviceEventDisconnect)
		t.Run("TransactionComplete", testListenerOnDeviceEventTransactionComplete)
		t.Run("Ping", testListenerOnDeviceEventPing)
		t.Run("Pong", testListenerOnDeviceEventPong)
	})
}