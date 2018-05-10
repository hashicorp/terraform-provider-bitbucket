package bitbucket

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
			},
		},
		ConfigureFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"bitbucket_access_key": dataAccessKey(),
			"bitbucket_user_v1":    dataSourceUserV1(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_access_key":        resourceAccessKey(),
			"bitbucket_hook":              resourceHook(),
			"bitbucket_default_reviewers": resourceDefaultReviewers(),
			"bitbucket_repository":        resourceRepository(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &BitbucketClient{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		HTTPClient: &http.Client{},
	}

	return client, nil
}
