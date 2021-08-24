package model

// DllProxyer 生成proxyer所需要的dto
type DllProxyer struct {
	// Shellcode This is a base64 encoded shellcode
	Shellcode string `json:"shellcode" binding:"required,base64" example:"MTIzemN4"`
	// X64 Whether the shellcode is x64
	X64 bool `json:"x64" example:"false"`
	// DllData the data from base64 encoded dll
	DllData string `json:"dll_data" binding:"required,base64" example:"MTIzemN4"`
	// the filename of origin dll
	DllName string `json:"dll_name" binding:"required,endswith=.dll" example:"add.dll"`
}
