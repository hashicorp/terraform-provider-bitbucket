package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

type GroupMember struct {
	Username    string `json:"username"`
	Firstname   string `json:"first_name,omitempty"`
	Lastname    string `json:"last_name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	ResourceURI string `json:"resource_uri,omitempty"`
}

type GroupMembers []GroupMember

type Group struct {
	AccountName string `json:"accountname,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Name        string `json:"name,omitempty"`
	AutoAdd     bool   `json:"auto_add,omitempty"`
	Permission  string `json:"permission,omitempty"`
	Members     GroupMembers
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
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"members": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	// Update group's main attributes
	accountName := d.Get("accountname").(string)
	groupSlug := d.Get("slug").(string)
	group := &Group{
		AccountName: accountName,
		Name:        d.Get("name").(string),
		Slug:        groupSlug,
		AutoAdd:     d.Get("auto_add").(bool),
		Permission:  d.Get("permission").(string),
	}

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(group)

	_, err := client.Put(
		fmt.Sprintf("1.0/groups/%s/%s", accountName, groupSlug),
		jsonpayload)
	if err != nil {
		return err
	}

	// Get desired group members
	membersData := d.Get("members").([]interface{})

	// Get current group members
	resourceGroupMembersRead(d, m)
	currentMembersData := d.Get("members").([]interface{})

	// Add missing group members
	for _, memberData := range membersData {
		username := memberData.(string)

		var found = false

		for i, currentMemberData := range currentMembersData {
			currentUsername := currentMemberData.(string)
			if currentUsername == username {
				found = true
				// Remove found member to collect undesired group members
				currentMembersData = append(currentMembersData[:i], currentMembersData[i+1:]...)
				break
			}
		}

		if !found {
			response, err := client.PutOnly(fmt.Sprintf("1.0/groups/%s/%s/members/%s",
				accountName,
				groupSlug,
				username,
			))

			if err != nil {
				return err
			}

			if response.StatusCode != 200 {
				return fmt.Errorf("Failed to add member %s", username)
			}
		}
	}

	// Remove undesired group members, i.e. those still in currentMembersData
	for _, currentMemberData := range currentMembersData {
		currentUsername := currentMemberData.(string)
		response, err := client.Delete(fmt.Sprintf("1.0/groups/%s/%s/members/%s",
			d.Get("accountname").(string),
			d.Get("slug").(string),
			currentUsername,
		))
		if err != nil {
			return err
		}

		if response.StatusCode != 204 {
			return fmt.Errorf("Failed to remove member %s", currentUsername)
		}
	}

	return resourceGroupRead(d, m)
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	// Only name can be passed at creation time, and not as JSON
	accountName := d.Get("accountname").(string)
	groupName := d.Get("name").(string)
	response, err := client.PostFormEncoded(
		fmt.Sprintf("1.0/groups/%s", accountName),
		bytes.NewBufferString(fmt.Sprintf("name=%s", groupName)))
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
		accountName,
		group.Slug,
	)))

	// Update the group after creation as not all settings can be passed at creation time
	return resourceGroupUpdate(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, "/")
		if len(idparts) == 2 {
			d.Set("accountname", idparts[0])
			d.Set("slug", idparts[1])
		} else {
			return fmt.Errorf("Incorrect ID format, should match `accountname/slug`")
		}
	}

	slug := d.Get("slug").(string)
	if slug == "" {
		slug = d.Get("name").(string)
	}

	client := m.(*BitbucketClient)
	response, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s",
		d.Get("accountname").(string),
		slug,
	))

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to read group")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var group Group
	err = json.Unmarshal(body, &group)
	if err != nil {
		return err
	}

	d.Set("account_name", group.AccountName)
	d.Set("slug", group.Slug)
	d.Set("name", group.Name)
	d.Set("auto_add", group.AutoAdd)
	d.Set("permission", group.Permission)

	return resourceGroupMembersRead(d, m)
}

func resourceGroupMembersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	response, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s/members",
		d.Get("accountname").(string),
		d.Get("slug").(string),
	))

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to read members")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var groupMembers []GroupMember
	err = json.Unmarshal(body, &groupMembers)
	if err != nil {
		return err
	}

	terraformMembers := make([]string, 0, len(groupMembers))

	for _, groupMember := range groupMembers {
		terraformMembers = append(terraformMembers, groupMember.Username)
	}

	d.Set("members", terraformMembers)
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
