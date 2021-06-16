/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/

package execution

import (
	"testing"

	"github.com/getgauge/gauge/execution/result"
	"github.com/getgauge/gauge/gauge"

	"github.com/getgauge/gauge-proto/go/gauge_messages"
)

func TestNotifyBeforeScenarioShouldAddBeforeScenarioHookMessages(t *testing.T) {
	r := &mockRunner{}
	h := &mockPluginHandler{NotifyPluginsfunc: func(m *gauge_messages.Message) {}, GracefullyKillPluginsfunc: func() {}}
	r.ExecuteAndGetStatusFunc = func(m *gauge_messages.Message) *gauge_messages.ProtoExecutionResult {
		if m.MessageType == gauge_messages.Message_ScenarioExecutionStarting {
			return &gauge_messages.ProtoExecutionResult{
				Message:       []string{"Before Scenario Called"},
				Failed:        false,
				ExecutionTime: 10,
			}
		}
		return &gauge_messages.ProtoExecutionResult{}
	}
	ei := &gauge_messages.ExecutionInfo{}
	sce := newScenarioExecutor(r, h, ei, nil, nil, nil, 0)
	scenario := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "A scenario"},
		Span:    &gauge.Span{Start: 2, End: 10},
	}
	scenarioResult := result.NewScenarioResult(gauge.NewProtoScenario(scenario))
	sce.notifyBeforeScenarioHook(scenarioResult)
	gotMessages := scenarioResult.ProtoScenario.PreHookMessages

	if len(gotMessages) != 1 {
		t.Errorf("Expected 1 message, got : %d", len(gotMessages))
	}
	if gotMessages[0] != "Before Scenario Called" {
		t.Errorf("Expected `Before Scenario Called` message, got : %s", gotMessages[0])
	}
}

func TestNotifyAfterScenarioShouldAddAfterScenarioHookMessages(t *testing.T) {
	r := &mockRunner{}
	h := &mockPluginHandler{NotifyPluginsfunc: func(m *gauge_messages.Message) {}, GracefullyKillPluginsfunc: func() {}}
	r.ExecuteAndGetStatusFunc = func(m *gauge_messages.Message) *gauge_messages.ProtoExecutionResult {
		if m.MessageType == gauge_messages.Message_ScenarioExecutionEnding {
			return &gauge_messages.ProtoExecutionResult{
				Message:       []string{"After Scenario Called"},
				Failed:        false,
				ExecutionTime: 10,
			}
		}
		return &gauge_messages.ProtoExecutionResult{}
	}
	ei := &gauge_messages.ExecutionInfo{}
	sce := newScenarioExecutor(r, h, ei, nil, nil, nil, 0)
	scenario := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "A scenario"},
		Span:    &gauge.Span{Start: 2, End: 10},
	}
	scenarioResult := result.NewScenarioResult(gauge.NewProtoScenario(scenario))
	sce.notifyAfterScenarioHook(scenarioResult)
	gotMessages := scenarioResult.ProtoScenario.PostHookMessages

	if len(gotMessages) != 1 {
		t.Errorf("Expected 1 message, got : %d", len(gotMessages))
	}
	if gotMessages[0] != "After Scenario Called" {
		t.Errorf("Expected `After Scenario Called` message, got : %s", gotMessages[0])
	}
}

func TestNotifyBeforeScenarioShouldAddBeforeScenarioHookScreenshots(t *testing.T) {
	r := &mockRunner{}
	h := &mockPluginHandler{NotifyPluginsfunc: func(m *gauge_messages.Message) {}, GracefullyKillPluginsfunc: func() {}}
	r.ExecuteAndGetStatusFunc = func(m *gauge_messages.Message) *gauge_messages.ProtoExecutionResult {
		if m.MessageType == gauge_messages.Message_ScenarioExecutionStarting {
			return &gauge_messages.ProtoExecutionResult{
				ScreenshotFiles: []string{"screenshot1.png", "screenshot2.png"},
				Failed:          false,
				ExecutionTime:   10,
			}
		}
		return &gauge_messages.ProtoExecutionResult{}
	}
	ei := &gauge_messages.ExecutionInfo{}
	sce := newScenarioExecutor(r, h, ei, nil, nil, nil, 0)
	scenario := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "A scenario"},
		Span:    &gauge.Span{Start: 2, End: 10},
	}
	scenarioResult := result.NewScenarioResult(gauge.NewProtoScenario(scenario))
	sce.notifyBeforeScenarioHook(scenarioResult)
	beforeScenarioScreenShots := scenarioResult.ProtoScenario.PreHookScreenshotFiles
	expected := []string{"screenshot1.png", "screenshot2.png"}

	if len(beforeScenarioScreenShots) != len(expected) {
		t.Errorf("Expected 2 screenshots, got : %d", len(beforeScenarioScreenShots))
	}

	for i, e := range expected {
		if string(beforeScenarioScreenShots[i]) != e {
			t.Errorf("Expected `%s` screenshot, got : %s", e, beforeScenarioScreenShots[i])
		}
	}
}

