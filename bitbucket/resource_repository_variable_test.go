package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketRepositoryVariable_basic(t *testing.T) {

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketRepositoryVariableConfig := fmt.Sprintf(`
		resource "bitbucket_repository" "test_repo" {
			owner = "%s"
			name = "test-repo-default-reviewers"
		}

		resource "bitbucket_repository_variable" "testvar" {
			key = "test"
			value = "test"
			repository = "${bitbucket_repository.test_repo.id}"
			secured = false
		  }
	`, testUser)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketRepositoryVariableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryVariableExists("bitbucket_repository_variable.testvar", "test", "test"),
				),
			},
		},
	})
}

func testAccCheckBitbucketRepositoryVariableDestroy(s *terraform.State) error {
	_, ok := s.RootModule().Resources["bitbucket_repository_variable.testvar"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_repository_variable.testvar")
	}
	return nil
}

func testAccCheckBitbucketRepositoryVariableExists(n, key, value string) resource.TestCheckFunc {
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
