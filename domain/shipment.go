package domain

import "errors"

type Shipment struct {
	ID          ShipmentID
	State       ShipmentState
	Origin      string
	Destination string
}

type ShipmentID int

type ShipmentState string

var Created = ShipmentState("Created")
var Handled = ShipmentState("Handled")
var Shipped = ShipmentState("Shipped")
var Cancelled = ShipmentState("Cancelled")
var Delivered = ShipmentState("Delivered")

var InvalidOrigin = errors.New("Invalid Origin")
var InvalidDestination = errors.New("Invalid Destination")
var InvalidState = errors.New("Invalid State")
var ShipmentAlreadyCreated = errors.New("Shipment already has a state")
var InvalidStateForDeliver = errors.New("Shipment is not shipped")
var ShipmentAlreadyDelivered = errors.New("Shipment is already delivered")

func NewShipment(id ShipmentID, origin string, destination string) (Shipment, error) {
	if origin == "" {
		return Shipment{}, InvalidOrigin
	}
	if destination == "" {
		return Shipment{}, InvalidDestination
	}

	return Shipment{
		id,
		"",
		origin,
		destination,
	}, nil
}

func (s *Shipment) Create() error {
	if s.State != "" {
		return ShipmentAlreadyCreated
	}

	s.State = Created

	return nil
}

func (s *Shipment) Deliver() error {
	if s.State == Delivered {
		return ShipmentAlreadyDelivered
	}
	if s.State != Shipped {
		return InvalidStateForDeliver
	}

	s.State = Delivered

	return nil
}

func (s *Shipment) IsNil() bool {
	return s.ID == 0 &&
		s.State == "" &&
		s.Origin == "" &&
		s.Destination == ""
}
