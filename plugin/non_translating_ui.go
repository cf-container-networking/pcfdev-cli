package plugin

import "strings"

type NonTranslatingUI struct {
	UI
}

func (ui *NonTranslatingUI) Confirm(message string, args ...interface{}) bool {
	response := ui.Ask(message)
	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	}
	return false
}
