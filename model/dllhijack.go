package model

type DllHijack struct {
	// Shellcode This is a base64 encoded shellcode
	Shellcode string `json:"shellcode" binding:"required,base64" example:"MTIzemN4"`
	// Type This is dll hijack type
	Type int `json:"type" binding:"required" example:"1"`
}