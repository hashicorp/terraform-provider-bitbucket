package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"io/ioutil"
	"log"
	"net/url"
)

type BranchRestriction struct {
	Id      int    `json:"id,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Pattern string `json:"pattern,omitempty"`
	Value   int    `json:"value,omitempty"`
}

func resourceBranchRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceBranchRestrictionsCreate,
		Read:   resourceBranchRestrictionsRead,
		Update: resourceBranchRestrictionsUpdate,
		Delete: resourceBranchRestrictionsDelete,
		Exists: resourceBranchRestrictionsExists,

		Schema: map[string]*schema.Schema{
			"owner": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"require_tasks_to_be_completed",
					"require_passing_builds_to_merge",
					"force",
					"require_all_dependencies_merged",
					"push",
					"require_approvals_to_merge",
					"enforce_merge_checks",
					"restrict_merges",
					"reset_pullrequest_approvals_on_change",
					"delete",
				},
					false),
			},
			"pattern": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func createBranchRestriction(d *schema.ResourceData) *BranchRestriction {
	return &BranchRestriction{
		Kind:    d.Get("kind").(string),
		Pattern: d.Get("pattern").(string),
		Value:   d.Get("value").(int),
	}
}

func resourceBranchRestrictionsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	branchRestriction := createBranchRestriction(d)

	bytedata, err := json.Marshal(branchRestriction)

	if err != nil {
		return err
	}

	branchRestrictionReq, err := client.Post(fmt.Sprintf("2.0/repositories/%s/%s/branch-restrictions",
		d.Get("owner").(string),
		d.Get("repository").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	body, readerr := ioutil.ReadAll(branchRestrictionReq.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &branchRestriction)
	if decodeerr != nil {
		return decodeerr
	}

	d.SetId(string(fmt.Sprintf("%v", branchRestriction.Id)))

	return resourceBranchRestrictionsRead(d, m)
}

func resourceBranchRestrictionsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)

	branchRestrictionsReq, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/%s/branch-restrictions/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id()),
	))

	log.Printf("ID: %s", url.PathEscape(d.Id()))

	if branchRestrictionsReq.StatusCode == 200 {
		var branchRestriction BranchRestriction
		body, readerr := ioutil.ReadAll(branchRestrictionsReq.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &branchRestriction)
		if decodeerr != nil {
			return decodeerr
		}

		d.SetId(string(fmt.Sprintf("%v", branchRestriction.Id)))
		d.Set("kind", branchRestriction.Kind)
		d.Set("pattern", branchRestriction.Pattern)
		d.Set("value", branchRestriction.Value)
	}

	return nil
}

func resourceBranchRestrictionsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	branchRestriction := createBranchRestriction(d)
	payload, err := json.Marshal(branchRestriction)
	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("2.0/repositories/%s/%s/branch-restrictions/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id()),
	), bytes.NewBuffer(payload))

	if err != nil {
		return err
	}

	return resourceBranchRestrictionsRead(d, m)
}

func resourceBranchRestrictionsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	_, err := client.Delete(fmt.Sprintf("2.0/repositories/%s/%s/branch-restrictions/%s",
		d.Get("owner").(string),
		d.Get("repository").(string),
		url.PathEscape(d.Id()),
	))

	return err
}

func resourceBranchRestrictionsExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*BitbucketClient)

	if v := d.Id(); v != "" {
		branchRestrictionsReq, err := client.Get(fmt.Sprintf("2.0/repositories/%s/%s/branch-restrictions/%s",
			d.Get("owner").(string),
			d.Get("repository").(string),
			url.PathEscape(d.Id()),
		))
		if err != nil {
			panic(err)
		}

		if branchRestrictionsReq.StatusCode != 200 {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
