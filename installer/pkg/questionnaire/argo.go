package questionnaire

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"regexp"
)

type ArgoQuestionnaire struct {
	prompt prompt.Prompt
}

func NewArgoQuestionnaire() *ArgoQuestionnaire {
	return &ArgoQuestionnaire{prompt: prompt.NewPrompt()}
}

func retrieveHostFromLB(installOptions *entity.InstallCmdOptions, kubeClient kube.Kube) error {
	kubeOptions := installOptions.Kube
	argoServerSvc, err := kubeClient.GetArgoServerSvc(kubeOptions.Namespace)

	if err != nil {
		msg := fmt.Sprintf("We didn't find ArgoCD on \"%s/%s\"", installOptions.Kube.ClusterName, kubeOptions.Namespace)
		return errors.New(msg)
	}

	if kube.IsLoadBalancer(argoServerSvc) {
		balancerHost, err := kubeClient.GetLoadBalancerHost(argoServerSvc)
		if err != nil {
			return err
		}
		if balancerHost != "" {
			installOptions.Argo.Host = balancerHost
		}
		return nil
	}
	return errors.New("Failed to retrieve LoadBalancer information from argocd-server svc")
}

// AskAboutArgoCredentials request argocd credentials if it wasnt passed in cli during installation
func (argoQuestionnaire *ArgoQuestionnaire) AskAboutArgoCredentials(installOptions *entity.InstallCmdOptions, kubeClient kube.Kube) error {

	if installOptions.Argo.Host == "" {
		err := retrieveHostFromLB(installOptions, kubeClient)
		if err != nil {
			logger.Warning(err.Error())
		}
	}

	if installOptions.Argo.Host == "" {
		err := argoQuestionnaire.prompt.InputWithDefault(&installOptions.Argo.Host, "ArgoCD host", "https://example.com")
		if err != nil {
			return err
		}
	}

	withProtocol, err := regexp.MatchString("^https?://", installOptions.Argo.Host)
	if err != nil {
		return err
	}

	// customer not put protocol during installation
	if !withProtocol {
		installOptions.Argo.Host = "https://" + installOptions.Argo.Host
	}

	// removing / in the end
	installOptions.Argo.Host = regexp.MustCompile("/+$").ReplaceAllString(installOptions.Argo.Host, "")

	if (installOptions.Argo.Token != "") || ((installOptions.Argo.Username != "") && (installOptions.Argo.Password != "")) {
		return nil
	}

	//err, useArgocdToken := prompt.Confirm("Choose an authentication method")
	useArgocdToken := "Auth token - Recommended [https://codefresh.io/docs/docs/ci-cd-guides/gitops-deployments/]"
	useUserAndPass := "Username and password"
	authenticationMethodOptions := []string{useArgocdToken, useUserAndPass}
	err, authenticationMethod := argoQuestionnaire.prompt.Select(authenticationMethodOptions, "Choose an authentication method")
	if err != nil {
		return err
	}

	if authenticationMethod == useArgocdToken {
		err = argoQuestionnaire.prompt.InputWithDefault(&installOptions.Argo.Token, "Argo token", "")
		if err != nil {
			return err
		}
	} else if authenticationMethod == useUserAndPass {
		err = argoQuestionnaire.prompt.InputWithDefault(&installOptions.Argo.Username, "Argo username", "admin")
		if err != nil {
			return err
		}

		err = argoQuestionnaire.prompt.InputPassword(&installOptions.Argo.Password, "Argo password")
		if err != nil {
			return err
		}
	}

	return nil
}
