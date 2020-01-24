package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"strings"
)

type DeployKey struct {
	Id         json.Number `json:"id,omitempty,Number"`
	Key        string      `json:"key,omitempty"`
	Comment    string      `json:"comment,omitempty"`
	Label      string      `json:"label,omitempty"`
	Repository string      `json:"-"`
	Owner      string      `json:"-"`
}

func (k *DeployKey) getAuthorizedKey() string {
	return fmt.Sprintf("%s %s", k.Key, k.Comment)
}

func resourceRepositoryDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryDeployKeyCreate,
		Read:   resourceRepositoryDeployKeyRead,
		Delete: resourceRepositoryDeployKeyDelete,
		Update: resourceRepositoryDeployKeyUpdate,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateKey,
				Description:  "The key in OpenSSH authorized_keys format",
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The owner of the repository (user or team)",
			},
			"repo_slug": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The repository slug",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The label for the key to be shown in Bitbucket",
			},
		},
	}
}

func createDeployKey(d *schema.ResourceData) (*DeployKey, error) {
	splitKey := strings.Split(d.Get("key").(string), " ")

	key := strings.Join(splitKey[:2], " ")
	comment := ""
	if len(splitKey) > 2 {
		comment = splitKey[2]
	}

	return &DeployKey{
		Id:         json.Number(d.Id()),
		Label:      d.Get("label").(string),
		Key:        key,
		Comment:    comment,
		Owner:      d.Get("owner").(string),
		Repository: d.Get("repo_slug").(string),
	}, nil
}

func resourceRepositoryDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createDeployKey(d)
	if err != nil {
		return err
	}

	keyReq, err := client.Post(
		fmt.Sprintf("2.0/repositories/%s/%s/deploy-keys", key.Owner, key.Repository),
		bytes.NewBufferString(fmt.Sprintf(`{"key":"%s","label":"%s"}`, key.getAuthorizedKey(), key.Label)),
	)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(keyReq.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &key)
	if err != nil {
		return err
	}

	d.SetId(string(key.Id))

	return resourceRepositoryDeployKeyRead(d, m)
}

func resourceRepositoryDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createDeployKey(d)
	if err != nil {
		return err
	}

	keyReq, err := client.Get(
		fmt.Sprintf("2.0/repositories/%s/%s/deploy-keys/%s", key.Owner, key.Repository, d.Id()),
	)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(keyReq.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &key)
	if err != nil {
		return err
	}

	_ = d.Set("key", key.getAuthorizedKey())
	_ = d.Set("label", key.Label)

	return nil
}

func resourceRepositoryDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createDeployKey(d)
	if err != nil {
		return err
	}

	_, err = client.Delete(
		fmt.Sprintf("2.0/repositories/%s/%s/deploy-keys/%s", key.Owner, key.Repository, d.Id()),
	)

	return err
}

func resourceRepositoryDeployKeyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createDeployKey(d)
	if err != nil {
		return err
	}

	_, err = client.Put(
		fmt.Sprintf("2.0/repositories/%s/%s/deploy-keys/%s", key.Owner, key.Repository, d.Id()),
		bytes.NewBufferString(fmt.Sprintf(`{"label":"%s"}`, key.Label)),
	)

	if err != nil {
		return err
	}

	return resourceRepositoryDeployKeyRead(d, m)
}
