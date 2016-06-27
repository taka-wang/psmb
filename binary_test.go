/*
func main() {
	var i int64 = 2323
	buf := Int64ToBytes(i)
	fmt.Println(buf)
	fmt.Println(BytesToInt64(buf))

	var arr []uint16 = []uint16{4396, 79, 4660, 22136}
	buf2 := Uint16ToBytes(arr)
	buf3 := Uint16To2Bytes(arr)
	fmt.Println(buf2)
	fmt.Println(buf3)
	fmt.Println(hex.EncodeToString(buf2))
	fmt.Println(hex.EncodeToString(buf3))
	fmt.Println(Float32frombytes(buf2))
}
*/

package psmb

import "testing"

func Test_Division_1(t *testing.T) {
	HexStrToRegs("AB12EF34")
	t.Log("pass")
}

/*

func Test_Division_1(t *testing.T) {
	t.Log("pass")
}

func Test_Division_2(t *testing.T) {
	t.Error("fail")
}

func Test_Division_3(t *testing.T) {
	var arr []uint16 = []uint16{4396, 79, 4660, 22136}
	buf := new(bytes.Buffer)
	err := Write(buf, MidBigEndian, arr)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	t.Log(buf.Bytes())
}

type ByteOrder interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}

// Write my writer
func Write(w io.Writer, order ByteOrder, data interface{}) error {
	return binary.Write(w, order.(binary.ByteOrder), data)
}


// BADC
// DCBA
func (midBigEndian) Uint32(b []byte) uint32 {
	//return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2]) | uint32(b[3])<<24

	return uint32(b[0])<<16 | uint32(b[1])<<24 | uint32(b[2]) | uint32(b[3])<<8
	//A             //B             //C                 //D
}

// Mid-Big Endian (BADC)

var MidBigEndian midBigEndian

type midBigEndian struct{}

// Mid-Little Endian (CDAB)

var MidLittle midLittleEndian

type midLittleEndian struct{}
*/
/*
// LittleEndian is the little-endian implementation of ByteOrder.
var LittleEndian littleEndian

type bigEndian struct{}

// BigEndian is the big-endian implementation of ByteOrder.
var BigEndian bigEndian

type littleEndian struct{}
*/

/*

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Uint16To2Bytes(data []uint16) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func Uint16ToBytes(data []uint16) []byte {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}
	return buf.Bytes()
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

*/
