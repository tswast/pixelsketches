// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package strategy

import (
	"testing"

	"github.com/tswast/pixelsketches/village/gui"
)

func TestSimPaint(t *testing.T) {
	app := gui.NewAppState()
	got := simPaint(app, gui.Action{Horizontal: 1})
	if got.rate > 0 || got.reason != "no-different-colors-found" {
		t.Errorf("simPaint(NewAppState(), TO_THE_RIGHT) => %#v, want rate: 0, reason: no-different-colors-found", got)
	}
}
