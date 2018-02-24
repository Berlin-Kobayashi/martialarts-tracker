package uuid

import (
	"testing"
	"regexp"
)

func TestV4(t *testing.T) {
	uuid := V4()

	uuidRegex := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")

	if !uuidRegex.Match([]byte(uuid)) {
		t.Errorf("%s aint an UUID V4", uuid)
	}
}