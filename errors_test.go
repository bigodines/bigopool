package async

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	ee := errs{}
	assert.Empty(t, ee.All())

	ee.append(errors.New("error 1"))
	ee.append(errors.New("error 2"))
	assert.Equal(t, 2, len(ee.All()))
}

func TestToError(t *testing.T) {
	ee := errs{}
	assert.Nil(t, ee.ToError())

	ee.append(errors.New("error 1"))
	ee.append(errors.New("error 2"))
	assert.NotNil(t, ee.ToError())
}

func TestIsEmpty(t *testing.T) {
	ee := errs{}
	assert.True(t, ee.IsEmpty())

	ee.append(errors.New("error 1"))
	ee.append(errors.New("error 2"))
	assert.False(t, ee.IsEmpty())
}
