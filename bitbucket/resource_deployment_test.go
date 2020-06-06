package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketDeployment_basic(t *testing.T) {
	var repo Deployment

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketDeploymentConfig := fmt.Sprintf(`
		resource "bitbucket_repository "test_repo" {
			owner = "%s"
			name = "test-repo-for-deployment-test"
		}
		resource "bitbucket_deployment" "test_deploy" {
			name = "test_deploy"
            stage = "Staging"
            repository = bitbucket_repository.test_repo.id
		}
	`, testUser)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDeploymentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeploymentExists("bitbucket_deployment.test_deploy", &repo),
				),
			},
		},
	})
}

func testAccCheckBitbucketDeploymentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	rs, ok := s.RootModule().Resources["bitbucket_deployment.test_deploy"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_deployment.test_deploy")
	}

	response, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/%s", rs.Primary.Attributes["owner"], rs.Primary.Attributes["name"]))

	if response.StatusCode != 404 {
		return fmt.Errorf("Deployment still exists")
	}

	return nil
}

func testAccCheckBitbucketDeploymentExists(n string, deployment *Deployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No deployment ID is set")
		}
		return nil
	}
}
