package main

import (
	"context"
	"flag"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"os"
	"path/filepath"
	"time"

	//"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNode() {
	clientset, err := createClient()
	if err != nil {
		klog.Error(err)
	}
	opts := v1.ListOptions{
		Limit: 100,
	}
	nodes,err := clientset.CoreV1().Nodes().List(context.TODO(),opts)
	if err != nil {
		klog.Error(err)
	}
	for _, v := range nodes.Items {
		klog.Info(v.Name)
	}
}

func GetPod() {
	clientset, err := createClient()
	if err != nil {
		klog.Error(err)
	}
	opts := v1.ListOptions{
		Limit: 100,
	}
	podwatch, err := clientset.CoreV1().Pods("kube-system").Watch(context.TODO(),opts)
	if err != nil {
		klog.Error(err)
	}
	for {
		select {
		case e, ok := <-podwatch.ResultChan():
			if !ok {
				klog.Warning("podWatch chan has been close!!!")
				klog.Info("clean chan over!")
				time.Sleep(time.Second * 5)
			}
			if e.Object != nil {
				klog.Info("chan is ok")
				klog.Info(e.Type)
				klog.Info(e.Object.DeepCopyObject())
			}
		}
	}
}

func GetDeployment(ns string) {
	clientset, err := createClient()
	if err != nil {
		klog.Error(err)
	}
	opts := v1.ListOptions{
		Watch: false,
	}
	dp, err := clientset.AppsV1().Deployments(ns).List(context.TODO(), opts)
	for _, v := range dp.Items {
		klog.Info(v.Name, v.Status, v.APIVersion)
	}
}

func GetService(ns string) {
	clientset, err := createClient()
	if err != nil {
		klog.Error(err)
	}
	opts := v1.ListOptions{
		Watch: false,
		Limit: 100,
	}
	svc, err := clientset.CoreV1().Services(ns).List(context.TODO(), opts)
	if err != nil {
		klog.Error(err)
	}
	for _, v := range svc.Items {
		klog.Info(v.Name, v.Namespace, v.Spec, v.APIVersion)
	}

}

func createClient() (clientset *kubernetes.Clientset, err error) {
	var kubeconfig *string
	if home := homeDir();home != "" {
		kubeconfig = flag.String("kubeconfig",filepath.Join(home,".kube","config"),"(optional) absolute path to kubeconfig file")
	}
	config, err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	return

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}



func main() {
	//GetNode()
	//GetPod()
	//GetDeployment("kube-system")
	GetService("kube-system")
}


