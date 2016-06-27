package psmb

import (
	"strings"
	"testing"

	"github.com/marksalpeter/sugar"
)

func Test16bitConvert(t *testing.T) {
	s := sugar.New(nil)
	arr := []uint16{4396, 79, 4660, 22136} // 112C004F12345678

	s.Title("Bytes to 16-bit integer array")

	s.Assert("`byte to uint16` in big endian order", func(log sugar.Log) bool {
		desire := []uint16{4396, 79, 4660, 22136}
		bytes := RegsToBytes(arr)
		result, _ := BytesToUInt16s(bytes, 0)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to uint16` in little endian order", func(log sugar.Log) bool {
		desire := []uint16{11281, 20224, 13330, 30806}
		bytes := RegsToBytes(arr)
		result, _ := BytesToUInt16s(bytes, 1)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to int16` in big endian order", func(log sugar.Log) bool {
		desire := []int16{4396, 79, 4660, 22136}
		bytes := RegsToBytes(arr)
		result, _ := BytesToInt16s(bytes, 0)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to int16` in little endian order", func(log sugar.Log) bool {
		desire := []int16{11281, 20224, 13330, 30806}
		bytes := RegsToBytes(arr)
		result, _ := BytesToInt16s(bytes, 1)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Title("Bytes to 32-bit integer array")

	s.Assert("`byte to uint32` in (ABCD) Big Endian order", func(log sugar.Log) bool {
		desire := []uint32{288096335, 305419896}
		bytes := RegsToBytes(arr)
		result, _ := BytesToUInt32s(bytes, 0)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to uint32` in (DCBA) Little Endian order", func(log sugar.Log) bool {
		desire := []uint32{1325411345, 2018915346}
		bytes := RegsToBytes(arr)
		result, _ := BytesToUInt32s(bytes, 1)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to uint32` in (BADC) Mid-Big Endian order", func(log sugar.Log) bool {
		desire := []uint32{739331840, 873625686}
		bytes := RegsToBytes(arr)
		result, _ := BytesToUInt32s(bytes, 2)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to uint32` in (CDAB) Mid-Little Endian order", func(log sugar.Log) bool {
		desire := []uint32{5181740, 1450709556}
		bytes := RegsToBytes(arr)
		result, _ := BytesToUInt32s(bytes, 3)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to int32` in (ABCD) Big Endian  order", func(log sugar.Log) bool {
		desire := []int32{288096335, 305419896}
		bytes := RegsToBytes(arr)
		result, _ := BytesToInt32s(bytes, 0)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to int32` in (DCBA) Little Endian order", func(log sugar.Log) bool {
		desire := []int32{1325411345, 2018915346}
		bytes := RegsToBytes(arr)
		result, _ := BytesToInt32s(bytes, 1)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to int32` in (BADC) Mid-Big Endian order", func(log sugar.Log) bool {
		desire := []int32{739331840, 873625686}
		bytes := RegsToBytes(arr)
		result, _ := BytesToInt32s(bytes, 2)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to int32` in (CDAB) Mid-Little Endian order", func(log sugar.Log) bool {
		desire := []int32{5181740, 1450709556}
		bytes := RegsToBytes(arr)
		result, _ := BytesToInt32s(bytes, 3)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Title("Bytes to 32-bit float array")

	s.Assert("`byte to float32` in (ABCD) Big Endian order", func(log sugar.Log) bool {
		arr2 := []uint16{17820, 16863, 17668, 46924} // 459C41DF4504B74C
		desire := []float32{5000.234, 2123.456}
		bytes := RegsToBytes(arr2)
		result, _ := BytesToFloat32s(bytes, 0)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to float32` in (DCBA) Little Endian order", func(log sugar.Log) bool {
		desire := []float32{2150371580, 1.73782444e34}
		bytes := RegsToBytes(arr)
		result, _ := BytesToFloat32s(bytes, 1)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to float32` in (BADC) Mid-Big Endian order", func(log sugar.Log) bool {
		desire := []float32{2.06495931e-12, 1.36410875e-7}
		bytes := RegsToBytes(arr)
		result, _ := BytesToFloat32s(bytes, 2)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`byte to float32` in (CDAB) Mid-Little Endian order", func(log sugar.Log) bool {
		desire := []float32{7.261164e-39, 68189266400000}
		bytes := RegsToBytes(arr)
		result, _ := BytesToFloat32s(bytes, 3)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Title("Bytes utility test")

	s.Assert("`RegsToBytes` test", func(log sugar.Log) bool {
		desire := []byte{0x11, 0x2C, 0x00, 0x4F, 0x12, 0x34, 0x56, 0x78}
		result := RegsToBytes(arr)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToHexStr` test", func(log sugar.Log) bool {
		arr2 := []byte{0x11, 0x2C, 0x00, 0x4F, 0x12, 0x34, 0x56, 0x78}
		desire := "112C004F12345678"
		result := BytesToHexStr(arr2)
		log("desire:%s, result:%s", desire, result)
		if !strings.EqualFold(result, desire) {
			return false
		}
		return true
	})

	s.Assert("`DecStrToRegs` test", func(log sugar.Log) bool {
		input := "4396,79,4660,22136"
		result, _ := DecStrToRegs(input)
		desire := []uint16{4396, 79, 4660, 22136}
		log(result)
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

}
