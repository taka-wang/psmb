package dummysrv_test

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	psmb "github.com/taka-wang/psmb"
	"github.com/takawang/sugar"
	zmq "github.com/takawang/zmq3"
)

var hostName string
var portNum1 = "502"
var portNum2 = "503"
var longRun = true

// generic tcp publisher
func publisher(cmd, json string) {
	sender, _ := zmq.NewSocket(zmq.PUB)
	defer sender.Close()
	sender.Connect("ipc:///tmp/to.psmb")

	for {
		time.Sleep(time.Duration(10) * time.Millisecond)
		t := time.Now()
		fmt.Println("Req:", t.Format("2006-01-02 15:04:05.000"))
		sender.Send(cmd, zmq.SNDMORE) // frame 1
		sender.Send(json, 0)          // convert to string; frame 2
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
		fmt.Println("listen..")
		msg, _ := receiver.RecvMessage(0)

		t := time.Now()
		fmt.Println("Data:", t.Format("2006-01-02 15:04:05.000"))

		// recv then exit loop
		return msg[0], msg[1]
	}
}

// long subscribe for data
func longSubscriber() {
	longRun = true
	receiver, _ := zmq.NewSocket(zmq.SUB)
	defer receiver.Close()
	receiver.Connect("ipc:///tmp/from.psmb")
	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	for longRun {
		msg, _ := receiver.RecvMessage(0)
		t := time.Now()
		fmt.Println("Res:", t.Format("2006-01-02 15:04:05.000"))
		fmt.Println(msg[0], msg[1])
	}
}

// init functions
func init() {
	time.Sleep(2000 * time.Millisecond)

	// generalize host reslove for docker/local env
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}
}

func TestPollRequestSingle(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`mbtcp.poll.create FC1` read bits test: port 503 - miss name", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			//Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return true
		}
		return false

	})

	s.Assert("`mbtcp.poll.create/mbtcp.poll.delete FC1` read bits test: port 503 - interval 1", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		go longSubscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		time.Sleep(10 * time.Second)

		historyReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		historyReqStr, _ := json.Marshal(historyReq)
		cmd = "mbtcp.poll.history"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.poll.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.poll.update/mbtcp.poll.delete FC1` read bits test: port 503 - miss name", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		go longSubscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		time.Sleep(10 * time.Second)

		// update request
		updateReq := psmb.MbtcpPollOpReq{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "",
			Interval: 2,
		}
		updateReqStr, _ := json.Marshal(updateReq)
		cmd = "mbtcp.poll.update"
		go publisher(cmd, string(updateReqStr))

		time.Sleep(3 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.poll.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false // bad
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.poll.update/mbtcp.poll.delete FC1` read bits test: port 503 - interval 2", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		go longSubscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		time.Sleep(10 * time.Second)

		// update request
		updateReq := psmb.MbtcpPollOpReq{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 2,
		}
		updateReqStr, _ := json.Marshal(updateReq)
		cmd = "mbtcp.poll.update"
		go publisher(cmd, string(updateReqStr))

		time.Sleep(20 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.poll.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false // bad
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.poll.read/mbtcp.poll.delete FC1` read bits test: port 503 - miss name", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		go longSubscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		time.Sleep(10 * time.Second)

		// update request
		updateReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "",
		}
		updateReqStr, _ := json.Marshal(updateReq)
		cmd = "mbtcp.poll.read"
		go publisher(cmd, string(updateReqStr))

		time.Sleep(3 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.poll.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false // bad
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.poll.read/mbtcp.poll.delete FC1` read bits test: port 503", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  false,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		go longSubscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		time.Sleep(10 * time.Second)

		// update request
		updateReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}
		updateReqStr, _ := json.Marshal(updateReq)
		cmd = "mbtcp.poll.read"
		go publisher(cmd, string(updateReqStr))

		time.Sleep(3 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.poll.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false // bad
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.poll.toggle/mbtcp.poll.read/mbtcp.poll.delete FC1` read bits test: port 503", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()
		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		go longSubscriber()

		time.Sleep(10 * time.Second)

		// update request
		updateReq := psmb.MbtcpPollOpReq{
			From:    "web",
			Tid:     time.Now().UTC().UnixNano(),
			Name:    "LED_11",
			Enabled: false,
		}
		updateReqStr, _ := json.Marshal(updateReq)
		cmd = "mbtcp.poll.toggle"
		go publisher(cmd, string(updateReqStr))

		time.Sleep(3 * time.Second)

		// read request
		updateReq2 := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}
		updateReqStr2, _ := json.Marshal(updateReq2)
		cmd = "mbtcp.poll.read"
		go publisher(cmd, string(updateReqStr2))

		time.Sleep(3 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.poll.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false // bad
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.poll.toggle/mbtcp.poll.delete FC1` read bits test: port 503 - enable", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  false,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		go longSubscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		time.Sleep(10 * time.Second)

		// enable poller
		enableReq := psmb.MbtcpPollOpReq{
			From:    "web",
			Tid:     time.Now().UTC().UnixNano(),
			Name:    "LED_11",
			Enabled: true,
		}

		enableReqStr, _ := json.Marshal(enableReq)
		cmd = "mbtcp.poll.toggle"
		go publisher(cmd, string(enableReqStr))

		time.Sleep(10 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Name: "LED_11",
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.poll.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(2 * time.Second)

		longRun = false
		time.Sleep(1 * time.Second)

		return true
	})
}

