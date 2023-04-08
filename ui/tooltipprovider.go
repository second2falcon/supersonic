package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
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

	toolTipLayer fyne.CanvasObject
}

// Returns the CanvasObject into which tool tips are rendered. This should be
// placed as a layer on top of the main application window's content in a MaxLayout.
func (t *ToolTipProvider) ToolTipLayer() fyne.CanvasObject {
	return t.toolTipLayer
}

func (t *ToolTipProvider) MakeToolTippable(obj fyne.CanvasObject, toolTipFn func() string) fyne.CanvasObject {
	toolTippable := &toolTippableObject{object: obj, toolTipFn: toolTipFn, toolTipProvider: t}
	toolTippable.ExtendBaseWidget(toolTippable)
	return toolTippable
}

type toolTippableObject struct {
	widget.BaseWidget

	object          fyne.CanvasObject
	toolTipFn       func() string
	toolTipProvider *ToolTipProvider
}

var _ desktop.Hoverable = (*toolTippableObject)(nil)

func (t *toolTippableObject) MouseIn(*desktop.MouseEvent) {

}

func (t *toolTippableObject) MouseOut() {

}

func (t *toolTippableObject) MouseMoved(*desktop.MouseEvent) {

}

func (t *toolTippableObject) CreateRenderer() fyne.WidgetRenderer {
	if t.object == nil {
		return nil
	}
	if o, ok := t.object.(fyne.Widget); ok {
		return o.CreateRenderer()
	}
	return widget.NewSimpleRenderer(t.object)
}
