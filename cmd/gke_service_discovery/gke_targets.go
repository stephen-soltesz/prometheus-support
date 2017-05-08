package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	container "google.golang.org/api/container/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typesv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type GKESource struct {
	project string

	client           *http.Client
	computeService   *compute.Service
	containerService *container.Service
	targets          []interface{}
}

var (
	gkeScopes = []string{compute.CloudPlatformScope}
)

// NewGKESource returns a new GKESource object with an authenticated clients for Compute & Container APIs.
func NewGKESource(project string) (*GKESource, error) {
	gke := &GKESource{
		project: project,
	}

	var err error

	// Create a new authenticated HTTP client.
	gke.client, err = google.DefaultClient(oauth2.NoContext, gkeScopes...)
	if err != nil {
		return nil, fmt.Errorf("Error setting up Compute client: %s", err)
	}

	// Create a new Compute service instance.
	gke.computeService, err = compute.New(gke.client)
	if err != nil {
		return nil, fmt.Errorf("Error setting up Compute client: %s", err)
	}

	// Create a new Container Engine service object.
	gke.containerService, err = container.New(gke.client)
	if err != nil {
		return nil, fmt.Errorf("Error setting up Container Engine client: %s", err)
	}

	// Allocate space for the list of targets.
	gke.targets = make([]interface{}, 0)
	return gke, nil
}

// Saves collected targets to the given name.
func (gke *GKESource) Save(name string) error {
	// Convert the targets to JSON.
	data, err := json.MarshalIndent(gke.targets, "", "    ")
	if err != nil {
		log.Printf("Failed to Marshal JSON: %s", err)
		log.Printf("Pretty data: %s", pretty.Sprint(gke.targets))
		return err
	}

	// Save targets to output file.
	err = ioutil.WriteFile(name, data, 0644)
	if err != nil {
		log.Printf("Failed to write %s: %s", "output.json", err)
		return err
	}
	return nil
}

// Collect sources.
func (gke *GKESource) Collect() error {
	// Get all zones in a project.
	zoneListCall := gke.computeService.Zones.List("mlab-sandbox")
	err := zoneListCall.Pages(nil, func(zones *compute.ZoneList) error {
		for _, zone := range zones.Items {

			// Get all clusters in a zone.
			clusterList, err := gke.containerService.Projects.Zones.Clusters.List(
				gke.project, zone.Name).Do()
			if err != nil {
				return err
			}

			// Look for targets from every cluster.
			for _, cluster := range clusterList.Clusters {
				targets, err := checkCluster(zone, cluster)
				if err != nil {
					return err
				}
				gke.targets = append(gke.targets, targets...)
			}
		}
		return nil
	})
	return err
}

// checkCluster uses the kubernetes API to search for GKE targets.
func checkCluster(zone *compute.Zone, cluster *container.Cluster) ([]interface{}, error) {
	targets := []interface{}{}
	// Use information about the GKE cluster to create a k8s API client.
	kubeClient, err := gkeClusterToKubeClient(cluster)
	if err != nil {
		return nil, err
	}

	// Fetch all services in the k8s cluster.
	services, err := kubeClient.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	log.Printf("%s - %s - There are %d services in the cluster\n", zone.Name, cluster.Name, len(services.Items))

	// Check each service, and
	for _, service := range services.Items {
		// Federation scraping is opt-in only.
		if service.ObjectMeta.Annotations["gke-prometheus-federation/scrape"] != "true" {
			continue
		}
		values := findTargetAndLables(zone, cluster, service)
		if values != nil {
			targets = append(targets, values)
		}
	}
	return targets, nil
}

// findTargetAndLables identifies at most one target per service and returns a
// target configuration for use with Prometheus file service discovery.
func findTargetAndLables(zone *compute.Zone, cluster *container.Cluster, service typesv1.Service) interface{} {
	var target string

	if len(service.Spec.ExternalIPs) > 0 && len(service.Spec.Ports) > 0 {
		// Static IP addresses appear in the Service.Spec.
		// ---
		//    Spec: v1.ServiceSpec{
		//        ExternalIPs:              {"104.196.164.214"},
		//    },
		target = fmt.Sprintf("%s:%d",
			service.Spec.ExternalIPs[0],
			service.Spec.Ports[0].Port)
	} else if len(service.Status.LoadBalancer.Ingress) > 0 {
		// Ephemeral IP addresses appear in the Service.Status field.
		// ---
		//    Status: v1.ServiceStatus{
		//        LoadBalancer: v1.LoadBalancerStatus{
		//            Ingress: {
		//                {IP:"104.197.220.28", Hostname:""},
		//            },
		//        },
		//    },
		target = fmt.Sprintf("%s:%d",
			service.Status.LoadBalancer.Ingress[0].IP,
			service.Spec.Ports[0].Port)
	}
	if target == "" {
		return nil
	}
	values := map[string]interface{}{
		"labels": map[string]string{
			"service": service.ObjectMeta.Name,
			"cluster": cluster.Name,
			"zone":    zone.Name,
		},
		"targets": []string{target},
	}
	return values
}

// gkeClusterToKubeClient converts a container engine API Cluster object into
// a kubernetes API client configuration object.
func gkeClusterToKubeClient(c *container.Cluster) (*kubernetes.Clientset, error) {
	rawCaCert, err := base64.URLEncoding.DecodeString(c.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, err
	}
	clusterClient := api.Config{
		Clusters: map[string]*api.Cluster{
			// Define the cluster address and CA Certificate.
			"cluster": &api.Cluster{
				Server:                   fmt.Sprintf("https://%s", c.Endpoint),
				InsecureSkipTLSVerify:    false, // Require a valid CA Certificate.
				CertificateAuthorityData: rawCaCert,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			// Define the user credentials with access to the API.
			"user": &api.AuthInfo{
				Username: c.MasterAuth.Username,
				Password: c.MasterAuth.Password,
			},
		},
		Contexts: map[string]*api.Context{
			// Define a context that refers to the above cluster and user.
			"cluster-user": &api.Context{
				Cluster:  "cluster",
				AuthInfo: "user",
			},
		},
		// Use the above context.
		CurrentContext: "cluster-user",
	}
	restConfig, err := clientcmd.NewDefaultClientConfig(
		clusterClient, &clientcmd.ConfigOverrides{ClusterInfo: api.Cluster{Server: ""}}).ClientConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}
