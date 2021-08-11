package controller

import (
	"encoding/base64"
	"mime"
	"toolset/httputil"
	"toolset/model"
	"toolset/plugins"

	"github.com/gin-gonic/gin"
)

// ShellcodeRunner godoc
// @Summary generate a shellcode runner
// @Description Generate Runner according to shellcode provided by the user
// @Tags loader
// @Accept json
// @Produce octet-stream
// @Produce json
// @Param shellcode body model.ShellcodeRunner true "shellcode" Format(base64)
// @Success 200 {file} binary "shellcode runner"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /loader/shellcode-runner [post]
func (c *Controller) ShellcodeRunner(ctx *gin.Context) {
	var shellcodeRunner model.ShellcodeRunner

	if err := ctx.ShouldBindJSON(&shellcodeRunner); err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	shellcode, err := base64.StdEncoding.DecodeString(shellcodeRunner.Shellcode)
	if err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	tool := plugins.NewShellcodeLoader(shellcode, shellcodeRunner.ReGen, shellcodeRunner.X64)
	data, err := tool.Run()
	if err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	ctx.Data(200, mime.TypeByExtension(".zip"), data)
}
