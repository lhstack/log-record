package utils

import (
	"bytes"
	"errors"
)

//ReadInt 读取buffer中的4个字节，用大端解码
func ReadInt(buffer *bytes.Buffer) (int32, error) {
	if buffer.Len() >= 4 {
		byteArray := make([]byte, 4)
		_, err := buffer.Read(byteArray)
		if err != nil {
			return 0, err
		}
		return int32(byteArray[0])<<24 | int32(byteArray[1])<<16 | int32(byteArray[2])<<8 | int32(byteArray[3]), nil
	}
	return 0, errors.New("buffer length  lt 4,cannot convert int32")
}

//ReadMiddle 读取buffer中的3个字节，用大端解码
func ReadMiddle(buffer *bytes.Buffer) (int32, error) {
	if buffer.Len() >= 3 {
		byteArray := make([]byte, 3)
		_, err := buffer.Read(byteArray)
		if err != nil {
			return 0, err
		}
		return int32(byteArray[0])<<16 | int32(byteArray[0])<<8 | int32(byteArray[1]), nil
	}
	return 0, errors.New("buffer length  lt 3,cannot convert int32")
}

//ReadShort 读取buffer中的2个字节，用大端解码
func ReadShort(buffer *bytes.Buffer) (int16, error) {
	if buffer.Len() >= 2 {
		byteArray := make([]byte, 2)
		_, err := buffer.Read(byteArray)
		if err != nil {
			return 0, err
		}
		return int16(byteArray[0])<<8 | int16(byteArray[1]), nil
	}
	return 0, errors.New("buffer length  lt 2,cannot convert int16")
}
