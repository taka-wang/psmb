package tcp

// command table for Downstream services

// MbCmdType defines modbus tcp command type
type MbCmdType string

// command table for modbusd
const (
	fc1  MbCmdType = "1"
	fc2  MbCmdType = "2"
	fc3  MbCmdType = "3"
	fc4  MbCmdType = "4"
	fc5  MbCmdType = "5"
	fc6  MbCmdType = "6"
	fc15 MbCmdType = "15"
	fc16 MbCmdType = "16"
	// setMbTimeout set TCP connection timeout
	setMbTimeout MbCmdType = "50"
	// getMbTimeout get TCP connection timeout
	getMbTimeout MbCmdType = "51"
)
