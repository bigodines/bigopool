package bigopool

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

func TestError(t *testing.T) {
	ee := errs{}

	ee.append(errors.New("one"))

	assert.Equal(t, "one", ee.Error())

	ee.append(errors.New("two"))

	assert.Equal(t, "one; two", ee.Error())
}

func TestUnwrapping(t *testing.T) {
	notFoundErr := errors.New("not found")
	badRequestErr := errors.New("bad request")

	ee := errs{}

	ee.append(notFoundErr)

	assert.True(t, errors.Is(ee.ToError(), notFoundErr))
	assert.False(t, errors.Is(ee.ToError(), badRequestErr))

	ee.append(badRequestErr)
	assert.True(t, errors.Is(ee.ToError(), badRequestErr))
}
