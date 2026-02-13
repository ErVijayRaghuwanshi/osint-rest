package jaco

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RegisterRoutes(rg *gin.RouterGroup, svc *Service, log zerolog.Logger) {

	h := NewHandler(svc, log)

	rg.GET("/ping", h.Ping)
	rg.GET("/userinfo", h.GetUserInfo)
	rg.GET("/reel", h.GetUserReel)
	rg.GET("/comments", h.GetReelComments)
}
