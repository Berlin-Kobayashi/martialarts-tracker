package uuid

import (
	"crypto/rand"
	"fmt"
)

func V4() string {
	b := make([]byte, 16)
	rand.Read(b)

	newUUID := []byte(fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
	newUUID[14] = '4'
	newUUID[19] = 'a'

	return string(newUUID)
}
