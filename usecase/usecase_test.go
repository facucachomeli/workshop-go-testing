package usecase_test

import (
	"errors"
	"testing"

	"github.com/facucachomeli/workshop-go-testing/domain"
	"github.com/facucachomeli/workshop-go-testing/usecase"
)

type getterMock struct {
	mock func(domain.ShipmentID) (domain.Shipment, error)
}

func (m getterMock) GetByID(id domain.ShipmentID) (domain.Shipment, error) {
	return m.mock(id)
}

func TestShipmentUseCase_Create_CouldNotCreateShipment(t *testing.T) {
	sequence := func() domain.ShipmentID {
		return domain.ShipmentID(1)
	}
	uc := usecase.NewShipmentUseCase(nil, nil, sequence)

	s, err := uc.Create("", "")
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != usecase.CouldNotCreateShipment {
		t.Errorf("expected '%s' error but got '%s'", usecase.CouldNotCreateShipment, err)
	}
	if !s.IsNil() {
		t.Errorf("expected shipment to be nil but got %#v", s)
	}
}

func TestShipmentUseCase_Create_GetterError(t *testing.T) {
	sequence := func() domain.ShipmentID {
		return domain.ShipmentID(1)
	}
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return domain.Shipment{}, errors.New("Getter error")
		},
	}
	uc := usecase.NewShipmentUseCase(nil, getter, sequence)

	s, err := uc.Create("valid origin", "valid destination")
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != usecase.CouldNotCheckExistingShipment {
		t.Errorf("expected '%s' error but got '%s'", usecase.CouldNotCheckExistingShipment, err)
	}
	if !s.IsNil() {
		t.Errorf("expected shipment to be nil but got %#v", s)
	}
}

func TestShipmentUseCase_Create_SaveError(t *testing.T) {
	sequence := func() domain.ShipmentID {
		return domain.ShipmentID(1)
	}
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return domain.Shipment{}, nil
		},
	}
	save := func(*domain.Shipment) error {
		return errors.New("Save error")
	}

	uc := usecase.NewShipmentUseCase(save, getter, sequence)

	s, err := uc.Create("valid origin", "valid destination")
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != usecase.CouldNotCreateShipment {
		t.Errorf("expected '%s' error but got '%s'", usecase.CouldNotCreateShipment, err)
	}
	if !s.IsNil() {
		t.Errorf("expected shipment to be nil but got %#v", s)
	}
}

func TestShipmentUseCase_Create_OK(t *testing.T) {
	id := domain.ShipmentID(1)
	origin := "valid origin"
	destination := "valid destination"
	sequence := func() domain.ShipmentID {
		return id
	}
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return domain.Shipment{}, nil
		},
	}
	save := func(*domain.Shipment) error {
		return nil
	}

	uc := usecase.NewShipmentUseCase(save, getter, sequence)

	s, err := uc.Create(origin, destination)
	if err != nil {
		t.Fatalf("expected error to be nil but got '%s'", err)
	}
	if s.IsNil() {
		t.Fatalf("expected shipment not to be nil")
	}
	if s.ID != id {
		t.Errorf("expected ID to be '%v' but got '%v'", id, s.ID)
	}
	if s.Origin != origin {
		t.Errorf("expected Origin to be '%v' but got '%v'", origin, s.Origin)
	}
	if s.Destination != destination {
		t.Errorf("expected ID to be '%v' but got '%v'", destination, s.Destination)
	}
}

func TestShipmentUseCase_Deliver_CouldNotCheckExistingShipment(t *testing.T) {
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return domain.Shipment{}, errors.New("Get error")
		},
	}
	uc := usecase.NewShipmentUseCase(nil, getter, nil)

	s, err := uc.Deliver(domain.ShipmentID(1))
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != usecase.CouldNotCheckExistingShipment {
		t.Errorf("expected '%s' error but got '%s'", usecase.CouldNotCheckExistingShipment, err)
	}
	if !s.IsNil() {
		t.Errorf("expected shipment to be nil but got %#v", s)
	}
}

func TestShipmentUseCase_Deliver_ShipmentDoesNotExist(t *testing.T) {
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			return domain.Shipment{}, nil
		},
	}
	uc := usecase.NewShipmentUseCase(nil, getter, nil)

	s, err := uc.Deliver(domain.ShipmentID(1))
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != usecase.ShipmentDoesNotExist {
		t.Errorf("expected '%s' error but got '%s'", usecase.ShipmentDoesNotExist, err)
	}
	if !s.IsNil() {
		t.Errorf("expected shipment to be nil but got %#v", s)
	}
}

func TestShipmentUseCase_Deliver_ShipmentCanNotBeDelivered(t *testing.T) {
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			s := domain.Shipment{
				domain.ShipmentID(1),
				domain.Created,
				"valid origin",
				"valid destination",
			}
			return s, nil
		},
	}
	uc := usecase.NewShipmentUseCase(nil, getter, nil)

	s, err := uc.Deliver(domain.ShipmentID(1))
	if err == nil {
		t.Errorf("expected error but found none")
	} else if err != usecase.ShipmentCanNotBeDelivered {
		t.Errorf("expected '%s' error but got '%s'", usecase.ShipmentCanNotBeDelivered, err)
	}
	if s.IsNil() {
		t.Errorf("expected shipment not to be nil")
	}
}

func TestShipmentUseCase_Deliver_OK(t *testing.T) {
	id := domain.ShipmentID(1)
	origin := "valid origin"
	destination := "valid destination"
	sequence := func() domain.ShipmentID {
		return id
	}
	getter := getterMock{
		mock: func(domain.ShipmentID) (domain.Shipment, error) {
			s := domain.Shipment{
				domain.ShipmentID(1),
				domain.Shipped,
				"valid origin",
				"valid destination",
			}
			return s, nil
		},
	}
	save := func(*domain.Shipment) error {
		return nil
	}

	uc := usecase.NewShipmentUseCase(save, getter, sequence)

	s, err := uc.Deliver(id)
	if err != nil {
		t.Fatalf("expected error to be nil but got '%s'", err)
	}
	if s.IsNil() {
		t.Fatalf("expected shipment not to be nil")
	}
	if s.ID != id {
		t.Errorf("expected ID to be '%v' but got '%v'", id, s.ID)
	}
	if s.Origin != origin {
		t.Errorf("expected Origin to be '%v' but got '%v'", origin, s.Origin)
	}
	if s.Destination != destination {
		t.Errorf("expected ID to be '%v' but got '%v'", destination, s.Destination)
	}
}

// getter := getterMock{
// 	mock: func(domain.ShipmentID) (domain.Shipment, error) {
// 		s, _ := domain.NewShipment(domain.ShipmentID(1), "valid origin", "valid destination")
// 		return s, nil
// 	},
// }
