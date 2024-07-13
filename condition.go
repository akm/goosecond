package goosecond

import (
	"runtime"

	"github.com/pressly/goose/v3"
)

type Condition struct {
	predicate  func() bool
	CallerSkip int
}

func NewCondition(predicate func() bool) *Condition {
	return &Condition{predicate: predicate, CallerSkip: 1}
}

func (c *Condition) AddFileNameMigrationContext(up, down goose.GoMigrationContext) {
	_, filename, _, _ := runtime.Caller(c.CallerSkip)
	if c.predicate() {
		goose.AddNamedMigrationContext(filename, up, down)
	} else {
		goose.AddNamedMigrationContext(filename, EmptyMigrationContext, EmptyMigrationContext)
	}
}
