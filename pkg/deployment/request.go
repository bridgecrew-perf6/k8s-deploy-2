package deployment

import "github.com/spf13/viper"

type Request struct {
	AppName       string            `json:"appName"`   // 应用名
	Replicas      *int32             `json:"replicas,omitempty"`  // 副本数量
	Image         string            `json:"image,omitempty"`     // 镜像
	Namespace     string            `json:"namespace"` // 命名空间
	LimitCpu      string            `json:"limitCpu,omitempty"`
	LimitMemory   string            `json:"limitMemory,omitempty"`
	RequestCpu    string            `json:"requestCpu,omitempty"`
	RequestMemory string            `json:"requestMemory,omitempty"`
	Env           map[string]string `json:"env,omitempty"`     // 环境变量
	Command       []string          `json:"command,omitempty"` // 启动命令
	Args          []string          `json:"args,omitempty"`    // 启动参数
	Labels        map[string]string `json:"labels"`
	Selector      map[string]string `json:"selector"`
	Annotation    map[string]string `json:"annotation,omitempty"` //标注
	//Volume        []Volume          // 挂载目录
	//VolumeMount   []VolumeMount
	//Ports         []PortMap // 多端口映射
	//Probe         Probe
}

type Probe struct {
	DelaySeconds     int32
	PeriodSeconds    int32
	FailureThreshold int32
	SuccessThreshold int32
	Timeout          int32
}

type PortMap struct {
	Port       int32
	TargetPort int
	Type       string
}

type Volume struct {
	Name     string
	HostPath string
}

type VolumeMount struct {
	Name      string
	MountPath string
}

func InitRequest() Request {
	var request Request
	var rep *int32
	var i int32 = 1
	rep = &i
	request.Replicas = rep
	request.RequestCpu =  viper.GetString("resources.requests.cpu")
	request.LimitCpu =  viper.GetString("resources.limit.cpu")
	request.RequestMemory = viper.GetString("resources.requests.memory")
	request.LimitMemory =  viper.GetString("resources.limits.memory")
	request.Env = map[string]string{}
	//request.Volume = []Volume{}
	//request.VolumeMount = []VolumeMount{}
	//request.Ports = []PortMap{}
	request.Command = []string{}
	request.Args = []string{}
	request.Selector = map[string]string{}
	request.Labels = map[string]string{}
	request.Annotation = map[string]string{}
	//request.Probe = Probe{
	//	DelaySeconds:     60,
	//	Timeout:          5,
	//	PeriodSeconds:    5,
	//	FailureThreshold: 1,
	//	SuccessThreshold: 1,
	//}
	return request
}
