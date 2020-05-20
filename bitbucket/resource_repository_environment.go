package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// RepositoryEnvironment structure for handling key info
type RepositoryEnvironment struct {
	Name  string `json:"name"`
	Stage *Stage `json:"environment_type"`
	UUID  string `json:"uuid,omitempty"`
}

type Stage struct {
	Name string `json:"name"`
}

func resourceRepositoryEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryEnvironmentCreate,
		Update: resourceRepositoryEnvironmentUpdate,
		Read:   resourceRepositoryEnvironmentRead,
		Delete: resourceRepositoryEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stage": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Test",
					"Staging",
					"Production",
				},
					false),
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func newRepositoryEnvironmentFromResource(d *schema.ResourceData) *RepositoryEnvironment {
	dk := &RepositoryEnvironment{
		Name: d.Get("name").(string),
		Stage: &Stage{
			Name: d.Get("stage").(string),
		},
	}
	return dk
}

func resourceRepositoryEnvironmentCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	rvcr := newRepositoryEnvironmentFromResource(d)
	bytedata, err := json.Marshal(rvcr)

	if err != nil {
		return err
	}
	req, err := client.Post(fmt.Sprintf("2.0/repositories/%s/environments/",
		d.Get("repository").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	var repositoryEnvironment RepositoryEnvironment

	body, readerr := ioutil.ReadAll(req.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &repositoryEnvironment)
	if decodeerr != nil {
		return decodeerr
	}
	d.Set("uuid", repositoryEnvironment.UUID)
	d.SetId(fmt.Sprintf("%s:%s", d.Get("repository"), repositoryEnvironment.UUID))

	return resourceRepositoryEnvironmentRead(d, m)
}

func resourceRepositoryEnvironmentRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	rvReq, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	))

	log.Printf("ID: %s", url.PathEscape(d.Id()))

	if rvReq.StatusCode == 200 {
		var repositoryEnvironment RepositoryEnvironment
		body, readerr := ioutil.ReadAll(rvReq.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &repositoryEnvironment)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("uuid", repositoryEnvironment.UUID)
		d.Set("name", repositoryEnvironment.Name)
		d.Set("stage", repositoryEnvironment.Stage.Name)
	}

	if rvReq.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceRepositoryEnvironmentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	rvcr := newRepositoryEnvironmentFromResource(d)
	bytedata, err := json.Marshal(rvcr)

	if err != nil {
		return err
	}
	req, err := client.Put(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	if req.StatusCode != 200 {
		return nil
	}

	return resourceRepositoryEnvironmentRead(d, m)
}

func resourceRepositoryEnvironmentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	)))
	return err
}
