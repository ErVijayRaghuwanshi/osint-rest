package twitter

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, svc *Service) {

    h := NewHandler(svc)

    rg.GET("/ping", h.Ping)
}
