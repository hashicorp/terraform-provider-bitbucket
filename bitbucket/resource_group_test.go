package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketGroup_basic(t *testing.T) {
	var group Group

	rInt := acctest.RandInt()

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketGroupConfig := fmt.Sprintf(`
		resource "bitbucket_group" "test_group" {
			accountname = "%s"
			name = "test group for group test %d"
			auto_add = true
			permission = "read"
		}
	`, testUser, rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBitbucketGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketGroupExists("bitbucket_group.test_group", &group),
				),
			},
		},
	})
}

func testAccCheckBitbucketGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*BitbucketClient)
	rs, ok := s.RootModule().Resources["bitbucket_group.test_group"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_group.test_group")
	}

	response, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s", rs.Primary.Attributes["accountname"], rs.Primary.Attributes["slug"]))

	if response.StatusCode != 404 {
		return fmt.Errorf("Group still exists")
	}

	return nil
}

func testAccCheckBitbucketGroupExists(n string, group *Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No group ID is set")
		}
		return nil
	}
}
