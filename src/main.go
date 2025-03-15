package main

import (
	"bufio"
	"fmt"
	"log"
	"path/filepath"

	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Menu struct {
	header  string
	footer  string
	choices []string
	cursor  int
}

type AliasAndCmd struct {
	alias string
	cmd   string
}

var currAliases []AliasAndCmd

func getBashAliases() error {
	var curr AliasAndCmd

	// Getting the dir for .bash_aliases
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get UserHomeDir")
		return nil
	}

	// Opening file
	dir = filepath.Join(dir, ".bash_aliases")
	file, err := os.Open(dir)
	if err != nil {
		log.Fatalf("Failed to open .bash_aliases")
		return err
	}

	defer file.Close() // Makes sure that file is closed at the end of the function

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.SplitN(line, " ", 2)[0] == "alias" { // We only care about lines that start with alias
			aliasAndCmd := strings.SplitN(line, "=", 2) // Only split on the first occurance of =
			curr.alias = aliasAndCmd[0]                 // alias = first element
			curr.cmd = aliasAndCmd[1]                   // cmd = second element
			currAliases = append(currAliases, curr)     // add to current aliases
		}
	}

	// If any errors occured during the reading of the file
	// log an error and return
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read .bash_aliases: %s", err)
		return err
	}

	return nil
}

func initialModel() Menu {
	return Menu{
		choices: []string{"Create New Alias", "Delete Aliases", "View Current Aliases", "Exit"},
		cursor:  0, // Cursor starts on Create New Alias
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			// Do something on enter

		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Menu) View() string {
	// The header
	s := m.header + "\n\n"
	// Iterate over our choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// Send the UI for rendering
	return s
}

func main() {

	// Pre-load all bash aliases
	err := getBashAliases()
	if err != nil {
		log.Fatalf("%s", err)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There has been error: %v", err)
		os.Exit(1)
	}
}
