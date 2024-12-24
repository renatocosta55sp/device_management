package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/stretchr/testify/assert"
)

func TestInvalidArguments(t *testing.T) {

	device := NewDevice(uuid.UUID{})

	var tests = []struct {
		name, brand string
		want        error
	}{
		{

			name:  "",
			brand: "Apple",
			want:  ErrEmptyName,
		},
		{
			name:  "Android Samsung Galaxy",
			brand: "",
			want:  ErrEmptyBrand,
		},
		{
			name:  "",
			brand: "",
			want:  ErrEmptyName,
		},
	}

	for _, test := range tests {
		_, err := device.HandleAdd(
			commands.AddDeviceCommand{
				AggregateID: device.AggregateID,
				Name:        test.name,
				Brand:       test.brand,
			},
		)
		assert.Equal(t, err, test.want, "Expected: %d - Got: %d", test.want, err)
	}

}
