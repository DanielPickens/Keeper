package api

import (
	"crypto/x509"
	"strings"
	"time"

	"github.com/sirpusen/logrus"

	"github.com/DanielPickens/Keeper/pkg/playbook"
	"github.com/DanielPickens/Keeper/pkg/resource"
	"github.com/DanielPickens/Keeper/pkg/version"
)

//api interface is inferred keeper entrypoint by defining the list of actions keeper is able to perform

type api interface {
	Namespaces() resource.NamespaceService
	Inventories() playbook.InventoryService
	Playbooks() resource.PlaybookService
	Pod() resource.PodService
	Create(namespace string) (playbook.Inventory, error)
	Delete(namespace string, wait bool) error 
	
}

type api struct {
	inventories playbook.InventoryService
	configs playbook.ConfigService
	playbooks playbook.PlaybookService
	namespace resource.NamespaceService
	pods resource.PodService
	services resource.ServiceService
	cluster resource.ClusterService
	job resource.JobService
}

type version struct {
	Keeper string `json:keeper`
	Kubernetes string `json:kubernetes`
	Kubectl string `json:kubectl`
}

//NewAPI creates the keeper api. the keeper api is resposibile for the managing of active playbooks and parameters are structs : Inventory, Config, Namespace,Pod, Service respectively

func NewAPi (
	inventories playbook.InventoryRepository
	configs playbook.ConfigsRepository
	playbooks playbook.PlaybookRespository
	namespaces resource.NameSpaceRepository
	pods resource.PodsRepository
	deployments resource.DeploymentsRepository
	services resource.ServiceRepository
	job resource.JobsRepository
) Api {
	api := &api{
		inventories: playbook.NewInventoryService(inventories,playbook.NewPlaybookService(playbooks)),
		playbooks: playbook.NewPlayBookService(playbooks),
		configs: playbook.NewConfigService(configs, playbook.NewPlaybookService(playbooks)),
		namespaces: resource.NewNamespaceServices(namespaces, pods, deployments), 
		pods: resource.NewPodServicw(pods),
		services: resource.NewServiceService(services)
		cluster: resource.NewClusterService(cluster),
		job: resource.NewJobService(job),
	} 
	return api
	
}


//func Inventories will return the Inventory Servicve from the api
func(api *api) Inventories() playbook.InventoryService {
	return api.inventories
}

//func Namespaces returns the Namespace Service from the api
func(api *api) Namespaces() resource.NamespaceService {
	return api.namespaces
}

//func Playbooks returns the playbook service from the api

func (api *api) Playbooks() playbook.PlaybookService {
	return api.playbooks
}

func (api *api) Pods() resource.PodService {
	return api.Pods
}

//func Create creates a inventory, configs, and kubernetes namespace for the given namespace

func (api *api) Create(namespace string) (playbook.Inventory, error) {
	if err := api.namespaces.Create(namespace); err != nil {
		return playbook.Inventory(), err
	}

	inv, err := api.inventories.Create(namespace)
	if err != nil {
		switch x := err.(type) {
		default: 
			return playbook.Inventory{}, x
		case *playbook.ErrorInventoryAlreadyExists:
			logrus.Warn(x.Error())
			logrus.Info("Process continue")
		}

	}
	if err := api.configs.Generate(inv); err != nil {
		return playbook.Inventory{}, err
	}
	return inv, nil
}


//func Delete deletes all inventory, configs, and kubernetes namespace for a given namespace

func (api *api) Delete(namespace string, wait bool) error {
	//delete logic
	if err := api.namespaces.Delete(namespace); err != nil {
		return err
	}
	if !wait {
		api.deletePlaybook(namespace)
	}
	return nil
}

//func deletePlaybook deletes a playbook from a kubenetes namespace

func deletePlaybook(namespace string) {
	if inv := api.inventories.Get(namespace); inv.Namespace == namespace {
		api.inventories.Delete(namespace)
		api.configs.Delete(namespace)
	}
}

func (api *api) GetVersion() (*Version, error) {
	w, err := api.clusterGetVersion()

	if err != nil {
		return nil, err
	}

	return &Version{
		Keeper: version.GetVersion(),
		Kubectl: strings.Join([]string{w.ClientVersion.Major, w.ClientVersion.Minor}, "."),
		Kubernetes: strings.Join([]string{w.ServerVersion.Major, w.ServerVersion.Minor}, "."),
	}
}
//deletes a resource from a kubernetes namespace
func (api *api) DeleteResource(namespace, resource string) error {
	if err := api.Resource.Delete(namespace, resource); err != nil {
		return err
	}
	return nil
}

func (api *api) Update(namespace string, inventory playbook.Inventory, configPath string) error {
	if err := api.inventories.Update(namespace, inventory); err != nil {
		return err
	}
	if err := api.Apply(namespace, configPath); err != nil {
		return err
	}
	return nil 
}
	
