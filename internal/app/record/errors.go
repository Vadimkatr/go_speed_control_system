package record

import (
	"errors"
	"fmt"
)

var (
	ErrValidateRecSpeed = errors.New(
		fmt.Sprintf("speed cannot be less than %f and more then %f", minSpeed, maxSpeed),
	)
	ErrValidateRecVehNum   = errors.New("vehicle number cannot be empty")
	ErrValidateRecDatetime = errors.New(
		fmt.Sprintf("datetime cannot be before %v and after %v", minDate, maxDate),
	)
)
