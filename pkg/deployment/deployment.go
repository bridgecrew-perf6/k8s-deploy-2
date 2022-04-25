package deployment

import (
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"strings"
)

func createDeployment(clientset *kubernetes.Clientset, request Request) error {

	deploymentsClient := clientset.AppsV1().Deployments(request.Namespace)
	createOption := metav1.CreateOptions{}

	_, err := deploymentsClient.Create(context.Background(), GetDeployment(request), createOption)

	return err
}

func deleteDeployment(clientset *kubernetes.Clientset, request Request) error {
	deploymentsClient := clientset.AppsV1().Deployments(request.Namespace)

	//deletePolicy := metav1.DeletePropagationForeground
	deleteOption := metav1.DeleteOptions{}

	err := deploymentsClient.Delete(context.Background(), request.AppName, deleteOption)
	return err
}

func getK8sDeployment(clientset *kubernetes.Clientset, request Request) (*appsv1.Deployment, error) {
	deploymentsClient := clientset.AppsV1().Deployments(request.Namespace)
	getOption := metav1.GetOptions{}
	return deploymentsClient.Get(context.Background(), request.AppName, getOption)
}

func updateDeployment(clientset *kubernetes.Clientset, request Request) error {
	deploymentsClient := clientset.AppsV1().Deployments(request.Namespace)
	updateOption := metav1.UpdateOptions{}

	_, err := deploymentsClient.Update(context.Background(), GetDeployment(request), updateOption)
	return err
}

func patchDeployment(clientset *kubernetes.Clientset, request Request) error {
	curJson, modJson, err := getPatchDeployment(clientset, request)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, appsv1.Deployment{})
	if err != nil {
		return err
	}

	deploymentsClient := clientset.AppsV1().Deployments(request.Namespace)
	patchOption := metav1.PatchOptions{}

	_, err = deploymentsClient.Patch(context.Background(), request.AppName, types.StrategicMergePatchType, patch, patchOption)
	if err != nil {
		return err
	}
	return nil
}

func getPatchDeployment(clientset *kubernetes.Clientset, request Request) ([]byte, []byte, error) {
	curDeployment, err := getK8sDeployment(clientset, request)
	if err != nil {
		return nil, nil, err
	}
	curDep := &appsv1.Deployment{
		ObjectMeta: curDeployment.ObjectMeta,
		Spec: curDeployment.Spec,
	}
	curJson, err := json.Marshal(curDep)
	if err != nil {
		return nil, nil, err
	}
	mod := curDep.DeepCopy()
	if request.Image != "" {
		mod.Spec.Template.Spec.Containers[0].Image = request.Image
	}
	if request.Replicas != nil {
		mod.Spec.Replicas = request.Replicas
	}
	if request.LimitCpu != "" {
		mod.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse(request.LimitCpu),
		}
	}
	if request.LimitMemory != "" {
		mod.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(request.LimitMemory),
		}
	}
	if request.RequestCpu != "" {
		mod.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse(request.RequestCpu),
		}
	}
	if request.RequestMemory != "" {
		mod.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(request.RequestMemory),
		}
	}
	if request.Env != nil {
		envMap := request.Env
		var envs = []corev1.EnvVar{}

		for entity := range envMap {
			envs = append(envs, corev1.EnvVar{
				Name:  entity,
				Value: envMap[entity],
			})
		}
		mod.Spec.Template.Spec.Containers[0].Env = envs
	}
	if request.Command != nil {
		mod.Spec.Template.Spec.Containers[0].Command = request.Command
	}
	if request.Args != nil {
		mod.Spec.Template.Spec.Containers[0].Args = request.Args
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, nil, err
	}
	return curJson, modJson, nil
}


func GetDeployment(request Request) *appsv1.Deployment {
	var r corev1.ResourceRequirements
	//j := `{"limits": {"cpu":"2000m", "memory": "1Gi"}, "requests": {"cpu":"2000m", "memory": "1Gi"}}`
	j := `{"limits": {"cpu":"` + request.LimitCpu + `", "memory": "` + request.LimitMemory + `"}, "requests": {"cpu":"` + request.RequestCpu + `", "memory": "` + request.RequestMemory + `"}}`

	envMap := request.Env
	var envs = []corev1.EnvVar{}

	for entity := range envMap {
		envs = append(envs, corev1.EnvVar{
			Name:  entity,
			Value: envMap[entity],
		})
	}

	json.Unmarshal([]byte(j), &r)
	// 向label中添加app
	if request.Labels == nil {
		request.Labels = map[string]string{}
	}
	request.Selector["app"] = request.AppName
	request.Labels["app"] = request.AppName
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.AppName,
			Labels:      request.Labels,
			Annotations: request.Annotation,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: request.Selector,
			},
			Replicas: request.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      request.Labels,
					Annotations: request.Annotation,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:           request.AppName,
							Image:          request.Image,
							Env:            envs,
							Resources:      r,
							Command:        request.Command,
							//VolumeMounts:   getVolumeMount(request),
							Args:           request.Args,
							//LivenessProbe:  getProbe(request),
							//ReadinessProbe: getProbe(request),
						},
					},
					ImagePullSecrets: getImagePullSecrets(),
					//Volumes:          getVolume(request),
				},
			},
		},
	}
	return deployment
}

// 镜像私钥
func getImagePullSecrets() []corev1.LocalObjectReference {
	imagePullSecrets := viper.GetString("imagePullSecrets")
	secretArr := strings.Split(imagePullSecrets, ",")
	var result = []corev1.LocalObjectReference{}
	for index := range secretArr {
		result = append(result, corev1.LocalObjectReference{Name: secretArr[index]})
	}
	return result
}

//获取挂载卷
//func getVolume(request Request) []corev1.Volume {
//	var result = []corev1.Volume{}
//	for index := range request.Volume {
//		result = append(result, corev1.Volume{
//			Name: request.Volume[index].Name,
//			VolumeSource: corev1.VolumeSource{
//				HostPath: &corev1.HostPathVolumeSource{
//					Path: request.Volume[index].HostPath,
//				},
//			},
//		})
//	}
//	return result
//}

//将挂载卷挂载到容器的某目录
//func getVolumeMount(request Request) []corev1.VolumeMount {
//	var result = []corev1.VolumeMount{}
//	for index := range request.VolumeMount {
//		result = append(result, corev1.VolumeMount{
//			Name:      request.VolumeMount[index].Name,
//			MountPath: request.VolumeMount[index].MountPath,
//		})
//	}
//	return result
//}

//获取探活
//func getProbe(request Request) *corev1.Probe {
//	probe :=
//		corev1.Probe{
//			InitialDelaySeconds: request.Probe.DelaySeconds,
//			TimeoutSeconds:      request.Probe.Timeout,
//			PeriodSeconds:       request.Probe.PeriodSeconds,
//			FailureThreshold:    request.Probe.FailureThreshold,
//			SuccessThreshold:    request.Probe.SuccessThreshold,
//			Handler: corev1.Handler{
//				TCPSocket: &corev1.TCPSocketAction{Port: intstr.IntOrString{Type: 0, IntVal: request.Port}},
//			},
//		}
//
//	return &probe
//}
