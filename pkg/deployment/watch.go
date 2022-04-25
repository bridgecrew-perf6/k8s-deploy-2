package deployment

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	deploymentutil "k8s.io/kubectl/pkg/util/deployment"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (d *Deployment) WatchDeploy(startWatchTime time.Time) {
	var stopper = make(chan struct{})
	var finish = make(chan struct{})
	clientset := d.client
	request := d.request
	log.Infof("Start watch deploy in namespace=%s for app=%s", request.Namespace, request.AppName)
	//var wg sync.WaitGroup
	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Apps().V1().Deployments().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			d := obj.(*appsv1.Deployment)
			if d.GetCreationTimestamp().After(startWatchTime) {
				log.Infof("Starting deployment in namespace=%s for app=%s at %s", d.GetNamespace(),
					d.Name, d.GetCreationTimestamp())
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			d := newObj.(*appsv1.Deployment)
			if d.Generation <= d.Status.ObservedGeneration {
				cond := deploymentutil.GetDeploymentCondition(d.Status, appsv1.DeploymentProgressing)
				if cond != nil && cond.Reason == deploymentutil.TimedOutReason {
					log.Errorf("deployment %s/%q exceeded its progress deadline", d.Namespace, d.Name)
					return
				}
				if d.Spec.Replicas != nil && d.Status.UpdatedReplicas < *d.Spec.Replicas {
					log.Infof("Waiting for deployment %s/%q rollout to finish: %d out of %d new replicas have been updated...", d.Namespace, d.Name, d.Status.UpdatedReplicas, *d.Spec.Replicas)
					return
				}
				if d.Status.Replicas > d.Status.UpdatedReplicas {
					log.Infof("Waiting for deployment %s/%q rollout to finish: %d old replicas are pending termination...", d.Namespace, d.Name, d.Status.Replicas-d.Status.UpdatedReplicas)
					return
				}
				if d.Status.AvailableReplicas < d.Status.UpdatedReplicas {
					log.Infof("Waiting for deployment %s/%q rollout to finish: %d of %d updated replicas are available...", d.Namespace, d.Name, d.Status.AvailableReplicas, d.Status.UpdatedReplicas)
					return
				}
				log.Infof("deployment %s/%q successfully rolled out", d.Namespace, d.Name)
				err := callback(request, true)
				if err != nil {
					log.Errorf("Call back error:%s", err.Error())
				}
				finish <- struct{}{}
				return
			}
		},
	})

	//stopper := make(chan struct{})
	go func() {
		informer.Run(stopper)
		//defer close(stopper)
		//wg.Wait()
	}()

	wt := viper.GetDuration("watchTimeout")
	select {
	case <-finish:
		close(stopper)
	case <-time.After(wt * time.Minute):
		log.Infof("Deployment timeout in namespace=%s for app=%s timed out, exiting", request.Namespace, request.AppName)
		err := callback(request, false)
		if err != nil {
			log.Errorf("Send call back data error for %s/%s: %s", request.Namespace, request.AppName, err.Error())
		}
		close(stopper)
		return
	}
}

func callback(request Request, finished bool) error {
	type callData struct {
		AppName   string `json:"appName"`
		Namespace string `json:"namespace"`
		Finished  bool   `json:"finished"`
	}
	url := viper.GetString("callBackURL")

	rb := callData{
		AppName:   request.AppName,
		Namespace: request.Namespace,
		Finished:  finished,
	}
	bodyStr, err := json.Marshal(&rb)
	if err != nil {
		return err
	}
	bodyBuffer := bytes.NewBuffer([]byte(bodyStr))
	req, err := http.NewRequest(http.MethodPost, url, bodyBuffer)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	defer req.Body.Close()

	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Second
	_, err = httpClient.Do(req)
	if err != nil {
		return err
	}
	log.Infof("Send call back data success: %s/%s", rb.Namespace, rb.AppName)
	return nil
}
