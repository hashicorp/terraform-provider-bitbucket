package bitbucket

import "testing"

var validKeys = []string{
	"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKqP3Cr632C2dNhhgKVcon4ldUSAeKiku2yP9O9/bDtY user@myhost",
	"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKqP3Cr632C2dNhhgKVcon4ldUSAeKiku2yP9O9/bDtY",
	"ssh-rsa AAAAC3NzaC1lZDI1NTE5AAAAIKqP3Cr632C2dNhhgKVcon4ldUSAeKiku2yP9O9/bDtY",
	"ssh-ecdsa AAAAC3NzaC1lZDI1NTE5AAAAIKqP3Cr632C2dNhhgKVcon4ldUSAeKiku2yP9O9/bDtY",
}
var invalidKeys = []string{
	"",
	"ssh-notakeytype AAAAC3NzaC1lZDI1NTE5AAAAIKqP3Cr632C2dNhhgKVcon4ldUSAeKiku2yP9O9/bDtY user@myhost",
	"ssh-ed25519 @@@@@@@@@@@ user@myhost",
}

func TestValidateKey(t *testing.T) {
	for _, key := range validKeys {
		_, errs := validateKey(key, "key")
		if len(errs) > 0 {
			t.Logf("Key '%s' should be valid", key)
			t.Fail()
		}
	}

	for _, key := range invalidKeys {
		_, errs := validateKey(key, "key")
		if len(errs) < 1 {
			t.Logf("Key '%s' should not be valid", key)
			t.Fail()
		}
	}
}
