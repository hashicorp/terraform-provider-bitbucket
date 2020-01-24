package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"strings"
)

type SshKey struct {
	UUID    string `json:"uuid,omitempty"`
	Key     string `json:"key,omitempty"`
	Comment string `json:"comment,omitempty"`
	Label   string `json:"label,omitempty"`
	Owner   string `json:"-"`
}

func (k *SshKey) getAuthorizedKey() string {
	return fmt.Sprintf("%s %s", k.Key, k.Comment)
}

func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSshKeyCreate,
		Read:   resourceSshKeyRead,
		Delete: resourceSshKeyDelete,
		Update: resourceSshKeyUpdate,

		Schema: map[string]*schema.Schema{
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Username or UUID of the owner of the key`,
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateKey,
				Description:  `SSH public key in OpenSSH 'authorized_keys' format`,
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Managed by Terraform",
				Description: `Label for this key`,
			},
		},
	}
}

func createSshKey(d *schema.ResourceData) (*SshKey, error) {
	splitKey := strings.Split(d.Get("key").(string), " ")

	key := strings.Join(splitKey[:2], " ")
	comment := ""
	if len(splitKey) > 2 {
		comment = splitKey[2]
	}

	return &SshKey{
		UUID:    d.Id(),
		Label:   d.Get("label").(string),
		Key:     key,
		Comment: comment,
		Owner:   d.Get("owner").(string),
	}, nil
}

func resourceSshKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createSshKey(d)
	if err != nil {
		return err
	}

	keyReq, err := client.Post(
		fmt.Sprintf("2.0/users/%s/ssh-keys", key.Owner),
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

	d.SetId(key.UUID)

	return resourceSshKeyRead(d, m)
}

func resourceSshKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createSshKey(d)
	if err != nil {
		return err
	}

	keyReq, err := client.Get(
		fmt.Sprintf("2.0/users/%s/ssh-keys/%s", key.Owner, key.UUID),
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

func resourceSshKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createSshKey(d)
	if err != nil {
		return err
	}

	_, err = client.Delete(
		fmt.Sprintf("2.0/users/%s/ssh-keys/%s", key.Owner, key.UUID),
	)

	return err
}

func resourceSshKeyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	key, err := createSshKey(d)
	if err != nil {
		return err
	}

	_, err = client.Put(
		fmt.Sprintf("2.0/users/%s/ssh-keys/%s", key.Owner, key.UUID),
		bytes.NewBufferString(fmt.Sprintf(`{"label":"%s"}`, key.Label)),
	)

	if err != nil {
		return err
	}

	return resourceSshKeyRead(d, m)
}
