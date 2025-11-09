package main

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/ambersignal/blacksunrising/pkg/loader"
)

type RenderContext struct {
	Loader *loader.Loader
}

type Widget interface {
	Render(rndrCtx RenderContext) (widget.PreferredSizeLocateableWidget, error)
}

type Root []Widget

func (r Root) Build(rndrCtx RenderContext) (*ebitenui.UI, error) {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(5))),
		),
	)

	var errs []error

	for _, item := range r {
		widget, err := item.Render(rndrCtx)

		if err != nil {
			errs = append(errs, err)
		}

		if widget != nil {
			rootContainer.AddChild(widget)
		}
	}

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return &ebitenui.UI{
		Container: rootContainer,
	}, nil
}

type Menu []Widget

func (m Menu) Render(rndrCtx RenderContext) (widget.PreferredSizeLocateableWidget, error) {
	containerImg, err := rndrCtx.Loader.LoadImage("container.png")
	if err != nil {
		return nil, fmt.Errorf("load container image: %w", err)
	}

	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(&widget.Insets{
				Left:   16,
				Right:  16,
				Top:    10,
				Bottom: 10,
			})),
		),
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceSimple(containerImg, 16, 64),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)

	var errs []error
	for _, item := range m {
		widget, err := item.Render(rndrCtx)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		container.AddChild(widget)
	}

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return container, nil
}

type MenuButton struct {
	Text    string
	OnPress widget.ButtonClickedHandlerFunc
}

func menuButton(text string) MenuButton {
	return MenuButton{
		Text: text,
	}
}

func (mb MenuButton) Render(rndrCtx RenderContext) (widget.PreferredSizeLocateableWidget, error) {
	buttonReleasedImg, err := rndrCtx.Loader.LoadImage("button_released.png")
	if err != nil {
		return nil, fmt.Errorf("load button image: %w", err)
	}

	buttonPressedImg, err := rndrCtx.Loader.LoadImage("button_pressed.png")
	if err != nil {
		return nil, fmt.Errorf("load button image: %w", err)
	}

	face, err := rndrCtx.Loader.LoadFont("laika-14.ttf", 14)
	if err != nil {
		return nil, err
	}

	buttonOpts := []widget.ButtonOpt{
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSlice(buttonReleasedImg, [3]int{7, 82, 7}, [3]int{13, 6, 13}),
			Pressed: image.NewNineSlice(buttonPressedImg, [3]int{7, 82, 7}, [3]int{13, 6, 13}),
		}),
		widget.ButtonOpts.Text(mb.Text, &face, &widget.ButtonTextColor{
			Idle:    color.White,
			Pressed: color.Black,
		}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   16,
			Right:  16,
			Top:    10,
			Bottom: 10,
		}),
	}

	if mb.OnPress != nil {
		buttonOpts = append(buttonOpts, widget.ButtonOpts.ClickedHandler(mb.OnPress))
	}

	button := widget.NewButton(buttonOpts...)

	return button, nil
}
