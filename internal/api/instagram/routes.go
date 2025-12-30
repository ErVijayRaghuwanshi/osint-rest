package instagram

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, svc *Service) {

	h := NewHandler(svc)

	rg.GET("/ping", h.Ping)
	rg.GET("/userinfo", h.GetUserInfo)
	rg.GET("/timeline", h.GetTimeline)
	rg.GET("/hashtags", h.SearchHashtags)
}
