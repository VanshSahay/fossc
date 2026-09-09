package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/VanshSahay/fossc/internal/club"
)

// viewDef is one entry of the sidebar navigation.
type viewDef struct {
	id    string // stable id, also used by announcement action links
	label string // sidebar label
	tag   string // banner label
}

var viewDefs = []viewDef{
	{"frontpage", "Front page", "WEEKLY EDITION • FRONT PAGE"},
	{"weekly", "Weekly edition", "WEEKLY EDITION • THIS WEEK"},
	{"manpage", "man fossc(1)", "FOSSC(1) · LINUX REFERENCE MANUAL"},
	{"bootcamp", "Bootcamp (Sep 9)", "EVENT DOSSIER // LINUX BASICS & OPEN-SOURCE HISTORY BOOTCAMP"},
	{"members", "Members", "MEMBERS"},
	{"discord", "Discord Community", "COMMUNITY // DISCORD HUB"},
}

const (
	sidebarWidth = 30
	footerHeight = 1
	mastheadH    = 4 // title, note, rule, blank
)

// layout returns the working page geometry: full terminal width.
func (m model) layout() (pageW, sidebarW, mainW, bodyH int) {
	pageW = m.w
	sidebarW = 0
	if pageW >= 72 {
		sidebarW = sidebarWidth
	}
	mainW = max(pageW-sidebarW, 20)
	bodyH = max(m.h-mastheadH-footerHeight, 1)
	return pageW, sidebarW, mainW - 1, bodyH // -1: gutter column
}

// View assembles the full page: masthead, sidebar + scrolling main, footer.
func (m model) View() string {
	if !m.ready {
		return "Loading the gazette…"
	}
	pageW, sidebarW, _, bodyH := m.layout()

	masthead := m.renderMasthead(pageW)

	var body string
	if sidebarW > 0 {
		sidebar := m.renderSidebar(sidebarW - 2)
		body = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, " ", m.vp.View())
	} else {
		body = m.vp.View()
	}
	body = lipgloss.NewStyle().Height(bodyH).MaxHeight(bodyH).Render(body)

	return lipgloss.JoinVertical(lipgloss.Left, masthead, body, m.renderFooter(pageW))
}

// renderMasthead is the sticky header: brand, tagline, date, theme toggle.
func (m model) renderMasthead(w int) string {
	s := m.sty()
	date := time.Now().Format("January 2, 2006")

	toggle := "🌙 Dark"
	if m.dark {
		toggle = "☀️ Light"
	}
	chip := s.themeChip().Render(toggle)

	right := lipgloss.JoinHorizontal(lipgloss.Center,
		s.mastheadMeta().Render(date), " ", chip)
	left := s.mastheadTitle().Render("FOSS CLUB KIET")
	gap := max(w-lipgloss.Width(left)-lipgloss.Width(right), 1)
	head := lipgloss.JoinHorizontal(lipgloss.Bottom, left, strings.Repeat(" ", gap), right)

	note := s.mastheadNote().Render(club.MastheadNote)
	if lipgloss.Width(note) > w {
		note = wrapAnsi(note, w)
	}
	rule := s.rule().Render(strings.Repeat("━", w))

	return lipgloss.JoinVertical(lipgloss.Left, head, note, rule, "")
}

// renderSidebar draws the LWN-style widget stack: nav, hours, quote, CTA.
func (m model) renderSidebar(w int) string {
	s := m.sty()
	boxW := w - 2 // inner content width inside the border

	// Box 1: club navigation
	var nav []string
	for i, vd := range viewDefs {
		st := s.navItem()
		marker := "  "
		if i == m.cur {
			st = s.navActive()
			marker = "▸ "
		}
		label := vd.label
		if vd.id == "bootcamp" {
			label = lipgloss.NewStyle().Bold(true).Render(label)
		}
		nav = append(nav, marker+st.Render(label))
	}
	navBox := joinBox(s, boxW, "FOSS CLUB KIET", strings.Join(nav, "\n"))

	// Box 2: hours & venue
	hours := strings.Join([]string{
		bold("Daily Open Hours:") + "\n" + club.DailyHours,
		bold("Weekly Meetup:") + "\n" + club.WeeklyMeetup,
		bold("Club Room Location:") + "\n" + club.ClubRoom,
		bold("Paper Reading:") + "\n" + club.PaperReading,
	}, "\n\n")
	hoursBox := joinBox(s, boxW, "HOURS & VENUE", hours)

	// Box 3: quote of the day
	q := club.QuoteOfTheDay()
	quote := wrapAnsi("“"+q.Text+"”", boxW)
	quoteBox := joinBox(s, boxW, "QUOTE OF THE DAY",
		s.mastheadNote().Render(quote)+"\n"+rightAlign(bold("— "+q.Author), boxW))

	// Box 4: discord CTA
	cta := strings.Repeat(" ", max(0, (boxW-lipgloss.Width("Join Discord ↗"))/2)) +
		link(s.chip(), club.DiscordInvite, "Join Discord ↗")
	ctaBox := joinBox(s, boxW, "", cta)

	return lipgloss.JoinVertical(lipgloss.Left,
		navBox, "", hoursBox, "", quoteBox, "", ctaBox)
}

