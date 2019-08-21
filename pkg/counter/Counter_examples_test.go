// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package counter

import (
	"fmt"
)

// This example shows you can setup a counter and get all the values that have occured.
func ExampleCounter_all() {
	c := New()
	for i := 0; i < 10; i++ {
		c.Increment("foo")
	}
	for i := 0; i < 5; i++ {
		c.Increment("bar")
	}
	values := c.All(true) // get all values
	fmt.Println(values)
	// Output: [bar foo]
}

// This example shows you can setup a counter and find the value with the most number of occurences.
func ExampleCounter_top() {
	c := New()
	for i := 0; i < 10; i++ {
		c.Increment("foo")
	}
	for i := 0; i < 5; i++ {
		c.Increment("bar")
	}
	values := c.Top(1, 0, true) // get top value, greater than zero, and in sorted order
	if len(values) > 0 {
		fmt.Println(values[0])
	}
	// Output: foo
}

// This example shows you can setup a counter and find the value with the least number of occurences.
func ExampleCounter_bar() {
	c := New()
	for i := 0; i < 10; i++ {
		c.Increment("foo")
	}
	for i := 0; i < 5; i++ {
		c.Increment("bar")
	}
	values := c.Bottom(1, -1, true) // get bottom value, no maximum, and in sorted order
	if len(values) > 0 {
		fmt.Println(values[0])
	}
	// Output: bar
}