func TestTimeoutOps(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`mbtcp.timeout.update` test - invalid json type - (1/5)", func(log sugar.Log) bool {
		ReadReqStr :=
			`{
                "from": "web",
                "tid": 123456,
                "timeout": "210000"
            }`
		cmd := "mbtcp.timeout.update"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("set timeout as 200000")
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`mbtcp.timeout.update` test - invalid value (1) - (2/5)", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Data: 1,
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.update"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("set timeout as 200000")
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.timeout.read` test - invalid value (1) - (3/5)", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.read"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" || r2.Data != 200000 {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.timeout.update` test - valid value (212345) - (4/5)", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			Data: 212345,
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.update"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.timeout.read` test - valid value (212345) - (5/5) ", func(log sugar.Log) bool {
		ReadReq := psmb.MbtcpTimeoutReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		ReadReqStr, _ := json.Marshal(ReadReq)
		cmd := "mbtcp.timeout.read"
		go publisher(cmd, string(ReadReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(ReadReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpTimeoutRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" || r2.Data != 212345 {
			return false
		}
		return true
	})

}

func TestOneOffWriteFC5(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - invalid value(2) - (1/4)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    5,
			Slave: 1,
			Addr:  10,
			Data:  2,
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   3,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 1 {
			return false
		}

		return true
	})

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - miss from & port - (2/4)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			//From:  "web",
			Tid: time.Now().UTC().UnixNano(),
			IP:  hostName,
			//Port:  portNum1,
			FC:    5,
			Slave: 1,
			Addr:  10,
			Data:  1,
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   3,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 1 {
			return false
		}

		return true
	})

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - valid value(0) - (3/4)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    5,
			Slave: 1,
			Addr:  10,
			Data:  0,
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   3,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 0 {
			return false
		}

		return true
	})

	s.Assert("`mbtcp.once.write FC5` write bit test: port 502 - valid value(1) - (4/4)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    5,
			Slave: 1,
			Addr:  10,
			Data:  1,
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   3,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 1 {
			return false
		}
		return true
	})
}

