package views

import tea "github.com/charmbracelet/bubbletea"

type BaseModel struct {
	Parent     tea.Model
	Breadcrumb []string
	Width      int
	Height     int
}
