package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"encoding/json"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
)

func TestAccBitbucketGroupMember_basic(t *testing.T) {
	var groupMember GroupMember

	rInt := acctest.RandInt()

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketGroupMemberConfigStep0 := fmt.Sprintf(`
		resource "bitbucket_group" "test_group" {
			accountname = "%[1]s"
			name = "test group for group member test %[2]d"
			auto_add = true
			permission = "read"
		}

		resource "bitbucket_group_member" "test_group_member" {
			accountname = "%[1]s"
			group_slug = "${bitbucket_group.test_group.slug}"
			username = "%[1]s"
		}
	`, testUser, rInt)

	testAccBitbucketGroupMemberConfigStep1 := fmt.Sprintf(`
		resource "bitbucket_group" "test_group" {
			accountname = "%[1]s"
			name = "test group for group member test %[2]d"
			auto_add = true
			permission = "read"
		}
	`, testUser, rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketGroupMemberDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBitbucketGroupMemberConfigStep0,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketGroupMemberExists("bitbucket_group_member.test_group_member", &groupMember),
				),
			},
			resource.TestStep{
				Config: testAccBitbucketGroupMemberConfigStep1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketGroupMemberDoesNotExists("bitbucket_group_member.test_group_member", &groupMember),
				),
			},
		},
	})
}

func testAccCheckBitbucketGroupMemberDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bitbucket_group_member" {
			continue
		}

		client := testAccProvider.Meta().(*BitbucketClient)
		response, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s/members", rs.Primary.Attributes["accountname"], rs.Primary.Attributes["group_slug"]))

		if response.StatusCode == 200 {
			return fmt.Errorf("Group still exists, should have been destroyed")
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Unexpected error while getting group members")
		}
	}

	return nil
}

func testAccCheckBitbucketGroupMemberExists(n string, groupMember *GroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No group member ID is set")
		}

		client := testAccProvider.Meta().(*BitbucketClient)
		response, _ := client.Get(fmt.Sprintf("1.0/groups/%s/%s/members", rs.Primary.Attributes["accountname"], rs.Primary.Attributes["group_slug"]))

		if response.StatusCode == 200 {
			var groupMembers []GroupMember

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			err = json.Unmarshal(body, &groupMembers)
			if err != nil {
				return err
			}

			for _, groupMember := range groupMembers {
				if groupMember.Username == rs.Primary.Attributes["username"] {
					return nil
				}
			}
		} else {
			return fmt.Errorf("Unexpected error while getting group members")
		}

		return fmt.Errorf("Group member does not exists")
	}
}

func testAccCheckBitbucketGroupMemberDoesNotExists(n string, groupMember *GroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found %s", n)
		}
		return nil
	}
}
