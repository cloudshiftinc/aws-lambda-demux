package demux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockEvent EventStruct = map[string]any{
	"abc":           "def",
	"something":     true,
	"somethingElse": 5,
	"nullValue":     nil,
}

func TestHasAttributeWorks(t *testing.T) {
	assert.True(t, mockEvent.HasAttribute("abc"))
	assert.False(t, mockEvent.HasAttribute("def"))
}

func TestStringAttributeValueWorks(t *testing.T) {
	v := mockEvent.StringAttributeValue("abc")
	assert.NotNil(t, v)
	assert.Equal(t, "def", *v)

	v = mockEvent.StringAttributeValue("something")
	assert.Nil(t, v)

	v = mockEvent.StringAttributeValue("missingAttr")
	assert.Nil(t, v)
}

func TestStringAttributeMatches(t *testing.T) {
	assert.True(t, mockEvent.StringAttributeMatches("abc", "def"))
	assert.False(t, mockEvent.StringAttributeMatches("abc", "xxxf"))
	assert.False(t, mockEvent.StringAttributeMatches("something", "true"))
}
