package schema_test

import (
	"errors"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/schema"
)

func TestCheckAcceptsEveryVersionUpToCurrent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		stored, current schema.Version
		wantErr         bool
	}{
		{0, 1, false},
		{1, 1, false},
		{0, 2, false},
		{1, 2, false},
		{2, 2, false},
		{2, 1, true},
		{3, 2, true},
		{-1, 2, true},
	}
	for _, tc := range tests {
		err := schema.Check(tc.stored, tc.current)
		if (err != nil) != tc.wantErr {
			t.Errorf("Check(%d, %d) = %v, wantErr %v", tc.stored, tc.current, err, tc.wantErr)
		}
		if err != nil && !errors.Is(err, schema.ErrUnsupported) {
			t.Errorf("Check(%d, %d) error does not wrap ErrUnsupported: %v", tc.stored, tc.current, err)
		}
	}
}

func TestCheckTellsAnOlderRskToUpgrade(t *testing.T) {
	t.Parallel()
	err := schema.Check(2, 1)
	if err == nil || err.Error() == "" {
		t.Fatal("expected an error")
	}
	want := "unsupported schema version: file version 2, rsk supports 1 — upgrade rsk or re-run 'rsk init'"
	if err.Error() != want {
		t.Errorf("error = %q\nwant    %q", err, want)
	}
}
