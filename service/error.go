package service

import "errors"

var UnsupportedMethod = errors.New("method is not supported")
var UnsupportedEntity = errors.New("entity is not supported")
