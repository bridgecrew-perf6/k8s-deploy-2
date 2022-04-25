package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"k8s-deploy/pkg/deployment"
	"k8s-deploy/pkg/result"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	log "github.com/sirupsen/logrus"
)

func CreateDeployment(c *gin.Context) {
	req := deployment.InitRequest()
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Create request param error:%s", err.Error())
		c.JSON(http.StatusOK, result.Err.WithMsg(err.Error()))
		return
	}
	log.Infof("请求创建deployment: %s/%s", req.Namespace, req.AppName)

	if req.Image == ""{
		log.Error("Request param image is empty")
		c.JSON(http.StatusOK, result.ErrImageParam)
		return
	}
	if req.AppName == ""{
		log.Error("Request param appname is empty")
		c.JSON(http.StatusOK, result.ErrAppNameParam)
		return
	}
	if req.Namespace == "" {
		log.Error("Request param namespace is empty")
		c.JSON(http.StatusOK, result.ErrNamespaceParam)
		return
	}

	d, err:= deployment.New(req)
	if err != nil {
		log.Errorf("New deployment error: %s", err.Error())
		c.JSON(http.StatusOK, result.ErrDeploymentCreate)
		return
	}

	createTime := time.Now()

	err = d.Create()
	if err != nil {
		log.Errorf("Create deployment error: %s", err.Error())
		c.JSON(http.StatusOK, result.Err.WithMsg(err.Error()))
		return
	}

	go d.WatchDeploy(createTime)

	c.JSON(http.StatusOK, result.OK)
}

func DeleteDeployment(c *gin.Context) {
	req := deployment.InitRequest()
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Delete request param error:%s", err.Error())
		c.JSON(http.StatusOK, result.Err.WithMsg(err.Error()))
		return
	}
	log.Infof("请求删除eployment: %s/%s", req.Namespace, req.AppName)

	if req.AppName == ""{
		log.Error("Request param appname is empty")
		c.JSON(http.StatusOK, result.ErrAppNameParam)
		return
	}
	if req.Namespace == "" {
		log.Error("Request param namespace is empty")
		c.JSON(http.StatusOK, result.ErrNamespaceParam)
		return
	}

	d, err:= deployment.New(req)
	if err != nil {
		log.Errorf("New deployment error: %s", err.Error())
		c.JSON(http.StatusOK, result.ErrDeploymentCreate)
		return
	}

	err = d.Delete()
	if err != nil {
		log.Errorf("Delete deployment error: %s", err.Error())
		c.JSON(http.StatusOK, result.Err.WithMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.OK)
}

func PatchDeployment(c *gin.Context) {
	//req := deployment.InitRequest()
	req := deployment.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Update request param error:%s", err.Error())
		c.JSON(http.StatusOK, result.Err.WithMsg(err.Error()))
		return
	}
	log.Infof("请求更新deployment: %s/%s", req.Namespace, req.AppName)

	if req.AppName == ""{
		c.JSON(http.StatusOK, result.ErrAppNameParam)
		return
	}
	if req.Namespace == "" {
		c.JSON(http.StatusOK, result.ErrNamespaceParam)
		return
	}

	d, err:= deployment.New(req)
	if err != nil {
		log.Errorf("New deployment error: %s", err.Error())
		c.JSON(http.StatusOK, result.ErrDeploymentCreate)
		return
	}

	patchTime := time.Now()

	err = d.Patch()
	if err != nil {
		log.Errorf("Patch deployment error: %s", err.Error())
		c.JSON(http.StatusOK, result.Err.WithMsg(err.Error()))
		return
	}

	go d.WatchDeploy(patchTime)

	c.JSON(http.StatusOK, result.OK)
}

func GetDeployment(c *gin.Context)  {
	namespace := c.Param("namespace")
	appname := c.Param("appname")
	if namespace == "" || appname == "" {
		log.Errorf("Get deployment error: %s", errors.New("namespace or appname is empty"))
		c.JSON(http.StatusOK, result.ErrParam)
	}

	req := deployment.Request{
		AppName: appname,
		Namespace: namespace,
	}

	d, err:= deployment.New(req)
	if err != nil {
		log.Errorf("New deployment error: %s", err.Error())
		c.JSON(http.StatusOK, result.ErrDeploymentCreate)
		return
	}

	deploy, err := d.Get()
	if k8serrors.IsNotFound(err) {
		c.JSON(http.StatusOK, result.ErrDeploymentNotFound)
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, result.Err.WithMsg(err.Error()))
		return
	}
	log.Infof("Get Deployment %s/%s", deploy.GetNamespace(), deploy.GetName())
	c.JSON(http.StatusOK, result.OK)
}

func CallBack(c *gin.Context)  {
	type request struct {
		AppName string `json:"appName"`
		Namespace string `json:"namespace"`
		Finished bool `json:"finished"`
	}
	var r request

	err := c.BindJSON(&r)
	if err != nil {
		fmt.Println(err.Error())
	}
	req, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(req))
}
