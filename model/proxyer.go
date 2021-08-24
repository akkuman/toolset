package model

// DllProxyer 生成proxyer所需要的dto
type DllProxyer struct {
	// Shellcode This is a base64 encoded shellcode
	Shellcode string `json:"shellcode" binding:"base64" example:"MTIzemN4"`
	// ReGen Whether to regenerate the loader, if true, it will remove cache, this may be beneficial for bypass AV
	ReGen bool `json:"regen" example:"true"`
	// X64 Whether the shellcode is x64
	X64 bool `json:"x64" binding:"required" example:"false"`
	// DllData the data from base64 encoded dll
	DllData string `json:"dll_data" binding:"base64" example:"MTIzemN4"`
	// the dll name wich
	DllName string `json:"dll_name" binding:"endswith=.dll" example:"add.dll"`
}
