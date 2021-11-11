#include <windows.h>

extern void OnProcessAttach();
void WINAPI func_wrapper();

BOOL WINAPI DllMain(
    HINSTANCE _hinstDLL,  // handle to DLL module
    DWORD _fdwReason,     // reason for calling function
    LPVOID _lpReserved    // reserved
);
