package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

type Project struct {
	Key         string `json:"key,omitempty"`
	IsPrivate   bool   `json:"is_private,omitempty"`
	Owner       string `json:"owner.username,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Update: resourceProjectUpdate,
		Read:   resourceProjectRead,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
		},
	}
}

func newProjectFromResource(d *schema.ResourceData) *Project {
	project := &Project{
		Name:        d.Get("name").(string),
		IsPrivate:   d.Get("is_private").(bool),
		Description: d.Get("description").(string),
		Key:         d.Get("key").(string),
	}

	return project
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	project := newProjectFromResource(d)

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(project)

	var projectKey string
	projectKey = d.Get("key").(string)
	if projectKey == "" {
		projectKey = d.Get("key").(string)
	}

	_, err := client.Put(fmt.Sprintf("2.0/teams/%s/projects/%s",
		d.Get("owner").(string),
		projectKey,
	), jsonpayload)

	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	project := newProjectFromResource(d)

	bytedata, err := json.Marshal(project)

	if err != nil {
		return err
	}

	var projectKey string
	projectKey = d.Get("key").(string)
	if projectKey == "" {
		projectKey = d.Get("key").(string)
	}

	_, err = client.Post(fmt.Sprintf("2.0/teams/%s/projects/",
		d.Get("owner").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	d.SetId(string(fmt.Sprintf("%s/%s", d.Get("owner").(string), projectKey)))

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, "/")
		if len(idparts) == 2 {
			d.Set("owner", idparts[0])
			d.Set("key", idparts[1])
		} else {
			return fmt.Errorf("Incorrect ID format, should match `owner/key`")
		}
	}

	var projectKey string
	projectKey = d.Get("key").(string)
	if projectKey == "" {
		projectKey = d.Get("key").(string)
	}

	client := m.(*BitbucketClient)
	project_req, _ := client.Get(fmt.Sprintf("2.0/teams/%s/projects/%s",
		d.Get("owner").(string),
		projectKey,
	))

	if project_req.StatusCode == 200 {

		var project Project

		body, readerr := ioutil.ReadAll(project_req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &project)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("key", project.Key)
		d.Set("is_private", project.IsPrivate)
		d.Set("name", project.Name)
		d.Set("description", project.Description)
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {

	var projectKey string
	projectKey = d.Get("key").(string)
	if projectKey == "" {
		projectKey = d.Get("key").(string)
	}

	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("2.0/teams/%s/projects/%s",
		d.Get("owner").(string),
		projectKey,
	))

	return err
}
