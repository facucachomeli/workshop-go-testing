package usecase

import (
	"errors"
	"reflect"

	"github.com/facucachomeli/workshop-go-testing/domain"
)

type shipmentUseCase struct {
	save     func(*domain.Shipment) error
	getter   Getter
	sequence func() domain.ShipmentID
}

type Getter interface {
	GetByID(domain.ShipmentID) (domain.Shipment, error)
}

func NewShipmentUseCase(save func(*domain.Shipment) error, getter Getter, sequence func() domain.ShipmentID) shipmentUseCase {
	return shipmentUseCase{save, getter, sequence}
}

var CouldNotCreateShipment = errors.New("Could not create shipment")
var CouldNotCheckExistingShipment = errors.New("Could not check existing shipment")
var ShipmentAlreadyExists = errors.New("Shipment already exists")
var ShipmentDoesNotExist = errors.New("Shipment does not exist")
var ShipmentCanNotBeDelivered = errors.New("Shipement can not be delivered")

func (uc shipmentUseCase) Create(origin string, destination string) (domain.Shipment, error) {

	s, err := domain.NewShipment(uc.sequence(), origin, destination)
	if err != nil {
		return domain.Shipment{}, CouldNotCreateShipment
	}

	if err := uc.canCreateShipment(s); err != nil {
		return domain.Shipment{}, err
	}

	if err := uc.save(&s); err != nil {
		return domain.Shipment{}, CouldNotCreateShipment
	}

	return s, nil
}

func (uc shipmentUseCase) canCreateShipment(s domain.Shipment) error {
	s, err := uc.getter.GetByID(s.ID)
	if err != nil {
		return CouldNotCheckExistingShipment
	}

	if !reflect.DeepEqual(s, domain.Shipment{}) {
		return ShipmentAlreadyExists
	}

	return nil
}

func (uc shipmentUseCase) Deliver(id domain.ShipmentID) (domain.Shipment, error) {
	s, err := uc.getter.GetByID(id)
	if err != nil {
		return domain.Shipment{}, CouldNotCheckExistingShipment
	}

	if s.IsNil() {
		return domain.Shipment{}, ShipmentDoesNotExist
	}

	err = s.Deliver()
	if err != nil && err != domain.ShipmentAlreadyDelivered {
		return s, ShipmentCanNotBeDelivered
	}

	return s, nil
}
