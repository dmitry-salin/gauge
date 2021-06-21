/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/

package filter

import (
	"os"
	"testing"

	"github.com/getgauge/gauge/env"
	"github.com/getgauge/gauge/gauge"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestFilterInvalidScenarios(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
	}
	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
	}
	var scenarios = []string{"First Scenario", "Third Scenario"}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	c.Assert(len(specs[0].Scenarios), Equals, 2)
	filteredScenarios := filterValidScenarios(specs, scenarios)
	c.Assert(len(filteredScenarios), Equals, 1)
	c.Assert(filteredScenarios[0], Equals, "First Scenario")
}

// ByScenarioName

func (s *MySuite) TestFilterScenariosByName(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
	}
	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
	}
	var scenarios = []string{"First Scenario"}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	c.Assert(len(specs[0].Scenarios), Equals, 2)
	specs = filterSpecsByScenarioName(specs, scenarios)
	c.Assert(len(specs[0].Scenarios), Equals, 1)
}

func (s *MySuite) TestFilterScenarioByNameWhichDoesNotExists(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
	}
	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
	}
	var scenarios = []string{"Third Scenario"}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	c.Assert(len(specs[0].Scenarios), Equals, 2)
	specs = filterSpecsByScenarioName(specs, scenarios)
	c.Assert(len(specs), Equals, 0)
}

func (s *MySuite) TestFilterMultipleScenariosByName(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
	}
	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
	}
	var scenarios = []string{"First Scenario", "Second Scenario"}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	c.Assert(len(specs[0].Scenarios), Equals, 2)
	specs = filterSpecsByScenarioName(specs, scenarios)
	c.Assert(len(specs[0].Scenarios), Equals, 2)
}

// FilterBasedOnSpan

func (s *MySuite) TestScenarioSpanFilter(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}
	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 7, End: 10},
	}
	scenario4 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Fourth Scenario"},
		Span:    &gauge.Span{Start: 11, End: 15},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3, scenario4},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3, scenario4},
	}

	specWithFilteredItems, specWithOtherItems := spec.Filter(NewScenarioFilterBasedOnSpan([]int{8}))

	c.Assert(len(specWithFilteredItems.Scenarios), Equals, 1)
	c.Assert(specWithFilteredItems.Scenarios[0], Equals, scenario3)

	c.Assert(len(specWithOtherItems.Scenarios), Equals, 3)
	c.Assert(specWithOtherItems.Scenarios[0], Equals, scenario1)
	c.Assert(specWithOtherItems.Scenarios[1], Equals, scenario2)
	c.Assert(specWithOtherItems.Scenarios[2], Equals, scenario4)
}

func (s *MySuite) TestScenarioSpanFilterLastScenario(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}
	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 7, End: 10},
	}
	scenario4 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Fourth Scenario"},
		Span:    &gauge.Span{Start: 11, End: 15},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3, scenario4},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3, scenario4},
	}

	specWithFilteredItems, specWithOtherItems := spec.Filter(NewScenarioFilterBasedOnSpan([]int{13}))
	c.Assert(len(specWithFilteredItems.Scenarios), Equals, 1)
	c.Assert(specWithFilteredItems.Scenarios[0], Equals, scenario4)

	c.Assert(len(specWithOtherItems.Scenarios), Equals, 3)
	c.Assert(specWithOtherItems.Scenarios[0], Equals, scenario1)
	c.Assert(specWithOtherItems.Scenarios[1], Equals, scenario2)
	c.Assert(specWithOtherItems.Scenarios[2], Equals, scenario3)

}

func (s *MySuite) TestScenarioSpanFilterFirstScenario(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}
	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 7, End: 10},
	}
	scenario4 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Fourth Scenario"},
		Span:    &gauge.Span{Start: 11, End: 15},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3, scenario4},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3, scenario4},
	}

	specWithFilteredItems, specWithOtherItems := spec.Filter(NewScenarioFilterBasedOnSpan([]int{2}))

	c.Assert(len(specWithFilteredItems.Scenarios), Equals, 1)
	c.Assert(specWithFilteredItems.Scenarios[0], Equals, scenario1)

	c.Assert(len(specWithOtherItems.Scenarios), Equals, 3)
	c.Assert(specWithOtherItems.Scenarios[0], Equals, scenario2)
	c.Assert(specWithOtherItems.Scenarios[1], Equals, scenario3)
	c.Assert(specWithOtherItems.Scenarios[2], Equals, scenario4)

}

func (s *MySuite) TestScenarioSpanFilterForSingleScenarioSpec(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1},
		Scenarios: []*gauge.Scenario{scenario1},
	}

	specWithFilteredItems, specWithOtherItems := spec.Filter(NewScenarioFilterBasedOnSpan([]int{3}))
	c.Assert(len(specWithFilteredItems.Scenarios), Equals, 1)
	c.Assert(specWithFilteredItems.Scenarios[0], Equals, scenario1)

	c.Assert(len(specWithOtherItems.Scenarios), Equals, 0)
}

