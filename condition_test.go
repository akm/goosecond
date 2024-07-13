package goosecond

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
)

func TestConditionAddFileNameMigrationContext(t *testing.T) {
	type givenArgSet struct {
		filename string
		up, down GoMigrationContext
	}
	dummyFunc := func(name string) func(context.Context, *sql.Tx) error {
		return func(context.Context, *sql.Tx) error { return fmt.Errorf(name) }
	}

	t.Run("true case", func(t *testing.T) {
		givenArgSets := []*givenArgSet{}
		defer setAddNamedMigrationContext(func(filename string, up, down GoMigrationContext) {
			givenArgSets = append(givenArgSets, &givenArgSet{filename, up, down})
		})()

		cond := NewCondition(func() bool { return true })
		cond.AddFileNameMigrationContext(dummyFunc("up"), dummyFunc("down"))

		if len(givenArgSets) != 1 {
			t.Errorf("len(givenArgSets) = %v, want %v", len(givenArgSets), 1)
		}
		args := givenArgSets[0]
		if filename := filepath.Base(args.filename); filename != "condition_test.go" {
			t.Errorf("filename = %v, want %v", filename, "condition_test.go")
		}
		// 関数は直接比較できないので、エラーメッセージで比較
		if args.up(nil, nil).Error() != "up" {
			t.Errorf("up = %v, want %v", args.up, "up")
		}
		if args.down(nil, nil).Error() != "down" {
			t.Errorf("down = %v, want %v", args.down, "down")
		}
	})

	t.Run("false case", func(t *testing.T) {
		givenArgSets := []*givenArgSet{}
		defer setAddNamedMigrationContext(func(filename string, up, down GoMigrationContext) {
			givenArgSets = append(givenArgSets, &givenArgSet{filename, up, down})
		})()

		cond := NewCondition(func() bool { return false })
		cond.AddFileNameMigrationContext(dummyFunc("up"), dummyFunc("down"))

		if len(givenArgSets) != 1 {
			t.Errorf("len(givenArgSets) = %v, want %v", len(givenArgSets), 1)
		}
		args := givenArgSets[0]
		if filename := filepath.Base(args.filename); filename != "condition_test.go" {
			t.Errorf("filename = %v, want %v", filename, "condition_test.go")
		}
		// 関数は直接比較できないので、エラーを返さないことを確認
		if err := args.up(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := args.down(nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func setAddNamedMigrationContext(fn func(filename string, up, down GoMigrationContext)) func() {
	var original func(filename string, up, down GoMigrationContext)
	AddNamedMigrationContext, original = fn, AddNamedMigrationContext
	return func() {
		AddNamedMigrationContext = original
	}
}
