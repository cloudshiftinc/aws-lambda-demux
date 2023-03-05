package demux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockEvent = map[string]any{
	"abc":           "def",
	"something":     true,
	"somethingElse": 5,
	"nullValue":     nil,
}

func TestHasAttributeWorks(t *testing.T) {
	assert.True(t, HasAttribute(mockEvent, "abc"))
	assert.False(t, HasAttribute(mockEvent, "def"))
}

func TestStringAttributeValueWorks(t *testing.T) {
	v := StringAttributeValue(mockEvent, "abc")
	assert.NotNil(t, v)
	assert.Equal(t, "def", *v)

	v = StringAttributeValue(mockEvent, "something")
	assert.Nil(t, v)

	v = StringAttributeValue(mockEvent, "missingAttr")
	assert.Nil(t, v)
}

func TestStringAttributeMatches(t *testing.T) {
	assert.True(t, StringAttributeMatches(mockEvent, "abc", "def"))
	assert.False(t, StringAttributeMatches(mockEvent, "abc", "xxxf"))
	assert.False(t, StringAttributeMatches(mockEvent, "something", "true"))
}
