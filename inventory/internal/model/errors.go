package model

import "errors"

// ErrPartNotFound — детали с таким UUID нет в каталоге.
var ErrPartNotFound = errors.New("part not found")
