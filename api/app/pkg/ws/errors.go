package ws

import "errors"

var (
	ErrManagerClosed = errors.New("manager is closed")
	ErrCloseProperly = errors.New("error close properly")
)
