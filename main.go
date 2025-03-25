package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"time"
)

const (
	workTimeInMinutes = 25
)

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("ugh... the tomatoes are rotten, %v", err)
		os.Exit(1)
	}
}

type model struct {
	totalSeconds     int
	remainingSeconds int
	isPaused         bool
	message          string
	progress         progress.Model
	boxStyle         lipgloss.Style
}

func initialModel() model {
	return model{
		totalSeconds:     workTimeInMinutes * 60,
		remainingSeconds: workTimeInMinutes * 60,
		isPaused:         false,
		message:          "running",
		progress:         progress.New(progress.WithDefaultGradient()),
	}
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return time.Time(t)
	})
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	{
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "p", " ":
				m.isPaused = !m.isPaused
				if m.isPaused == true {
					m.message = "Paused"
				} else {
					m.message = "Resumed"
				}

			case "R", "r":
				m.remainingSeconds = m.totalSeconds
			}
		case time.Time:
			{
				if !m.isPaused {
					m.remainingSeconds = m.remainingSeconds - 1
					percent := 1 - (float64(m.remainingSeconds) / float64(m.totalSeconds))
					cmd := m.progress.SetPercent(percent)
					return m, tea.Batch(doTick(), cmd)
				} else {
					return m, doTick()
				}

			}

		case progress.FrameMsg:
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return m, cmd

		}
	}
	return m, nil
}

func (m model) View() string {
	var status string
	if m.isPaused {
		status = "WORK MODE - PAUSED 🟠"
	} else {
		status = "WORK MODE - RUNNING 🟢"
	}

	mins := m.remainingSeconds / 60
	secs := m.remainingSeconds % 60

	timeDisplay := fmt.Sprintf("[REMAINING TIME %02d:%02d OFF %02d:00]", mins, secs, workTimeInMinutes)
	controls := "Press q to quit"
	controls += "\nPress r to reset"
	controls += "\nPress space to pause/resume"

	content := status + "\n\n" +
		timeDisplay + "\n\n" +
		m.progress.View() + "\n\n" +
		controls

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 3).
		Width(60).
		Align(lipgloss.Center)

	return boxStyle.Render(content)

}

///
/// Previous Code using the GoRoutines
///

//func DetectEachSecond(
//	seconds int,
//	c chan int,
//) {
//	for i := 0; i < seconds; i++ {
//		time.Sleep(1 * time.Second)
//		c <- 1
//	}
//	close(c)
//}

//func main(){
//eachSecond := make(chan int)
//totalMinutes := 5
//seconds := totalMinutes * 60
//
//go DetectEachSecond(seconds, eachSecond)
//
//for {
//	_, open := <-eachSecond
//	if !open {
//		break
//	}
//	seconds = seconds - 1
//	mins := seconds / 60
//	secs := seconds % 60
//
//	fmt.Printf("\r%02d:%02d", mins, secs)
//}
//}
