package service

import "errors"

var NotSupportedMethod = errors.New("method is not supported")
var NotSupportedEntity = errors.New("entity is not supported")
