package router

import (
	"github.com/gin-gonic/gin"
	"k8s-deploy/pkg/api"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/deployment", api.CreateDeployment)
		apiv1.DELETE("/deployment", api.DeleteDeployment)
		apiv1.PATCH("/deployment", api.PatchDeployment)
		apiv1.GET("/deployment/:namespace/:appname", api.GetDeployment)

		apiv1.POST("/callback", api.CallBack)
	}

	return r
}
