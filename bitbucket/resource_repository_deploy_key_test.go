package bitbucket

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const publicKeyContents = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD55FjmCr3Zi9IRlEMpLsdQ4lys9MgLZk8YXeEk31ZEwDtgbJ+RSCtHHeP+tTqfu9iNYGA+UUNE3HEO9PEg8vY307CEoQ4SrA4NmccrRnw960P54Zp+YvVTol2LfOn7hICMPJ/Wl4aSJI8jpK3bZc4iL6C/+Z3EbfE1Vm2cTafrLQy95pM2M5sZfCs+4k5UE7TzXjuE+qaSww0ZVVcQl4H5PzNNFW2yFtAUSjtzJt80kmYmue1dFDsUgEjqhqEp2/5XVzfKIhHkjIHYNulNouTbPGbZ6wuhfO5SikDApWgMkRbu386eN/H0wa0h4Jvg8Iy64Mrv3qZcpb7AXhWZLngR me@example.com"

func TestAccBitbucketDeployKey_basic(t *testing.T) {
	var key Key

	testUser := os.Getenv("BITBUCKET_USERNAME")
	config := fmt.Sprintf(`
		resource "bitbucket_repository" "test_repo" {
			owner = "%s"
			name = "test-repo-for-deploy-key-test"
		}
		resource "bitbucket_repository_deploy_key" "test_key" {
			owner = "%s"
			repository = "${bitbucket_repository.test_repo.name}"
			public_key_contents = "%s"
			label = "Test deploy key for terraform"
		}
	`, testUser, testUser, publicKeyContents)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDeployKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryDeployKeyExists(
						"bitbucket_repository_deploy_key.test_key", &key),
				),
			},
		},
	})
}

func testAccCheckBitbucketRepositoryDeployKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketClient)
	rs, ok := s.RootModule().Resources["bitbucket_repository_deploy_key.test_key"]
	if !ok {
		return fmt.Errorf("Not found %s",
			"bitbucket_repository_deploy_key.test_key")
	}

	response, err := client.Get(
		fmt.Sprintf("1.0/repositories/%s/%s/deploy-key/%s",
			rs.Primary.Attributes["owner"],
			rs.Primary.Attributes["repository"],
			url.PathEscape(rs.Primary.Attributes["uuid"])))

	if err == nil {
		return fmt.Errorf("The resource was found should have errored")
	}

	if response.StatusCode != 404 {
		return fmt.Errorf("Deploy key still exists")
	}

	return nil
}

func testAccCheckBitbucketRepositoryDeployKeyExists(n string,
	key *Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Deploy Key ID is set")
		}
		return nil
	}
}
