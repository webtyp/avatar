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

	icon := svg.Icon("test-hair")
	cat = cat.WithHair(icon)

	icon2 := svg.Icon("test-torso")
	cat = cat.WithTorso(icon2)

	icon3 := svg.Icon("test-head")
	cat = cat.WithHead(icon3)

	av := avatar.New(cat)
	out := av.Render().String()

	if !strings.Contains(out, "#test-hair") {
		t.Errorf("expected Hair #test-hair, got %v", out)
	}
	if !strings.Contains(out, "#test-torso") {
		t.Errorf("expected Torso #test-torso, got %v", out)
	}
	if !strings.Contains(out, "#test-head") {
		t.Errorf("expected Head #test-head, got %v", out)
	}
}

func TestAvatar_WorkedExamples(t *testing.T) {
	tests := []struct {
		name     string
		calls    func(av *avatar.Avatar)
		contains []string
	}{
		{
			name: "New(NewCatalog(t))",
			calls: func(av *avatar.Avatar) {},
			contains: []string{
				`href='#t-brows-normal'`,
				`href='#t-eyes-normal'`,
				`href='#t-mouth-normal'`,
			},
		},
		{
			name: "SetEmotion(Sad)",
			calls: func(av *avatar.Avatar) {
				av.SetEmotion(avatar.Sad)
			},
			contains: []string{
				`#t-brows-sad`,
				`#t-eyes-sad`,
				`#t-mouth-sad`,
				`#t-nose-sad`,
				`#t-extras-drop`,
			},
		},
		{
			name: "SetEmotion(Angry); SetActivity(Speak)",
			calls: func(av *avatar.Avatar) {
				av.SetEmotion(avatar.Angry)
				av.SetActivity(avatar.Speak)
			},
			contains: []string{
				`#t-brows-anger`,
				`display='inline'`, // speech group display (dom.String uses single quotes)
				`#t-mouth-anger`, // rest frame
			},
		},
		{
			name: "SetActivity(Think) (Neutral)",
			calls: func(av *avatar.Avatar) {
				av.SetActivity(avatar.Think)
			},
			contains: []string{
				`#t-brows-think`,
				`#t-mouth-think`,
				`display='inline'`, // think group
			},
		},
		{
			name: "SetEmotion(Happy); SetActivity(Think)",
			calls: func(av *avatar.Avatar) {
				av.SetEmotion(avatar.Happy)
				av.SetActivity(avatar.Think)
			},
			contains: []string{
				`#t-brows-greeting`, // emotion wins
				`#t-mouth-think`,    // activity wins
			},
		},
		{
			name: "SetActivity(Fatigue); SetActivity(Greeting)",
			calls: func(av *avatar.Avatar) {
				av.SetActivity(avatar.Fatigue)
				av.SetActivity(avatar.Greeting)
			},
			contains: []string{
				`#t-mouth-greeting`,
			},
		},
		{
			name: "NewCatalog(t).WithHair(u-hair-short)",
			calls: func(av *avatar.Avatar) {
				cat := avatar.NewCatalog("t").WithHair(svg.Icon("u-hair-short"))
				*av = *avatar.New(cat)
			},
			contains: []string{
				`#u-hair-short`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cat := avatar.NewCatalog("t")
			av := avatar.New(cat)
			tc.calls(av)
			out := av.Render().String()

			for _, want := range tc.contains {
				if !strings.Contains(out, want) {
					t.Errorf("expected out to contain %q, but got:\n%s", want, out)
				}
			}
		})
	}
}

func TestAvatar_Speak_Stop(t *testing.T) {
	cat := avatar.NewCatalog("test")
	av := avatar.New(cat)
	av.Speak("hola")

	// During speak, group should have display='inline'
	outSpeak := av.Render().String()
	if !strings.Contains(outSpeak, `class='avatar__talk-hands' display='inline'`) {
		t.Errorf("expected talk-hands group to be visible while speaking, got %v", outSpeak)
	}

	av.Stop()
	outStop := av.Render().String()
	if strings.Contains(outStop, `class='avatar__talk-hands' display='inline'`) {
		t.Errorf("expected talk-hands group to be hidden after stop, got %v", outStop)
	}
}

func TestAvatar_SetLayerVisible(t *testing.T) {
	cat := avatar.NewCatalog("test")
	av := avatar.New(cat)

	av.SetLayerVisible(avatar.LayerHair, false)
	out := av.Render().String()

	// When layer hair is hidden, its group gets display='none'
	if !strings.Contains(out, `display='none'`) {
		t.Errorf("expected layer group to have display='none', got %v", out)
	}
}

func TestAvatar_RenderCSS(t *testing.T) {
	cat := avatar.NewCatalog("test")
	av := avatar.New(cat)

	sheet := av.RenderCSS()
	if sheet == nil {
		t.Errorf("expected valid stylesheet, got nil")
	}

	if errs := sheet.Validate(); len(errs) > 0 {
		t.Errorf("expected valid stylesheet, got errors: %v", errs)
	}
}
