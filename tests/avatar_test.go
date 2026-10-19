//go:build !wasm

package tests

import (
	"strings"
	"testing"
	"webtyp.com/avatar"
	"webtyp.com/svg"
)

func TestCatalog(t *testing.T) {
	cat := avatar.NewCatalog("test")
	if cat.Prefix != "test" {
		t.Errorf("expected Prefix 'test', got %v", cat.Prefix)
	}

	icon := svg.Icon("test-hair")
	cat.WithHair(icon)
	if cat.Hair != icon {
		t.Errorf("expected Hair %v, got %v", icon, cat.Hair)
	}

	cat.WithTorso(icon)
	if cat.Torso != icon {
		t.Errorf("expected Torso %v, got %v", icon, cat.Torso)
	}

	cat.WithHead(icon)
	if cat.Head != icon {
		t.Errorf("expected Head %v, got %v", icon, cat.Head)
	}
}

func TestAvatar_State(t *testing.T) {
	cat := avatar.NewCatalog("test")
	av := avatar.New(cat)

	av.SetState(avatar.Speak)
	av.Speak("hello")
	av.Stop()
}

func TestAvatar_Render(t *testing.T) {
	cat := avatar.NewCatalog("test")
	av := avatar.New(cat)

	// Test rendering valid SVG
	el := av.Render()
	markup := el.String()

	if !strings.Contains(markup, "viewBox='0 0 512 512'") {
		t.Errorf("expected viewBox attribute, got %v", markup)
	}
	if !strings.Contains(markup, "class='avatar'") {
		t.Errorf("expected class attribute, got %v", markup)
	}
	if !strings.Contains(markup, "<use ") {
		t.Errorf("expected <use> tags, got %v", markup)
	}
}

func TestAvatar_RenderCSS(t *testing.T) {
	cat := avatar.NewCatalog("test")
	av := avatar.New(cat)

	style := av.RenderCSS()
	if style == nil {
		t.Errorf("expected valid stylesheet, got nil")
	}
}
