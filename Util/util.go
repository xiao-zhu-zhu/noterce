package Util

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func ParseAnArg(buf *bytes.Buffer) ([]byte, error) {
	argLenBytes := make([]byte, 4)
	_, err := buf.Read(argLenBytes)
	if err != nil {
		return nil, err
	}
	argLen := binary.BigEndian.Uint32(argLenBytes)
	if argLen != 0 {
		arg := make([]byte, argLen)
		_, err = buf.Read(arg)
		if err != nil {
			return nil, err
		}
		args := strings.Split(strings.TrimRight(string(arg), "\x00"), "\x00")
		return []byte(args[0]), nil
	} else {
		return nil, err
	}

}

func ConvertChinese(byte []byte) []byte {
	result, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
	return result
}

func RandString(arr []string) string {
	n := len(arr)
	if n == 0 {
		return ""
	}
	b := make([]byte, 1)
	rand.Read(b)
	i := int(b[0]) % n
	return arr[i]
}

func DebugError() {

}