func (s *MySuite) TestScenarioSpanFilterWithWrongScenarioIndex(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1},
		Scenarios: []*gauge.Scenario{scenario1},
	}

	specWithFilteredItems, specWithOtherItems := spec.Filter(NewScenarioFilterBasedOnSpan([]int{5}))
	c.Assert(len(specWithFilteredItems.Scenarios), Equals, 0)

	c.Assert(len(specWithOtherItems.Scenarios), Equals, 1)
	c.Assert(specWithOtherItems.Scenarios[0], Equals, scenario1)
}

func (s *MySuite) TestScenarioSpanFilterWithMultipleLineNumbers(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}
	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 7, End: 10},
	}
	scenario4 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Fourth Scenario"},
		Span:    &gauge.Span{Start: 11, End: 15},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3, scenario4},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3, scenario4},
	}

	specWithFilteredItems, specWithOtherItems := spec.Filter(NewScenarioFilterBasedOnSpan([]int{3, 13}))

	c.Assert(len(specWithFilteredItems.Scenarios), Equals, 2)
	c.Assert(specWithFilteredItems.Scenarios[0], Equals, scenario1)
	c.Assert(specWithFilteredItems.Scenarios[1], Equals, scenario4)

	c.Assert(len(specWithOtherItems.Scenarios), Equals, 2)
	c.Assert(specWithOtherItems.Scenarios[0], Equals, scenario2)
	c.Assert(specWithOtherItems.Scenarios[1], Equals, scenario3)

}

// FilterBasedOnTags

func (s *MySuite) TestToEvaluateTagExpressionWithTwoTags(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "tag1 & tag3"}
	c.Assert(filter.FilterTags([]string{"tag1", "tag2"}), Equals, false)
}

func (s *MySuite) TestToEvaluateTagExpressionWithComplexTagExpression(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "tag1 & ((tag3 | tag2) & (tag5 | tag4 | tag3) & tag7) | tag6"}
	c.Assert(filter.FilterTags([]string{"tag1", "tag2", "tag7", "tag4"}), Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionWithFailingTagExpression(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "tag1 & ((tag3 | tag2) & (tag5 | tag4 | tag3) & tag7) & tag6"}
	c.Assert(filter.FilterTags([]string{"tag1", "tag2", "tag7", "tag4"}), Equals, false)
}

