#include "dllmain.h"
#include "windows.h"

void WINAPI workman() {
    DWORD baseAddress = (DWORD)GetModuleHandleA(NULL);
    PIMAGE_DOS_HEADER dosHeader = (PIMAGE_DOS_HEADER)baseAddress;
    PIMAGE_NT_HEADERS32 ntHeader = (PIMAGE_NT_HEADERS32)(baseAddress + dosHeader->e_lfanew);
    DWORD entryPoint = (DWORD)baseAddress + ntHeader->OptionalHeader.AddressOfEntryPoint;
    DWORD old;
    /*
    68 56341200     push 0x123456
    58              pop eax
    FFE0            jmp eax

    也可利用 https://bbs.pediy.com/thread-266711-1.htm 中的方法
    */
    BYTE shellcode[] = {0x68,0x00,0x00,0x00,0x00,0x58,0xff,0xe0};
    int size = sizeof(shellcode) / sizeof(BYTE);
    // 下面的DOWRD是32位程序的指针长度
    *(DWORD *)(shellcode+1) = (DWORD)OnProcessAttach;
    // 将上面的shellcdoe拷贝到ep
    VirtualProtect((LPVOID)entryPoint, size, PAGE_READWRITE, &old);
    for (int i = 0; i < size; i++) {
        *((PBYTE)entryPoint + i) = shellcode[i];
    }
    VirtualProtect((LPVOID)entryPoint, size, old, &old);
}

BOOL WINAPI DllMain(
    HINSTANCE _hinstDLL,  // handle to DLL module
    DWORD _fdwReason,     // reason for calling function
    LPVOID _lpReserved)   // reserved
{
    switch (_fdwReason) {
	case DLL_PROCESS_ATTACH:
		// Initialize once for each new process.
        // Return FALSE to fail DLL load.
        workman();
        break;
    case DLL_PROCESS_DETACH:
        // Perform any necessary cleanup.
        break;
    case DLL_THREAD_DETACH:
        // Do thread-specific cleanup.
        break;
    case DLL_THREAD_ATTACH:
		// Do thread-specific initialization.
        break;
    }
    return TRUE; // Successful.
}

void WINAPI func_wrapper() {};