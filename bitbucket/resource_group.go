package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

type Group struct {
	AccountName string `json:"accountname,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Name        string `json:"name,omitempty"`
	AutoAdd     bool   `json:"auto_add,omitempty"`
	Permission  string `json:"permission,omitempty"`
}

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Update: resourceGroupUpdate,
		Read:   resourceGroupRead,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"accountname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_add": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permission": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func newGroupFromResource(d *schema.ResourceData) *Group {
	group := &Group{
		AccountName: d.Get("accountname").(string),
		Name:        d.Get("name").(string),
		Slug:        d.Get("slug").(string),
		AutoAdd:     d.Get("auto_add").(bool),
		Permission:  d.Get("permission").(string),
	}
	return group
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	group := newGroupFromResource(d)

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(group)

	_, err := client.Put(
		fmt.Sprintf("1.0/groups/%s/%s", d.Get("accountname").(string), d.Get("slug").(string)),
		jsonpayload)
	if err != nil {
		return err
	}

	return resourceGroupRead(d, m)
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	// Only name can be passed at creation time, and not as JSON
	response, err := client.PostFormEncoded(
		fmt.Sprintf("1.0/groups/%s", d.Get("accountname").(string)),
		bytes.NewBufferString(fmt.Sprintf("name=%s", d.Get("name").(string))))
	if err != nil {
		return err
	}

	// Get the group's slug from the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	group := Group{}
	err = json.Unmarshal(body, &group)
	if err != nil {
		return err
	}
	d.Set("slug", group.Slug)

	d.SetId(string(fmt.Sprintf("%s/%s",
		d.Get("accountname").(string),
		d.Get("slug").(string),
	)))

	// Update the group after creation as not all settings can be passed at creation time
	return resourceGroupUpdate(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	response, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s",
		d.Get("accountname").(string),
		d.Get("slug").(string),
	))

	if response.StatusCode == 200 {

		var group Group

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &group)
		if err != nil {
			return err
		}

		d.Set("account_name", group.AccountName)
		d.Set("slug", group.Slug)
		d.Set("name", group.Name)
		d.Set("auto_add", group.AutoAdd)
		d.Set("permission", group.Permission)
	}

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("1.0/groups/%s/%s",
		d.Get("accountname").(string),
		d.Get("slug").(string),
	))

	return err
}
