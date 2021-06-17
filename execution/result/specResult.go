/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/

package result

import (
	"github.com/getgauge/gauge-proto/go/gauge_messages"
)

// SpecResult represents the result of spec execution
type SpecResult struct {
	ProtoSpec            *gauge_messages.ProtoSpec
	ScenarioFailedCount  int
	ScenarioCount        int
	IsFailed             bool
	FailedDataTableRows  []int32
	SkippedDataTableRows []int32
	ExecutionTime        int64
	Skipped              bool
	ScenarioSkippedCount int
	Errors               []*gauge_messages.Error
}

// SetFailure sets the result to failed
func (specResult *SpecResult) SetFailure() {
	specResult.IsFailed = true
}

func (specResult *SpecResult) SetSkipped(skipped bool) {
	specResult.Skipped = skipped
}

func (specResult *SpecResult) AddSpecItems(resolvedItems []*gauge_messages.ProtoItem) {
	specResult.ProtoSpec.Items = append(specResult.ProtoSpec.Items, resolvedItems...)
}

// AddScenarioResults adds the result of each scenario to spec result.
func (specResult *SpecResult) AddScenarioResults(scenarioResults []*ScenarioResult) {
	tableDrivenScenariosMap := make(map[int64]bool)

	for _, scenarioResult := range scenarioResults {
		specResult.AddExecTime(scenarioResult.ExecTime())
		item := scenarioResult.ConvertToProtoItem()

		isScenarioTableRelated := false
		isScenarioTableDriven := false
		if item.GetItemType() == gauge_messages.ProtoItem_TableDrivenScenario {
			if item.GetTableDrivenScenario().IsSpecTableDriven {
				isScenarioTableRelated = true
			}
			isScenarioTableDriven = item.GetTableDrivenScenario().IsScenarioTableDriven
		}
		if isScenarioTableDriven {
			if _, ok := tableDrivenScenariosMap[scenarioResult.SpanStart()]; !ok {
				tableDrivenScenariosMap[scenarioResult.SpanStart()] = true
				specResult.ScenarioCount++
			}
		} else {
			specResult.ScenarioCount++
		}
		if scenarioResult.GetFailed() {
			specResult.IsFailed = true
			specResult.ScenarioFailedCount++
			if isScenarioTableRelated {
				specResult.FailedDataTableRows = append(specResult.FailedDataTableRows, int32(scenarioResult.SpecDataTableRowIndex))
			}
		} else if scenarioResult.GetSkipped() {
			specResult.ScenarioSkippedCount++
			if isScenarioTableRelated {
				specResult.SkippedDataTableRows = append(specResult.SkippedDataTableRows, int32(scenarioResult.SpecDataTableRowIndex))
			}
		}
		specResult.ProtoSpec.Items = append(specResult.ProtoSpec.Items, item)
	}
}

func (specResult *SpecResult) AddExecTime(execTime int64) {
	specResult.ExecutionTime += execTime
}

func (specResult *SpecResult) GetPreHook() []*gauge_messages.ProtoHookFailure {
	return specResult.ProtoSpec.PreHookFailures
}

func (specResult *SpecResult) GetPostHook() []*gauge_messages.ProtoHookFailure {
	return specResult.ProtoSpec.PostHookFailures
}

func (specResult *SpecResult) AddPreHook(f ...*gauge_messages.ProtoHookFailure) {
	specResult.ProtoSpec.PreHookFailures = append(specResult.ProtoSpec.PreHookFailures, f...)
}

func (specResult *SpecResult) AddPostHook(f ...*gauge_messages.ProtoHookFailure) {
	specResult.ProtoSpec.PostHookFailures = append(specResult.ProtoSpec.PostHookFailures, f...)
}

func (specResult *SpecResult) ExecTime() int64 {
	return specResult.ExecutionTime
}

// GetFailed returns the state of the result
func (specResult *SpecResult) GetFailed() bool {
	return specResult.IsFailed
}

func (specResult *SpecResult) Item() interface{} {
	return specResult.ProtoSpec
}
