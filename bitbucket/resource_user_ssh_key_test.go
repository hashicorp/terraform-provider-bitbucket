package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"net/url"
	"os"
	"testing"
)

const sshKeyResource = "bitbucket_user_ssh_key.test_ssh_key"

func TestAccBitbucketSshKey_basic(t *testing.T) {
	testUser := os.Getenv("BITBUCKET_USERNAME")

	template := `
		resource "bitbucket_user_ssh_key" "test_ssh_key" {
			owner = "%s"
			key   = "%s"
			label = "%s"
		}
	`

	key := generateKey(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketSshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:             fmt.Sprintf(template, testUser, key, "test"),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketSshKeyExists(sshKeyResource),
				),
			},
			{
				Config:             fmt.Sprintf(template, testUser, key, "test2"),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketSshKeyExists(sshKeyResource),
				),
			},
			{
				Config:             fmt.Sprintf(template, testUser, generateKey(t), "test2"),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketSshKeyExists(sshKeyResource),
				),
			},
		},
	})
}

func testAccCheckBitbucketSshKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	rs, ok := s.RootModule().Resources[sshKeyResource]
	if !ok {
		return fmt.Errorf("not found %s", sshKeyResource)
	}

	response, err := client.Get(fmt.Sprintf(
		"2.0/users/%s/ssh-keys/%s",
		rs.Primary.Attributes["owner"],
		url.PathEscape(rs.Primary.Attributes["id"]),
	))

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

func testAccCheckBitbucketSshKeyExists(name string) resource.TestCheckFunc {
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
