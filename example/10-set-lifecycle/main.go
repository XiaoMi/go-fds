package main

import (
	"log"
	"os"

	"github.com/hujianxin/go-fds/fds"
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
		"rules":[
			{
				"prefix":"helloword",
				"id":"1",
				"actions":{
					"expiration":{
						"days":164
					}
				},
				"enabled":false
			},
			{
				"prefix":"helloword",
				"id":"2",
				"actions":{
					"expiration":{
						"days":164
					}
				},
				"enabled":true
			},
			{
				"prefix":"helloword",
				"id":"3",
				"actions":{
					"expiration":{
						"days":164
					}
				},
				"enabled":false
			}
		]
	}
	`
	lifecycle, err := fds.NewLifecycleConfigFromJSON([]byte(content))
	if err != nil {
		log.Fatal(err)
	}

	err = fdsClient.SetLifecycleConfig("bucketname", lifecycle)
	if err != nil {
		log.Fatal(err)
	}
}
