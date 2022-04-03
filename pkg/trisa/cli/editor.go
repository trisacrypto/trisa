package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type EditorMode uint8

const (
	EditInstall EditorMode = iota + 1
	EditCreate
	EditUpdate
)

// Editor is an interactive application that prompts for input from the user.
type Editor struct {
	name     string
	scanner  *bufio.Reader
	profiles *Profiles
	mode     EditorMode
}

func (p *Profiles) NewEditor(name string, mode EditorMode) *Editor {
	return &Editor{
		name:     name,
		profiles: p,
		scanner:  bufio.NewReader(os.Stdin),
		mode:     mode,
	}
}

func (e *Editor) Edit() (err error) {
	// Step One: are we overwriting an existing file or profile?
	switch e.mode {
	case EditInstall:
		if err = e.CheckOverwrite(); err != nil {
			return err
		}
	case EditCreate:
		if err = e.CheckProfileExists(); err != nil {
			return err
		}
	}

	// If a profile with that name doesn't exist, create it
	if _, ok := e.profiles.Profiles[e.name]; !ok {
		e.profiles.Profiles[e.name] = &Profile{folder: e.profiles.folder}
	}

	// Get the pointer to the profile
	profile := e.profiles.Profiles[e.name]

	if e.PromptBool("lookup in directory service (or enter manually)?", false) {
		if err = e.DirectoryLookup(profile); err != nil {
			return err
		}
	} else {
		if err = e.ManualEntry(profile); err != nil {
			return err
		}
	}

	// If the profile being edited isn't active, prompt to activate it
	if e.name != e.profiles.Active {
		if e.PromptBool(fmt.Sprintf("activate %s?", e.name), true) {
			e.profiles.Active = e.name
		}
	}

	// Save the specified profile to disk
	profile.LastUpdated = time.Now()
	return nil
}

//====================================================================================
// User Input
//====================================================================================

func (e *Editor) PromptString(prompt, defaultText string) string {
	return e.prompt(prompt, defaultText)
}

func (e *Editor) PromptBool(prompt string, defaultBool bool) bool {
	var rep, defaultText string
	if defaultBool {
		defaultText = "Y/n"
	} else {
		defaultText = "y/N"
	}

	rep = e.prompt(prompt, defaultText)
	if rep == defaultText {
		return defaultBool
	}

	switch strings.ToLower(rep) {
	case "y", "ye", "yes", "t", "true", "on", "1":
		return true
	case "n", "no", "f", "false", "off", "0":
		return false
	default:
		fmt.Println("please specify yes or no")
		return e.PromptBool(prompt, defaultBool)
	}
}

func (e *Editor) prompt(prompt, defaultText string) string {
	if defaultText == "" {
		fmt.Print(prompt + ": ")
	} else {
		fmt.Printf("%s [%s]: ", prompt, defaultText)
	}

	text, _ := e.scanner.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultText
	}
	return text
}

//====================================================================================
// Workflow Commands
//====================================================================================

func (e *Editor) CheckOverwrite() (err error) {
	var path string
	if path, err = e.profiles.Path(); err != nil {
		return err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	if !e.PromptBool("a configuration file already exists, overwrite it?", false) {
		return ErrDoNotOverwrite
	}

	return nil
}

func (e *Editor) CheckProfileExists() (err error) {
	if _, ok := e.profiles.Profiles[e.name]; ok {
		if !e.PromptBool("a profile already exists with that name, update it?", false) {
			return ErrDoNotOverwrite
		}
	}
	return nil
}

func (e *Editor) DirectoryLookup(profile *Profile) (err error) {
	return nil
}

func (e *Editor) ManualEntry(profile *Profile) (err error) {
	// Prompt for the endpoint
	profile.Endpoint = e.PromptString("node endpoint with port", profile.Endpoint)

	return nil
}
