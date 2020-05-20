package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/url"
	"os"
	"testing"
)

func TestAccBitbucketBranchRestriction_basic(t *testing.T) {
	var branchRestriction BranchRestriction

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketBranchRestrictionConfig := fmt.Sprintf(`
		resource "bitbucket_repository" "test_repo" {
			owner = "%s"
			name = "test-repo-for-branch-restriction-test"
		}
		resource "bitbucket_branch_restriction" "test_repo_branch_restriction" {
			owner = "%s"
			repository = "${bitbucket_repository.test_repo.name}"
 			kind = "force"
 			pattern = "master"
		}
	`, testUser, testUser)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketBranchRestrictionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketBranchRestrictionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketBranchRestrictionExists("bitbucket_branch_restriction.test_repo_branch_restriction", &branchRestriction),
				),
			},
		},
	})
}

func testAccCheckBitbucketBranchRestrictionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	rs, ok := s.RootModule().Resources["bitbucket_branch_restriction.test_repo_branch_restriction"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_branch_restriction.test_repo_branch_restriction")
	}

	response, err := client.Get(fmt.Sprintf("2.0/repositories/%s/%s/branch-restrictions/%s", rs.Primary.Attributes["owner"], rs.Primary.Attributes["repository"], url.PathEscape(rs.Primary.Attributes["id"])))

	if err == nil {
		return fmt.Errorf("The resource was found should have errored")
	}

	if response.StatusCode != 404 {
		return fmt.Errorf("BranchRestriction still exists")
	}

	return nil
}

func testAccCheckBitbucketBranchRestrictionExists(n string, branchRestriction *BranchRestriction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No BranchRestriction ID is set")
		}
		return nil
	}
}
