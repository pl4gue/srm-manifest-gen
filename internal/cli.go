package srm

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func Spin(title string, action func()) error {
	return spinner.New().
		Title(title).
		Action(action).
		Run()
}

func Input(title string, validationFunction func(string) error) (string, error) {
	var choice string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(title).
				Value(&choice).
				Validate(validationFunction),
		),
	).Run()

	return choice, err
}

func Select[T comparable](title string, options []huh.Option[T]) (T, error) {
	var choice T

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[T]().
				Title(title).
				Options(options...).
				Value(&choice),
		),
	).Run()

	return choice, err
}

func Confirm(title string) (bool, error) {
	var confirm bool

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Value(&confirm),
		),
	).Run()

	return confirm, err
}
