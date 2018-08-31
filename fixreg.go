package vkapi

import (
	"regexp"
)

var (
	fixTopicIDStr *regexp.Regexp
)

func init() {
	fixTopicIDStr = regexp.MustCompile(`"topic_id":"[0-9]+"`)
}
