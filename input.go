package main

import (
	"regexp"
	"strings"
)

func sanitizeInput(line string) string {
	// Remove proto
	trimProto := strings.Index(line, "://")
	if trimProto != -1 {
		line = line[trimProto+3:]
	}
	// strip anything after and including first '/'
	trimPath := strings.Index(line, "/")
	if trimPath != -1 {
		line = line[:trimPath]
	}

	// Use regex to check if target ends with a port specification, otherwise assume 443
	// Props: https://stackoverflow.com/questions/12968093/regex-to-validate-port-number
	// There may be corner cases like 00000 being accepted, you may subit a pull request
	// if this is an issue for you and you have a better solution
	match, _ := regexp.MatchString(":([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$", line)
	if !match {
		line = line + ":443"
	}
	return line
}
