package goclient

import (
	"testing"
	"time"

	"github.com/marksalpeter/sugar"
	zmq "github.com/taka-wang/zmq3"
)

var hostName, portNum string

// generic tcp publisher
func publisher(cmd string) {

	sender, _ := zmq.NewSocket(zmq.PUB)
	defer sender.Close()
	sender.Connect("ipc:///tmp/to.psmb")

	for {
		time.Sleep(time.Duration(1) * time.Second)
		sender.Send("tcp", zmq.SNDMORE) // frame 1
		sender.Send(cmd, 0)             // convert to string; frame 2
		// send the exit loop
		break
	}
}

// generic subscribe
func subscriber() (string, string) {
	receiver, _ := zmq.NewSocket(zmq.SUB)
	defer receiver.Close()
	receiver.Connect("ipc:///tmp/from.psmb")
	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	for {
		msg, _ := receiver.RecvMessage(0)
		// recv then exit loop
		return msg[0], msg[1]
	}
}

func TestPsmb(t *testing.T) {
	s := sugar.New(nil)
	s.Title("psmb test")
}
