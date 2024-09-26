// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"database/sql/driver"
	"fmt"
)

type RoleEnum string

const (
	RoleEnumInviter RoleEnum = "inviter"
	RoleEnumInvitee RoleEnum = "invitee"
)

func (e *RoleEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RoleEnum(s)
	case string:
		*e = RoleEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for RoleEnum: %T", src)
	}
	return nil
}

type NullRoleEnum struct {
	RoleEnum RoleEnum
	Valid    bool // Valid is true if RoleEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoleEnum) Scan(value interface{}) error {
	if value == nil {
		ns.RoleEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RoleEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoleEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RoleEnum), nil
}

type Connection struct {
	ConnectionID string
	ID           int64
	Alias        string
	MyRole       RoleEnum
}