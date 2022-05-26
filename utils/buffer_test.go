package utils

import (
	bytes2 "bytes"
	"fmt"
	"math"
	"testing"
)

func TestReadInt(t *testing.T) {
	value := math.MaxInt32
	byts := make([]byte, 4)
	byts[0] = byte(value >> 24)
	byts[1] = byte(value >> 16)
	byts[2] = byte(value >> 8)
	byts[3] = byte(value)
	readInt, err := ReadInt(bytes2.NewBuffer(byts))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(readInt)
}

func TestReadShort(t *testing.T) {
	value := math.MaxInt16
	byts := make([]byte, 2)
	byts[0] = byte(value >> 8)
	byts[1] = byte(value)
	readShort, err := ReadShort(bytes2.NewBuffer(byts))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(readShort)
}

func TestReadMiddle(t *testing.T) {
	value := 1<<24 - 1
	byts := make([]byte, 3)
	byts[0] = byte(value >> 16)
	byts[0] = byte(value >> 8)
	byts[1] = byte(value)
	readShort, err := ReadMiddle(bytes2.NewBuffer(byts))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(readShort)
}