func TestNotifyAfterScenarioShouldAddAfterScenarioHookScreenshots(t *testing.T) {
	r := &mockRunner{}
	h := &mockPluginHandler{NotifyPluginsfunc: func(m *gauge_messages.Message) {}, GracefullyKillPluginsfunc: func() {}}
	r.ExecuteAndGetStatusFunc = func(m *gauge_messages.Message) *gauge_messages.ProtoExecutionResult {
		if m.MessageType == gauge_messages.Message_ScenarioExecutionEnding {
			return &gauge_messages.ProtoExecutionResult{
				ScreenshotFiles: []string{"screenshot1.png", "screenshot2.png"},
				Failed:          false,
				ExecutionTime:   10,
			}
		}
		return &gauge_messages.ProtoExecutionResult{}
	}
	ei := &gauge_messages.ExecutionInfo{}
	sce := newScenarioExecutor(r, h, ei, nil, nil, nil, 0)
	scenario := &gauge.Scenario{
		Heading: &gauge.Heading{Value: "A scenario"},
		Span:    &gauge.Span{Start: 2, End: 10},
	}
	scenarioResult := result.NewScenarioResult(gauge.NewProtoScenario(scenario))
	sce.notifyAfterScenarioHook(scenarioResult)
	afterScenarioScreenShots := scenarioResult.ProtoScenario.PostHookScreenshotFiles
	expected := []string{"screenshot1.png", "screenshot2.png"}

	if len(afterScenarioScreenShots) != len(expected) {
		t.Errorf("Expected 2 screenshots, got : %d", len(afterScenarioScreenShots))
	}

	for i, e := range expected {
		if string(afterScenarioScreenShots[i]) != e {
			t.Errorf("Expected `%s` screenshot, got : %s", e, afterScenarioScreenShots[i])
		}
	}
}

func TestToEvaluateTableFilter(t *testing.T) {
	cell1 := gauge.TableCell{Value: "john", CellType: gauge.Static}
	cell2 := gauge.TableCell{Value: "mike", CellType: gauge.Static}
	cell3 := gauge.TableCell{Value: "tag3, tag2", CellType: gauge.Static}
	cell4 := gauge.TableCell{Value: "tag4, tag10", CellType: gauge.Static}

	headers := []string{"name", "tags"}
	cols1 := [][]gauge.TableCell{{cell1}, {cell3}}
	cols2 := [][]gauge.TableCell{{cell2}, {cell4}}

	table1 := gauge.NewTable(headers, cols1, 1)
	table2 := gauge.NewTable(headers, cols2, 1)

	scenario1 := &gauge.Scenario{
		Heading:             &gauge.Heading{Value: "First Scenario"},
		Span:                &gauge.Span{Start: 1, End: 3},
		Tags:                &gauge.Tags{RawValues: [][]string{}},
		SpecDataTableRow:    *table1,
		SpecDataTableFilter: "tag3 & tag2",
	}
	scenario2 := &gauge.Scenario{
		Heading:             &gauge.Heading{Value: "Second Scenario"},
		Span:                &gauge.Span{Start: 4, End: 6},
		Tags:                &gauge.Tags{RawValues: [][]string{}},
		SpecDataTableRow:    *table2,
		SpecDataTableFilter: "tag4 & tag5",
	}

	if !shouldExecuteForSpecDataTable(scenario1) {
		expr := scenario1.SpecDataTableFilter
		specDataTags := scenario1.SpecDataTableRow.Rows()[0][1]
		t.Errorf("Expected that scenario with spec_table_filter: `%s` should excecute for spec table row with tags: `%s`", expr, specDataTags)
	}
	if shouldExecuteForSpecDataTable(scenario2) {
		expr := scenario2.SpecDataTableFilter
		specDataTags := scenario2.SpecDataTableRow.Rows()[0][1]
		t.Errorf("Expected that scenario with spec_table_filter: `%s` should not excecute for spec table row with tags: `%s`", expr, specDataTags)
	}
}
