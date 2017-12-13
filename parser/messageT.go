package parser

type MessageType int

//go:generate stringer -type=MessageType

const (
	ISAUTHENTICATED MessageType = iota
	PUBLISH
	REMOVETOKEN
	SETTOKEN
	EVENT
	ACKRECEIVE
)
