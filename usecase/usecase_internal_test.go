package usecase

import (
	"errors"
	"testing"

	"github.com/facucachomeli/workshop-go-testing/domain"
)

type getterMock struct {
	mock func(domain.ShipmentID) (domain.Shipment, error)
}

func (m getterMock) GetByID(id domain.ShipmentID) (domain.Shipment, error) {
	return m.mock(id)
}

func TestShipmentUseCase_CanCreateShipment_CouldNotCheckExistingShipment(t *testing.T) {
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return domain.Shipment{}, errors.New("Getter error")
		},
	}
	uc := shipmentUseCase{
		nil,
		getter,
		nil,
	}
	s := domain.Shipment{}
	err := uc.canCreateShipment(s)
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != CouldNotCheckExistingShipment {
		t.Errorf("expected '%s' error but got '%s'", CouldNotCheckExistingShipment, err)
	}
	if !s.IsNil() {
		t.Errorf("expected shipment to be nil but got %#v", s)
	}
}

func TestShipmentUseCase_CanCreateShipment_ShipmentAlreadyExists(t *testing.T) {
	s := domain.Shipment{
		domain.ShipmentID(1),
		domain.Created,
		"valid origin",
		"valid destination",
	}
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return s, nil
		},
	}
	uc := shipmentUseCase{
		nil,
		getter,
		nil,
	}
	err := uc.canCreateShipment(s)
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != ShipmentAlreadyExists {
		t.Errorf("expected '%s' error but got '%s'", ShipmentAlreadyExists, err)
	}
	if s.State != domain.Created {
		t.Errorf("expected shipment to be Shipped but got %s", s.State)
	}
}

func TestShipmentUseCase_CanCreateShipment_OK(t *testing.T) {
	s := domain.Shipment{}
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return s, nil
		},
	}
	uc := shipmentUseCase{
		nil,
		getter,
		nil,
	}
	err := uc.canCreateShipment(s)
	if err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	}
}
