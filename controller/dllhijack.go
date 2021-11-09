package controller

import (
	"toolset/plugins"

	"github.com/gin-gonic/gin"
)

// DllHijackConfig godoc
// @Summary get dll hijack option
// @Description get all valid dll hijack option
// @Tags loader
// @Accept json
// @Produce json
// @Success 200 {array} plugins.DllHijackOptionItem
// @Router /loader/dllhijack/config [post]
func (c *Controller) DllHijackConfig(ctx *gin.Context) {
	ctx.JSON(200, plugins.DllHijackConfig())
}