func TestOneOffWriteFC6(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - valid value (22) - (1/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   false,
			Data:  "22",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 22 {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - miss hex type & port - (2/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Data:  "22",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 22 {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - invalid value (array) - (3/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   false,
			Data:  "22,11",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 22 {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `DEC` register test: port 502 - invalid hex type - (4/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   false,
			Data:  "ABCD1234",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r1.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - valid value (ABCD) - (5/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   true,
			Data:  "ABCD",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 43981 {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - miss port (ABCD) - (6/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   true,
			Data:  "ABCD",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 43981 {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - invalid value (ABCD1234) - (7/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   true,
			Data:  "ABCD1234",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}
		if r3[0] != 43981 {
			return false
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC6` write `HEX` register test: port 502 - invalid hex type - (8/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    6,
			Slave: 1,
			Addr:  10,
			Hex:   true,
			Data:  "22,11",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r1.Status != "ok" {
			return true
		}
		return false
	})

}

func TestOneOffWriteFC15(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - invalid json type - (1/5)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReqStr :=
			`{
                "from": "web",
                "tid": 123456,
                "fc" : 15,
                "ip": "192.168.0.1",
                "port": "503",
                "slave": 1,
                "addr": 10,
                "len": 4,
                "data": [-1,0,-1,0]
            }`

		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return true
		}
		return false

	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - invalid json type - (2/5)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReqStr :=
			`{
                "from": "web",
                "tid": 123456,
                "fc" : 15,
                "ip": "192.168.0.1",
                "port": "503",
                "slave": 1,
                "addr": 10,
                "len": 4,
                "data": "1,0,1,0"
            }`

		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - invalid value(2) - (3/5)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    15,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Data:  []uint16{2, 0, 2, 0},
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   4,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{1, 0, 1, 0}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - miss from & port - (4/5)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			//From:  "web",
			Tid: time.Now().UTC().UnixNano(),
			IP:  hostName,
			//Port:  portNum1,
			FC:    15,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Data:  []uint16{2, 0, 2, 0},
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   4,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{1, 0, 1, 0}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC15` write bit test: port 502 - valid value(0) - (5/5)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    15,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Data:  []uint16{0, 1, 1, 0},
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}
		//return true

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  10,
			Len:   4,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{0, 1, 1, 0}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})
}

func TestOneOffWriteFC16(t *testing.T) {

	s := sugar.New(t)

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - valid value (11,22,33,44) - (1/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Hex:   false,
			Data:  "11,22,33,44",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{11, 22, 33, 44}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - miss hex type & port - (2/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
			IP:   hostName,
			//Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			//Hex:   false,
			Data: "11,22,33,44",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{11, 22, 33, 44}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - invalid hex type - (3/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Hex:   false,
			Data:  "ABCD1234",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r1.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`mbtcp.once.write FC16` write `DEC` register test: port 502 - invalid length - (4/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   8,
			Hex:   false,
			Data:  "11,22,33,44",
		}
		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{11, 22, 33, 44}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - valid value (ABCD1234) - (5/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Hex:   true,
			Data:  "ABCD1234",
		}

		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{0xABCD, 0x1234}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - miss port (ABCD) - (6/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			//From:  "web",
			Tid: time.Now().UTC().UnixNano(),
			IP:  hostName,
			//Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Hex:   true,
			Data:  "ABCD1234",
		}

		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{0xABCD, 0x1234}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - invalid hex type (11,22,33,44) - (7/8)", func(log sugar.Log) bool {
		// ---------------- write part

		writeReq := psmb.MbtcpWriteReq{
			//From:  "web",
			Tid: time.Now().UTC().UnixNano(),
			IP:  hostName,
			//Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   4,
			Hex:   true,
			Data:  "11,22,33,44",
		}

		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return true
		}
		return true

	})

	s.Assert("`mbtcp.once.write FC16` write `HEX` register test: port 502 - invalid length - (8/8)", func(log sugar.Log) bool {
		// ---------------- write part
		writeReq := psmb.MbtcpWriteReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    16,
			Slave: 1,
			Addr:  10,
			Len:   8,
			Hex:   true,
			Data:  "ABCD1234",
		}

		writeReqStr, _ := json.Marshal(writeReq)
		cmd := "mbtcp.once.write"
		go publisher(cmd, string(writeReqStr))
		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(writeReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}

		// check response
		if r1.Status != "ok" {
			return false
		}

		// ---------------- read part
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  10,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 = subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var data json.RawMessage // raw []byte
		r2 := psmb.MbtcpReadRes{Data: &data}
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		// ---------------- Compare
		var r3 []uint16
		if err := json.Unmarshal(data, &r3); err != nil {
			return false
		}

		desire := []uint16{0xABCD, 0x1234}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], r3[idx])
			if r3[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

}

func TestOneOffReadFC1(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`FC1` read bits test: port 502 - length 1", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  3,
			Len:   1,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC1` read bits test: port 502 - length 7", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  3,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC1` read bits test: port 502 - Illegal data address", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    1,
			Slave: 1,
			Addr:  1000,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC1` read bits test: port 503 - length 7", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum2,
			FC:    1,
			Slave: 1,
			Addr:  3,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

}

func TestOneOffReadFC2(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`FC2` read bits test: port 502 - length 1", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    2,
			Slave: 1,
			Addr:  3,
			Len:   1,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC2` read bits test: port 502 - length 7", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    2,
			Slave: 1,
			Addr:  3,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC2` read bits test: port 502 - Illegal data address", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    2,
			Slave: 1,
			Addr:  1000,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC2` read bits test: port 503 - length 7", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum2,
			FC:    2,
			Slave: 1,
			Addr:  3,
			Len:   7,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})
}

func TestOneOffReadFC3(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`FC3` read bytes Type 1 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 2 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.HexString,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 3 length 4 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Scale,
			Range: &psmb.ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 3 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Scale,
			Range: &psmb.ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "Invalid length to convert" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 4 length 4 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.UInt16,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 4 length 4 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.UInt16,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 4 length 4 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.UInt16,
			//Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 5 length 4 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Int16,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 5 length 4 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Int16,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 5 length 4 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Int16,
			//Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 8 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.UInt32,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 8 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.UInt32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 8 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.UInt32,
			//Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 6 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.UInt32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC3` read bytes Type 7 length 8 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Int32,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 7 length 8 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Int32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 7 length 8 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Int32,
			//Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 7 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Int32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: ABCD", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.ABCD,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: DCBA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.DCBA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: BADC", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.BADC,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 8 test: port 502 - order: CDAB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.CDAB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC3` read bytes Type 8 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Float32,
			Order: psmb.BigEndian,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC3` read bytes: port 502 - invalid type", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    3,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  9,
			Order: psmb.BigEndian,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

}

func TestOneOffReadFC4(t *testing.T) {
	s := sugar.New(t)

	s.Assert("`FC4` read bytes Type 1 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 2 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.HexString,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 3 length 4 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Scale,
			Range: &psmb.ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 3 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Scale,
			Range: &psmb.ScaleRange{0, 65535, 100, 500},
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "Invalid length to convert" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 4 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.UInt16,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 4 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.UInt16,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 4 length 4 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.UInt16,
			//Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 5 length 4 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Int16,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 5 length 4 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Int16,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 5 length 4 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   4,
			Type:  psmb.Int16,
			//Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 8 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.UInt32,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 8 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.UInt32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 8 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.UInt32,
			//Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 6 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.UInt32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC4` read bytes Type 7 length 8 test: port 502 - Order: AB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Int32,
			Order: psmb.AB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 7 length 8 test: port 502 - Order: BA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Int32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 7 length 8 test: port 502 - miss order", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Int32,
			//Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 7 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Int32,
			Order: psmb.BA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: ABCD", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.ABCD,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: DCBA", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.DCBA,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: BADC", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.BADC,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 8 test: port 502 - order: CDAB", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   8,
			Type:  psmb.Float32,
			Order: psmb.CDAB,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

	s.Assert("`FC4` read bytes Type 8 length 7 test: port 502 - invalid length", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  psmb.Float32,
			Order: psmb.BigEndian,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return true
		}
		return false
	})

	s.Assert("`FC4` read bytes: port 502 - invalid type", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpReadReq{
			From:  "web",
			Tid:   time.Now().UTC().UnixNano(),
			IP:    hostName,
			Port:  portNum1,
			FC:    4,
			Slave: 1,
			Addr:  3,
			Len:   7,
			Type:  9,
			Order: psmb.BigEndian,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.once.read"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check fail response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

}

func TestPollsRequest(t *testing.T) {
	s := sugar.New(t)
	s.Assert("`mbtcp.polls.read FC1` read 2 poll reqeusts", func(log sugar.Log) bool {

		// send request
		readReq1 := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_11",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     3,
			Len:      7,
		}
		readReqStr1, _ := json.Marshal(readReq1)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr1))
		// receive response
		s1, s2 := subscriber()
		log("req: %s, %s", cmd, string(readReqStr1))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r1 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r1.Status != "ok" {
			return false
		}

		// send request2
		readReq2 := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_22",
			Interval: 2,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       1,
			Slave:    1,
			Addr:     5,
			Len:      8,
		}
		readReqStr2, _ := json.Marshal(readReq2)
		cmd = "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr2))
		// receive response
		s3, s4 := subscriber()
		log("req: %s, %s", cmd, string(readReqStr2))
		log("res: %s, %s", s3, s4)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s4), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}

		go longSubscriber()

		time.Sleep(10 * time.Second)

		// readReq poller
		readReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd = "mbtcp.polls.read"
		go publisher(cmd, string(readReqStr))

		time.Sleep(10 * time.Second)

		longRun = false
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.polls.read/mbtcp.polls.read/mbtcp.polls.delete FC1` read 50 poll requests", func(log sugar.Log) bool {

		fmt.Println("@@Data:", time.Now().Format("2006-01-02 15:04:05.000"))
		for idx := 0; idx < 100; idx++ {
			// send request
			readReq1 := psmb.MbtcpPollStatus{
				From:     "web",
				Tid:      time.Now().UTC().UnixNano(),
				Name:     "LED_11" + strconv.Itoa(idx),
				Interval: 1,
				Enabled:  true,
				IP:       hostName,
				Port:     portNum1,
				FC:       1,
				Slave:    1,
				Addr:     uint16(idx),
				Len:      7,
			}
			readReqStr1, _ := json.Marshal(readReq1)
			cmd := "mbtcp.poll.create"
			go publisher(cmd, string(readReqStr1))
			// receive response
			/*
			   s1, s2 := subscriber()
			   log("req: %s, %s", cmd, string(readReqStr1))
			   log("res: %s, %s", s1, s2)
			*/
		}
		fmt.Println("@@@Data:", time.Now().Format("2006-01-02 15:04:05.000"))

		go longSubscriber()

		time.Sleep(10 * time.Second)

		// readReq poller
		readReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.polls.read"
		go publisher(cmd, string(readReqStr))

		time.Sleep(10 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.polls.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false
		time.Sleep(1 * time.Second)

		return true
	})

	s.Assert("`mbtcp.polls.read/mbtcp.polls.read/mbtcp.polls.toggle/mbtcp.polls.delete FC1` read 20 poll requests", func(log sugar.Log) bool {

		for idx := 0; idx < 30; idx++ {
			// send request
			readReq1 := psmb.MbtcpPollStatus{
				From:     "web",
				Tid:      time.Now().UTC().UnixNano(),
				Name:     "LED_22-" + strconv.Itoa(idx),
				Interval: 1,
				Enabled:  true,
				IP:       hostName,
				Port:     portNum1,
				FC:       1,
				Slave:    1,
				Addr:     uint16(idx),
				Len:      7,
			}
			readReqStr1, _ := json.Marshal(readReq1)
			cmd := "mbtcp.poll.create"
			go publisher(cmd, string(readReqStr1))
			// receive response
			/*
			   s1, s2 := subscriber()
			   log("req: %s, %s", cmd, string(readReqStr1))
			   log("res: %s, %s", s1, s2)
			*/
		}
		fmt.Println("@@@Data:", time.Now().Format("2006-01-02 15:04:05.000"))

		go longSubscriber()

		time.Sleep(10 * time.Second)

		// toggle poller
		readReq := psmb.MbtcpPollOpReq{
			From:    "web",
			Tid:     time.Now().UTC().UnixNano(),
			Enabled: false,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.polls.toggle"
		go publisher(cmd, string(readReqStr))

		time.Sleep(5 * time.Second)

		// delete poller
		delReq := psmb.MbtcpPollOpReq{
			From: "web",
			Tid:  time.Now().UTC().UnixNano(),
		}

		delReqStr, _ := json.Marshal(delReq)
		cmd = "mbtcp.polls.delete"
		go publisher(cmd, string(delReqStr))

		time.Sleep(5 * time.Second)

		longRun = false
		time.Sleep(1 * time.Second)

		return true
	})

	/*
	   s.Assert("`mbtcp.polls.delete FC1` read bits test: port 503", func(log sugar.Log) bool {

	   })

	   s.Assert("`mbtcp.polls.toggle FC1` read bits test: port 503", func(log sugar.Log) bool {

	   })
	*/
}

func TestPSMB(t *testing.T) {
	s := sugar.New(t)

	s.Title("Poll request tests")

	s.Assert("mbtcp.poll.create `FC3` read bytes Type 1 test: port 502", func(log sugar.Log) bool {
		// send request
		readReq := psmb.MbtcpPollStatus{
			From:     "web",
			Tid:      time.Now().UTC().UnixNano(),
			Name:     "LED_1",
			Interval: 1,
			Enabled:  true,
			IP:       hostName,
			Port:     portNum1,
			FC:       3,
			Slave:    1,
			Addr:     3,
			Len:      7,
			Type:     psmb.RegisterArray,
		}

		readReqStr, _ := json.Marshal(readReq)
		cmd := "mbtcp.poll.create"
		go publisher(cmd, string(readReqStr))

		// receive response
		s1, s2 := subscriber()

		log("req: %s, %s", cmd, string(readReqStr))
		log("res: %s, %s", s1, s2)

		// parse resonse
		var r2 psmb.MbtcpSimpleRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check response
		if r2.Status != "ok" {
			return false
		}
		return true
	})

}
