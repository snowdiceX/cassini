package event

import (
	"testing"

	qostypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/assert"
)

func Test_ParseCoins(t *testing.T) {
	v, n, err := qostypes.ParseCoins("99HEHE")
	assert.NoError(t, err)

	assert.Equal(t, int64(0), v.Int64(), "wrong qos value")
	assert.Equal(t, int64(99), n[0].GetAmount().Int64(), "wrong qsc value")
	assert.Equal(t, "HEHE", n[0].GetName(), "wrong qsc name")
}
