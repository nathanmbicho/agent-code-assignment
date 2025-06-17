package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// color constants
const (
	ErrorColor   = "196" // Bright red
	SuccessColor = "70"  // Green
	HeaderColor  = "33"  // Blue
	CodeColor    = "141" // Purple
	BorderColor  = "39"  // Cyan
	TextColor    = "250" // Light gray
	InfoColor    = "214" // Orange
)

var TextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(TextColor))

// ErrorStyle - error message style - bright red with bold
var ErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ErrorColor)).
	Bold(true).
	Margin(1, 0)

// SuccessStyle - success message style - green with bold
var SuccessStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(SuccessColor)).
	Bold(true).
	Margin(1, 0)

// SuccessStyle2 - success message style - green with bold
var SuccessStyle2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color(SuccessColor)).Bold(true)

// HeaderStyle - blue with underline and padding
var HeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(HeaderColor)).
	Bold(true).
	Padding(1)

// InfoStyle info/help text style
var InfoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(InfoColor)).
	Italic(true)

// CodeStyle code/filename style - purple with a background
var CodeStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(CodeColor)).
	Background(lipgloss.Color("235")).
	Padding(0, 1).
	Bold(true)

// CuteBorder - cli border style
var CuteBorder = lipgloss.Border{
	Top:         "._.:*:",
	Bottom:      "._.:*:",
	Left:        "|*",
	Right:       "|*",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}

// CLIStyle CLI container style with the cute border
var CLIStyle = lipgloss.NewStyle().
	Border(CuteBorder).
	BorderForeground(lipgloss.Color(BorderColor)).
	Padding(1, 2).
	Margin(1)

// InputStyle Input box style
var InputStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62")).
	Foreground(lipgloss.Color(TextColor)).
	Padding(0, 1).
	Width(60)

// File extension specific styles
var (
	GoFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")). // Go cyan
			Bold(true)

	JSFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")). // JavaScript yellow
			Bold(true)

	PythonFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")). // Python blue
			Bold(true)

	PHPFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")). // PHP purple
			Bold(true)
)

// Helper functions for common UI elements

// RenderError renders an error message with icon
func RenderError(message string) string {
	return ErrorStyle.Render("❌ " + message)
}

// RenderSuccess renders a success message with icon
func RenderSuccess(message string) string {
	return SuccessStyle.Render("✅ " + message)
}

// RenderHeader renders a header with styling
func RenderHeader(title string) string {
	return HeaderStyle.Render(title)
}

// RenderCode renders code/filename with appropriate styling
func RenderCode(code string) string {
	return CodeStyle.Render(code)
}

// RenderInfo renders info/help text
func RenderInfo(info string) string {
	return InfoStyle.Render("💡 " + info)
}

// GetFileStyle returns appropriate style based on file extension
func GetFileStyle(extension string) lipgloss.Style {
	switch extension {
	case ".go":
		return GoFileStyle
	case ".js":
		return JSFileStyle
	case ".py":
		return PythonFileStyle
	case ".php":
		return PHPFileStyle
	default:
		return CodeStyle
	}
}

// RenderFileWithExtension renders filename with extension-specific styling
func RenderFileWithExtension(filename string, extension string) string {
	style := GetFileStyle(extension)
	return style.Render(filename)
}

// RenderProgressBar renders a simple progress bar
func RenderProgressBar(progress int, total int, width int) string {
	if total == 0 {
		return ""
	}

	filled := (progress * width) / total
	bar := ""

	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	return progressStyle.Render(bar)
}
