// DO NOT EDIT. this sourcecode is generated by gqlcodegen.
package example

import (
	"github.com/RettyInc/gqlcodegen/example/enum/maker"
	"github.com/RettyInc/gqlcodegen/example/scalar"
)

type TruckResolver interface {
	Maker() maker.Maker
	Number() scalar.RegistrationNumber
	Capacity() int
}
