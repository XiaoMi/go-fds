package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hujianxin/galaxy-fds-sdk-go/fds"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	fdsConf.EnableHTTPS = false
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

	content := `
			{
				"prefix":"other",
				"id":"1",
				"actions":{
					"expiration":{
						"days":164
					}
				},
				"enabled":false
			}
	`
	rule, err := fds.NewLifecycleRuleFromJSON([]byte(content))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", rule)

	err = fdsClient.SetLifecycleRule("bucketname", rule)
	if err != nil {
		log.Fatal(err)
	}
}