// joinBox wraps body text in a bordered sidebar box with a title bar.
func joinBox(s styles, w int, title, body string) string {
	if title != "" {
		body = s.sidebarTitle(w).Render(title) + "\n" + body
	}
	inner := wrapAnsi(body, w)
	return s.sidebarBox(w).Render(inner)
}

// renderFooter is the key-hint status line; transient messages take over.
func (m model) renderFooter(w int) string {
	s := m.sty()
	hints := "1-6 view · ←/→ switch · ↑/↓ scroll · o open · y copy · d theme · q quit"
	if m.status != "" {
		return s.statusChip().Render(" " + m.status + " ")
	}
	right := "fossc(1) · " + link(s.hint(), club.Website, club.Website)
	gap := w - lipgloss.Width(hints) - lipgloss.Width(right)
	if gap < 1 {
		return s.hint().Render(hints)
	}
	return s.hint().Render(hints + strings.Repeat(" ", max(gap, 0)) + right)
}

// renderMain builds the scrollable document for the active view.
func (m model) renderMain(w int) string {
	s := m.sty()
	id := viewDefs[m.cur].id
	switch id {
	case "frontpage", "weekly":
		return m.renderFrontPage(s, w, viewDefs[m.cur].tag)
	case "manpage":
		return m.renderManPage(s, w)
	case "bootcamp":
		return m.renderBootcamp(s, w)
	case "members":
		return m.renderMembers(s, w)
	case "discord":
		return m.renderDiscord(s, w)
	}
	return ""
}

// banner renders the full-width section strip with left/right spans.
func banner(s styles, w int, left, right string) string {
	left = strings.ToUpper(left)
	inner := w - 2 // padding
	gap := inner - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		return s.banner(w).Render(left)
	}
	return s.banner(w).Render(lipgloss.JoinHorizontal(lipgloss.Top,
		left, strings.Repeat(" ", gap), right))
}

func (m model) renderFrontPage(s styles, w int, tag string) string {
	var b strings.Builder
	b.WriteString(banner(s, w, tag, time.Now().Format("January 2, 2006")))
	b.WriteString("\n\n")

	for i, a := range club.Announcements {
		catColor := lipgloss.Color(a.CategoryLight)
		if m.dark {
			catColor = lipgloss.Color(a.CategoryDark)
		}
		b.WriteString(s.kicker(catColor).Render(strings.ToUpper(a.Category)))
		b.WriteString("\n")
		b.WriteString(s.headline().Render(wrapAnsi(a.Title, w)))
		b.WriteString("\n")

		byline := "[" + a.Date + "]"
		if a.Author != "" {
			byline += " By " + a.Author
		}
		if a.VenueOrDetails != "" {
			byline += " • " + a.VenueOrDetails
		}
		b.WriteString(s.byline().Render(wrapAnsi(byline, w)))
		b.WriteString("\n\n")

		for _, p := range a.Content {
			b.WriteString(wrapAnsi(p, w))
			b.WriteString("\n\n")
		}
		if a.ActionText != "" {
			// On the site these switch views; as clickable links they lead to
			// the corresponding page — the discord action goes to the invite.
			url := club.Website
			if a.ActionView == "discord" {
				url = club.DiscordInvite
			}
			b.WriteString(link(s.actionLink(), url, a.ActionText))
			b.WriteString("\n\n")
		}
		if i < len(club.Announcements)-1 {
			b.WriteString(s.hint().Render(strings.Repeat("─", w)))
			b.WriteString("\n\n")
		}
	}
	return b.String()
}

