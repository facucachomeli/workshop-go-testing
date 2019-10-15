package domain_test

import (
	"testing"

	"github.com/facucachomeli/workshop-go-testing/domain"
	"github.com/stretchr/testify/assert"
)

func TestShipment_NewShipment_Error(t *testing.T) {
	cases := []struct {
		name          string
		id            domain.ShipmentID
		origin        string
		destination   string
		expectedError error
	}{
		{
			name:          "Invalid Origin",
			id:            1,
			origin:        "",
			destination:   "",
			expectedError: domain.InvalidOrigin,
		},
		{
			name:          "Invalid Destination",
			id:            1,
			origin:        "valid origin",
			destination:   "",
			expectedError: domain.InvalidDestination,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := domain.NewShipment(c.id, c.origin, c.destination)
			if err == nil {
				t.Errorf("expected error but found none")
			} else if c.expectedError != err {
				t.Errorf("expected '%s' error but got '%s'", c.expectedError, err)
			}
			if !s.IsNil() {
				t.Errorf("expected shipment to be nil but got %#v", s)
			}
		})
	}
}

func TestShipment_NewShipment_OK(t *testing.T) {
	id := domain.ShipmentID(1)
	origin := "valid origin"
	destination := "valid destination"

	s, err := domain.NewShipment(id, origin, destination)

	assert.Nilf(t, err, "expected error to be nil but got '%s'", err)
	assert.Equal(t, id, s.ID, "expected ID to be '%v' but got '%v'", id, s.ID)
	assert.Equal(t, origin, s.Origin, "expected Origin to be '%s' but got '%s'", origin, s.Origin)
	assert.Equal(t, destination, s.Destination, "expected Destination to be '%s' but got '%s'", destination, s.Destination)
	assert.Zero(t, s.State, "expected shipment state to empty but got %s", s.State)
}

func TestShipment_Create_Error(t *testing.T) {
	cases := []struct {
		name          string
		state         domain.ShipmentState
		expectedError error
	}{
		{
			name:          "Invalid State Cancelled",
			state:         domain.Cancelled,
			expectedError: domain.ShipmentAlreadyCreated,
		},
		{
			name:          "Invalid State Created",
			state:         domain.Created,
			expectedError: domain.ShipmentAlreadyCreated,
		},
		{
			name:          "Invalid State Handled",
			state:         domain.Handled,
			expectedError: domain.ShipmentAlreadyCreated,
		},
	}

	var s domain.Shipment
	for _, c := range cases {
		s = domain.Shipment{
			State: c.state,
		}
		t.Run(c.name, func(t *testing.T) {
			err := s.Create()
			assert.Equal(t, c.state, s.State)
			assert.NotNilf(t, err, "expected error but found none")
			assert.Equal(t, c.expectedError, err)
		})
	}
}

func TestShipment_Create_OK(t *testing.T) {
	s := domain.Shipment{
		State: "",
	}
	err := s.Create()

	assert.Nilf(t, err, "expected error to be nil but got '%s'", err)
	assert.Equal(t, domain.Created, s.State)
}

func TestShipment_Deliver_Error(t *testing.T) {
	cases := []struct {
		name          string
		state         domain.ShipmentState
		expectedError error
	}{
		{
			name:          "Invalid State Cancelled",
			state:         domain.Cancelled,
			expectedError: domain.InvalidStateForDeliver,
		},
		{
			name:          "Invalid State Created",
			state:         domain.Created,
			expectedError: domain.InvalidStateForDeliver,
		},
		{
			name:          "Invalid State Delivered",
			state:         domain.Delivered,
			expectedError: domain.ShipmentAlreadyDelivered,
		},
	}

	var s domain.Shipment
	for _, c := range cases {
		s = domain.Shipment{
			State: c.state,
		}
		t.Run(c.name, func(t *testing.T) {
			err := s.Deliver()
			assert.Equal(t, c.state, s.State)
			assert.Nilf(t, err)
			if err == nil {
				t.Errorf("expected error but found none")
			} else if c.expectedError != err {
				t.Errorf("expected '%s' error but got '%s'", c.expectedError, err)
			}
		})
}

func TestShipment_Deliver_OK(t *testing.T) {
	s := domain.Shipment{
		State: domain.Shipped,
	}
	err := s.Deliver()

	if err != nil {
		t.Fatalf("expected error to be nil but got '%s'", err)
	}
	if s.State != domain.Delivered {
		t.Errorf("expected shipment state to be '%s' but got '%s'", domain.Delivered, s.State)
	}
}
