package Util

import (
	"bufio"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"hack8-note_rce/mode"
	"os"
)

func Jsonp(data interface{}) mode.Teample {
	teample := mode.Teample{Data: data}

	return teample
}

func WriteFile(filename, content string) error {
	// 创建一个文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入字符串
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// 检测页面编码
func DeterminePageEncoding(r *bufio.Reader) encoding.Encoding {
	//使用peek读取十分关键，只是偷看一下，不会移动读取位置，否则其他地方就没法读取了
	bytes, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

func IsGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return true
			}
		}
	}
	return false
}
