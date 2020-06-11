// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"testing"
	"time"

	"golang.org/x/build/maintner"
)

func TestNewIntervalFromCL(t *testing.T) {
	var (
		t0 = time.Now()
		t1 = t0.Add(1 * time.Hour)
	)
	testCases := []struct {
		cl         *maintner.GerritCL
		start, end int64
	}{
		{
			cl: &maintner.GerritCL{
				Created: t0,
				Status:  "new",
			},
			start: t0.Unix(),
			end:   math.MaxInt64,
		},
		{
			cl: &maintner.GerritCL{
				Created: t0,
				Status:  "merged",
				Metas: []*maintner.GerritMeta{
					{
						Commit: &maintner.GitCommit{
							Msg:        "autogenerated:gerrit:merged",
							CommitTime: t1,
						},
					},
				},
			},
			start: t0.Unix(),
			end:   t1.Unix(),
		},
		{
			cl: &maintner.GerritCL{
				Created: t0,
				Status:  "abandoned",
				Metas: []*maintner.GerritMeta{
					{
						Commit: &maintner.GitCommit{
							Msg:        "autogenerated:gerrit:abandon",
							CommitTime: t1,
						},
					},
				},
			},
			start: t0.Unix(),
			end:   t1.Unix(),
		},
	}

	for _, tc := range testCases {
		ival := newIntervalFromCL(tc.cl)
		if got, want := ival.start, tc.start; got != want {
			t.Errorf("start: got %d; want %d", got, want)
		}
		if got, want := ival.end, tc.end; got != want {
			t.Errorf("end: got %d; want %d", got, want)
		}
		if got, want := ival.cl, tc.cl; got != want {
			t.Errorf("cl: got %+v; want %+v", got, want)
		}
	}
}

func TestIntervalIntersection(t *testing.T) {
	testCases := []struct {
		interval   *clInterval
		t0, t1     time.Time
		intersects bool
	}{
		{
			&clInterval{start: 0, end: 5},
			time.Unix(0, 0),
			time.Unix(10, 0),
			true,
		},
		{
			&clInterval{start: 10, end: 20},
			time.Unix(0, 0),
			time.Unix(10, 0),
			true,
		},
		{
			&clInterval{start: 10, end: 20},
			time.Unix(0, 0),
			time.Unix(9, 0),
			false,
		},
		{
			&clInterval{start: 0, end: 5},
			time.Unix(6, 0),
			time.Unix(10, 0),
			false,
		},
	}

	for _, tc := range testCases {
		if got, want := tc.interval.intersects(tc.t0, tc.t1), tc.intersects; got != want {
			t.Errorf("(%v).intersects(%v, %v): got %v; want %v", tc.interval, tc.t0, tc.t1, got, want)
		}
	}
}
