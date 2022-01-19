// Package encoding decodes the JSON-encoded workflow,
package encoding

import (
	"encoding/json"
)

type I18nField map[string]string

type Vertex interface {
	ID() string
}

type WatchedDirectory struct {
	ChainID  string `json:"chain_id"`
	OnlyDirs bool   `json:"only_dirs"`
	Path     string `json:"path"`
	UnitType string `json:"unit_type"`
}

func (v WatchedDirectory) ID() string {
	return v.Path
}

var _ Vertex = WatchedDirectory{}

type Chain struct {
	id          string
	Description I18nField `json:"description"`
	LinkID      string    `json:"link_id"`
}

func (v Chain) ID() string {
	return v.id
}

var _ Vertex = Chain{}

type Link struct {
	id                string
	Config            LinkConfig            `json:"config"`
	Description       I18nField             `json:"description"`
	ExitCodes         map[int]*LinkExitCode `json:"exit_codes"`
	FallbackJobStatus string                `json:"fallback_job_status"`
	FallbackLinkID    string                `json:"fallback_link_id"`
	Group             I18nField             `json:"group"`

	// Start and end of the processing workflow.
	Start bool `json:"start"`
	End   bool `json:"end"`
}

func (v Link) ID() string {
	return v.id
}

var _ Vertex = WatchedDirectory{}

type LinkConfig struct {
	Manager string `json:"@manager"`
	Model   string `json:"@model"`

	// StandardTaskConfig
	Arguments          string `json:"arguments"`
	Execute            string `json:"execute"`
	FilterFileEnd      string `json:"filter_file_end"`
	FilterSubdir       string `json:"filter_subdir"`
	RequiresOutputLock bool   `json:"requires_output_lock"`
	StderrFile         string `json:"stderr_file"`
	StdoutFile         string `json:"stdout_file"`

	// MicroServiceChainChoice
	ChainChoices []string     `json:"chain_choices"`
	LinkChoices  []LinkChoice `json:"choices"`

	// TaskConfigSetUnitVariable
	Variable      string `json:"variable"`
	VariableValue string `json:"variable_value"`
	ChainID       string `json:"chain_id"`

	// MicroServiceChoiceReplacementDic
	Replacements []LinkConfigReplacement `json:"replacements"`

	// TaskConfigUnitVariableLinkPull
	// @ Variable, ChainID
}

type LinkChoice struct {
	Value  bool   `json:"value"`
	LinkID string `json:"link_id"`
}

type LinkConfigReplacement struct {
	ID          string            `json:"id"`
	Description I18nField         `json:"description"`
	Items       map[string]string `json:"items"`
}

type LinkExitCode struct {
	JobStatus string `json:"job_status"`
	LinkID    string `json:"link_id"`
}

type WorkflowData struct {
	WatchedDirectories []*WatchedDirectory `json:"watched_directories"`
	Chains             map[string]*Chain   `json:"chains"`
	Links              map[string]*Link    `json:"links"`
}

func New() *WorkflowData {
	return &WorkflowData{
		WatchedDirectories: []*WatchedDirectory{},
		Chains:             map[string]*Chain{},
		Links:              map[string]*Link{},
	}
}

func LoadWorkflowData(stream []byte) (*WorkflowData, error) {
	wd := New()
	if err := json.Unmarshal(stream, wd); err != nil {
		return nil, err
	}
	// This can probably be done by implementing Unmarshaller.
	for id, item := range wd.Chains {
		item.id = id
	}
	for id, item := range wd.Links {
		item.id = id
	}
	return wd, nil
}
