package argo

import (
	"fmt"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	store2 "github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	applicationCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}

	projectCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "appprojects",
	}
)

func buildConfig() *rest.Config {
	inCluster, _ := strconv.ParseBool(os.Getenv("IN_CLUSTER"))
	if inCluster {
		config, _ := rest.InClusterConfig()
		return config
	} else {
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}
		return config
	}
}

func watchApplicationChanges() {
	clientset, err := dynamic.NewForConfig(buildConfig())
	if err != nil {
		glog.Errorln(err)
	}

	kubeInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(clientset, time.Minute*30)
	applicationInformer := kubeInformerFactory.ForResource(applicationCRD).Informer()

	codefreshConfig := store2.GetStore().Codefresh
	api := codefresh2.Api{
		Token:       codefreshConfig.Token,
		Host:        codefreshConfig.Host,
		Integration: codefreshConfig.Integration,
	}

	//queue := Get()

	applicationInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			env := PrepareEnvironment(obj)
			api.SendEnvironment(env)

			//queue.Push(env)

			applications := GetApplications()
			err := api.SendResources("applications", AdaptArgoApplications(applications))
			if err != nil {
				fmt.Print(err)
			}

			log.Println(applications)
		},
		DeleteFunc: func(obj interface{}) {
			applications := GetApplications()
			err := api.SendResources("applications", AdaptArgoApplications(applications))
			if err != nil {
				fmt.Print(err)
			}

			log.Println(applications)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			env := PrepareEnvironment(newObj)
			api.SendEnvironment(env)

			//queue.Push(env)
		},
	})

	projectInformer := kubeInformerFactory.ForResource(projectCRD).Informer()

	projectInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("project added: %s \n", obj)
			projects := GetProjects()
			err := api.SendResources("projects", AdaptArgoProjects(projects))
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println(projects)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("project deleted: %s \n", obj)
			projects := GetProjects()
			err := api.SendResources("projects", AdaptArgoProjects(projects))
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println(projects)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("project updated: %s \n", newObj)
		},
	})

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)

	for {
		time.Sleep(time.Second)
	}

}

func Watch() {
	watchApplicationChanges()

}
