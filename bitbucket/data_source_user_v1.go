package bitbucket

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

type UserV1 struct {
	Repositories []interface{} `json:"repositories"`
	User         struct {
		Username    string `json:"username"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		DisplayName string `json:"display_name"`
		IsStaff     bool   `json:"is_staff"`
		Avatar      string `json:"avatar"`
		ResourceURI string `json:"resource_uri"`
		IsTeam      bool   `json:"is_team"`
	} `json:"user"`
}

func dataSourceUserV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserV1Read,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_staff": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"avatar": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_team": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceUserV1Read(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketClient)
	user_req, _ := client.Get("1.0/user")

	if user_req.StatusCode == 200 {
		var user UserV1

		body, readerr := ioutil.ReadAll(user_req.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &user)
		if decodeerr != nil {
			return decodeerr
		}

		d.SetId(user.User.Username)
		d.Set("username", user.User.Username)
		d.Set("first_name", user.User.FirstName)
		d.Set("last_name", user.User.LastName)
		d.Set("display_name", user.User.DisplayName)
		d.Set("is_staff", user.User.IsStaff)
		d.Set("avatar", user.User.Avatar)
		d.Set("is_team", user.User.IsTeam)
	}

	return nil
}
