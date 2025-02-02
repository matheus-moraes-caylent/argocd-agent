package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

type GitContextQuestionnaire struct {
	codefreshApi codefresh.CodefreshApi
}

func NewGitContextQuestionnaire() *GitContextQuestionnaire {
	return &GitContextQuestionnaire{codefreshApi: codefresh.GetInstance()}
}

// AskAboutGitContext request git integration , should be selected from list of codefresh git contexts
func (questionnaire *GitContextQuestionnaire) AskAboutGitContext(installOptions *entity.InstallCmdOptions) error {
	if installOptions.Git.Integration != "" { // Integration is passed
		err, _ := questionnaire.codefreshApi.GetGitContextByName(installOptions.Git.Integration)
		if err != nil {
			return err
		}
		return nil
	}

	err, contexts := questionnaire.codefreshApi.GetGitContexts()
	if err != nil {
		return err
	}

	//var selectedContext string
	var values = make(map[string]string)
	var list []string
	for _, v := range *contexts {
		values[v.Metadata.Name] = v.Spec.Data.Auth.Password
		list = append(list, v.Metadata.Name)
	}

	if len(list) == 1 {
		installOptions.Git.Integration = list[0]
	} else {
		err, installOptions.Git.Integration = prompt.NewPrompt().Select(list, "Select Git/GithubApp context (Please create a dedicated context for the agent to  avoid hitting the Github rate limits or use github app integration)")
	}

	return err
}
