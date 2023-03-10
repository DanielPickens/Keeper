package kubernetes

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Danielpickens/Keeper/pkg/resource"
)

const (
	configDir  = ".kube"
	configFile = "config"
)

type Client struct {
	kubernetes   kubernetes.Interface
	namespaces   resource.NamespaceRepository
	pods         resource.PodRepository
	deployments  resource.DeploymentRepository
	statefulsets resource.StatefulsetRepository
	services     resource.ServiceRepository
	cluster      resource.ClusterRepository
	jobs         resource.JobRepository
}

// NewClient return a new kubernetes client
func NewClient(configFilePath string) (*Client, error) {

	config, err := clientcmd.BuildConfigFromFlags("", configFilePath)
	if err != nil {
		return &Client{}, fmt.Errorf("kubernetes client build config : %s", err.Error())
	}

	config.QPS = float32(250)
	config.Burst = 500

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &Client{}, fmt.Errorf("kubernetes new client for config : %s", err.Error())
	}

	return &Client{
		kubernetes:   clientSet,
		namespaces:   NewNamespaceRepository(clientSet),
		pods:         NewPodRepository(clientSet),
		deployments:  NewDeploymentRepository(clientSet),
		statefulsets: NewStatefulsetRepository(clientSet),
		services:     NewServiceRepository(clientSet, GetKubernetesHost(configFilePath)),
		cluster:      NewClusterRepository(),
		jobs:         NewJobRepository(clientSet),
	}, nil
}

func (c *Client) Jobs() resource.JobRepository {
	return c.jobs
}

func (c *Client) Namespaces() resource.NamespaceRepository {
	return c.namespaces
}

func (c *Client) Pods() resource.PodRepository {
	return c.pods
}

func (c *Client) Services() resource.ServiceRepository {
	return c.services
}

func (c *Client) Cluster() resource.ClusterRepository {
	return c.cluster
}

func (c *Client) Deployments() resource.DeploymentRepository {
	return c.deployments
}

func (c *Client) Statefulsets() resource.StatefulsetRepository {
	return c.statefulsets
}

// KubeConfigDefaultPath return the kubernetes default config path
func KubeConfigDefaultPath() string {
	return filepath.Join(homeDir(), configDir, configFile)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// GetKubernetesHost return the kubernetes cluster domain name used in the ~/.kube/config file
// The returned host takes the form : mydomainname.com
// Notice : this is just the host, without any schema or port.
func GetKubernetesHost(configFilePath string) string {

	config, _ := clientcmd.BuildConfigFromFlags("", configFilePath)

	u, err := url.Parse(config.Host)
	if err != nil {
		logrus.Fatalf("Impossible to get K8s host : %s", err.Error())
	}

	return strings.Split(u.Host, ":")[0]
}
