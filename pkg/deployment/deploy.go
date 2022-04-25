package deployment

import (
	"errors"
	"k8s-deploy/pkg/kube"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

type Deployment struct {
	client *kubernetes.Clientset
	request Request
}

func New(request Request) (*Deployment, error) {
	env := request.Namespace
	client, err := kube.GetKubernetesClient(env)
	if err != nil {
		return nil, err
	}
	d := &Deployment{
		client: client,
		request: request,
	}
	return d, nil
}

func (d *Deployment) Create() error {
	deployment, err := getK8sDeployment(d.client, d.request)
	if err ==  nil{
		return errors.New(deployment.Name+"应用已存在")
	}

	err = createDeployment(d.client, d.request);
	if err != nil {
		return err
	}

	return nil
}

func (d *Deployment) Delete() error {
	err := deleteDeployment(d.client, d.request);
	if err != nil {
		return err
	}
	return nil
}

func (d *Deployment) Patch() error {
	err := patchDeployment(d.client, d.request)
	if err!=nil{
		return err
	}
	return nil
}

func (d *Deployment) Get() (*appsv1.Deployment,error) {
	deployment, err := getK8sDeployment(d.client, d.request)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}
