package main

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func run(k8sConfigFile string) {
	// creates the in-cluster config
	var config *rest.Config
	var err error

	if k8sConfigFile != "" {
		config, err = clientcmd.BuildConfigFromFlags("", k8sConfigFile)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	w, err := clientset.CoreV1().ServiceAccounts("").Watch(context.TODO(), metav1.ListOptions{})

	for event := range w.ResultChan() {
		//fmt.Printf("Type: %v\n", event.Type)
		p, ok := event.Object.(*v1.ServiceAccount)
		if !ok {
			log.Infof("unexpected type")
		}
		for ak, _ := range p.Annotations {
			if strings.EqualFold(ak, "azure.pod.identity/use") {
				s := v1.Secret{}
				s.Name = "arc-" + p.Name
				if event.Type == watch.Deleted {
					err = clientset.CoreV1().Secrets(p.Namespace).Delete(context.TODO(), s.Name, metav1.DeleteOptions{})
					if err != nil && !errors.IsNotFound(err) {
						log.Errorf("Unable to delete secret %s:%s. %v ", p.Namespace, s.Name, err.Error())
					} else if err == nil {
						log.Infof("Deleted secret %s:%s", p.Namespace, s.Name)
					}
				} else {
					s.StringData = make(map[string]string)
					token, err := getSAToken(p.Namespace, s.Name)
					if err != nil {
						panic(err)
					}
					s.StringData["token"] = token
					_, err = clientset.CoreV1().Secrets(p.Namespace).Create(context.TODO(), &s, metav1.CreateOptions{})
					if err != nil && !errors.IsAlreadyExists(err) {
						log.Errorf("Unable to add secret %s:%s. %v", p.Namespace, s.Name, err.Error())
					} else if err == nil {
						log.Infof("Created secret %s:%s", p.Namespace, s.Name)
					}
				}
				break
			}
		}
	}

	return
}

func loadSigningKey() error {
	s, err := clientset.CoreV1().Secrets("azure-arc").Get(context.TODO(), )
}