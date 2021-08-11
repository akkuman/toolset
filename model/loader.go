package model

// ShellcodeRunner 生成runner所需要的dto
type ShellcodeRunner struct {
	// Shellcode This is a base64 encoded shellcode
	Shellcode string `json:"shellcode" binding:"base64" example:"MTIzemN4"`
	// ReGen Whether to regenerate the loader, if true, it will remove cache, this may be beneficial for bypass AV
	ReGen bool `json:"regen" example:"true"`
	// X64 Whether the shellcode is x64
	X64 bool `json:"x64" example:"false"`
}