func (m model) renderManPage(s styles, w int) string {
	var b strings.Builder

	// FOSSC(1)   Linux Reference Manual   FOSSC(1)
	hdrL, hdrM := "FOSSC(1)", "Linux Reference Manual"
	gap := max(w-lipgloss.Width(hdrL)*2-lipgloss.Width(hdrM), 2)
	b.WriteString(s.manHeader().Render(hdrL + strings.Repeat(" ", gap) + hdrM + strings.Repeat(" ", gap) + hdrL))
	b.WriteString("\n")
	b.WriteString(s.rule().Render(strings.Repeat("─", w)))
	b.WriteString("\n\n")

	section := func(name, body string) {
		b.WriteString(s.manSection().Render(name))
		b.WriteString("\n")
		b.WriteString(body)
		b.WriteString("\n\n")
	}

	section("NAME", "    "+bold("fossc")+" - FOSS Club KIET operational handbook and gateway")

	syn := "    " + bold("fossc") + " [" + bold("--discord") + "] [" + bold("--bootcamp") + "] [" +
		bold("--meeting") + " tuesday-5pm] [" + bold("--research-paper") + "] [" + italic("command") + "]"
	b.WriteString(s.manSection().Render("SYNOPSIS"))
	b.WriteString("\n")
	b.WriteString(s.manSynopsis(w - 4).Render(wrapAnsi(syn, w-6)))
	b.WriteString("\n\n")

	section("DESCRIPTION", "    "+wrapAnsi(
		bold("FOSS Club KIET")+" is the student-led software freedom collective at KIET Deemed To Be University. We run open labs in "+bold("Room H808 or H108")+" (confirm on Discord once temporarily) every day after class hours conclude.", w-4))

	b.WriteString(s.manSection().Render("MEETINGS & TIMINGS"))
	b.WriteString("\n")
	b.WriteString(renderScheduleTable(w))
	b.WriteString("\n\n")

	comm := "    Discord: " + link(s.actionLink(), club.DiscordInvite, club.DiscordInvite) +
		"\n" + wrapAnsi("    Venue: "+club.ClubRoom, w)
	section("COMMUNICATION", comm)

	return b.String()
}

// renderScheduleTable draws the MEETINGS & TIMINGS table with box characters,
// wrapping cell contents like an HTML table would.
func renderScheduleTable(w int) string {
	const pad = 1
	// Row width = 4 border chars + 3 cells (each +2 padding) → the column
	// total must leave 10 columns of chrome inside w.
	total := w - 10
	// Proportional column widths: cadence / time / venue.
	c1 := max(10, total*18/100)
	c3 := max(18, total*42/100)
	c2 := max(10, total-c1-c3)
	widths := []int{c1, c2, c3}

	header := []string{"Cadence", "Time & Schedule", "Venue / Details"}
	rows := make([][]string, 0, len(club.Schedule)+1)
	rows = append(rows, header)
	for _, r := range club.Schedule {
		rows = append(rows, []string{r.Cadence, r.Time, r.Venue})
	}

	renderRow := func(row []string) []string {
		wrapped := make([][]string, 3)
		nlines := 1
		for c, cell := range row {
			wrapped[c] = strings.Split(wrapAnsi(cell, widths[c]), "\n")
			nlines = max(nlines, len(wrapped[c]))
		}
		var lines []string
		for ln := range nlines {
			var cols []string
			for c := range 3 {
				text := ""
				if ln < len(wrapped[c]) {
					text = wrapped[c][ln]
				}
				cols = append(cols, " "+text+strings.Repeat(" ", max(0, widths[c]-lipgloss.Width(text)))+" ")
			}
			lines = append(lines, "│"+strings.Join(cols, "│")+"│")
		}
		return lines
	}

	sep := func(l, m, r string) string {
		parts := make([]string, 0, 3)
		for _, cw := range widths {
			parts = append(parts, strings.Repeat("─", cw+pad*2))
		}
		return l + strings.Join(parts, m) + r
	}

	var out strings.Builder
	out.WriteString(sep("┌", "┬", "┐") + "\n")
	for i, row := range rows {
		for _, ln := range renderRow(row) {
			out.WriteString(ln + "\n")
		}
		if i < len(rows)-1 {
			out.WriteString(sep("├", "┼", "┤") + "\n")
		}
	}
	out.WriteString(sep("└", "┴", "┘"))
	return out.String()
}

