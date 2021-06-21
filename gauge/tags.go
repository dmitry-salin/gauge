/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/

package gauge

import "strings"

type Tags struct {
	RawValues [][]string
}

func (tags *Tags) Add(values []string) {
	tags.RawValues = append(tags.RawValues, values)
}

func (tags *Tags) Values() (val []string) {
	for i := range tags.RawValues {
		val = append(val, tags.RawValues[i]...)
	}
	return val
}

func (tags *Tags) Kind() TokenKind {
	return TagKind
}

func SplitAndTrimTags(tag string) []string {
	listOfTags := strings.Split(tag, ",")
	for i, aTag := range listOfTags {
		listOfTags[i] = strings.TrimSpace(aTag)
	}
	return listOfTags
}
