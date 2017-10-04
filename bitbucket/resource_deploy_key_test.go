package bitbucket

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketDeployKey_basic(t *testing.T) {
	var deployKey DeployKey

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketDeployKeyConfig := fmt.Sprintf(`
    resource "bitbucket_repository" "test_repo" {
      owner = "%s"
      name  = "test-repo-for-deploy-key-test"
    }

    resource "bitbucket_deploy_key" "test_deploy_key" {
      owner      = "%s"
      repository = "${bitbucket_repository.test_repo.name}"
      label      = "test-deploy-key"
      key        = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEZBW8HOt+CKpEoFtF5q2NJNUyb+Z3wFcmZBX1UtW6CROJDz8AZfQTWQCi/7pz5+K1iqkZ7VvEg153MxbJXXa2sbpzeqLTuZdk8dGumhGxOGua6oLWLqO51k3H/dgK/tF4IQJqTe8p7XaolL4dnz87MU9GdDL1JV+ctdWH96+lX+9XGyC3momWNGCUxtGWwJAeyU0PSwcNmUjqqAryKMCrPtajKRjcjKS2WMpG1RML9nlkV4JLljof4wDo9aDxMhYSMyV1FQryUMcrOBaVbmP8AKru2AipHY89gReRG3pLgJrCe4Fi+d+BTqmMoJ2Sa8+RPPZA72sKg91+0KigIsl7 test@test"
    }
  `, testUser, testUser)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDeployKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBitbucketDeployKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeployKeyExists("bitbucket_deploy_key.test_deploy_key", &deployKey),
				),
			},
		},
	})
}

func testAccCheckBitbucketDeployKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketClient)
	rs, ok := s.RootModule().Resources["bitbucket_deploy_key.test_deploy_key"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_deploy_key.test_deploy_key")
	}

	response, err := client.Get(fmt.Sprintf("1.0/repositories/%s/%s/deploy-keys/%s",
		rs.Primary.Attributes["owner"],
		rs.Primary.Attributes["repository"],
		url.PathEscape(rs.Primary.Attributes["pk"])))

	if err == nil {
		return fmt.Errorf("The resource was found should have errored")
	}

	if response.StatusCode != 404 {
		return fmt.Errorf("DeployKey still exists")
	}

	return nil
}

func testAccCheckBitbucketDeployKeyExists(n string, deployKey *DeployKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No DeployKey ID is set")
		}
		return nil
	}
}
