package kubectl

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = "/root/.kube/config"
	namespace  = "default"
)

// Context holds the namespace and KubeClient
type Context struct {
	namespace  string
	kubeclient *kubernetes.Clientset
}

// NewContext return new Context
func NewContext() *Context {
	return &Context{}
}

// GetNamespace return the namespace
func (c *Context) GetNamespace() string {
	return c.namespace
}

// GetClient returns the kubeClient
func (c *Context) GetClient() *kubernetes.Clientset {
	return c.kubeclient
}

// Make prepares the kubeClient with kubeConfig
func (c *Context) Make() {

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	kubeclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	c.kubeclient = kubeclient
}
