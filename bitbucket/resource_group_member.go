package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

type GroupMember struct {
	Username    string `json:"username"`
	Firstname   string `json:"first_name,omitempty"`
	Lastname    string `json:"last_name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	ResourceURI string `json:"resource_uri,omitempty"`
}

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMemberCreate,
		Read:   resourceGroupMemberRead,
		Delete: resourceGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"accountname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_slug": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"avatar": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_uri": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupMemberCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	response, err := client.PutOnly(fmt.Sprintf("1.0/groups/%s/%s/members/%s",
		d.Get("accountname").(string),
		d.Get("group_slug").(string),
		d.Get("username").(string),
	))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		fmt.Errorf("Failed to add member")
	}

	d.SetId(string(fmt.Sprintf("%s/%s/%s",
		d.Get("accountname").(string),
		d.Get("group_slug").(string),
		d.Get("username").(string),
	)))

	return resourceGroupMemberRead(d, m)
}

func resourceGroupMemberRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	response, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s/members",
		d.Get("accountname").(string),
		d.Get("group_slug").(string),
	))

	if response.StatusCode == 200 {

		var groupMembers []GroupMember

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &groupMembers)
		if err != nil {
			return err
		}

		//FIXME: maybe cache groupMembers response to avoid too many API calls
		for _, groupMember := range groupMembers {
			if groupMember.Username == d.Get("username").(string) {
				d.Set("first_name", groupMember.Firstname)
				d.Set("last_name", groupMember.Lastname)
				d.Set("avatar", groupMember.Avatar)
				d.Set("resource_uri", groupMember.ResourceURI)
			}
		}
	}

	return nil
}

func resourceGroupMemberDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("1.0/groups/%s/%s/members/%s",
		d.Get("accountname").(string),
		d.Get("group_slug").(string),
		d.Get("username").(string),
	))

	return err
}
