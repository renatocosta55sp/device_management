package domain

import (
	"testing"

	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/modeling/domain"
	"github.com/stretchr/testify/assert"
)

func TestInvalidArguments(t *testing.T) {

	device := NewDeviceAggregate([]domain.Event{})

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
		_, err := device.Add(
			commands.AddDeviceCommand{
				AggregateID: device.AggregateID,
				Name:        test.name,
				Brand:       test.brand,
			},
		)
		assert.Equal(t, err, test.want, "Expected: %d - Got: %d", test.want, err)
	}

}
