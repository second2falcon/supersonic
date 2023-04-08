package ui

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ToolTipProvider is a component to allow fyne.CanvasObjects to have tool tips,
// which are shown on hover events (desktop only).
// It renders tool tips into the object returned by ToolTipLayer(),
// which should be stacked in a MaxLayout on top of the app's main window content.
type ToolTipProvider struct {

	// Amount of delay after a tool tippable item is hovered before the tool tip is shown,
	// if there has *not* been a tooltip shown within the last BetweenToolTipTime.
	InitialToolTipDelay time.Duration

	// Amount of delay after a tool tippable item is hovered before the tool tip is shown,
	// if there *has* been a tooltip shown within the last BetweenToolTipTime.
	SubsequentToolTipDelay time.Duration

	// Amount of time after a tool tip is dismissed before switching back to
	// the InitialToolTipDelay for the next tooltip that is shown.
	BetweenToolTipTime time.Duration

	mutex              sync.Mutex
	pendingToolTipObj  *toolTippableObject
	showToolTipAfter   time.Time
	lastToolTipShownAt time.Time
	pendingToolTipFn   func()

	toolTip      *canvas.Text
	toolTipLayer *fyne.Container
}

func NewToolTipProvider() *ToolTipProvider {
	toolTip := canvas.NewText("", theme.ForegroundColor())
	toolTip.TextSize = theme.CaptionTextSize()
	return &ToolTipProvider{
		toolTip:      toolTip,
		toolTipLayer: container.NewWithoutLayout(toolTip),
	}
}

// Returns the CanvasObject into which tool tips are rendered. This should be
// placed as a layer on top of the main application window's content in a MaxLayout.
func (t *ToolTipProvider) ToolTipLayer() fyne.CanvasObject {
	return t.toolTipLayer
}

func (t *ToolTipProvider) MakeToolTippable(obj fyne.CanvasObject, toolTipFn func() string) fyne.CanvasObject {
	toolTippable := &toolTippableObject{object: obj, toolTipFn: toolTipFn, provider: t}
	toolTippable.ExtendBaseWidget(toolTippable)
	return toolTippable
}

func (t *ToolTipProvider) newWaitToolTipFn() func() {
	return func() {
		// we only return when we've either shown a tool tip
		// or there is no tool tip queued to show
		// else we keep waiting
		for {
			t.mutex.Lock()
			if t.pendingToolTipObj == nil {
				// no tool tip queued to be shown
				t.pendingToolTipFn = nil
				t.mutex.Unlock()
				return
			}

			now := time.Now()
			if now.After(t.showToolTipAfter) {
				// show tool tip
				t.toolTip.Text = t.pendingToolTipObj.toolTipFn()
				t.pendingToolTipObj = nil
				t.pendingToolTipFn = nil
				t.mutex.Unlock()
				t.toolTip.Refresh()
				return
			} else {
				// wait to check after the next showToolTipAfter time
				newDelay := t.showToolTipAfter.Sub(now)
				t.mutex.Unlock()
				<-time.After(newDelay)
			}
		}
	}
}

type toolTippableObject struct {
	widget.BaseWidget

	object    fyne.CanvasObject
	toolTipFn func() string
	provider  *ToolTipProvider
}

var _ desktop.Hoverable = (*toolTippableObject)(nil)

func (t *toolTippableObject) MouseIn(*desktop.MouseEvent) {
	t.provider.mutex.Lock()
	defer t.provider.mutex.Unlock()
	t.provider.pendingToolTipObj = t

	now := time.Now()
	var delay time.Duration
	if time.Since(t.provider.lastToolTipShownAt) > t.provider.BetweenToolTipTime {
		delay = t.provider.InitialToolTipDelay
	} else {
		delay = t.provider.SubsequentToolTipDelay
	}
	t.provider.showToolTipAfter = now.Add(delay)
	if t.provider.pendingToolTipFn == nil {
		t.provider.pendingToolTipFn = t.provider.newWaitToolTipFn()
		_ = time.AfterFunc(delay, t.provider.pendingToolTipFn)
	}
}

func (t *toolTippableObject) MouseOut() {
	var needRefresh bool

	t.provider.mutex.Lock()
	t.provider.pendingToolTipObj = nil
	if t.provider.toolTip.Text != "" {
		needRefresh = true
		t.provider.toolTip.Text = ""
	}
	t.provider.mutex.Unlock()

	if needRefresh {
		t.provider.toolTip.Refresh()
	}
}

func (t *toolTippableObject) MouseMoved(*desktop.MouseEvent) {}

func (t *toolTippableObject) CreateRenderer() fyne.WidgetRenderer {
	if t.object == nil {
		return nil
	}
	if o, ok := t.object.(fyne.Widget); ok {
		return o.CreateRenderer()
	}
	return widget.NewSimpleRenderer(t.object)
}
