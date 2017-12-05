package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

// See https://confluence.atlassian.com/bitbucket/deploy-keys-resource-296095243.html
type Key struct {
	Id                int    `json:"pk,omitempty"`
	PublicKeyContents string `json:"key,omitempty"`
	Label             string `json:"label,omitempty"`
}

func resourceRepositoryDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryDeployKeyCreate,
		Read:   resourceRepositoryDeployKeyRead,
		Delete: resourceRepositoryDeployKeyDelete,
		Exists: resourceRepositoryDeployKeyExists,

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
			"public_key_contents": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func createKey(d *schema.ResourceData) *Key {
	return &Key{
		PublicKeyContents: d.Get("public_key_contents").(string),
		Label:             d.Get("label").(string),
	}
}

func resourceRepositoryDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	key := createKey(d)

	payload, err := json.Marshal(key)
	if err != nil {
		return err
	}

	keyRequest, err := client.Post(
		fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys",
			d.Get("owner").(string),
			d.Get("repository").(string),
		), bytes.NewBuffer(payload))

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(keyRequest.Body)
	if err != nil {
		return err
	}

	if keyRequest.StatusCode != 200 {
		return fmt.Errorf("Got non-200 response from POST: %s", body)
	}

	err = json.Unmarshal(body, &key)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(key.Id))

	return resourceRepositoryDeployKeyRead(d, m)
}

func resourceRepositoryDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	keyRequest, _ := client.Get(
		fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
			d.Get("owner").(string),
			d.Get("repository").(string),
			url.PathEscape(d.Id()),
		))

	body, err := ioutil.ReadAll(keyRequest.Body)
	if err != nil {
		return err
	}

	if keyRequest.StatusCode != 200 {
		return fmt.Errorf("Got non-200 response from read: %s", body)
	}

	var key Key
	err = json.Unmarshal(body, &key)
	if err != nil {
		return err
	}

	d.Set("public_key_contents", key.PublicKeyContents)
	d.Set("label", key.Label)

	return nil
}

func resourceRepositoryDeployKeyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketClient)
	keyRequest, err := client.Get(
		fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
			d.Get("owner").(string),
			d.Get("repository").(string),
			url.PathEscape(d.Id()),
		))

	if err != nil {
		return false, err
	}

	if keyRequest.StatusCode == 404 {
		return false, nil
	} else if keyRequest.StatusCode != 200 {
		return false, fmt.Errorf("Unexpected StatusCode %d from GET: %s",
			keyRequest.StatusCode, err)
	}

	return true, nil
}

func resourceRepositoryDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	keyRequest, err := client.Delete(
		fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
			d.Get("owner").(string),
			d.Get("repository").(string),
			url.PathEscape(d.Id()),
		))

	if err != nil {
		return err
	}

	if keyRequest.StatusCode != 204 {
		return fmt.Errorf("Unsuccessful response %d from DELETE",
			keyRequest.StatusCode)
	}

	return nil
}
