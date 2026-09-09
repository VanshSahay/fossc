# fossc

> FOSS Club KIET operational handbook and gateway — the terminal gazette.

`fossc` is a Bubble Tea TUI that mirrors [fossclubkiet.org](https://www.fossclubkiet.org/)
as a full-screen newspaper: masthead, sidebar widgets (nav, hours & venue,
quote of the day, Discord CTA) and editorial views — front page, weekly
edition, `man fossc(1)`, the bootcamp dossier, the member roster and the
Discord hub.

Content is ported verbatim from the site's source (`data/announcements/*.json`,
`app/GazetteClient.tsx`), and the **quote of the day uses the site's exact
golden-ratio picker** over the full 360-quote corpus — so `fossc --quote`
matches what the website renders today.

## Install

<details open>
<summary>macOS — Homebrew</summary>

```sh
brew install VanshSahay/tap/fossc
```

</details>

<details>
<summary>Debian / Ubuntu — apt repo (self-hosted)</summary>

```sh
curl -fsSL https://vanshsahay.github.io/fossc/fossc-archive-keyring.asc \
  | sudo tee /usr/share/keyrings/fossc-archive-keyring.asc >/dev/null
echo "deb [signed-by=/usr/share/keyrings/fossc-archive-keyring.asc] https://vanshsahay.github.io/fossc stable main" \
  | sudo tee /etc/apt/sources.list.d/fossc.list
sudo apt update && sudo apt install fossc
```

</details>

<details>
<summary>Arch — PKGBUILD</summary>

AUR submissions are paused upstream, so install from the in-repo PKGBUILD:

```sh
git clone https://github.com/VanshSahay/fossc && cd fossc/packaging/AUR
updpkgsums && makepkg -si
```

</details>

<details>
<summary>RPM / other Linux — from a release</summary>

Grab the `.rpm` (or a tarball for any linux/darwin/windows machine) from
[Releases](https://github.com/VanshSahay/fossc/releases), or:

```sh
sudo apt install ./fossc_*_amd64.deb    # from a release, no repo
```

</details>

<details>
<summary>Go</summary>

```sh
go install github.com/VanshSahay/fossc@latest
```

Requires Go 1.25+. Built with [bubbletea v1.3.10](https://github.com/charmbracelet/bubbletea),
[bubbles v1.0.0](https://github.com/charmbracelet/bubbles) and
[lipgloss v1.1.0](https://github.com/charmbracelet/lipgloss) — check
`fossc --version` for the exact build.

</details>

## Usage

```sh
fossc                 # the full terminal gazette
fossc bootcamp        # open the gazette straight on a view
fossc --discord       # open the Discord invite in your browser
fossc --meeting tuesday-5pm   # schedule + countdown to the next sync
```

### Gateway flags

| Flag | What it does |
| --- | --- |
| `--discord` | Print + open the Discord invite |
| `--bootcamp` | Print the full bootcamp dossier |
| `--meeting tuesday-5pm` | Print the meetings table and countdown (also `daily`, `paper`) |
| `--research-paper` | Print the paper reading circle notice |
| `--quote` | Print today's quote of the day |
| `--members` | Print the member roster |
| `--version` | Version + linked bubbletea version |

Positional commands `frontpage weekly man bootcamp members` start the TUI on
that view; `discord` acts as a gateway and opens the invite. Piping `fossc`
prints the front page as plain text.

### Keys inside the gazette

| Key | Action |
| --- | --- |
| `1`–`6` | Jump to a view |
| `←`/`→`, `h`/`l`, `tab` | Switch views |
| `↑`/`↓`, `k`/`j`, `pgup`/`pgdn`, mouse wheel | Scroll the articles |
| `o` | Open the Discord invite (`O` opens the website) |
| `y` | Copy the bootcamp announcement text |
| `d` | Toggle light/dark (follows the terminal background on start) |
| `q` | Quit |

All URLs — action links, Discord invites, GitHub profiles, the footer site
link — are OSC 8 hyperlinks, so clicking them opens your browser in terminals
that support clickable links (iTerm2, kitty, WezTerm, VS Code, GNOME
Terminal, …). macOS's default Terminal.app doesn't support OSC 8; there, use
the `o` key instead.

## Layout

```
internal/club/    site content: announcements, bootcamp, members, schedule,
                  quotes corpus + the golden-ratio daily picker
internal/tui/     bubbletea model, lipgloss styles, view renderers
main.go           gateway flags, commands, non-TTY fallback
```

No backgrounds are painted — theming is foreground-only, so the gazette sits
on whatever your terminal already looks like.
