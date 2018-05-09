package bitbucket

//https://confluence.atlassian.com/bitbucket/deploy-keys-resource-296095243.html

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
)

type AccessKey struct {
	Pk    uint64 `json:"pk,omitempty"`
	Key   string `json:"key,omitempty"`
	Label string `json:"label"`
}

func resourceAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessKeyCreate,
		Delete: resourceAccessKeyDetele,
		Read:   resourceAccessKeyRead,
		Update: resourceAccessKeyUpdate,

		Schema: map[string]*schema.Schema{
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAccessKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	akReq, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		d.Get("owner"),
		d.Get("repository"),
		url.PathEscape(d.Id())),
	)

	log.Printf("ID: %s", url.PathEscape(d.Id()))

	if err != nil {
		return err
	}

	if akReq.StatusCode == 200 {
		var ak AccessKey

		body, err := ioutil.ReadAll(akReq.Body)
		if err != nil {
			return err
		}

		decodingerr := json.Unmarshal(body, &ak)
		if decodingerr != nil {
			return decodingerr
		}

		d.SetId(fmt.Sprintf("%d", ak.Pk))
		d.Set("key", ak.Key)
		d.Set("label", ak.Label)
	}

	return nil
}

func resourceAccessKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	endpoint := fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys",
		d.Get("owner").(string),
		d.Get("repository").(string),
	)

	jsonpayload, err := json.Marshal(&AccessKey{
		Label: d.Get("label").(string),
		Key:   d.Get("key").(string),
	})
	if err != nil {
		return err
	}

	akResp, err := client.Post(endpoint, bytes.NewBuffer(jsonpayload))
	if err != nil {
		return err
	}

	if akResp.StatusCode == 200 {
		var ak AccessKey

		body, err := ioutil.ReadAll(akResp.Body)
		if err != nil {
			return err
		}

		decodingerr := json.Unmarshal(body, &ak)
		if decodingerr != nil {
			return decodingerr
		}

		log.Printf("[DEBUG] Access Key [%v]", ak)

		d.SetId(fmt.Sprintf("%d", ak.Pk))
		d.Set("key", ak.Key)
		d.Set("label", ak.Label)
	}

	return nil
}

func resourceAccessKeyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	endpoint := fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id()),
	)

	jsonpayload, err := json.Marshal(&AccessKey{
		Label: d.Get("label").(string),
	})
	if err != nil {
		return err
	}

	akResp, err := client.Put(endpoint, bytes.NewBuffer(jsonpayload))
	if err != nil {
		return err
	}

	if akResp.StatusCode == 200 {
		var ak AccessKey

		body, err := ioutil.ReadAll(akResp.Body)
		if err != nil {
			return err
		}

		decodingerr := json.Unmarshal(body, &ak)
		if decodingerr != nil {
			return decodingerr
		}

		log.Printf("[DEBUG] Access Key [%v]", ak)

		d.SetId(fmt.Sprintf("%d", ak.Pk))
		d.Set("key", ak.Key)
		d.Set("label", ak.Label)
	}

	return nil
}

func resourceAccessKeyDetele(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	_, err := client.Delete(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id())),
	)

	return err
}
