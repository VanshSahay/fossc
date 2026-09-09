package tui

import "github.com/charmbracelet/lipgloss"

// theme mirrors the site's light/dark palettes (see app/GazetteClient.tsx and
// app/globals.css), reduced to foreground accents: this TUI paints no
// backgrounds at all — general text sits on the terminal's own background.
type theme struct {
	Strong lipgloss.Color // strong rules, table borders (#333333 / #666666)
	Soft   lipgloss.Color // box borders            (#999999 / #444444)
	Text   lipgloss.Color // body text              (#111111 / #e0e0e0)
	Muted  lipgloss.Color // bylines, hints         (#666666 / #999999)
	Link   lipgloss.Color // .lwn-link              (#003399 / #66b2ff)
}

func newTheme(dark bool) theme {
	if dark {
		return theme{
			Strong: "#666666", Soft: "#444444",
			Text: "#e0e0e0", Muted: "#999999", Link: "#66b2ff",
		}
	}
	return theme{
		Strong: "#333333", Soft: "#999999",
		Text: "#111111", Muted: "#666666", Link: "#003399",
	}
}

// styles are rebuilt whenever the theme changes; cheap enough to rebuild per
// render.
type styles struct {
	t theme
}

func (m model) sty() styles {
	return styles{t: newTheme(m.dark)}
}

// mastheadTitle is the serif-bold "FOSS CLUB KIET" h1.
func (s styles) mastheadTitle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Text)
}

func (s styles) mastheadNote() lipgloss.Style {
	return lipgloss.NewStyle().Italic(true).Foreground(s.t.Muted)
}

func (s styles) mastheadMeta() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(s.t.Muted)
}

// themeChip is the "🌙 Dark" toggle label in the masthead.
func (s styles) themeChip() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(s.t.Muted)
}

func (s styles) rule() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(s.t.Strong)
}

// sidebarBox renders ".lwn-sidebar-box": 1px border, no fill.
func (s styles) sidebarBox(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.t.Soft).
		Width(w)
}

// sidebarTitle renders ".lwn-sidebar-title" (a dark bar on the site), as bold
// text here.
func (s styles) sidebarTitle(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(s.t.Text).
		Width(w).
		Padding(0, 1)
}

func (s styles) navActive() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Text)
}

func (s styles) navItem() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(s.t.Link)
}

// banner is the full-width "Weekly Edition • Front Page" strip.
func (s styles) banner(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.t.Text).
		Border(lipgloss.NormalBorder(), true, false).
		BorderForeground(s.t.Strong).
		Bold(true).
		Padding(0, 1).
		Width(w)
}

func (s styles) kicker(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(color)
}

func (s styles) headline() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Text)
}

func (s styles) byline() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(s.t.Muted)
}

func (s styles) actionLink() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Link).Underline(true)
}

func (s styles) moduleBox(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.t.Soft).
		Padding(0, 1).
		Width(w)
}

func (s styles) prereqBox(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.t.Strong).
		Padding(0, 1).
		Width(w)
}

// chip is an emphasis label ("Join Discord ↗", "[y] Copy …").
func (s styles) chip() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Text)
}

func (s styles) discordChip() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Link)
}

func (s styles) hint() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(s.t.Muted)
}

func (s styles) statusChip() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Link)
}

func (s styles) manSection() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Underline(true).Foreground(s.t.Text)
}

// manSynopsis is the left-barred SYNOPSIS block.
func (s styles) manSynopsis(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(s.t.Strong).
		Padding(0, 1).
		Width(w)
}

func (s styles) manHeader() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(s.t.Text)
}
