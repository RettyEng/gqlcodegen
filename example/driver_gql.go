// DO NOT EDIT. this sourcecode is generated by gqlcodegen.
package example

import (
	"github.com/RettyInc/gqlcodegen/example/enum/class"
)

type DriverResolver interface {
	LicenceNumber() *string
	Name() string
	MiddleName() *string
	FamilyName() string
	IsOnDuty() bool
	Class() class.Class
}