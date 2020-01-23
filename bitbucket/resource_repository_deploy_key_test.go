package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"net/url"
	"os"
	"testing"
)

const deployKeyResource = "bitbucket_repository_deploy_key.test_deploy_key"

func TestAccBitbucketDeployKey_basic(t *testing.T) {
	testUser := os.Getenv("BITBUCKET_USERNAME")
	template := `
		resource bitbucket_repository test_repo {
			owner = "%s"
			name  = "test-repo-for-deployment-key-test"
		}
		resource bitbucket_repository_deploy_key test_deploy_key {
			owner      = bitbucket_repository.test_repo.owner
			repository = bitbucket_repository.test_repo.slug
			key        = "%s"
			label      = "%s"
		}
	`

	key := generateKey(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDeployKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:             fmt.Sprintf(template, testUser, key, "test"),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeployKeyExists(deployKeyResource),
				),
			},
			{
				Config:             fmt.Sprintf(template, testUser, key, "test2"),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeployKeyExists(deployKeyResource),
				),
			},
			{
				Config:             fmt.Sprintf(template, testUser, generateKey(t), "test2"),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeployKeyExists(deployKeyResource),
				),
			},
		},
	})
}

func testAccCheckBitbucketDeployKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	rs, ok := s.RootModule().Resources[deployKeyResource]
	if !ok {
		return fmt.Errorf("not found %s", deployKeyResource)
	}

	response, err := client.Get(fmt.Sprintf(
		"2.0/repositories/%s/%s/deploy-keys/%s",
		rs.Primary.Attributes["owner"],
		rs.Primary.Attributes["repository"],
		url.PathEscape(rs.Primary.Attributes["uuid"])))

	if err == nil {
		return fmt.Errorf("%v the resource was found, should have errored", response.Status)
	}

	if response == nil {
		return err
	}

	if response.StatusCode != 404 {
		return fmt.Errorf("ssh key still exists")
	}

	return nil
}

func testAccCheckBitbucketDeployKeyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no SSH Key ID is set")
		}
		return nil
	}
}
