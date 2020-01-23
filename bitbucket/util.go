package bitbucket

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"golang.org/x/crypto/ssh"
	"regexp"
	"strings"
	"testing"
)

// generateKey generates an SSH key
func generateKey(t *testing.T) string {
	rsa2, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	key2, err := ssh.NewPublicKey(&rsa2.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	key2Marshalled := ssh.MarshalAuthorizedKey(key2)
	return string(key2Marshalled[:len(key2Marshalled)-1])
}

// validateKey checks that an SSH key meets the OpenSSH authorized_keys format
func validateKey(val interface{}, key string) (warns []string, errs []error) {
	splitKey := strings.Split(val.(string), " ")
	if len(splitKey) < 2 {
		return nil, []error{fmt.Errorf("%s should be in OpenSSH authorized_keys format", key)}
	}

	if splitKey[0] != "ssh-ecdsa" && splitKey[0] != "ssh-rsa" && splitKey[0] != "ssh-ed25519" {
		errs = append(errs, fmt.Errorf("%s should start with 'ssh-ecdsa', 'ssh-rsa' or 'ssh-ed25519'", key))
	}

	matched, err := regexp.MatchString("^[-A-Za-z0-9+/=]+$", splitKey[1])
	if err != nil || !matched {
		errs = append(errs, fmt.Errorf("%s is not Base64-encoded", key))
	}

	return warns, errs
}
