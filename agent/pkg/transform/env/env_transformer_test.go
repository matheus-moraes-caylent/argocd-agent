package env

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type MockArgoApi struct {
}

func (api *MockArgoApi) CheckToken() error {
	panic("implement me")
}

func (m MockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (m MockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (m MockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	panic("implement me")
}

func (m MockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
	panic("implement me")
}

func (m MockArgoApi) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	var nodes = make([]argoSdk.Node, 0)
	nodes = append(nodes, argoSdk.Node{
		Kind: "Deploy",
		Uid:  "Uid",
		Health: argoSdk.Health{
			Status: "Health",
		},
	})

	nodes = append(nodes, argoSdk.Node{
		Kind: "Deploy",
		Uid:  "Uid2",
		Health: argoSdk.Health{
			Status: "Unhealth",
		},
	})

	return &argoSdk.ResourceTree{
		Nodes: nodes,
	}, nil
}

func (m MockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	var result []interface{}
	item := map[string]interface{}{
		"kind": "Application",
		"name": "app-name",
	}
	result = append(result, item)
	return result, nil
}

func (m MockArgoApi) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	liveState := "{\"kind\":\"Service\",\"metadata\":{ \"name\":\"test-api\",\"namespace\":\"andrii\",\"uid\":\"46263671-f290-11ea-8d49-42010a8001b0\"},\"spec\":{ \"template\": { \"spec\": { \"containers\":[{\"image\":\"andriicodefresh/test:v7\",\"name\":\"test-api\"}] } }, \"clusterIP\":\"10.27.251.224\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":1700}]}}"

	var resourceItems = make([]argoSdk.ManagedResourceItem, 0)
	resourceItems = append(resourceItems, argoSdk.ManagedResourceItem{
		Kind:        "Deployment",
		TargetState: "",
		LiveState:   liveState,
		Name:        "Test",
	})

	resourceItems = append(resourceItems, argoSdk.ManagedResourceItem{
		Kind:        "Application",
		TargetState: "",
		LiveState:   liveState,
		Name:        "RootApp",
	})

	return &argoSdk.ManagedResource{
		Items: resourceItems,
	}, nil
}

func (api *MockArgoApi) GetClusters() ([]argoSdk.ClusterItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetApplications() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetRepositories() ([]argoSdk.RepositoryItem, error) {
	panic("implement me")
}

func TestPrepareEnvironment(t *testing.T) {

	envTransformer := EnvTransformer{argoApi: &MockArgoApi{}}

	services, err := envTransformer.prepareEnvironmentActivity("test")
	if err != nil {
		t.Error(err)
	}

	if len(services) != 2 {
		t.Errorf("We should prepare 2 services for send to codefresh")
	}
	labels := map[string]interface{}{"app.kubernetes.io/instance": "apps-root"}
	status := map[string]interface{}{
		"operationState": map[string]interface{}{
			"syncResult": map[string]interface{}{"revision": "some revision"},
		},
	}
	envItem := map[string]interface{}{
		"status": status,
		"metadata": struct {
			name   string
			labels map[string]interface{}
		}{
			labels: labels,
		},
	}

	var env argoSdk.ArgoApplication

	util.Convert(envItem, &env)

	err, _ = envTransformer.PrepareEnvironment(env, 0)
	if err != nil {
		t.Errorf("Should successful finish PrepareEnvironment")
	}
}
