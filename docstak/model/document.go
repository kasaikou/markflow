package model

import "context"

type Document struct {
	Title       string
	Description string
	Tasks       map[string]DocumentTask
	GlobalEnvs  map[string]string
}

type DocumentConfig struct {
	Document         Document
	ExecPathResolver map[string]string
}

type DocumentTask struct {
	Title       string
	Call        string
	Description string
	Scripts     []DocumentTaskScript
	Envs        map[string]string
}

type DocumentTaskScript struct {
	ExecPath string
	Script   string
}

type NewDocumentOption func(ctx context.Context, d *DocumentConfig) error

func NewDocument(ctx context.Context, options ...NewDocumentOption) (Document, error) {
	document := DocumentConfig{
		Document: Document{
			Tasks: map[string]DocumentTask{},
		},
	}

	for i := range options {
		if err := options[i](ctx, &document); err != nil {
			return document.Document, err
		}
	}

	return document.Document, nil
}
