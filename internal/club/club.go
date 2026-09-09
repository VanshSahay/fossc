// Package club holds the FOSS Club KIET operational data mirrored from the
// official website, https://www.fossclubkiet.org/ (a Next.js "gazette").
// Content here tracks data/announcements/*.json and app/GazetteClient.tsx.
package club

// DiscordInvite is the club's official Discord gateway.
const DiscordInvite = "https://discord.gg/JK272Ef8Pm"

// Website is the official site this tool mirrors.
const Website = "https://www.fossclubkiet.org/"

// Announcement mirrors one data/announcements/*.json entry.
type Announcement struct {
	ID             string
	Title          string
	Category       string
	CategoryDark   string // dark-mode color for the category kicker
	CategoryLight  string // light-mode color for the category kicker
	Author         string
	Date           string
	VenueOrDetails string
	Content        []string
	ActionText     string
	ActionView     string // view id the action link opens
}

// Announcements are the front-page articles, ordered like the site (order asc).
var Announcements = []Announcement{
	{
		ID:            "sync",
		Title:         "Weekly Club Sync Every Tuesday at 5:00 PM",
		Category:      "Club Notice & General Sync",
		CategoryLight: "#8a1f11",
		CategoryDark:  "#ff7777",
		Date:          "August 22, 2026",
		VenueOrDetails: "Venue: Room H808 or H108 (confirm on Discord once temporarily), CSE-AI / AI&ML Dept",
		Content: []string{
			"Welcome to all new and returning members of FOSS Club KIET! Our primary weekly meeting takes place every Tuesday at 5:00 PM in Room H808 or H108 (confirm on Discord once temporarily, CSE-AI / AI&ML Department). During these sessions, members share lightning talks on recent open source discoveries, review pull requests, and plan upcoming project sprints.",
			"Additionally, our space is open daily in Room H808 or H108 (confirm on Discord once temporarily) right after your classes conclude. Drop in anytime to work, discuss or ideate.",
		},
		ActionText: "[View Room Schedule]",
		ActionView: "manpage",
	},
	{
		ID:            "bootcamp",
		Title:         "Linux Basics & Open-Source History Bootcamp Scheduled for September 9",
		Category:      "Flagship Event Announcement",
		CategoryLight: "#1b4332",
		CategoryDark:  "#52c48a",
		Author:        "FOSS Club KIET Committee",
		Date:          "August 20, 2026",
		VenueOrDetails: "Date: September 9, 2026 (17:00 - 20:30) • Venue: Room H106, KIET",
		Content: []string{
			"FOSS Club KIET is proud to announce our upcoming Linux Basics & Open-Source History Bootcamp on September 9, 2026 (17:00 - 20:30) in Room H106, KIET. A beginner-friendly session where you won't just learn Linux, you'll play it. Explore an abandoned system in our browser adventure, uncover hidden files, unlock gates, and learn how to start contributing to open source this semester.",
			"Modules will cover Linux Fundamentals (core CLI navigation, filesystem hierarchy, and essential utilities), Open-Source History, and our browser-based Linux quest puzzle adventure. Absolute beginners welcome - bring a laptop with a modern web browser and a Google account. Nothing to install.",
		},
		ActionText: "[View Full Bootcamp Schedule & Modules →]",
		ActionView: "bootcamp",
	},
	{
		ID:            "paper-reading",
		Title:         "Monthly Research Paper Reading Circle",
		Category:      "Research Group",
		CategoryLight: "#194569",
		CategoryDark:  "#7bb0e0",
		Date:          "August 15, 2026",
		VenueOrDetails: "Cadence: Monthly • Announced via Discord #paper-reading",
		Content: []string{
			"Our monthly reading circle gathers to dissect seminal research papers in operating systems, distributed consensus algorithms (Raft, Paxos), kernel tracing with eBPF, and cryptographic protocols.",
			"Dates and paper selections are posted in advance on our Discord server. Anyone interested in presenting or discussing a paper is welcome to propose one in the #paper-reading channel.",
		},
		ActionText: "[Discuss in Discord #paper-reading →]",
		ActionView: "discord",
	},
}

// Module is one bootcamp workshop module.
type Module struct {
	Name        string
	Description string
}

// BootcampModules is the "Curriculum & Workshop Modules" dossier.
var BootcampModules = []Module{
	{
		Name:        "Module 1: Linux Fundamentals",
		Description: "Core CLI navigation, filesystem hierarchy exploration, and essential command-line utilities.",
	},
	{
		Name:        "Module 2: Open-Source History",
		Description: "The evolution of FOSS culture, open collaboration philosophy, and community onboarding for the semester.",
	},
	{
		Name:        "Module 3: Browser-based Linux Quest",
		Description: "A browser-based terminal puzzle adventure: explore an abandoned system, spot misleading clues, unlock gates, and mint a completion certificate.",
	},
}

// BootcampPitch is the dossier's opening paragraph.
const BootcampPitch = "A beginner-friendly session where you won't just learn Linux, you'll play it. Explore an abandoned system in our browser adventure, uncover hidden files, unlock gates, and learn how to start contributing to open source this semester."

// BootcampPrereqs is the prerequisites box.
const BootcampPrereqs = "Prerequisites: Absolute beginners welcome. Bring a laptop with a modern web browser and a Google account. Nothing to install."

// BootcampBlurb is the shareable announcement text ("Copy Bootcamp Announcement").
const BootcampBlurb = "Linux Basics & Open-Source History Bootcamp by FOSS CLUB KIET\nDate: September 9, 2026\nVenue: Room H106, KIET\nCurriculum: Linux Fundamentals, Open-Source History, Browser-based Linux Quest\nDiscord: https://discord.gg/JK272Ef8Pm"

// Member is one club member card.
type Member struct {
	Name    string
	Discord string
	GitHub  string
}

// Members is the "Academic Year 2026-27" roster.
var Members = []Member{
	{Name: "Deepak Anand", Discord: "arcceus", GitHub: "arcceus"},
	{Name: "Vansh Sahay", Discord: "vansh_sahay", GitHub: "VanshSahay"},
	{Name: "Nikhil", Discord: "badnikhil", GitHub: "badnikhil"},
}

// ScheduleRow is one row of the man-page MEETINGS & TIMINGS table.
type ScheduleRow struct {
	Cadence string
	Time    string
	Venue   string
}

// Schedule is the MEETINGS & TIMINGS table from man fossc(1).
var Schedule = []ScheduleRow{
	{"Daily Hours", "Everyday after classes", "Room H808 or H108 (confirm on Discord once temporarily)"},
	{"Weekly Sync", "Every Tuesday @ 5:00 PM", "General meetings & lightning talks"},
	{"Paper Group", "Monthly", "Systems & architecture review"},
}

// Sidebar hour/venue facts, mirroring the "HOURS & VENUE" widget.
const (
	DailyHours    = "Everyday after classes (Room H808 or H108, confirm on Discord once temporarily)"
	WeeklyMeetup  = "Tuesdays @ 5:00 PM IST"
	ClubRoom      = "Room H808 or H108 (confirm on Discord once temporarily), CSE-AI / AI&ML Dept, KIET"
	PaperReading  = "Monthly discussion date varies"
	MastheadNote  = "Published weekly by the FOSS Club at KIET Deemed to be University • Room H808, CSE-AI/AI&ML Dept"
	DiscordBlurb  = "Get immediate assistance with Linux installations, receive workshop announcements, collaborate on open source repositories, and participate in paper reading discussions."
	MembersBlurb  = "Details for the members of FOSS Club KIET."
	MemberEdition = "Academic Year 2026-27"
)
