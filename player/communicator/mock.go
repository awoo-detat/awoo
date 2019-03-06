package communicator

import (
	"github.com/stretchr/testify/mock"
)

func NewMock() *MockCommunicator {
	return &MockCommunicator{}
}

type MockCommunicator struct {
	mock.Mock
}

func (mc *MockCommunicator) ReadMessage() (int, []byte, error) {
	args := mc.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

func (mc *MockCommunicator) WriteMessage(messageType int, data []byte) error {
	args := mc.Called(messageType, data)
	return args.Error(0)
}

func (mc *MockCommunicator) Close() error {
	args := mc.Called()
	return args.Error(0)
}
