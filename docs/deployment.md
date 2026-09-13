# Deployment

Team Impact Scorecard is a standalone desktop application. There is no server, no cloud, and no network component. Deployment means distributing a compiled binary to the user's machine. The user installs it once; all data is stored locally in a SQLite file.

---

## Deployment Model

- **Release unit:** one compiled binary per target platform — macOS (`.app` bundle), Windows (`.exe`), Linux (ELF binary).
- **Distribution:** GitHub Releases. Each release publishes pre-built binaries as downloadable assets.
- **Runtime dependencies:** none. The binary is self-contained. No Go runtime, no database server, no browser required.
- **Data storage:** a SQLite file in the OS user config directory (`os.UserConfigDir()/leadpulse/data.db`). The user owns their data file. It is not synced or backed up by the application.
- **Network access at runtime:** none. The application makes no outbound connections.

This model is appropriate because the PRD requires data privacy (no external data), a single-user tool (one Team Lead), and zero operational overhead for the user.

---

## Release Triggers and Versioning

- **Trigger:** a maintainer tags a commit with `v<major>.<minor>.<patch>` and pushes to `origin`. GitHub Actions picks up the tag and runs the release workflow.
- **Versioning:** Semantic Versioning (SemVer).
  - Patch: bug fixes, formula corrections, UI adjustments with no schema change.
  - Minor: new features, new screens, additive schema migrations.
  - Major: breaking schema changes that require a migration tool, or significant behavioral rewrites.
- **Ownership:** any maintainer with write access to the repository may tag a release. No automated release from `main`.
- **Required evidence before tagging:** CI is green on the commit being tagged (`go test -race ./...` passes on all three platforms).

---

## Environments

There is one environment: the user's local machine. There is no staging or production server.

For development, the application runs directly with `go run .` or `fyne run`. The data file used during development is the same SQLite file the built binary uses — developers should use a separate test data file to avoid corrupting real data.

**Platform targets:**
| Platform | Artifact | Build command |
|---|---|---|
| macOS (arm64 + amd64) | `.app` bundle | `fyne package -os darwin` |
| Windows (amd64) | `.exe` | `fyne package -os windows` |
| Linux (amd64) | binary | `fyne package -os linux` |

---

## Deployment Procedure

Release is performed by a maintainer. CI handles the builds.

**Prerequisites:**
- Go 1.21+ installed locally (for verification builds only; CI does the release builds)
- `gh` CLI authenticated to the repository
- All tests passing on the target commit

**Steps:**
1. Verify CI is green on the commit you intend to release: check GitHub Actions on `main`.
2. Confirm `go test -race ./...` passes locally.
3. Update `docs/roadmap.md` — move completed features to Done, add evidence links.
4. Tag the release:
   ```bash
   git tag v<major>.<minor>.<patch>
   git push origin v<major>.<minor>.<patch>
   ```
5. GitHub Actions runs the release workflow: builds binaries for all three platforms, creates a GitHub Release, and attaches the binaries as assets.
6. Verify the release on GitHub: confirm all three platform artifacts are attached and the release notes are correct.
7. If the release notes need editing, edit them directly on the GitHub Release page.

**Release notes format:**
- One-line summary of what changed.
- List of user-visible changes (not commit hashes).
- Known issues or upgrade notes if any.

---

## Rollback Procedure

There is no server to roll back. A user on a bad version downloads the previous release binary from GitHub Releases and replaces their current binary.

**Schema migrations:**
- All schema migrations in v1 must be **additive only** (add columns, add tables). No column drops, no renames, no data transforms.
- A user who downgrades to a previous binary may encounter unknown columns in their SQLite file. The application must handle unknown columns gracefully (ignore them on read; do not fail).
- If a destructive migration becomes unavoidable, it requires a major version bump and a migration tool shipped alongside the binary. This is deferred beyond v1.

**If a release is found to be broken:**
1. Create a GitHub Release marked as "Pre-release" or add a prominent note to the release description warning users not to download it.
2. Tag and release a patch version with the fix as quickly as possible.
3. There is no automated rollback — users must manually download the corrected binary.

---

## Configuration and Secrets

There are no secrets and no runtime configuration beyond the data file path. The data file path resolves automatically from `os.UserConfigDir()`. There is no `.env` file, no config file, and no environment variable required to run the application.

If a future version requires user-configurable settings (e.g., billability target), those settings are stored in the same SQLite database, not in a separate config file.

---

## Health Signals and Operational Risks

The application has no server-side health signals. The meaningful operational risks are:

| Risk | Signal | Action |
|---|---|---|
| SQLite file corruption | Application fails to open or crashes on start | User restores from a manual backup; we cannot provide automated recovery in v1 |
| Binary won't launch on target OS | User report | Reproduce on the target platform, patch, and re-release |
| Formula regression | User reports incorrect scores | Verify against PRD formulas, add a failing engine test, fix, release patch |

**Backup guidance (for users):** The data file at `os.UserConfigDir()/leadpulse/data.db` can be copied manually. The application does not back it up automatically in v1.
