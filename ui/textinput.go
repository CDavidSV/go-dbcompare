package ui

import (
	"fmt"
	"os"

	"github.com/CDavidSV/go-dbcompare/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModel struct {
	textInput      textinput.Model
	label          string
	err            error
	output         *TextInputValue
	req            bool
	validationFunc func(value string) error
}

type TextInputOptions struct {
	Label              string
	Placeholder        string
	CharLimit          int
	Required           bool
	MaskInput          bool
	ValidationFunction func(value string) error
}

type TextInputValue struct {
	Value string
}

type errMsg error

var labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#fffff")).Background(lipgloss.Color("99")).Padding(0, 1)

func InitialTextInputModel(options TextInputOptions, output *TextInputValue) TextInputModel {
	ti := textinput.New()

	ti.Placeholder = options.Placeholder
	ti.CharLimit = options.CharLimit
	ti.Focus()

	if options.MaskInput {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	if options.CharLimit < 1 {
		ti.CharLimit = 5000
	}

	output.Value = ""

	return TextInputModel{
		textInput:      ti,
		err:            nil,
		label:          options.Label,
		output:         output,
		req:            options.Required,
		validationFunc: options.ValidationFunction,
	}
}

func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.req && m.textInput.Value() == "" {
				return m, nil
			}
			if !m.req && m.textInput.Value() == "" {
				return m, tea.Quit
			}

			if m.validationFunc != nil {
				err := m.validationFunc(m.textInput.Value())

				if err != nil {
					m.err = err
					return m, textinput.Blink
				}
			}

			m.output.Value = m.textInput.Value()

			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(0)
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	errMsg := ""
	if m.err != nil {
		errMsg = m.err.Error()
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		labelStyle.Render(m.label),
		m.textInput.View(),
		config.ErrorStyle.Render(errMsg),
	)
}
