package crypt

import (
	"encoding/base64"
	"errors"
	"math/rand"
)

func XOR(text []byte, key []byte) []byte {
	for i := 0; i < len(text); i++ {
		text[i] = text[i] ^ key[i%len(key)]
	}
	return text
}

func Base64Encode(data []byte) []byte{
	result := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(result, data)
	return result
}

func Base64Decode(data []byte) ([]byte, error){
	result := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	_, err := base64.StdEncoding.Decode(result, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Base64URLEncode(data []byte) []byte{
	result := make([]byte, base64.RawURLEncoding.EncodedLen(len(data)))
	base64.RawURLEncoding.Encode(result, data)
	return result
}

func Base64URLDecode(data []byte) ([]byte, error){
	result := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	_, err := base64.RawURLEncoding.Decode(result, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NetbiosEncode(data []byte, key byte) []byte {
	var result []byte
	for _, value := range data {
		buf := make([]byte, 2)
		buf[0] = (value >> 4) + key
		buf[1] = value&0xf + key
		result = append(result, buf...)
	}
	return result
}

func NetbiosDecode(data []byte, key byte) []byte {
	var result []byte
	for i := 0; i < len(data); i += 2 {
		result = append(result, (data[i]-key)<<4+(data[i+1]-key)&0xf)
	}
	return result
}

func MaskEncode(data []byte) []byte {
	key := make([]byte, 4)
	rand.Read(key)
	return append(key, XOR(data, key)...)
}

func MaskDecode(data []byte, key []byte) []byte {
	return XOR(data, key)
}

func Encrypt(data []byte, EncryptType string) ([] byte, error){
	var result [] byte
	switch EncryptType {
	case "base64":
		result = Base64Encode(data)
	case "base64url":
		result = Base64URLEncode(data)
	case "mask":
		result = MaskEncode(data)
	case "netbios":
		result = NetbiosEncode(data, byte('a'))
	case "netbiosu":
		result = NetbiosEncode(data, byte('A'))
	default:
		return nil, errors.New("Wrong encryption type.")
	}
	return result, nil
}

func Decrypt(data []byte, DecryptType string) ([] byte, error){
	var result [] byte
	switch DecryptType {
	case "base64":
		result, err := Base64Decode(data)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "base64url":
		result, err := Base64URLDecode(data)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "mask":
		if len(data) <= 4 {
			return result, nil
		}
		result = MaskDecode(data[4:], data[0:4])
	case "netbios":
		result = NetbiosDecode(data, byte('a'))
	case "netbiosu":
		result = NetbiosDecode(data, byte('A'))
	default:
		return nil, errors.New("Wrong decryption type.")
	}
	return result, nil
}

func EncryptMultipleTypes(data []byte, types []string) ([]byte, error){
	var err error
	for _, value :=range types{
		data, err = Encrypt(data, value)
		if err != nil{
			return nil, err
		}
	}
	return data, nil
}

func DecryptMultipleTypes(data []byte, types []string) ([]byte, error){
	var err error
	for index, _ :=range types{
		data, err = Decrypt(data, types[len(types) - index - 1])
		if err != nil{
			return nil, err
		}
	}
	return data, nil
}
