package storage

import "errors"

var NotFound = errors.New("resource could not be found")
var Invalid = errors.New("resource is invalid")
