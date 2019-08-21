// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	c := New()
	assert.False(t, c.Has("foo"))
	c.Increment("foo")
	assert.True(t, c.Has("foo"))
	assert.Equal(t, 1, c.Count("foo"))
	c.Increment("foo")
	assert.Equal(t, 2, c.Count("foo"))
	c.Increment("bar")
	assert.True(t, c.Has("bar"))
	assert.Equal(t, 1, c.Count("bar"))
	assert.Equal(t, 2, c.Len())
}

func TestAll(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Increment("foo")
	}
	for i := 0; i < 5; i++ {
		c.Increment("bar")
	}
	assert.Equal(t, []string{"bar", "foo"}, c.All(true))
}

func TestTop(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Increment("foo")
	}
	for i := 0; i < 5; i++ {
		c.Increment("bar")
	}
	assert.Equal(t, []string{"foo"}, c.Top(1, 6, false))
	assert.Equal(t, []string{"foo"}, c.Top(1, 0, true))
	assert.Equal(t, []string{"foo", "bar"}, c.Top(2, 0, true))
	assert.Equal(t, []string{"foo", "bar"}, c.Top(3, 0, true))
}

func TestBottom(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Increment("foo")
	}
	for i := 0; i < 5; i++ {
		c.Increment("bar")
	}
	assert.Equal(t, []string{"bar"}, c.Bottom(1, 6, false))
	assert.Equal(t, []string{"bar"}, c.Bottom(1, -1, true))
	assert.Equal(t, []string{"bar", "foo"}, c.Bottom(2, -1, true))
	assert.Equal(t, []string{"bar", "foo"}, c.Bottom(3, -1, true))
}
