package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

var Banner = `
 ____ ___.__   __  .__                   ________        __                               
|    |   \  |_/  |_|__| _____ _____     /  _____/_____ _/  |_  ______  _  _______  ___.__.
|    |   /  |\   __\  |/     \\__  \   /   \  ___\__  \\   __\/ __ \ \/ \/ /\__  \<   |  |
|    |  /|  |_|  | |  |  Y Y  \/ __ \_ \    \_\  \/ __ \|  | \  ___/\     /  / __ \\___  |
|______/ |____/__| |__|__|_|  (____  /  \______  (____  /__|  \___  >\/\_/  (____  / ____|
                            \/     \/          \/     \/          \/             \/\/     

`
var termWidth, termHeight, _ = term.GetSize(os.Stdout.Fd())

var TextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#943ddb"))

var TextInput = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FFFFFF"))

var TextStyle1 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF")).
	Bold(true).
	Margin(1, 2)

var TextStyle2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF")).
	Bold(true).
	Margin(0, 2)

var Header = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(0, 2).
	MarginBottom(1)

var SubTitle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#A0A0A0")).
	MarginBottom(1)

var LinkStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00BFFF")).
	Underline(true)

var CenterView = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Width(termWidth).
	AlignHorizontal(lipgloss.Center).
	Padding(1, 2)

var CenterText = lipgloss.NewStyle().
	Bold(true).
	Width(termWidth).
	AlignHorizontal(lipgloss.Center)

var GreenCheck = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#2bd94b"))

var CenterBottomText = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8a8a8a")).
	Align(lipgloss.Center).
	Height(termHeight).
	AlignVertical(lipgloss.Center).
	Margin(0, 0, 5)

var TextWarning = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFCC00"))

func HelpText(key, info string) string {
	bold := lipgloss.NewStyle().
		Bold(true)
	inf := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8a8a8a"))
	format := fmt.Sprintf("\n\n%s %s\n\n", bold.Render(key), inf.Render(info))
	return format
}
