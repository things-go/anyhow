package paginator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParam(t *testing.T) {
	param := Param{}
	assert.Equal(t, Param{1, DefaultPageSize}, param.Inspect())

	assert.Equal(t, Param{1, DefaultMaxPageSize}, New(0, 20000))

	param = Param{}
	assert.EqualValues(t, Param{1, 20}, param.Inspect(20))
}
