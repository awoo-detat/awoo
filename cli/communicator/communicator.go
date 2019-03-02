package communicator

import (
	"log"
	"os"
)

type Communicator struct {
	Logger *log.Logger
}

func New(fileName string) *Communicator {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	l := log.New(file, "", log.Ldate|log.Ltime)
	//l := log.New(file, "", log.Ldate|log.Ltime)
	return &Communicator{
		Logger: l,
	}
}

// ReadMessage acts as the function that would receive data from the front end.
func (c *Communicator) ReadMessage() (messageType int, p []byte, err error) {
	return 0, []byte{}, nil
}

// WriteMessage acts as the function that would send data to the front end.
func (c *Communicator) WriteMessage(messageType int, data []byte) error {
	c.Logger.Printf("%s\n", data)
	return nil
}

// Close is a no op.
func (c *Communicator) Close() error {
	return nil
}
