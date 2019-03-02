package player

// A Communicator is an interface that is satisfied by
// github.com/gorilla/websocket and can be implemented by mocked/stubbed
// structs for testing and web-independent purposes.
type Communicator interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}
