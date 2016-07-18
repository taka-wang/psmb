package psmb

import (
	"strings"
	"testing"

	"github.com/takawang/sugar"
)

func TestBinaryOps(t *testing.T) {

	s := sugar.New(t)
	arr := []uint16{4396, 79, 4660, 22136} // 112C004F12345678

	// --------------------------------------------//
	s.Title("Bytes to 16-bit integer array tests")

	s.Assert("`BytesToUInt16s` in big endian order - (1/4)", func(log sugar.Log) bool {
		desire := []uint16{4396, 79, 4660, 22136}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToUInt16s(bytes, BigEndian)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToUInt16s` in little endian order - (2/4)", func(log sugar.Log) bool {
		desire := []uint16{11281, 20224, 13330, 30806}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToUInt16s(bytes, LittleEndian)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToInt16s` in big endian order - (3/4)", func(log sugar.Log) bool {
		desire := []int16{4396, 79, 4660, 22136}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToInt16s(bytes, BigEndian)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToInt16s` in little endian order - (4/4)", func(log sugar.Log) bool {
		desire := []int16{11281, 20224, 13330, 30806}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToInt16s(bytes, LittleEndian)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	// --------------------------------------------//
	s.Title("Bytes to 32-bit integer array tests")

	s.Assert("`BytesToUInt32s` in (ABCD) Big Endian order - (1/4)", func(log sugar.Log) bool {
		desire := []uint32{288096335, 305419896}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToUInt32s(bytes, ABCD)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToUInt32s` in (DCBA) Little Endian order - (2/4)", func(log sugar.Log) bool {
		desire := []uint32{1325411345, 2018915346}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToUInt32s(bytes, DCBA)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToUInt32s` in (BADC) Mid-Big Endian order - (3/4)", func(log sugar.Log) bool {
		desire := []uint32{739331840, 873625686}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToUInt32s(bytes, BADC)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToUInt32s` in (CDAB) Mid-Little Endian order - (4/4)", func(log sugar.Log) bool {
		desire := []uint32{5181740, 1450709556}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToUInt32s(bytes, CDAB)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToInt32s` in (ABCD) Big Endian order - (1/4)", func(log sugar.Log) bool {
		desire := []int32{288096335, 305419896}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToInt32s(bytes, ABCD)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToInt32s` in (DCBA) Little Endian order - (2/4)", func(log sugar.Log) bool {
		desire := []int32{1325411345, 2018915346}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToInt32s(bytes, DCBA)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToInt32s` in (BADC) Mid-Big Endian order - (3/4)", func(log sugar.Log) bool {
		desire := []int32{739331840, 873625686}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToInt32s(bytes, BADC)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToInt32s` in (CDAB) Mid-Little Endian order - (4/4)", func(log sugar.Log) bool {
		desire := []int32{5181740, 1450709556}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToInt32s(bytes, CDAB)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	// --------------------------------------------//
	s.Title("Bytes to 32-bit float array tests")

	s.Assert("`BytesToFloat32s` in (ABCD) Big Endian order - (1/4)", func(log sugar.Log) bool {
		arr2 := []uint16{17820, 16863, 17668, 46924} // 459C41DF4504B74C
		desire := []float32{5000.234, 2123.456}
		bytes, _ := RegistersToBytes(arr2)
		result, _ := BytesToFloat32s(bytes, ABCD)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToFloat32s` in (DCBA) Little Endian order - (2/4)", func(log sugar.Log) bool {
		desire := []float32{2150371580, 1.73782444e34}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToFloat32s(bytes, DCBA)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToFloat32s` in (BADC) Mid-Big Endian order - (3/4)", func(log sugar.Log) bool {
		desire := []float32{2.06495931e-12, 1.36410875e-7}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToFloat32s(bytes, BADC)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToFloat32s` in (CDAB) Mid-Little Endian order - (4/4)", func(log sugar.Log) bool {
		desire := []float32{7.261164e-39, 68189266400000}
		bytes, _ := RegistersToBytes(arr)
		result, _ := BytesToFloat32s(bytes, CDAB)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	// --------------------------------------------//
	s.Title("Bytes/registers utility tests")

	s.Assert("`BitStringToUInt8s` test", func(log sugar.Log) bool {
		input := "1,0,1,1,0,1"
		result, _ := BitStringToUInt8s(input)
		desire := []uint8{1, 0, 1, 1, 0, 1}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BitStringToUInt8s` test - left comma", func(log sugar.Log) bool {
		input := ",1,0,1,1,0,1"
		result, err := BitStringToUInt8s(input)
		if err != nil {
			log(err)
			return true
		}
		desire := []uint8{1, 0, 1, 1, 0, 1}
		log("input length: %d; desire length: %d", len(result), len(desire))

		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BitStringToUInt8s` test - right comma", func(log sugar.Log) bool {
		input := "1,0,1,1,0,1,"
		result, _ := BitStringToUInt8s(input)
		desire := []uint8{1, 0, 1, 1, 0, 1}
		log("input length: %d; desire length: %d", len(result), len(desire))

		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BitStringToUInt8s` test - left, right comma", func(log sugar.Log) bool {
		input := ",1,0,1,1,0,1,"
		result, _ := BitStringToUInt8s(input)
		desire := []uint8{1, 0, 1, 1, 0, 1}
		log("input length: %d; desire length: %d", len(result), len(desire))

		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`RegistersToBytes` test", func(log sugar.Log) bool {
		desire := []byte{0x11, 0x2C, 0x00, 0x4F, 0x12, 0x34, 0x56, 0x78}
		result, _ := RegistersToBytes(arr)
		for idx := 0; idx < len(result); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`BytesToHexString` test", func(log sugar.Log) bool {
		arr2 := []byte{0x11, 0x2C, 0x00, 0x4F, 0x12, 0x34, 0x56, 0x78}
		desire := "112C004F12345678"
		result := BytesToHexString(arr2)
		log("desire:%s, result:%s", desire, result)
		if !strings.EqualFold(result, desire) {
			return false
		}
		return true
	})

	s.Assert("`DecimalStringToRegisters` test", func(log sugar.Log) bool {
		input := "4396,79,4660,22136"
		result, _ := DecimalStringToRegisters(input)
		desire := []uint16{4396, 79, 4660, 22136}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`DecimalStringToRegisters` test - left comma", func(log sugar.Log) bool {
		input := ",4396,79,4660,22136"
		result, _ := DecimalStringToRegisters(input)
		desire := []uint16{4396, 79, 4660, 22136}
		log("input length: %d; desire length: %d", len(result), len(desire))
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`DecimalStringToRegisters` test - right comma", func(log sugar.Log) bool {
		input := "4396,79,4660,22136,"
		result, _ := DecimalStringToRegisters(input)
		desire := []uint16{4396, 79, 4660, 22136}
		log("input length: %d; desire length: %d", len(result), len(desire))
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`DecimalStringToRegisters` test - left, right comma", func(log sugar.Log) bool {
		input := ",4396,79,4660,22136,"
		result, _ := DecimalStringToRegisters(input)
		desire := []uint16{4396, 79, 4660, 22136}
		log("input length: %d; desire length: %d", len(result), len(desire))
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`HexStringToRegisters` test", func(log sugar.Log) bool {
		input := "112C004F12345678"
		result, _ := HexStringToRegisters(input)
		desire := []uint16{4396, 79, 4660, 22136}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%d, result:%d", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`HexStringToRegisters` test - wrong length", func(log sugar.Log) bool {
		input := "112C004F1234567"
		_, err := HexStringToRegisters(input)
		if err != nil {
			log(err)
			return true
		}
		return false

		/*
			desire := []uint16{4396, 79, 4660, 22136}
			for idx := 0; idx < len(desire); idx++ {
				log("desire:%d, result:%d", desire[idx], result[idx])
				if result[idx] != desire[idx] {
					return false
				}
			}
			return true
		*/
	})

	s.Assert("`LinearScalingRegisters` test", func(log sugar.Log) bool {
		// arr := []uint16{4396, 79, 4660, 22136} // 112C004F12345678
		result, _ := LinearScalingRegisters(arr, 0, 65535, 100, 200)
		desire := []float32{106.70786602578775, 100.1205462729839, 107.11070420386054, 133.7773708705272}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})

	s.Assert("`LinearScalingRegisters` test - (0,0,0,0)", func(log sugar.Log) bool {
		// arr := []uint16{4396, 79, 4660, 22136} // 112C004F12345678
		_, err := LinearScalingRegisters(arr, 0, 0, 0, 0)
		if err != nil {
			log(err)
			return true
		}
		return false
	})

	s.Assert("`LinearScalingRegisters` test - reverse", func(log sugar.Log) bool {
		// arr := []uint16{4396, 79, 4660, 22136} // 112C004F12345678
		result, _ := LinearScalingRegisters(arr, 65535, 0, 100, 200)
		desire := []float32{106.70786602578775, 100.1205462729839, 107.11070420386054, 133.7773708705272}
		for idx := 0; idx < len(desire); idx++ {
			log("desire:%f, result:%f", desire[idx], result[idx])
			if result[idx] != desire[idx] {
				return false
			}
		}
		return true
	})
}
