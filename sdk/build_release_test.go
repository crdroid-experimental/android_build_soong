// Copyright (C) 2021 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"fmt"
	"testing"

	"android/soong/android"
)

// Tests for build_release.go

var (
	// Some additional test specific releases that are added after the currently supported ones and
	// so are treated as being for future releases.
	buildReleaseFuture1 = initBuildRelease("F1")
	buildReleaseFuture2 = initBuildRelease("F2")
)

func TestNameToRelease(t *testing.T) {
	t.Run("single release", func(t *testing.T) {
		release, err := nameToRelease("S")
		android.AssertDeepEquals(t, "errors", nil, err)
		android.AssertDeepEquals(t, "release", buildReleaseS, release)
	})
	t.Run("invalid release", func(t *testing.T) {
		release, err := nameToRelease("A")
		android.AssertDeepEquals(t, "release", (*buildRelease)(nil), release)
		// Uses a wildcard in the error message to allow for additional build releases to be added to
		// the supported set without breaking this test.
		android.FailIfNoMatchingErrors(t, `unknown release "A", expected one of \[S,T.*,F1,F2\]`, []error{err})
	})
}

func TestParseBuildReleaseSet(t *testing.T) {
	t.Run("single release", func(t *testing.T) {
		set, err := parseBuildReleaseSet("S")
		android.AssertDeepEquals(t, "errors", nil, err)
		android.AssertStringEquals(t, "set", "[S]", set.String())
	})
	t.Run("open range", func(t *testing.T) {
		set, err := parseBuildReleaseSet("F1+")
		android.AssertDeepEquals(t, "errors", nil, err)
		android.AssertStringEquals(t, "set", "[F1,F2]", set.String())
	})
	t.Run("closed range", func(t *testing.T) {
		set, err := parseBuildReleaseSet("S-F1")
		android.AssertDeepEquals(t, "errors", nil, err)
		android.AssertStringEquals(t, "set", "[S,T,F1]", set.String())
	})
	invalidAReleaseMessage := `unknown release "A", expected one of ` + allBuildReleaseSet.String()
	t.Run("invalid release", func(t *testing.T) {
		set, err := parseBuildReleaseSet("A")
		android.AssertDeepEquals(t, "set", (*buildReleaseSet)(nil), set)
		android.AssertStringDoesContain(t, "errors", fmt.Sprint(err), invalidAReleaseMessage)
	})
	t.Run("invalid release in open range", func(t *testing.T) {
		set, err := parseBuildReleaseSet("A+")
		android.AssertDeepEquals(t, "set", (*buildReleaseSet)(nil), set)
		android.AssertStringDoesContain(t, "errors", fmt.Sprint(err), invalidAReleaseMessage)
	})
	t.Run("invalid release in closed range start", func(t *testing.T) {
		set, err := parseBuildReleaseSet("A-S")
		android.AssertDeepEquals(t, "set", (*buildReleaseSet)(nil), set)
		android.AssertStringDoesContain(t, "errors", fmt.Sprint(err), invalidAReleaseMessage)
	})
	t.Run("invalid release in closed range end", func(t *testing.T) {
		set, err := parseBuildReleaseSet("T-A")
		android.AssertDeepEquals(t, "set", (*buildReleaseSet)(nil), set)
		android.AssertStringDoesContain(t, "errors", fmt.Sprint(err), invalidAReleaseMessage)
	})
	t.Run("invalid closed range reversed", func(t *testing.T) {
		set, err := parseBuildReleaseSet("F1-S")
		android.AssertDeepEquals(t, "set", (*buildReleaseSet)(nil), set)
		android.AssertStringDoesContain(t, "errors", fmt.Sprint(err), `invalid closed range, start release "F1" is later than end release "S"`)
	})
}

func TestBuildReleaseSetContains(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		set, _ := parseBuildReleaseSet("F1-F2")
		android.AssertBoolEquals(t, "set contains F1", true, set.contains(buildReleaseFuture1))
		android.AssertBoolEquals(t, "set does not contain S", false, set.contains(buildReleaseS))
		android.AssertBoolEquals(t, "set contains F2", true, set.contains(buildReleaseFuture2))
		android.AssertBoolEquals(t, "set does not contain T", false, set.contains(buildReleaseT))
	})
}
