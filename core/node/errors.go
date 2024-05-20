package node

import "errors"

var (
	ErrNodeNotInitialized = errors.New("Node object is not initialized")
	ErrNodeNotExists      = errors.New("Node is not exists")
	ErrWrongPassedStruct  = errors.New("passed wrong struct")
	ErrWrongPassedMap     = errors.New("passed wrong map")
	ErrNodeAlreadyExists  = errors.New("Node already exists")
	ErrPathNotFound       = errors.New("Node path not found")
)
