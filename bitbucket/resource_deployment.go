package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// Deployment structure for handling key info
type Deployment struct {
	Name  string `json:"name"`
	Stage *Stage `json:"environment_type"`
	UUID  string `json:"uuid,omitempty"`
}

type Stage struct {
	Name string `json:"name"`
}

type Values struct {
	Values []Value `json:"values"`
}

type Value struct {
	Category Category `json:"category"`
	UUID     string   `json:"uuid,omitempty"`
}
type Category struct {
	Name string `json:"name"`
}

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Update: resourceDeploymentUpdate,
		Read:   resourceDeploymentRead,
		Delete: resourceDeploymentDelete,

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

func newDeploymentFromResource(d *schema.ResourceData) *Deployment {
	dk := &Deployment{
		Name: d.Get("name").(string),
		Stage: &Stage{
			Name: d.Get("stage").(string),
		},
	}
	return dk
}

func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	exists, err := checkIfNameAlreadyExists(d, m)
	if err != nil {
		return err
	}

	if !exists {
		client := m.(*Client)
		rvcr := newDeploymentFromResource(d)
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

		var deployment Deployment

		body, readerr := ioutil.ReadAll(req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &deployment)
		if decodeerr != nil {
			return decodeerr
		}
		d.Set("uuid", deployment.UUID)
		d.SetId(fmt.Sprintf("%s:%s", d.Get("repository"), deployment.UUID))

		return resourceDeploymentRead(d, m)
	}
	return nil
}

func checkIfNameAlreadyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	exists := false
	name := d.Get("name").(string)
	client := m.(*Client)
	req, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/environments/", d.Get("repository").(string)))

	if req.StatusCode == http.StatusOK {
		var values Values
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return false, err
		}

		err = json.Unmarshal(body, &values)
		if err != nil {
			return false, err
		}

		for _, x := range values.Values {
			if name == x.Category.Name {
				exists = true
				d.Set("uuid", x.UUID)
				d.SetId(fmt.Sprintf("%s:%s", d.Get("repository"), x.UUID))
			}
		}
	}

	if exists {
		return true, resourceDeploymentRead(d, m)
	}
	return false, nil
}

func resourceDeploymentRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	req, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	))

	log.Printf("ID: %s", url.PathEscape(d.Id()))

	if req.StatusCode == 200 {
		var Deployment Deployment
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &Deployment)
		if err != nil {
			return err
		}

		d.Set("uuid", Deployment.UUID)
		d.Set("name", Deployment.Name)
		d.Set("stage", Deployment.Stage.Name)
	}

	if req.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	rvcr := newDeploymentFromResource(d)
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

	return resourceDeploymentRead(d, m)
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	))
	return err
}