func (s *MySuite) TestToEvaluateTagExpressionWithWrongTagExpression(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "tag1 & (((tag3 | tag2) & (tag5 | tag4 | tag3) & tag7) & tag6)"}
	c.Assert(filter.FilterTags([]string{"tag1", "tag2", "tag7", "tag4"}), Equals, false)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingOfSpaces(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "tag 1 & tag3"}
	c.Assert(filter.FilterTags([]string{"tag 1", "tag3"}), Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingLogicalNotOperator(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "!tag 1 & tag3"}
	c.Assert(filter.FilterTags([]string{"tag2", "tag3"}), Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingManyLogicalNotOperator(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "!(!(tag 1 | !(tag6 | !(tag5))) & tag2)"}
	value := filter.FilterTags([]string{"tag2", "tag4"})
	c.Assert(value, Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingParallelLogicalNotOperator(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "!(tag1) & ! (tag3 & ! (tag3))"}
	value := filter.FilterTags([]string{"tag2", "tag4"})
	c.Assert(value, Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingComma(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "tag 1 , tag3"}
	c.Assert(filter.FilterTags([]string{"tag2", "tag3"}), Equals, false)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingCommaGivesTrue(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "tag 1 , tag3"}
	c.Assert(filter.FilterTags([]string{"tag1", "tag3"}), Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingTrueAndFalseAsTagNames(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "true , false"}
	c.Assert(filter.FilterTags([]string{"true", "false"}), Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingTrueAndFalseAsTagNamesWithNegation(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "!true"}
	c.Assert(filter.FilterTags(nil), Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionConsistingSpecialCharacters(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "a && b || c | b & b"}
	c.Assert(filter.FilterTags([]string{"a", "b"}), Equals, true)
}

func (s *MySuite) TestToEvaluateTagExpressionWhenTagIsSubsetOfTrueOrFalse(c *C) {
	// https://github.com/getgauge/gauge/issues/667
	filter := &ScenarioFilterBasedOnTags{expression: "b || c | b & b && a"}
	c.Assert(filter.FilterTags([]string{"a", "b"}), Equals, true)
}

func (s *MySuite) TestParseTagExpression(c *C) {
	filter := &ScenarioFilterBasedOnTags{expression: "b || c | b & b && a"}
	txps, tags := filter.parseTagExpression()

	expectedTxps := []string{"b", "|", "|", "c", "|", "b", "&", "b", "&", "&", "a"}
	expectedTags := []string{"b", "c", "b", "b", "a"}

	for i, t := range txps {
		c.Assert(expectedTxps[i], Equals, t)
	}
	for i, t := range tags {
		c.Assert(expectedTags[i], Equals, t)
	}
}

func (s *MySuite) TestToFilterSpecsByWrongTagExpression(c *C) {
	tagsMap := map[string]bool{"{tag1}": true, "{tag2}": true}
	filter := &ScenarioFilterBasedOnTags{expression: "(tag1 & tag2"}
	filter.replaceSpecialChar()
	_, err := filter.formatAndEvaluateExpression(tagsMap, filter.isTagPresent)

	c.Assert(err.Error(), Equals, "invalid expression: `(tag1 & tag2` \neval:1:5: expected ')', found newline")
}

func (s *MySuite) TestFilterTags(c *C) {
	specTags := []string{"abcd", "foo", "bar", "foo bar"}
	spec := &gauge.Specification{
		Tags:             &gauge.Tags{RawValues: [][]string{specTags}},
		FilterExpression: "",
	}
	tagFilter := NewScenarioFilterBasedOnTags(spec, "abcd & foo bar")
	evaluateTrue := tagFilter.FilterTags(specTags)
	c.Assert(evaluateTrue, Equals, true)
}

func (s *MySuite) TestSanitizeTags(c *C) {
	specTags := []string{"abcd", "foo", "bar", "foo bar"}
	spec := &gauge.Specification{
		Tags:             &gauge.Tags{RawValues: [][]string{specTags}},
		FilterExpression: "",
	}
	tagFilter := NewScenarioFilterBasedOnTags(spec, "abcd & foo bar | true")
	evaluateTrue := tagFilter.FilterTags(specTags)
	c.Assert(evaluateTrue, Equals, true)
}

func (s *MySuite) TestToFilterSpecsByTagExpOfTwoTags(c *C) {
	myTags := []string{"tag1", "tag2"}
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}
	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
		Tags:      &gauge.Tags{RawValues: [][]string{myTags}},
	}

	spec2 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1, spec2)

	c.Assert(specs[0].Tags.Values()[0], Equals, myTags[0])
	c.Assert(specs[0].Tags.Values()[1], Equals, myTags[1])

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & tag2")
	c.Assert(len(filteredSpecs), Equals, 1)

	c.Assert(len(otherSpecs), Equals, 1)
}

func (s *MySuite) TestToEvaluateTagExpression(c *C) {
	myTags := []string{"tag1", "tag2"}

	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
		Tags:    &gauge.Tags{RawValues: [][]string{{myTags[0]}}},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{{"tag3"}}},
	}

	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario4 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Fourth Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}

	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
	}

	spec2 := &gauge.Specification{
		Items:     []gauge.Item{scenario3, scenario4},
		Scenarios: []*gauge.Scenario{scenario3, scenario4},
		Tags:      &gauge.Tags{RawValues: [][]string{myTags}},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1, spec2)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & !(tag1 & tag4) & (tag2 | tag3)")
	c.Assert(len(filteredSpecs), Equals, 1)
	c.Assert(len(filteredSpecs[0].Scenarios), Equals, 2)
	c.Assert(filteredSpecs[0].Scenarios[0], Equals, scenario3)
	c.Assert(filteredSpecs[0].Scenarios[1], Equals, scenario4)

	c.Assert(len(otherSpecs), Equals, 1)
	c.Assert(len(otherSpecs[0].Scenarios), Equals, 2)
	c.Assert(otherSpecs[0].Scenarios[0], Equals, scenario1)
	c.Assert(otherSpecs[0].Scenarios[1], Equals, scenario2)
}

func (s *MySuite) TestToFilterMultipleScenariosByMultipleTags(c *C) {
	myTags := []string{"tag1", "tag2"}
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
		Tags:    &gauge.Tags{RawValues: [][]string{{"tag1"}}},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}

	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}
	scenario4 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Fourth Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{{"prod", "tag7", "tag1", "tag2"}}},
	}
	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3, scenario4},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3, scenario4},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	c.Assert(len(specs[0].Scenarios), Equals, 4)
	c.Assert(len(specs[0].Scenarios[0].Tags.Values()), Equals, 1)
	c.Assert(len(specs[0].Scenarios[1].Tags.Values()), Equals, 2)
	c.Assert(len(specs[0].Scenarios[2].Tags.Values()), Equals, 2)
	c.Assert(len(specs[0].Scenarios[3].Tags.Values()), Equals, 4)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & tag2")
	c.Assert(len(filteredSpecs[0].Scenarios), Equals, 3)
	c.Assert(filteredSpecs[0].Scenarios[0], Equals, scenario2)
	c.Assert(filteredSpecs[0].Scenarios[1], Equals, scenario3)
	c.Assert(filteredSpecs[0].Scenarios[2], Equals, scenario4)

	c.Assert(len(otherSpecs[0].Scenarios), Equals, 1)
	c.Assert(otherSpecs[0].Scenarios[0], Equals, scenario1)

}

