//go:build darwin
package sysinfo

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

func GetOSVersion() (string, error) {
	cmd := exec.Command("sw_vers", "-productVersion")
	out, _ := cmd.CombinedOutput()
	return string(out[:]), nil
}

func IsHighPriv() bool {
	fd, err := os.Open("/root")
	defer fd.Close()
	if err != nil {
		return false
	}
	return false
}

func IsOSX64() (bool, error) {
	cmd := exec.Command("sysctl", "hw.cpu64bit_capable")
	out, _ := cmd.CombinedOutput()
	out = bytes.ReplaceAll(out, []byte("hw.cpu64bit_capable: "), []byte(""))
	if string(out) == "1" {
		return true, nil
	}
	return false, nil
}

func IsProcessX64() (bool, error) {
	if runtime.GOARCH == "amd64" {
		return false, nil
	}
	return true, nil
}

func GetCodePageANSI() ([]byte, error) {
	//hardcode for test
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 936)
	return b, nil
}

func GetCodePageOEM() ([]byte, error) {
	//hardcode for test
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 936)
	return b, nil
}

func GetUsername() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", nil
	}
	usr := user.Username
	return usr, nil
}

func IsPidX64(pid uint32) (bool, error) {
	/*is64 := false

	hProcess, err := windows.OpenProcess(uint32(0x1000), false, pid)
	if err != nil {
		return IsProcessX64()
	}

	_ = windows.IsWow64Process(hProcess, &is64)*/

	return true, nil
}

