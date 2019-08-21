// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package counter provides functions to generate a frequency histogram of values.
package counter

import (
	"sort"
)

// Counter is used for calculating a frequency histogram of strings.
type Counter map[string]int

// New returns a new Counter.
func New() Counter {
	return Counter(map[string]int{})
}

// Len returns the current number of unique values.
func (c Counter) Len() int {
	return len(c)
}

// Has returns true if the counter contains the value.
func (c Counter) Has(value string) bool {
	_, ok := c[value]
	return ok
}

// Count returns the current count for a given value.
// Returns 0 if the value has not occured.
func (c Counter) Count(value string) int {
	i, ok := c[value]
	if ok {
		return i
	}
	return 0
}

// Increment increases the count for a given value by 1.
func (c Counter) Increment(value string) {
	if count, ok := c[value]; ok {
		c[value] = count + 1
	} else {
		c[value] = 1
	}
}

// All returns all the values as a slice of strings.
// If s is set to true, then the values are sorted in alphabetical order.
func (c Counter) All(s bool) []string {
	values := make([]string, 0, len(c))
	for v := range c {
		values = append(values, v)
	}

	if s {
		sort.SliceStable(values, func(i, j int) bool {
			return values[i] < values[j]
		})
	}

	return values
}

// Top returns at most "n" values that have occured at least "min" times as a slice of strings.
// If s is set to true, then the values are sorted in descending order before the "n" values are chosen.
// If you would want to get the single most frequent value then use Top(1, 0, true).
// If you want 2 values that occured at least ten times, but do not care if they are the 2 most frequent values, then use Top(2, 10, false).
func (c Counter) Top(n int, min int, s bool) []string {

	if n == 0 {
		return make([]string, 0)
	}

	items := make([]struct {
		Value     string
		Frequency int
	}, 0)
	for value, frequency := range c {
		if frequency >= min {
			items = append(items, struct {
				Value     string
				Frequency int
			}{Value: value, Frequency: frequency})
		}
	}

	if s {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Frequency > items[j].Frequency
		})
	}

	if n > 0 && n < len(items) {
		values := make([]string, 0, n)
		for _, item := range items {
			values = append(values, item.Value)
		}
		return values[:n]
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}
	return values
}

// Top returns at most "n" values that have occured at most "max" times as a slice of strings.
// If max is less than zero, then ignore "max" as a threshold.
// If max is set to zero, then the function will return an empty slice of strings.
// If s is set to true, then the values are sorted in ascending order before the "n" values are chosen.
func (c Counter) Bottom(n int, max int, s bool) []string {

	if n == 0 {
		return make([]string, 0)
	}

	items := make([]struct {
		Value     string
		Frequency int
	}, 0)
	for value, frequency := range c {
		if max < 0 || frequency <= max {
			items = append(items, struct {
				Value     string
				Frequency int
			}{Value: value, Frequency: frequency})
		}
	}

	if s {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Frequency < items[j].Frequency
		})
	}

	if n > 0 && n < len(items) {
		values := make([]string, 0, n)
		for _, item := range items {
			values = append(values, item.Value)
		}
		return values[:n]
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}
	return values
}