func (s *MySuite) TestToFilterScenariosByTagsAtSpecLevel(c *C) {
	myTags := []string{"tag1", "tag2"}

	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}

	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}

	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3},
		Tags:      &gauge.Tags{RawValues: [][]string{myTags}},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	c.Assert(len(specs[0].Scenarios), Equals, 3)
	c.Assert(len(specs[0].Tags.Values()), Equals, 2)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & tag2")
	c.Assert(len(filteredSpecs[0].Scenarios), Equals, 3)
	c.Assert(filteredSpecs[0].Scenarios[0], Equals, scenario1)
	c.Assert(filteredSpecs[0].Scenarios[1], Equals, scenario2)
	c.Assert(filteredSpecs[0].Scenarios[2], Equals, scenario3)

	c.Assert(len(otherSpecs), Equals, 0)

}

func (s *MySuite) TestToFilterScenariosByTagExpWithDuplicateTagNames(c *C) {
	myTags := []string{"tag1", "tag12"}
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{{"tag1"}}},
	}

	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{{"tag12"}}},
	}

	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1)
	c.Assert(len(specs), Equals, 1)

	c.Assert(len(specs[0].Scenarios), Equals, 3)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & tag12")
	c.Assert(len(filteredSpecs[0].Scenarios), Equals, 1)
	c.Assert(filteredSpecs[0].Scenarios[0], Equals, scenario1)

	c.Assert(len(otherSpecs), Equals, 1)
	c.Assert(otherSpecs[0].Scenarios[0], Equals, scenario2)
	c.Assert(otherSpecs[0].Scenarios[1], Equals, scenario3)
}

func (s *MySuite) TestToFilterSpecsByTags(c *C) {
	myTags := []string{"tag1", "tag2"}
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}
	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}

	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2},
		Scenarios: []*gauge.Scenario{scenario1, scenario2},
	}
	spec2 := &gauge.Specification{
		Items:     []gauge.Item{scenario2, scenario3},
		Scenarios: []*gauge.Scenario{scenario2, scenario3},
	}

	spec3 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario3},
		Scenarios: []*gauge.Scenario{scenario1, scenario3},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1, spec2, spec3)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & tag2")
	c.Assert(len(filteredSpecs), Equals, 2)
	c.Assert(len(filteredSpecs[0].Scenarios), Equals, 1)
	c.Assert(len(filteredSpecs[1].Scenarios), Equals, 1)
	c.Assert(filteredSpecs[0].Scenarios[0], Equals, scenario1)
	c.Assert(filteredSpecs[1].Scenarios[0], Equals, scenario1)

	c.Assert(len(otherSpecs), Equals, 3)
	c.Assert(len(otherSpecs[0].Scenarios), Equals, 1)
	c.Assert(len(otherSpecs[1].Scenarios), Equals, 2)
	c.Assert(len(otherSpecs[2].Scenarios), Equals, 1)
	c.Assert(otherSpecs[0].Scenarios[0], Equals, scenario2)
	c.Assert(otherSpecs[1].Scenarios[0], Equals, scenario2)
	c.Assert(otherSpecs[1].Scenarios[1], Equals, scenario3)
	c.Assert(otherSpecs[2].Scenarios[0], Equals, scenario3)

}

func (s *MySuite) TestToFilterScenariosByTag(c *C) {
	myTags := []string{"tag1", "tag2"}

	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}

	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}

	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & tag2")
	c.Assert(len(filteredSpecs[0].Scenarios), Equals, 1)
	c.Assert(filteredSpecs[0].Scenarios[0], Equals, scenario2)

	c.Assert(len(otherSpecs), Equals, 1)
	c.Assert(len(otherSpecs[0].Scenarios), Equals, 2)
	c.Assert(otherSpecs[0].Scenarios[0], Equals, scenario1)
	c.Assert(otherSpecs[0].Scenarios[1], Equals, scenario3)
}

func (s *MySuite) TestToFilterMultipleScenariosByTags(c *C) {
	myTags := []string{"tag1", "tag2"}

	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
		Tags:    &gauge.Tags{RawValues: [][]string{{"tag1"}}},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}

	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}

	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag1 & tag2")

	c.Assert(len(filteredSpecs[0].Scenarios), Equals, 2)
	c.Assert(filteredSpecs[0].Scenarios[0], Equals, scenario2)
	c.Assert(filteredSpecs[0].Scenarios[1], Equals, scenario3)

	c.Assert(len(otherSpecs), Equals, 1)
	c.Assert(len(otherSpecs[0].Scenarios), Equals, 1)
	c.Assert(otherSpecs[0].Scenarios[0], Equals, scenario1)
}

