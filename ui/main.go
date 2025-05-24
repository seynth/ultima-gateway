package ui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"ultima/auxiliary"
	"ultima/model"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type JSON struct {
	StartURL string `json:"startURL"`
}

type GenerateConfigKeyHash struct {
	pureJson   string
	configHash string
	plusUrl    string
	finalHash  string
}

type Ultima struct {
	input          textinput.Model
	spinner        spinner.Model
	genConfKeyHash GenerateConfigKeyHash
	jsonObject     JSON

	message       string
	isSubmit      bool
	completed     bool
	statusCode    int
	statusNext    string
	completedStep []string
}

type tickMsg time.Time

func UltimaInit() Ultima {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = " \"I need config file . . .\""
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = EnhanceText

	return Ultima{
		input:          ti,
		spinner:        s,
		genConfKeyHash: GenerateConfigKeyHash{},
		completedStep:  []string{},
	}
}

func (ult Ultima) Init() tea.Cmd {
	return nil
}

func (ult Ultima) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !ult.isSubmit {
			switch msg.Type {
			case tea.KeyEnter:
				if ult.input.Value() != "" {
					ult.isSubmit = true
					ult.statusCode = 0

					return ult, tea.Batch(
						ult.spinner.Tick,
						tea.Tick(time.Second, func(t time.Time) tea.Msg {
							return tickMsg(t)
						}),
					)
				}
			case tea.KeyCtrlC, tea.KeyEsc:
				return ult, tea.Quit
			}
		} else {
			if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc || msg.String() == "q" {
				return ult, tea.Quit
			}
		}

	case model.ReadAndConvert:
		ult.genConfKeyHash.pureJson = msg.EncodedContent
		if len(ult.genConfKeyHash.pureJson) > 0 {
			s := fmt.Sprintf("Ultima is %s", GreenColor.Render("online"))
			ult.completedStep = append(ult.completedStep, s)
			ult.statusCode = 2
			return ult, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	case tickMsg:
		if ult.isSubmit && !ult.completed {
			switch ult.statusCode {
			case 0:
				ult.statusCode = 1
				ult.statusNext = "Starting Ultima Gateway"
				return ult, tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
					return tickMsg(t)
				})
			case 1:
				file := strings.ReplaceAll(ult.input.Value(), `"`, "")
				ult.statusNext = "Searching startURL key"
				return ult, auxiliary.Convert(file)

			case 2:
				ult.statusCode = 3
				json.Unmarshal([]byte(ult.genConfKeyHash.pureJson), &ult.jsonObject)
				ult.statusNext = "Enhance Feature"
				ult.completedStep = append(ult.completedStep, "Found startURL")
				return ult, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
					return tickMsg(t)
				})
			case 3:
				ult.statusCode = 4
				ult.statusNext = "Starting chrome"
				ult.completedStep = append(ult.completedStep, EnhanceText.Render("Dynamic Header Enabled"))
			case 4:
				ult.genConfKeyHash.finalHash = auxiliary.Sha256Encode([]byte(ult.genConfKeyHash.plusUrl))
				ult.completed = true
				return ult, auxiliary.StartChrome(
					ult.jsonObject.StartURL,
					ult.genConfKeyHash.configHash,
					ult.genConfKeyHash.finalHash,
					auxiliary.Sha256Encode([]byte("random")),
				)

			}
			return ult, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	}

	if ult.completed {
		ult.statusNext = WhiteText.Render(fmt.Sprintf("You can %s this terminal but dont %s it!!", GreenColor.Render("minimize"), RedColor.Render("close")))
	}

	var cmd tea.Cmd
	if ult.isSubmit {
		ult.spinner, cmd = ult.spinner.Update(msg)
	} else {
		ult.input, cmd = ult.input.Update(msg)
	}
	return ult, cmd
}

func (ult Ultima) View() string {
	var builder strings.Builder
	subTitle := fmt.Sprintf("%s %s", EnhanceText.Render("[Enhance]"), WhiteText.Render("Bypass all you can imagine"))

	builder.WriteString(TextStyle.Render(Banner))
	builder.WriteString("\n")
	builder.WriteString(subTitle)
	builder.WriteString("\n")
	if ult.jsonObject.StartURL != "" {
		builder.WriteString(LinkStyle.Render(ult.jsonObject.StartURL))
	}
	builder.WriteString("\n")

	builder.WriteString("\n")
	if ult.isSubmit {
		if len(ult.completedStep) > 0 {
			builder.WriteString(
				Text.Render(
					strings.Join(ult.completedStep, "\n"),
				),
			)
			builder.WriteString("\n\n")
		}

		if !ult.completed {
			builder.WriteString(ult.spinner.View())
		}
		builder.WriteString(" ")
		builder.WriteString(WhiteText.Render(ult.statusNext))

		builder.WriteString(HelpText("esc", "quit"))
		return CenterView.Render(builder.String())
	}

	if len(strings.TrimSpace(ult.message)) > 0 {
		builder.WriteString(TextStyle.Render(ult.message))
		return builder.String()
	}

	builder.WriteString("Please specify SEB config path\n")
	builder.WriteString(TextInput.Render(ult.input.View()))
	if len(ult.input.Value()) > 0 {
		builder.WriteString("\n\npress enter to begin")
	}
	builder.WriteString(HelpText("esc", "quit"))
	return CenterView.Render(builder.String())
}
