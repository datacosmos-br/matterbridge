package config

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

const stringSlice2DConfig = `
[general]
ReplaceMessages=[ ["cat","dog"], ["foo","bar"] ]

[irc.bad]
ReplaceNicks=[ ["a","b"], "c" ]

[irc.badvalue]
ReplaceNicks=[ ["a", 1] ]

[irc.notarray]
ExtractNicks="nope"
`

func TestGetStringSlice2D(t *testing.T) {
	logger, hook := test.NewNullLogger()
	cfg := NewConfigFromString(logger, []byte(stringSlice2DConfig))

	errorCount := func() int {
		count := 0

		for _, entry := range hook.AllEntries() {
			if entry.Level == logrus.ErrorLevel {
				count++
			}
		}

		return count
	}

	val, ok := cfg.GetStringSlice2D("general.ReplaceMessages")
	assert.True(t, ok)
	assert.Equal(t, [][]string{{"cat", "dog"}, {"foo", "bar"}}, val)

	// An optional setting that is not configured is not an error.
	val, ok = cfg.GetStringSlice2D("general.ReplaceNicks")
	assert.False(t, ok)
	assert.Nil(t, val)
	assert.Equal(t, 0, errorCount())

	val, ok = cfg.GetStringSlice2D("irc.bad.ReplaceNicks")
	assert.False(t, ok)
	assert.Nil(t, val)
	assert.Equal(t, 1, errorCount())

	val, ok = cfg.GetStringSlice2D("irc.badvalue.ReplaceNicks")
	assert.False(t, ok)
	assert.Nil(t, val)
	assert.Equal(t, 2, errorCount())

	val, ok = cfg.GetStringSlice2D("irc.notarray.ExtractNicks")
	assert.False(t, ok)
	assert.Nil(t, val)
	assert.Equal(t, 3, errorCount())
}
