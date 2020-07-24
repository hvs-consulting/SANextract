package main

import (
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	// Remove proto
	r := sanitizeInput("https://www.google.com:443")
	if r != "www.google.com:443" {
		t.Errorf("Input sanitization tests failed")
	}
	r = sanitizeInput("http://google.com:443")
	if r != "google.com:443" {
		t.Errorf("Input sanitization tests failed")
	}
	r = sanitizeInput("smtp://mail.google.com:25")
	if r != "mail.google.com:25" {
		t.Errorf("Input sanitization tests failed")
	}

	// Do nothing
	r = sanitizeInput("www.google.com:443")
	if r != "www.google.com:443" {
		t.Errorf("Input sanitization tests failed")
	}
	r = sanitizeInput("mail.google.com:25")
	if r != "mail.google.com:25" {
		t.Errorf("Input sanitization tests failed")
	}

	// Remove path
	r = sanitizeInput("google.com:443/frontpage/index.html")
	if r != "google.com:443" {
		t.Errorf("Input sanitization tests failed")
	}
	r = sanitizeInput("mail.google.com:25/")
	if r != "mail.google.com:25" {
		t.Errorf("Input sanitization tests failed")
	}

	// Append default port
	r = sanitizeInput("https://google.com")
	if r != "google.com:443" {
		t.Errorf("Input sanitization tests failed")
	}
	r = sanitizeInput("mail.google.com")
	if r != "mail.google.com:443" {
		t.Errorf("Input sanitization tests failed")
	}

	// All together
	r = sanitizeInput("https://stackoverflow.com/questions/12968093/regex-to-validate-port-number")
	if r != "stackoverflow.com:443" {
		t.Errorf("Input sanitization tests failed")
	}
	r = sanitizeInput("https://www.google.com/search?sxsrf=ALeKk02rpH0l0fpZEl1FWQ3KrqZw2jCWyA%3A1590595850808&source=hp&ei=CpHOXtmSL4GWaq36gfgL&q=trololol&oq=trololol&gs_lcp=CgZwc3ktYWIQAzIECCMQJzICCAAyAggAMgQIABAKMgIIADIECAAQCjIECAAQCjIECAAQCjIECAAQCjIECAAQCjoHCCMQ6gIQJzoGCCMQJxATOgUIABCDAVCbQ1iYTGCXTWgCcAB4AIABZogB7wWSAQM3LjGYAQCgAQGqAQdnd3Mtd2l6sAEK&sclient=psy-ab&ved=0ahUKEwjZv6qlt9TpAhUBixoKHS19AL8Q4dUDCAk&uact=5%E2%80%8B%E2%80%8B%E2%80%8B%E2%80%8B%E2%80%8B%E2%80%8B%E2%80%8B")
	if r != "www.google.com:443" {
		t.Errorf("Input sanitization tests failed")
	}
	r = sanitizeInput("https://en.wikipedia.org:8080/wiki/SQL_injection")
	if r != "en.wikipedia.org:8080" {
		t.Errorf("Input sanitization tests failed")
	}

}
