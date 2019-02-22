// DO NOT EDIT. this sourcecode is generated by gqlcodegen.
package example

import (
	"context"

	"github.com/RettyInc/gqlcodegen/example/scalar"
)

type QueryResolver interface {
	Truck(context.Context, QueryResolver_Truck_Arg) TruckResolver
	Garage(context.Context, QueryResolver_Garage_Arg) GarageResolver
}

type QueryResolver_Truck_Arg struct {
	Number *scalar.RegistrationNumber
}

type QueryResolver_Garage_Arg struct {
	Id scalar.Uint32
}
