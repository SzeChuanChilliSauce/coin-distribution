package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TODO:外部出发调用
func (cds *CoinDisServer) DistributeCoins(ctx *gin.Context) {
	resp := struct {
		Status int `json:"status"`
	}{
		Status: 200,
	}

	ctx.JSON(http.StatusOK, resp)
}
