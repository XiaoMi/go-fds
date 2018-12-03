package main

import (
	"context"
	"log"
	"os"

	"github.com/hujianxin/galaxy-fds-sdk-go/fds"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

	acl := &fds.AccessControlList{
		Owner: fds.Owner{
			ID: "CI25efd5b0-5e83-4621-ab11-6ac863bcd164",
		},
	}

	grant := fds.Grant{
		Grantee: fds.GrantKey{
			ID: "ALL_USERS",
		},
		Permission: fds.GrantPermissionRead,
		Type:       fds.GrantTypeGroup,
	}

	acl.AddGrant(grant)

	err = fdsClient.SetBucketACLWithContext(context.Background(), "hellonihao", acl)
	if err != nil {
		log.Fatal(err)
	}
}