func (m model) renderBootcamp(s styles, w int) string {
	var b strings.Builder
	b.WriteString(banner(s, w, viewDefs[3].tag, "September 9, 2026"))
	b.WriteString("\n\n")

	b.WriteString(s.headline().Render("Linux Basics & Open-Source History Bootcamp (September 9)"))
	b.WriteString("\n")
	b.WriteString(s.byline().Render("Venue: " + bold("Room H106, KIET") + " • Date: " + bold("September 9, 2026 (17:00 - 20:30)")))
	b.WriteString("\n")
	b.WriteString(s.rule().Render(strings.Repeat("─", min(w, 60))))
	b.WriteString("\n\n")

	b.WriteString(wrapAnsi(club.BootcampPitch, w))
	b.WriteString("\n\n")

	b.WriteString(s.headline().Render("Curriculum & Workshop Modules"))
	b.WriteString("\n\n")
	boxW := w - 2
	for _, mod := range club.BootcampModules {
		b.WriteString(s.moduleBox(boxW).Render(
			bold(mod.Name) + "\n" + s.byline().Render(wrapAnsi(mod.Description, boxW-2))))
		b.WriteString("\n\n")
	}

	b.WriteString(s.prereqBox(boxW).Render(wrapAnsi(club.BootcampPrereqs, boxW-2)))
	b.WriteString("\n\n")

	b.WriteString(s.chip().Render("[y] Copy Bootcamp Announcement Text"))
	b.WriteString("\n")
	return b.String()
}

func (m model) renderMembers(s styles, w int) string {
	var b strings.Builder
	b.WriteString(banner(s, w, "MEMBERS", club.MemberEdition))
	b.WriteString("\n\n")

	b.WriteString(s.headline().Render("Members"))
	b.WriteString("\n")
	b.WriteString(s.byline().Render(club.MembersBlurb))
	b.WriteString("\n\n")

	boxW := w - 2
	for _, mem := range club.Members {
		b.WriteString(s.moduleBox(boxW).Render(
			bold(mem.Name) + "\n" +
				s.byline().Render("Discord: @"+mem.Discord) + "\n" +
				"GitHub: " + link(s.actionLink(), "https://github.com/"+mem.GitHub, "@"+mem.GitHub)))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m model) renderDiscord(s styles, w int) string {
	var b strings.Builder
	b.WriteString(banner(s, w, "COMMUNITY // DISCORD HUB", "Real-time Chat"))
	b.WriteString("\n\n")

	b.WriteString(s.headline().Render("FOSS Club KIET Discord Community"))
	b.WriteString("\n")
	b.WriteString(s.byline().Render("Connect with 200+ fellow KIET Linux enthusiasts and developers."))
	b.WriteString("\n\n")

	boxW := w - 2
	inner := boxW - 2
	content := s.headline().Render("Join the Official Server") + "\n" +
		wrapAnsi(club.DiscordBlurb, inner) + "\n\n" +
		strings.Repeat(" ", max(0, (inner-lipgloss.Width("discord.gg/JK272Ef8Pm ↗"))/2)) +
		link(s.discordChip(), club.DiscordInvite, "discord.gg/JK272Ef8Pm ↗")
	b.WriteString(s.prereqBox(boxW).Render(content))
	b.WriteString("\n\n")

	b.WriteString(s.hint().Render("Press o to open the invite in your browser."))
	b.WriteString("\n")
	return b.String()
}

// ---- small text helpers -----------------------------------------------------

// link wraps styled text in an OSC 8 hyperlink, so terminals that support
// clickable links (iTerm2, kitty, WezTerm, GNOME Terminal, …) open url when
// the text is clicked. Unsupported terminals just show the styled text.
func link(st lipgloss.Style, url, text string) string {
	return ansi.SetHyperlink(url) + st.Render(text) + ansi.ResetHyperlink()
}

func wrapAnsi(s string, w int) string {
	if w <= 0 {
		return s
	}
	return lipgloss.NewStyle().Width(w).Render(s)
}

func bold(s string) string  { return lipgloss.NewStyle().Bold(true).Render(s) }
func italic(s string) string { return lipgloss.NewStyle().Italic(true).Render(s) }

func rightAlign(s string, w int) string {
	return strings.Repeat(" ", max(w-lipgloss.Width(s), 0)) + s
}
