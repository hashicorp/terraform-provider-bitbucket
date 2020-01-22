package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

type apiUser struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Nickname    string `json:"nickname"`
}

func dataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataReadUser,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataReadUser(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	username := d.Get("username")
	if username == "" {
		return fmt.Errorf("username must not be blank")
	}

	r, err := c.Get(fmt.Sprintf("2.0/users/%s", username))
	if err != nil {
		return err
	}

	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if r.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching user")
	}

	var u apiUser

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return err
	}

	d.SetId(u.UUID)
	d.Set("uuid", u.UUID)
	d.Set("nickname", u.Nickname)
	d.Set("display_name", u.DisplayName)

	return nil
}
