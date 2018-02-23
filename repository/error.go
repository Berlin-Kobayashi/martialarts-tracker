package repository

import "errors"

var NotFound = errors.New("training unit could not be found")
var Invalid = errors.New("training unit is invalid")
