# Publishing fossc

How a release gets from a git tag to every package manager.

## One-time setup (already done for v1.0.0)

| Thing | Where |
| --- | --- |
| Source repo | `github.com/VanshSahay/fossc` |
| Homebrew tap | `github.com/VanshSahay/homebrew-tap` (goreleaser pushes `Formula/fossc.rb`) |
| apt repo | `https://vanshsahay.github.io/fossc/` (gh-pages branch, built by aptly) |

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

That's it. The `release` workflow then:

1. Builds binaries (linux/darwin/windows × amd64/arm64) with goreleaser,
   `-s -w -X main.version=<tag>`.
2. Publishes the GitHub Release with tarballs/zip + checksums.
3. Pushes `Formula/fossc.rb` to `VanshSahay/homebrew-tap`.
4. Builds `deb`, `rpm`, and `archlinux` packages via nfpm (man page, README and
   LICENSE included).

The `apt-repo` workflow then fires on the published release: downloads the
`.deb`, adds it to the aptly repo `fossc-stable`, signs and publishes the
`stable` distribution, and force-deploys to `gh-pages`.

## Verify after a release

```sh
gh run watch -R VanshSahay/fossc                                   # release
gh run watch -R VanshSahay/fossc $(gh run list -R VanshSahay/fossc -w apt-repo -L1 --json databaseId -q '.[0].databaseId')
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
