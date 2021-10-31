package common

//go:generate stringer -type=RoleTypeEnum -linecomment=true

type RoleTypeEnum int

const (
	RoleTypeMaster RoleTypeEnum = iota
	RoleTypeSlaver
)
