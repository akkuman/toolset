package model

// ShellcodeRunner 生成runner所需要的dto
type ShellcodeRunner struct {
	Shellcode string `json:"shellcode" binding:"base64" example:"MTIzemN4"`
	ReGen     bool   `json:"regen" example:"true"`
	X64       bool   `json:"x64" example:"false"`
}
