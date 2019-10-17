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

func Test_checkQscMax(t *testing.T) {
	qscStr := []string{
		"99HEHE",
		"3ppp,999aaa,11iii"}
	qscs := checkQscMax(qscStr)

	assert.Equal(t, 4, len(qscs), "parse error")
	assert.Equal(t, int64(99), qscs[0].GetAmount().Int64(), "wrong qsc value")
	assert.Equal(t, "HEHE", qscs[0].GetName(), "wrong qsc value")
}

func Test_checkQosMax(t *testing.T) {
	qosStr := []string{
		"11", "3", "99", "13"}
	max := checkQosMax(qosStr)

	assert.Equal(t, int64(99), max, "wrong qos value")
}
