package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// CloneURL is the internal struct we use to represent urls
type CloneURL struct {
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
}

// PipelinesEnabled is the struct we send to turn on or turn off pipelines for a repository
type PipelinesEnabled struct {
	Enabled bool `json:"enabled"`
}

// Repository is the struct we need to send off to the Bitbucket API to create a repository
type Repository struct {
	SCM         string `json:"scm,omitempty"`
	HasWiki     bool   `json:"has_wiki,omitempty"`
	HasIssues   bool   `json:"has_issues,omitempty"`
	Website     string `json:"website,omitempty"`
	IsPrivate   bool   `json:"is_private,omitempty"`
	ForkPolicy  string `json:"fork_policy,omitempty"`
	Language    string `json:"language,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	Project     struct {
		Key string `json:"key,omitempty"`
	} `json:"project,omitempty"`
	Links struct {
		Clone []CloneURL `json:"clone,omitempty"`
	} `json:"links,omitempty"`
}

var RepositorySchema = map[string]*schema.Schema{
	"scm": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "git",
	},
	"has_wiki": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"has_issues": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"website": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"clone_ssh": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"clone_https": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"project_key": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"is_private": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"pipelines_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"fork_policy": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "allow_forks",
	},
	"language": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"description": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Required: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"slug": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Update: resourceRepositoryUpdate,
		Read:   resourceRepositoryRead,
		Delete: resourceRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: RepositorySchema,
	}
}

func newRepositoryFromResource(d *schema.ResourceData) *Repository {
	repo := &Repository{
		Name:        d.Get("name").(string),
		Slug:        d.Get("slug").(string),
		Language:    d.Get("language").(string),
		IsPrivate:   d.Get("is_private").(bool),
		Description: d.Get("description").(string),
		ForkPolicy:  d.Get("fork_policy").(string),
		HasWiki:     d.Get("has_wiki").(bool),
		HasIssues:   d.Get("has_issues").(bool),
		SCM:         d.Get("scm").(string),
		Website:     d.Get("website").(string),
	}

	repo.Project.Key = d.Get("project_key").(string)
	return repo
}

func resourceRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	repository := newRepositoryFromResource(d)

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(repository)

	var repoSlug string
	repoSlug = d.Get("slug").(string)
	if repoSlug == "" {
		repoSlug = d.Get("name").(string)
	}

	_, err := client.Put(fmt.Sprintf("2.0/repositories/%s/%s",
		d.Get("owner").(string),
		repoSlug,
	), jsonpayload)

	if err != nil {
		return err
	}

	var pipelinesEnabled bool
	pipelinesEnabled = d.Get("pipelines_enabled").(bool)
	pipelinesConfig := &PipelinesEnabled{Enabled: pipelinesEnabled}

	bytedata, err := json.Marshal(pipelinesConfig)

	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config",
		d.Get("owner").(string),
		repoSlug), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}
	return resourceRepositoryRead(d, m)
}

func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	repo := newRepositoryFromResource(d)

	bytedata, err := json.Marshal(repo)

	if err != nil {
		return err
	}

	var repoSlug string
	repoSlug = d.Get("slug").(string)
	if repoSlug == "" {
		repoSlug = d.Get("name").(string)
	}

	_, err = client.Post(fmt.Sprintf("2.0/repositories/%s/%s",
		d.Get("owner").(string),
		repoSlug,
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}
	d.SetId(string(fmt.Sprintf("%s/%s", d.Get("owner").(string), repoSlug)))

	var pipelinesEnabled bool
	pipelinesEnabled = d.Get("pipelines_enabled").(bool)
	pipelinesConfig := &PipelinesEnabled{Enabled: pipelinesEnabled}

	bytedata, err = json.Marshal(pipelinesConfig)

	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config",
		d.Get("owner").(string),
		repoSlug), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, "/")
		if len(idparts) == 2 {
			d.Set("owner", idparts[0])
			d.Set("slug", idparts[1])
		} else {
			return fmt.Errorf("Incorrect ID format, should match `owner/slug`")
		}
	}

	var repoSlug string
	repoSlug = d.Get("slug").(string)
	if repoSlug == "" {
		repoSlug = d.Get("name").(string)
	}

	client := m.(*Client)
	repoReq, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/%s",
		d.Get("owner").(string),
		repoSlug,
	))

	if repoReq.StatusCode == 200 {

		var repo Repository

		body, readerr := ioutil.ReadAll(repoReq.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &repo)
		if decodeerr != nil {
			return decodeerr
		}

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
			d.Get("owner").(string),
			repoSlug))

		if err != nil {
			return err
		}

		if pipelinesConfigReq.StatusCode == 200 {
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

	}

	return nil
}

func resourceRepositoryDelete(d *schema.ResourceData, m interface{}) error {

	var repoSlug string
	repoSlug = d.Get("slug").(string)
	if repoSlug == "" {
		repoSlug = d.Get("name").(string)
	}

	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf("2.0/repositories/%s/%s",
		d.Get("owner").(string),
		repoSlug,
	))

	return err
}
