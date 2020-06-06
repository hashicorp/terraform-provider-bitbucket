package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccBitbucketDeploymentVariable_basic(t *testing.T) {

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketDeploymentVariableConfig := fmt.Sprintf(`
		resource "bitbucket_repository" "test_repo" {
			owner = "%s"
			name = "test-repo-default-reviewers"
		}
		resource "bitbucket_deployment" "test_deploy" {
			name = "testdeploy"
			stage = "Test"
			repository = bitbucket_repository.test_repo.id
		  }
		resource "bitbucket_deployment_variable" "testvar" {
			key = "test"
			value = "test"
			deployment = bitbucket_deployment.test_deploy.id
			secured = false
		  }
	`, testUser)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDeploymentVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDeploymentVariableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeploymentVariableExists("bitbucket_deployment_variable.testvar", "test", "test"),
				),
			},
		},
	})
}

func testAccCheckBitbucketDeploymentVariableDestroy(s *terraform.State) error {
	_, ok := s.RootModule().Resources["bitbucket_deployment_variable.testvar"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_deployment_variable.testvar")
	}
	return nil
}

func testAccCheckBitbucketDeploymentVariableExists(n, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found %s", n)
		}

		if rs.Primary.Attributes["key"] != key {
			return fmt.Errorf("Key not set")
		}

		if rs.Primary.Attributes["value"] != value {
			return fmt.Errorf("Key not set")
		}

		return nil
	}
}
