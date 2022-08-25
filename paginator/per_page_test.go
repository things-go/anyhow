package paginator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParam2(t *testing.T) {
	param := Param2{}
	assert.Equal(t, Param2{1, DefaultPageSize}, param.Inspect())

	assert.Equal(t, Param2{1, DefaultMaxPageSize}, New2(0, 20000))

	param = Param2{}
	assert.EqualValues(t, Param2{1, 20}, param.Inspect(20))
}