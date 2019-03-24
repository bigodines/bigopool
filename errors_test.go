package bigopool

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	ee := errs{}
	assert.Empty(t, ee.All())

	ee.Append(errors.New("error 1"))
	ee.Append(errors.New("error 2"))
	assert.Equal(t, 2, len(ee.All()))
}

func TestToError(t *testing.T) {
	ee := errs{}
	assert.Nil(t, ee.ToError())

	ee.Append(errors.New("error 1"))
	ee.Append(errors.New("error 2"))
	assert.NotNil(t, ee.ToError())
}

func TestIsEmpty(t *testing.T) {
	ee := errs{}
	assert.True(t, ee.IsEmpty())

	ee.Append(errors.New("error 1"))
	ee.Append(errors.New("error 2"))
	assert.False(t, ee.IsEmpty())
}

func TestError(t *testing.T) {
	ee := errs{}

	ee.Append(errors.New("ok"))

	assert.Equal(t, "\nok", ee.Error())

	ee.Append(errors.New("two"))

	assert.Equal(t, "\nok\ntwo", ee.Error())
}
