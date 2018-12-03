package fds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLifecycleConfigFromJSON(t *testing.T) {
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
	lifecycle, err := NewLifecycleConfigFromJSON([]byte(content))
	assert.Nil(t, err)
	assert.Equal(t, len(lifecycle.Rules), 3)

	assert.False(t, lifecycle.Rules[0].Enabled)
	assert.True(t, lifecycle.Rules[1].Enabled)

	assert.Equal(t, "3", lifecycle.Rules[2].ID)
}

func Test_NewLifecycleRuleFromJSON(t *testing.T) {
	content := `
{
    "prefix":"helloworld/",
    "id":"1",
    "actions":{
        "expiration":{
            "days":164
        }
    },
    "enabled":false
}
`
	rule, err := NewLifecycleRuleFromJSON([]byte(content))
	assert.Nil(t, err)
	assert.Equal(t, "1", rule.ID)
	assert.Equal(t, "helloworld/", rule.Prefix)

	assert.False(t, rule.Enabled)
	assert.Equal(t, float64(164), rule.Action["expiration"].Days)
}
