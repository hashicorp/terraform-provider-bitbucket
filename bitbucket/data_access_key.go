package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataAccessKey() *schema.Resource {
	return &schema.Resource{

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
				Type: schema.TypeString,
				//				Optional: true,
				Required: true,
			},
			"label": &schema.Schema{
				Type: schema.TypeString,
				//				Optional: true,
				Required: true,
			},
		},
	}
}

func dataAccessKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	akReq, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		d.Get("owner"),
		d.Get("repository"),
		url.PathEscape(d.Id()),
	))
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
