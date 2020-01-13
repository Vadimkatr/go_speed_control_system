package record

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecord_CreateRecord(t *testing.T) {
	type Params struct {
		dt time.Time
		vn string
		s  float32
	}
	testCases := []struct {
		name   string
		params Params
		err    error
	}{
		{
			name: "valid",
			params: Params{
				dt: time.Now(),
				vn: "1234 PP",
				s:  60.0,
			},
			err: nil,
		},
		{
			name: "invalid datetime",
			params: Params{
				dt: time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
				vn: "1234 PP",
				s:  60.0,
			},
			err: ErrValidateRecDatetime,
		},
		{
			name: "invalid vehicle number",
			params: Params{
				dt: time.Now(),
				vn: "",
				s:  60.0,
			},
			err: ErrValidateRecVehNum,
		},
		{
			name: "invalid speed",
			params: Params{
				dt: time.Now(),
				vn: "",
				s:  -1.0,
			},
			err: ErrValidateRecSpeed,
		},
		{
			name: "invalid speed",
			params: Params{
				dt: time.Now(),
				vn: "",
				s:  1000.0,
			},
			err: ErrValidateRecSpeed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CreateRecord(tc.params.dt, tc.params.vn, tc.params.s)
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
