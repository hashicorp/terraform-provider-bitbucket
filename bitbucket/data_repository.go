package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataReadRepository,

		Schema: RepositorySchema,
	}
}

func dataReadRepository(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	owner := d.Get("owner")
	if owner == "" {
		return fmt.Errorf("owner must not be blank")
	}

	var repoSlug string
	repoSlug = d.Get("slug").(string)
	if repoSlug == "" {
		repoSlug = d.Get("name").(string)
	}
	if repoSlug == "" {
		return fmt.Errorf("repo slug or name must not be blank")
	}

	repoReq, err := client.Get(fmt.Sprintf("2.0/repositories/%s/%s",
		owner,
		repoSlug,
	))
	if err != nil {
		return err
	}

	if repoReq.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repo not found")
	}

	if repoReq.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching repository")
	}

	var repo Repository

	body, readerr := ioutil.ReadAll(repoReq.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &repo)
	if decodeerr != nil {
		return decodeerr
	}

	d.SetId(string(fmt.Sprintf("%s/%s", owner, repoSlug)))
	d.Set("scm", repo.SCM)
	d.Set("is_private", repo.IsPrivate)
	d.Set("has_wiki", repo.HasWiki)
	d.Set("has_issues", repo.HasIssues)
	d.Set("name", repo.Name)
	if repo.Slug != "" && repo.Name != repo.Slug {
		d.Set("slug", repo.Slug)
	}
	d.Set("language", repo.Language)
	d.Set("fork_policy", repo.ForkPolicy)
	d.Set("website", repo.Website)
	d.Set("description", repo.Description)
	d.Set("project_key", repo.Project.Key)

	for _, cloneURL := range repo.Links.Clone {
		if cloneURL.Name == "https" {
			d.Set("clone_https", cloneURL.Href)
		} else {
			d.Set("clone_ssh", cloneURL.Href)
		}
	}
	pipelinesConfigReq, err := client.Get(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config",
		owner,
		repoSlug))

	if err == nil && pipelinesConfigReq.StatusCode == 200 {
		var pipelinesConfig PipelinesEnabled

		body, readerr := ioutil.ReadAll(pipelinesConfigReq.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &pipelinesConfig)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("pipelines_enabled", pipelinesConfig.Enabled)
	}

	return nil
}
