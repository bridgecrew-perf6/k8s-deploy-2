module k8s-deploy

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	k8s.io/kubectl v0.23.0
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2

replace k8s.io/api => k8s.io/api v0.19.2

replace k8s.io/kubectl => k8s.io/kubectl v0.19.2

replace k8s.io/apimachinery => k8s.io/apimachinery v0.19.2