func (s *MySuite) TestToFilterScenariosByUnavailableTags(c *C) {
	myTags := []string{"tag1", "tag2"}
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{myTags}},
	}

	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
	}

	spec1 := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3},
	}

	var specs []*gauge.Specification
	specs = append(specs, spec1)

	filteredSpecs, otherSpecs := filterByTags(specs, "tag3")

	c.Assert(len(filteredSpecs), Equals, 0)

	c.Assert(len(otherSpecs), Equals, 1)
	c.Assert(len(otherSpecs[0].Scenarios), Equals, 3)
	c.Assert(otherSpecs[0].Scenarios[0], Equals, scenario1)
	c.Assert(otherSpecs[0].Scenarios[1], Equals, scenario2)
	c.Assert(otherSpecs[0].Scenarios[2], Equals, scenario3)
}

func (s *MySuite) TestScenarioTagFilterShouldNotRemoveNonScenarioKindItems(c *C) {
	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Second Scenario"},
		Span:    &gauge.Span{Start: 4, End: 6},
		Tags:    &gauge.Tags{RawValues: [][]string{{"tag2"}}},
	}
	scenario3 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "Third Scenario"},
		Span:    &gauge.Span{Start: 7, End: 10},
		Tags:    &gauge.Tags{RawValues: [][]string{{"tag1"}}},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1, scenario2, scenario3, &gauge.Table{}, &gauge.Comment{Value: "Comment", LineNo: 1}, &gauge.Step{}},
		Scenarios: []*gauge.Scenario{scenario1, scenario2, scenario3},
	}

	specWithFilteredItems, specWithOtherItems := filterByTags([]*gauge.Specification{spec}, "tag1 | tag2")

	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithFilteredItems[0].Items), Equals, 5)

	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithOtherItems[0].Items), Equals, 4)
}

func (s *MySuite) TestToFilterScenariosBySuiteTags(c *C) {
	os.Clearenv()
	os.Setenv(env.GaugeSuiteTags, "suite_tag1")

	scenario1 := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "First Scenario"},
		Span:    &gauge.Span{Start: 1, End: 3},
	}
	spec := &gauge.Specification{
		Items:     []gauge.Item{scenario1, &gauge.Table{}, &gauge.Comment{Value: "Comment", LineNo: 1}, &gauge.Step{}},
		Scenarios: []*gauge.Scenario{scenario1},
	}

	specWithFilteredItems, specWithOtherItems := filterByTags([]*gauge.Specification{spec}, "!suite_tag1")
	c.Assert(len(specWithFilteredItems), Equals, 0)
	c.Assert(len(specWithOtherItems), Equals, 1)
}

func (s *MySuite) TestToFilterScenariosSpecTableByTags(c *C) {
	headers := []string{"app", "tags"}
	cell1 := gauge.TableCell{Value: "app1", CellType: gauge.Static}
	cell2 := gauge.TableCell{Value: "app2", CellType: gauge.Static}
	cell3 := gauge.TableCell{Value: "tag1, tag2", CellType: gauge.Static}
	cell4 := gauge.TableCell{Value: "tag3, tag4", CellType: gauge.Static}

	cols1 := [][]gauge.TableCell{{cell1}, {cell3}}
	cols2 := [][]gauge.TableCell{{cell2}, {cell4}}

	table1 := gauge.NewTable(headers, cols1, 1)
	table2 := gauge.NewTable(headers, cols2, 1)

	scenario11 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{}},
		FilterExpression:     "",
		SpecDataTableRow:     *table1,
		ScenarioDataTableRow: gauge.Table{},
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario12 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{}},
		FilterExpression:     "",
		SpecDataTableRow:     *table2,
		ScenarioDataTableRow: gauge.Table{},
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario21 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{}},
		FilterExpression:     "",
		SpecDataTableRow:     *table1,
		ScenarioDataTableRow: gauge.Table{},
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	scenario22 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{}},
		FilterExpression:     "",
		SpecDataTableRow:     *table2,
		ScenarioDataTableRow: gauge.Table{},
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	spec1 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec"},
		Scenarios:        []*gauge.Scenario{scenario11, scenario21},
		DataTable:        gauge.DataTable{Table: table1},
		FileName:         "/home/user/gauge/specs/test.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{{"tag5", "tag6"}}},
		FilterExpression: "",
		Items:            []gauge.Item{table1, scenario11, scenario21, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}
	spec2 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec"},
		Scenarios:        []*gauge.Scenario{scenario12, scenario22},
		DataTable:        gauge.DataTable{Table: table2},
		FileName:         "/home/user/gauge/specs/test.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{{"tag5", "tag6"}}},
		FilterExpression: "",
		Items:            []gauge.Item{table2, scenario12, scenario22, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}

	specWithFilteredItems, specWithOtherItems := filterByTags([]*gauge.Specification{spec1, spec2}, "tag1 && tag2")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario21})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario12, scenario22})

	c.Assert(specWithFilteredItems[0].Scenarios[0].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].SpecDataTableRowIndex, Equals, 0)

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec1, spec2}, "tag7")
	c.Assert(len(specWithFilteredItems), Equals, 0)
	c.Assert(len(specWithOtherItems), Equals, 2)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[1].Scenarios), Equals, 2)

	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario21})
	c.Assert(specWithOtherItems[1].Scenarios, DeepEquals, []*gauge.Scenario{scenario12, scenario22})

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec1, spec2}, "tag1 || tag3")
	c.Assert(len(specWithFilteredItems), Equals, 2)
	c.Assert(len(specWithOtherItems), Equals, 0)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithFilteredItems[1].Scenarios), Equals, 2)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario21})
	c.Assert(specWithFilteredItems[1].Scenarios, DeepEquals, []*gauge.Scenario{scenario12, scenario22})

	c.Assert(specWithFilteredItems[0].Scenarios[0].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[1].Scenarios[0].SpecDataTableRowIndex, Equals, 1)
	c.Assert(specWithFilteredItems[1].Scenarios[1].SpecDataTableRowIndex, Equals, 1)

	spec3 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec"},
		Scenarios:        []*gauge.Scenario{scenario11, scenario21},
		DataTable:        gauge.DataTable{Table: table1},
		FileName:         "/home/user/gauge/specs/test.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{{"tag5", "tag6"}}},
		FilterExpression: "!tag4",
		Items:            []gauge.Item{table1, scenario11, scenario21, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}
	spec4 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec"},
		Scenarios:        []*gauge.Scenario{scenario12, scenario22},
		DataTable:        gauge.DataTable{Table: table2},
		FileName:         "/home/user/gauge/specs/test.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{{"tag5", "tag6"}}},
		FilterExpression: "!tag4",
		Items:            []gauge.Item{table2, scenario12, scenario22, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec3, spec4}, "")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario21})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario12, scenario22})

	c.Assert(specWithFilteredItems[0].Scenarios[0].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].SpecDataTableRowIndex, Equals, 0)
}

