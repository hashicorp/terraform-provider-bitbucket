package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

type DeployKey struct {
	Pk    int    `json:"pk,omitempty"`
	Label string `json:"label,omitempty"`
	Key   string `json:"key,omitempty"`
}

func resourceDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeployKeyCreate,
		Read:   resourceDeployKeyRead,
		Update: resourceDeployKeyUpdate,
		Delete: resourceDeployKeyDelete,
		Exists: resourceDeployKeyExists,

		Schema: map[string]*schema.Schema{
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					pattern := "ssh-rsa AAAA[0-9A-Za-z+/]+[=]{0,3}( [^@]+@[^@]+)?"
					_, err := regexp.MatchString(pattern, value)
					if err != nil {
						errors = append(errors, fmt.Errorf(
							"%q must be a valid rsa key", k))
					}
					return
				},
			},
			"pk": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func createDeployKey(d *schema.ResourceData) *DeployKey {
	return &DeployKey{
		Pk:    d.Get("pk").(int),
		Label: d.Get("label").(string),
		Key:   d.Get("key").(string),
	}
}

func resourceDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	deployKey := createDeployKey(d)

	payload, err := json.Marshal(deployKey)
	if err != nil {
		return err
	}

	deployKey_req, err := client.Post(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys",
		d.Get("owner").(string),
		d.Get("repository").(string),
	), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	body, readerr := ioutil.ReadAll(deployKey_req.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &deployKey)
	if decodeerr != nil {
		return decodeerr
	}

	d.SetId(strconv.Itoa(deployKey.Pk))

	return resourceDeployKeyRead(d, m)
}

func resourceDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	deployKey_req, _ := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id()),
	))

	log.Printf("ID: %s", url.PathEscape(d.Id()))

	if deployKey_req.StatusCode == 200 {
		var deployKey DeployKey

		body, readerr := ioutil.ReadAll(deployKey_req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &deployKey)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("pk", deployKey.Pk)
		d.Set("label", deployKey.Label)
		d.Set("key", deployKey.Key)
	}

	return nil
}

func resourceDeployKeyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	deployKey := createDeployKey(d)

	payload, err := json.Marshal(deployKey)
	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id()),
	), bytes.NewBuffer(payload))

	if err != nil {
		return err
	}

	return resourceDeployKeyRead(d, m)
}

func resourceDeployKeyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketClient)
	if _, okay := d.GetOk("pk"); okay {
		deployKey_req, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
			d.Get("owner").(string),
			d.Get("repository").(string),
			url.PathEscape(d.Id()),
		))

		if err != nil {
			panic(err)
		}

		if deployKey_req.StatusCode != 200 {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func resourceDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id()),
	))

	return err
}
