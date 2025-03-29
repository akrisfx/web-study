package greetings

import (
	"regexp"
	"testing"
)

func TestFuckYou(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := FuckYou(name)
	if !want.MatchString(msg) || err != nil {
		t.Errorf(`FuckYou("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

func TestFuckYouEmpty(t *testing.T) {
	name := ""
	msg, err := FuckYou(name)
	if msg != "" || err == nil {
		t.Errorf(`FuckYou("") = %q, %v, want "", error`, msg, err)
	}
}
