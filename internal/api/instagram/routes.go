package instagram

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RegisterRoutes(rg *gin.RouterGroup, svc *Service, log zerolog.Logger) {

	h := NewHandler(svc, log)

	rg.GET("/ping", h.Ping)
	rg.GET("/userinfo", h.GetUserInfo)
	rg.GET("/timeline", h.GetTimeline)
	rg.GET("/hashtags", h.SearchHashtags)
}
