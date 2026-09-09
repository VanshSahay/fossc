# Publishing fossc

How a release gets from a git tag to every package manager.

## One-time setup (already done for v1.0.0)

| Thing | Where |
| --- | --- |
| Source repo | `github.com/VanshSahay/fossc` |
| Homebrew tap | `github.com/VanshSahay/homebrew-tap` (`scripts/build-formula.sh` pushes `Formula/fossc.rb`) |
| apt repo | `https://vanshsahay.github.io/fossc/` (gh-pages branch, built by aptly in the release workflow) |

GitHub secrets on `VanshSahay/fossc`:

| Secret | Purpose |
| --- | --- |
| `TAP_GITHUB_TOKEN` | PAT with `repo` scope, lets goreleaser push the formula to `homebrew-tap` |
| `GPG_PRIVATE_KEY` | Armored private key signing the apt repo (no passphrase) |
| `GPG_KEY_ID` | Long key ID of the above |

The matching **public** key is committed at `packaging/apt/fossc-archive-keyring.asc`
and served from the Pages root so users can add it to their keyrings.

## Cutting a release

```sh
# from a clean main
git tag v1.2.3
git push origin main v1.2.3
```

That's it. The single `release` workflow then:

1. Builds binaries (linux/darwin/windows × amd64/arm64) with goreleaser,
   `-s -w -X main.version=<tag>`.
2. Publishes the GitHub Release with tarballs/zip + checksums.
3. Generates `Formula/fossc.rb` from `checksums.txt` and pushes it to
   `VanshSahay/homebrew-tap` (goreleaser's `brews` is deprecated in favor of
   casks, which don't suit a pure CLI — hence the small generator script).
4. Builds `deb`, `rpm`, and `archlinux` packages via nfpm (man page, README and
   LICENSE included).
5. Adds the `.deb`s to the aptly repo `fossc-stable`, signs and publishes the
   `stable` distribution, and deploys `dists/` + `pool/` + the armored public
   key to `gh-pages`.

Note: the apt publish lives in the same workflow because releases authored
with `GITHUB_TOKEN` don't trigger other workflows.

## Verify after a release

```sh
gh run watch -R VanshSahay/fossc            # one run does release + tap + apt
curl -fsSL https://vanshsahay.github.io/fossc/dists/stable/Release | head
brew info VanshSahay/tap/fossc
```

## AUR

AUR account creation is temporarily paused upstream, so there is no
goreleaser `aurs` push. A ready-to-submit `packaging/AUR/PKGBUILD` ships in
the repo; when submissions reopen:

```sh
cd packaging/AUR
updpkgsums                      # real sha256 for the tag tarball
# web upload or `aurpublish fossc`
```

`fossc-bin` (precompiled, from the GitHub tarball) can follow the same way.

## Apt repo notes

- The signing key lives in `GNUPGHOME` only on the runner; the private key
  never leaves GitHub secrets + this machine's `/tmp`.
- `aptly publish snapshot -force-overwrite` + `keep_files: true` means older
  versions stay in `pool/`, so `apt upgrade` and downgrades both work.
- If Pages isn't enabled after the first apt-repo run:
  repo → Settings → Pages → Source: *Deploy from a branch* → `gh-pages / (root)`.
