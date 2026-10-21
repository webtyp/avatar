//go:build !wasm

package tests

import (
	"strings"
	"testing"
	"webtyp.com/avatar/reference"
)

func TestReference_IconSvg(t *testing.T) {
	out := reference.IconSvg().String()

	variants := []string{
		"normal",
		"greeting",
		"sad",
		"astonishment",
		"anger",
		"fatigue",
		"listen",
		"think",
		"talk-a",
		"talk-o",
		"closed",
	}

	// Check brows, eyes, mouth
	layers := []string{"brows", "eyes", "mouth"}
	for _, l := range layers {
		for _, v := range variants {
			id := reference.Prefix + "-" + l + "-" + v
			if !strings.Contains(out, `id="`+id+`"`) {
				t.Errorf("missing face variant ID: %s", id)
			}
		}
	}

	// noses
	noses := []string{"nose-sad", "nose-closed"}
	for _, n := range noses {
		id := reference.Prefix + "-" + n
		if !strings.Contains(out, `id="`+id+`"`) {
			t.Errorf("missing nose ID: %s", id)
		}
	}

	// poses
	poses := []string{
		"pose-greeting",
		"pose-think",
		"pose-talk-a",
		"pose-talk-o",
	}
	for _, p := range poses {
		id := reference.Prefix + "-" + p
		if !strings.Contains(out, `id="`+id+`"`) {
			t.Errorf("missing pose ID: %s", id)
		}
	}

	// extras
	extras := []string{
		"extras-drop",
		"extras-marks",
		"extras-wave",
		"extras-speech-a",
		"extras-speech-o",
		"thought-bubble-1",
		"thought-bubble-2",
		"thought-bubble-3",
		"thought-bubble-3-ring",
		"thought-dot-1",
		"thought-dot-2",
		"thought-dot-3",
		"sound-inner",
		"sound-mid",
		"sound-outer",
		"zzz-1",
		"zzz-2",
		"zzz-3",
	}
	for _, x := range extras {
		id := reference.Prefix + "-" + x
		if !strings.Contains(out, `id="`+id+`"`) {
			t.Errorf("missing extra ID: %s", id)
		}
	}

	// The drawing has no happy variants; Happy uses greeting.
}
