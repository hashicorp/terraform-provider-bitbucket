package bitbucket

//https://confluence.atlassian.com/bitbucket/deploy-keys-resource-296095243.html

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	//	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// Accessey foo bar
type AccessKey struct {
	Pk    int    `json:"pk"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

func resourceAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessKeyCreate,
		Delete: resourceAccessKeyDetele,
		Exists: resourceAccessKeyExists,
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
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAccessKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	f := url.Values{}
	f.Set("key", url.PathEscape(d.Get("key").(string)))
	f.Set("label", d.Get("label").(string))

	endpoint := fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys", d.Get("owner").(string), d.Get("repository").(string))
	buffer := bytes.NewBufferString(f.Encode())
	//	buffer := strings.NewReader(f.Encode())

	fmt.Println(buffer)

	akResp, err := client.Post(endpoint, buffer)

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

		d.SetId(string(ak.Pk))
		d.Set("key", ak.Key)
		d.Set("label", ak.Label)
	}

	return nil
}

func resourceAccessKeyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketClient)

	akResp, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s", d.Get("owner").(string), d.Get("repository").(string), d.Get("pk")))

	if err != nil {
		return false, err
	}

	if akResp.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}

func resourceAccessKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	akReq, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s", d.Get("owner"), d.Get("repository"), d.Id()))

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

		//		d.SetId(string(ak.Pk))
		d.Set("key", ak.Key)
		d.Set("label", ak.Label)
	}

	return nil
}

func resourceAccessKeyUpdate(d *schema.ResourceData, m interface{}) error {
	//	client := m.(*BitbucketClient)

	//	akReq, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s", d.Get("owner"), d.Get("repository"), d.Id()))

	//	if err != nil {
	//		return err
	//	}

	//	if akReq.StatusCode == 200 {
	//		var ak AccessKey

	//		body, err := ioutil.ReadAll(akReq.Body)
	//		if err != nil {
	//			return err
	//		}

	//		decodingerr := json.Unmarshal(body, &ak)
	//		if decodingerr != nil {
	//			return decodingerr
	//		}

	//		//		d.SetId(string(ak.Pk))
	//		d.Set("key", ak.Key)
	//		d.Set("label", ak.Label)
	//	}

	resourceAccessKeyRead(d, m)

	return nil
}

func resourceAccessKeyDetele(d *schema.ResourceData, m interface{}) error {
	//	client := m.(*BitbucketClient)

	//	akReq, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s", d.Get("owner"), d.Get("repository"), d.Id()))

	//	if err != nil {
	//		return err
	//	}

	//	if akReq.StatusCode == 200 {
	//		var ak AccessKey

	//		body, err := ioutil.ReadAll(akReq.Body)
	//		if err != nil {
	//			return err
	//		}

	//		decodingerr := json.Unmarshal(body, &ak)
	//		if decodingerr != nil {
	//			return decodingerr
	//		}

	//		//		d.SetId(string(ak.Pk))
	//		d.Set("key", ak.Key)
	//		d.Set("label", ak.Label)
	//	}

	return nil
}