func (s *MySuite) TestToFilterScenariosTableByTags(c *C) {
	headers := []string{"app", "tags"}
	cell1 := gauge.TableCell{Value: "app1", CellType: gauge.Static}
	cell2 := gauge.TableCell{Value: "app2", CellType: gauge.Static}
	cell3 := gauge.TableCell{Value: "tag1, tag2", CellType: gauge.Static}
	cell4 := gauge.TableCell{Value: "tag3, tag4", CellType: gauge.Static}

	cell5 := gauge.TableCell{Value: "app3", CellType: gauge.Static}
	cell6 := gauge.TableCell{Value: "app4", CellType: gauge.Static}
	cell7 := gauge.TableCell{Value: "tag5, tag6", CellType: gauge.Static}
	cell8 := gauge.TableCell{Value: "tag7, tag8", CellType: gauge.Static}

	cols1 := [][]gauge.TableCell{{cell1}, {cell3}}
	cols2 := [][]gauge.TableCell{{cell2}, {cell4}}

	cols3 := [][]gauge.TableCell{{cell5}, {cell7}}
	cols4 := [][]gauge.TableCell{{cell6}, {cell8}}

	table1 := gauge.NewTable(headers, cols1, 1)
	table2 := gauge.NewTable(headers, cols2, 1)

	table3 := gauge.NewTable(headers, cols3, 1)
	table4 := gauge.NewTable(headers, cols4, 1)

	scenario11 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag9"}}},
		FilterExpression:     "!tag11",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: *table1,
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario12 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag9"}}},
		FilterExpression:     "!tag11",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: *table2,
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario21 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag10"}}},
		FilterExpression:     "!tag12",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: *table3,
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	scenario22 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag10"}}},
		FilterExpression:     "!tag12",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: *table4,
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	spec1 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec1"},
		Scenarios:        []*gauge.Scenario{scenario11, scenario12, scenario21, scenario22},
		DataTable:        gauge.DataTable{Table: &gauge.Table{}},
		FileName:         "/home/user/gauge/specs/test1.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{}},
		FilterExpression: "",
		Items:            []gauge.Item{scenario11, scenario12, scenario21, scenario22, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}

	specWithFilteredItems, specWithOtherItems := filterByTags([]*gauge.Specification{spec1}, "tag1 || tag3")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario12})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario21, scenario22})

	c.Assert(specWithFilteredItems[0].Scenarios[0].ScenarioDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].ScenarioDataTableRowIndex, Equals, 1)

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec1}, "tag1 && tag3 && tag5 && tag7")
	c.Assert(len(specWithFilteredItems), Equals, 0)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 4)

	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario12, scenario21, scenario22})

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec1}, "tag9 || tag10")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 0)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 4)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario12, scenario21, scenario22})

	c.Assert(specWithFilteredItems[0].Scenarios[0].ScenarioDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].ScenarioDataTableRowIndex, Equals, 1)
	c.Assert(specWithFilteredItems[0].Scenarios[2].ScenarioDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[3].ScenarioDataTableRowIndex, Equals, 1)

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec1}, "tag2 || tag6")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario11, scenario21})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario12, scenario22})

	c.Assert(specWithFilteredItems[0].Scenarios[0].ScenarioDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].ScenarioDataTableRowIndex, Equals, 0)

	scenario31 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Third Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag20"}}},
		FilterExpression:     "!tag1",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: *table1,
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario41 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Fourth Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag21"}}},
		FilterExpression:     "tag6 && tag21",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: *table3,
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	spec2 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec2"},
		Scenarios:        []*gauge.Scenario{scenario31, scenario41},
		DataTable:        gauge.DataTable{Table: &gauge.Table{}},
		FileName:         "/home/user/gauge/specs/test2.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{}},
		FilterExpression: "",
		Items:            []gauge.Item{scenario31, scenario41, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec2}, "")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 1)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 1)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario41})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario31})

	c.Assert(specWithFilteredItems[0].Scenarios[0].ScenarioDataTableRowIndex, Equals, 0)
}

