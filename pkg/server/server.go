package server

import (
	"errors"
	"github.com/spf13/viper"
	"k8s-deploy/pkg/router"
	"net/http"
)

func Start() error {

	addr := viper.GetString("addr")
	if addr == "" {
		return errors.New("listen addr is not config")
	}
	return http.ListenAndServe(addr, router.Init())
}
