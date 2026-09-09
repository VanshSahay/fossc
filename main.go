// fossc - FOSS Club KIET operational handbook and gateway.
//
// Run with no arguments, fossc launches a full terminal gazette mirroring
// https://www.fossclubkiet.org/ — masthead, sidebar, articles and all. Flags
// act as a gateway: they print (and where sensible, open) the thing you asked
// for without entering the TUI.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/VanshSahay/fossc/internal/club"
	"github.com/VanshSahay/fossc/internal/tui"
)

// version is overridden at release time via -ldflags "-X main.version=…";
// source builds fall back to this default.
var version = "1.0.0"

const usage = `fossc - FOSS Club KIET operational handbook and gateway

Usage:
  fossc [command]
  fossc [flags]

Commands (open the terminal gazette on a view):
  frontpage          Front page (default)
  weekly             Weekly edition
  man                man fossc(1) — meetings, venue, communication
  bootcamp           Linux Basics & Open-Source History Bootcamp dossier
  members            Member roster, Academic Year 2026-27
  discord            Open the Discord invite in your browser

Flags:
  --discord          Open the Discord invite (https://discord.gg/JK272Ef8Pm)
  --bootcamp         Print the full bootcamp dossier
  --meeting <when>   Print the schedule; when = tuesday-5pm | daily | paper
  --research-paper   Print the research paper reading circle notice
  --quote            Print today's quote of the day
  --members          Print the member roster
  --version          Print version and exit
  -h, --help         This help

Keys inside the gazette:
  1-6 jump to view · ←/→ switch · ↑/↓ scroll · o open discord · y copy
  announcement · d toggle light/dark · q quit

FOSS Club KIET — Room H808 or H108 (confirm on Discord once temporarily),
CSE-AI / AI&ML Dept, KIET. Open labs daily after classes.`

func main() {
	var (
		fDiscord   = flag.Bool("discord", false, "open the Discord invite")
		fBootcamp  = flag.Bool("bootcamp", false, "print the bootcamp dossier")
		fMeeting   = flag.String("meeting", "", "print the schedule: tuesday-5pm | daily | paper")
		fPaper     = flag.Bool("research-paper", false, "print the research paper circle notice")
		fQuote     = flag.Bool("quote", false, "print today's quote of the day")
		fMembers   = flag.Bool("members", false, "print the member roster")
		fVersion   = flag.Bool("version", false, "print version")
	)
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	if *fVersion {
		fmt.Printf("fossc %s (bubbletea %s)\n", version, teaVersion)
		return
	}

	// Gateway flags print their payload and exit; several may be combined.
	handled := false
	if *fDiscord {
		handled = true
		fmt.Println("Discord: " + club.DiscordInvite)
		if err := openBrowser(club.DiscordInvite); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't open a browser:", err)
		}
	}
	if *fBootcamp {
		handled = true
		printBootcamp()
	}
	if *fMeeting != "" {
		handled = true
		if err := printMeeting(*fMeeting); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
	if *fPaper {
		handled = true
		printPaper()
	}
	if *fQuote {
		handled = true
		q := club.QuoteOfTheDay()
		fmt.Printf("“%s”\n  — %s\n", q.Text, q.Author)
	}
	if *fMembers {
		handled = true
		printMembers()
	}
	if handled {
		return
	}

	// Positional command: either a gateway (discord) or a TUI start view.
	startView := "frontpage"
	if cmd := flag.Arg(0); cmd != "" {
		switch cmd {
		case "discord":
			fmt.Println("Discord: " + club.DiscordInvite)
			if err := openBrowser(club.DiscordInvite); err != nil {
				fmt.Fprintln(os.Stderr, "couldn't open a browser:", err)
			}
			return
		case "help":
			fmt.Println(usage)
			return
		case "frontpage", "weekly", "man", "manpage", "bootcamp", "members":
			startView = cmd
		default:
			fmt.Fprintf(os.Stderr, "fossc: unknown command %q — try 'fossc help'\n", cmd)
			os.Exit(2)
		}
	}

	// No TTY on stdout (pipe, redirect): print the front page as plain text
	// instead of hanging on a full-screen program.
	if !isTTY(os.Stdout) {
		printFrontPage()
		return
	}

	p := tea.NewProgram(tui.New(startView), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "fossc:", err)
		os.Exit(1)
	}
}

