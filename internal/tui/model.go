package tui

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/VanshSahay/fossc/internal/club"
)

// Model is the gazette TUI state.
type model struct {
	w, h   int
	cur    int // index into viewDefs
	vp     viewport.Model
	dark   bool
	ready  bool
	status string // transient footer message
}

// New builds the initial model. dark mirrors the site's default: follow the
// terminal background, overridable with d (like the site's 🌙 toggle).
func New(startView string) model {
	m := model{
		dark: lipgloss.HasDarkBackground(),
		vp:   viewport.New(80, 24),
	}
	m.cur = viewIndex(startView)
	return m
}

func viewIndex(id string) int {
	for i, vd := range viewDefs {
		if vd.id == id {
			return i
		}
	}
	return 0
}

// ---- messages ----------------------------------------------------------------

type copiedMsg struct{}
type copyErrMsg struct{ err error }
type openErrMsg struct{ err error }
type statusClearMsg struct{}

func clearStatus() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return statusClearMsg{} })
}

// clipboardCmd picks the first available local clipboard tool.
func clipboardCmd(text string) *exec.Cmd {
	bin, arg := func() (string, []string) {
		switch runtime.GOOS {
		case "darwin":
			return "pbcopy", nil
		case "windows":
			return "cmd", []string{"/c", "clip"}
		default:
			if _, err := exec.LookPath("wl-copy"); err == nil {
				return "wl-copy", nil
			}
			return "xclip", []string{"-selection", "clipboard"}
		}
	}()
	c := exec.Command(bin, arg...)
	c.Stdin = strings.NewReader(text)
	return c
}

// browserCmd opens url in the user's default browser.
func browserCmd(url string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url)
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return exec.Command("xdg-open", url)
	}
}

// ---- tea.Model ---------------------------------------------------------------

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		_, _, mainW, bodyH := m.layout()
		m.vp.Width = mainW
		m.vp.Height = bodyH
		m.ready = true
		m.vp.SetContent(m.renderMain(mainW))

	case statusClearMsg:
		m.status = ""

	case copiedMsg:
		m.status = "[Copied to Clipboard!]"
		cmds = append(cmds, clearStatus())

	case copyErrMsg:
		m.status = "Clipboard unavailable: " + msg.err.Error()
		cmds = append(cmds, clearStatus())

	case openErrMsg:
		m.status = "Couldn't open browser — " + club.DiscordInvite
		cmds = append(cmds, clearStatus())

	case tea.MouseMsg:
		nv, cmd := m.vp.Update(msg)
		m.vp = nv
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "1", "2", "3", "4", "5", "6":
			m.jump(msg.String())
			cmds = append(cmds, m.rescroll())

		case "right", "l", "tab":
			m.cur = (m.cur + 1) % len(viewDefs)
			cmds = append(cmds, m.rescroll())

		case "left", "h", "shift+tab":
			m.cur = (m.cur - 1 + len(viewDefs)) % len(viewDefs)
			cmds = append(cmds, m.rescroll())

		case "d":
			m.dark = !m.dark
			cmds = append(cmds, m.rescroll())

		case "o":
			cmds = append(cmds, tea.ExecProcess(browserCmd(club.DiscordInvite),
				func(err error) tea.Msg { return openErrMsg{err} }))

		case "O":
			cmds = append(cmds, tea.ExecProcess(browserCmd(club.Website),
				func(err error) tea.Msg { return openErrMsg{err} }))

		case "y":
			cmds = append(cmds, tea.ExecProcess(clipboardCmd(club.BootcampBlurb),
				func(err error) tea.Msg {
					if err != nil {
						return copyErrMsg{err}
					}
					return copiedMsg{}
				}))

		default:
			nv, cmd := m.vp.Update(msg)
			m.vp = nv
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// jump switches to the view bound to a number key.
func (m *model) jump(key string) {
	for i := range viewDefs {
		if key == strconv.Itoa(i+1) {
			m.cur = i
			return
		}
	}
}

// rescroll re-renders the main document (theme/size/view changed) from the top,
// mirroring the site's scrollTo({top:0}) on navigation.
func (m *model) rescroll() tea.Cmd {
	_, _, mainW, bodyH := m.layout()
	m.vp.Width = mainW
	m.vp.Height = bodyH
	m.vp.SetContent(m.renderMain(mainW))
	m.vp.GotoTop()
	return nil
}