func (s *MySuite) TestToFilterSimpleScenariosByTags(c *C) {
	scenario1 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag1"}}},
		FilterExpression:     "!tag7",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: gauge.Table{},
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario2 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag2"}}},
		FilterExpression:     "!tag4",
		SpecDataTableRow:     gauge.Table{},
		ScenarioDataTableRow: gauge.Table{},
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	spec1 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec1"},
		Scenarios:        []*gauge.Scenario{scenario1, scenario2},
		DataTable:        gauge.DataTable{Table: &gauge.Table{}},
		FileName:         "/home/user/gauge/specs/test1.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{{"tag5", "tag6"}}},
		FilterExpression: "tag2",
		Items:            []gauge.Item{scenario1, scenario2, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}
	spec2 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec2"},
		Scenarios:        []*gauge.Scenario{scenario1, scenario2},
		DataTable:        gauge.DataTable{Table: &gauge.Table{}},
		FileName:         "/home/user/gauge/specs/test2.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{{"tag7", "tag8"}}},
		FilterExpression: "tag1",
		Items:            []gauge.Item{scenario1, scenario2, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}

	specWithFilteredItems, specWithOtherItems := filterByTags([]*gauge.Specification{spec1}, "")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 1)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 1)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario2})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario1})

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec2}, "")
	c.Assert(len(specWithFilteredItems), Equals, 0)
	c.Assert(len(specWithOtherItems), Equals, 1)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)

	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario1, scenario2})
}