// teaVersion resolves the linked bubbletea version via build info.
var teaVersion = func() string {
	for _, dep := range moduleDeps() {
		if dep.Name == "github.com/charmbracelet/bubbletea" {
			return dep.Version
		}
	}
	return "unknown"
}()

func isTTY(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func openBrowser(url string) error {
	var c *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		c = exec.Command("open", url)
	case "windows":
		c = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		c = exec.Command("xdg-open", url)
	}
	return c.Run()
}

// ---- plain-text payloads -----------------------------------------------------

func printFrontPage() {
	fmt.Println("FOSS CLUB KIET")
	fmt.Println(club.MastheadNote)
	fmt.Println()
	for _, a := range club.Announcements {
		fmt.Println("== " + a.Category + " ==")
		fmt.Println(a.Title)
		byline := "[" + a.Date + "]"
		if a.Author != "" {
			byline += " By " + a.Author
		}
		if a.VenueOrDetails != "" {
			byline += " • " + a.VenueOrDetails
		}
		fmt.Println(byline)
		fmt.Println()
		for _, p := range a.Content {
			fmt.Println(wrap(p, 78))
			fmt.Println()
		}
		if a.ActionText != "" {
			fmt.Println(a.ActionText)
			fmt.Println()
		}
	}
	fmt.Println("Discord: " + club.DiscordInvite)
	fmt.Println("Web:     " + club.Website + "  (run `fossc` in a terminal for the full gazette)")
}

func printBootcamp() {
	fmt.Println("EVENT DOSSIER // LINUX BASICS & OPEN-SOURCE HISTORY BOOTCAMP")
	fmt.Println("Venue: Room H106, KIET • Date: September 9, 2026 (17:00 - 20:30)")
	fmt.Println()
	fmt.Println(wrap(club.BootcampPitch, 78))
	fmt.Println()
	fmt.Println("Curriculum & Workshop Modules")
	for _, m := range club.BootcampModules {
		fmt.Println("  " + m.Name)
		fmt.Println("    " + wrap(m.Description, 74))
	}
	fmt.Println()
	fmt.Println(wrap(club.BootcampPrereqs, 78))
	fmt.Println()
	fmt.Println("Discord: " + club.DiscordInvite)
}

func printPaper() {
	a := club.Announcements[2]
	fmt.Println("== " + a.Category + " ==")
	fmt.Println(a.Title)
	fmt.Println("[" + a.Date + "] • " + a.VenueOrDetails)
	fmt.Println()
	for _, p := range a.Content {
		fmt.Println(wrap(p, 78))
		fmt.Println()
	}
	fmt.Println("Propose a paper in the Discord #paper-reading channel:")
	fmt.Println(club.DiscordInvite)
}

func printMembers() {
	fmt.Println("MEMBERS — " + club.MemberEdition)
	fmt.Println(club.MembersBlurb)
	fmt.Println()
	for _, m := range club.Members {
		fmt.Println("  " + m.Name)
		fmt.Println("    Discord: @" + m.Discord)
		fmt.Println("    GitHub:  https://github.com/" + m.GitHub)
	}
}

func printMeeting(when string) error {
	fmt.Println("MEETINGS & TIMINGS")
	for _, r := range club.Schedule {
		fmt.Printf("  %-13s %-26s %s\n", r.Cadence, r.Time, r.Venue)
	}
	fmt.Println()

	switch when {
	case "tuesday-5pm", "tuesday", "weekly", "sync":
		next := club.NextWeeklySync(time.Now())
		in := club.UntilString(time.Until(next))
		fmt.Printf("Next weekly sync: %s • 5:00 PM IST (in %s)\n", next.Format("Monday, January 2, 2006"), in)
		fmt.Println("Venue: " + club.ClubRoom)
	case "daily":
		fmt.Println("Daily open hours: " + club.DailyHours)
	case "paper", "research":
		fmt.Println("Paper reading: " + club.PaperReading)
		fmt.Println("Picks are posted in Discord #paper-reading.")
	default:
		return fmt.Errorf("unknown meeting %q — valid: tuesday-5pm, daily, paper", when)
	}
	fmt.Println("Confirm the room on Discord: " + club.DiscordInvite)
	return nil
}