func (s *MySuite) TestToFilterTableRelatedAndTableDrivenScenariosByTags(c *C) {
	headers1 := []string{"data", "tags"}
	cell11 := gauge.TableCell{Value: "data1", CellType: gauge.Static}
	cell12 := gauge.TableCell{Value: "data2", CellType: gauge.Static}
	cell13 := gauge.TableCell{Value: "tag100, tag200", CellType: gauge.Static}
	cell14 := gauge.TableCell{Value: "tag300, tag400", CellType: gauge.Static}

	cols11 := [][]gauge.TableCell{{cell11}, {cell13}}
	cols12 := [][]gauge.TableCell{{cell12}, {cell14}}

	table11 := gauge.NewTable(headers1, cols11, 1)
	table12 := gauge.NewTable(headers1, cols12, 1)

	headers2 := []string{"app", "tags"}
	cell21 := gauge.TableCell{Value: "app1", CellType: gauge.Static}
	cell22 := gauge.TableCell{Value: "app2", CellType: gauge.Static}
	cell23 := gauge.TableCell{Value: "tag1, tag2", CellType: gauge.Static}
	cell24 := gauge.TableCell{Value: "tag3, tag4", CellType: gauge.Static}

	cell25 := gauge.TableCell{Value: "app3", CellType: gauge.Static}
	cell26 := gauge.TableCell{Value: "app4", CellType: gauge.Static}
	cell27 := gauge.TableCell{Value: "tag5, tag6", CellType: gauge.Static}
	cell28 := gauge.TableCell{Value: "tag7, tag8", CellType: gauge.Static}

	cols21 := [][]gauge.TableCell{{cell21}, {cell23}}
	cols22 := [][]gauge.TableCell{{cell22}, {cell24}}

	cols23 := [][]gauge.TableCell{{cell25}, {cell27}}
	cols24 := [][]gauge.TableCell{{cell26}, {cell28}}

	table21 := gauge.NewTable(headers2, cols21, 1)
	table22 := gauge.NewTable(headers2, cols22, 1)

	table23 := gauge.NewTable(headers2, cols23, 1)
	table24 := gauge.NewTable(headers2, cols24, 1)

	scenario111 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag9"}}},
		FilterExpression:     "tag100 || tag300",
		SpecDataTableRow:     *table11,
		ScenarioDataTableRow: *table21,
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario112 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag9"}}},
		FilterExpression:     "tag100 || tag300",
		SpecDataTableRow:     *table11,
		ScenarioDataTableRow: *table22,
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario121 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag10"}}},
		FilterExpression:     "tag200 || tag400",
		SpecDataTableRow:     *table11,
		ScenarioDataTableRow: *table23,
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	scenario122 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag10"}}},
		FilterExpression:     "tag200 || tag400",
		SpecDataTableRow:     *table11,
		ScenarioDataTableRow: *table24,
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	scenario211 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag9"}}},
		FilterExpression:     "tag100 || tag300",
		SpecDataTableRow:     *table12,
		ScenarioDataTableRow: *table21,
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario212 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "First Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag9"}}},
		FilterExpression:     "tag100 || tag300",
		SpecDataTableRow:     *table12,
		ScenarioDataTableRow: *table22,
		Span:                 &gauge.Span{Start: 1, End: 3},
	}
	scenario221 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag10"}}},
		FilterExpression:     "tag200 || tag400",
		SpecDataTableRow:     *table12,
		ScenarioDataTableRow: *table23,
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	scenario222 := &gauge.Scenario{
		Heading:              &gauge.Heading{Value: "Second Scenario"},
		Steps:                []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
		Tags:                 &gauge.Tags{RawValues: [][]string{{"tag10"}}},
		FilterExpression:     "tag200 || tag400",
		SpecDataTableRow:     *table12,
		ScenarioDataTableRow: *table24,
		Span:                 &gauge.Span{Start: 4, End: 6},
	}
	spec1 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec"},
		Scenarios:        []*gauge.Scenario{scenario111, scenario112, scenario121, scenario122},
		DataTable:        gauge.DataTable{Table: table11},
		FileName:         "/home/user/gauge/specs/test.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{}},
		FilterExpression: "tag3 || tag7",
		Items:            []gauge.Item{table11, scenario111, scenario112, scenario121, scenario122, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}
	spec2 := &gauge.Specification{
		Heading:          &gauge.Heading{Value: "Spec"},
		Scenarios:        []*gauge.Scenario{scenario211, scenario212, scenario221, scenario222},
		DataTable:        gauge.DataTable{Table: table12},
		FileName:         "/home/user/gauge/specs/test.spec",
		Tags:             &gauge.Tags{RawValues: [][]string{}},
		FilterExpression: "tag3 || tag7",
		Items:            []gauge.Item{table12, scenario211, scenario212, scenario221, scenario222, &gauge.Step{}},
		TearDownSteps:    []*gauge.Step{{Args: []*gauge.StepArg{{Value: "", ArgType: gauge.Static}}}},
	}

	specWithFilteredItems, specWithOtherItems := filterByTags([]*gauge.Specification{spec1, spec2}, "tag4 || tag8")
	c.Assert(len(specWithFilteredItems), Equals, 2)
	c.Assert(len(specWithOtherItems), Equals, 2)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithFilteredItems[1].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[1].Scenarios), Equals, 2)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario112, scenario122})
	c.Assert(specWithFilteredItems[1].Scenarios, DeepEquals, []*gauge.Scenario{scenario212, scenario222})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario111, scenario121})
	c.Assert(specWithOtherItems[1].Scenarios, DeepEquals, []*gauge.Scenario{scenario211, scenario221})

	c.Assert(specWithFilteredItems[0].Scenarios[0].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[0].ScenarioDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].ScenarioDataTableRowIndex, Equals, 0)

	c.Assert(specWithFilteredItems[1].Scenarios[0].SpecDataTableRowIndex, Equals, 1)
	c.Assert(specWithFilteredItems[1].Scenarios[0].ScenarioDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[1].Scenarios[1].SpecDataTableRowIndex, Equals, 1)
	c.Assert(specWithFilteredItems[1].Scenarios[1].ScenarioDataTableRowIndex, Equals, 0)

	specWithFilteredItems, specWithOtherItems = filterByTags([]*gauge.Specification{spec1, spec2}, "tag100")
	c.Assert(len(specWithFilteredItems), Equals, 1)
	c.Assert(len(specWithOtherItems), Equals, 2)
	c.Assert(len(specWithFilteredItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[0].Scenarios), Equals, 2)
	c.Assert(len(specWithOtherItems[1].Scenarios), Equals, 4)

	c.Assert(specWithFilteredItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario112, scenario122})
	c.Assert(specWithOtherItems[0].Scenarios, DeepEquals, []*gauge.Scenario{scenario111, scenario121})
	c.Assert(specWithOtherItems[1].Scenarios, DeepEquals, []*gauge.Scenario{scenario211, scenario212, scenario221, scenario222})

	c.Assert(specWithFilteredItems[0].Scenarios[0].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[0].ScenarioDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].SpecDataTableRowIndex, Equals, 0)
	c.Assert(specWithFilteredItems[0].Scenarios[1].ScenarioDataTableRowIndex, Equals, 0)
}
