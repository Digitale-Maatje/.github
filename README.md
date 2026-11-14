# .github

Shared release notes and PR labels for Digitale Maatje repositories. Each repository holds one line of configuration and two small workflows; the rules live here.

## What a repository gets

- **PR labels** from [`.github/labeler.yml`](.github/labeler.yml): the version (`version: major|minor|patch`) and type (`feature`, `bug`, `dependencies`, …) from the conventional-commit title or the branch prefix, plus labels for changed files and PR size.
- **A draft release** after every merge to `main`, from [`.github/release-drafter.yml`](.github/release-drafter.yml): its version follows the merged PRs, and its notes group them by label.

Publishing the draft stays a person's step. In several repositories, publishing a release is what applies production.

### Release notes

Each merged PR is listed once, under the first category that matches. The conventional-commit prefix comes off its title.

| Category | Labels (or title) |
|---|---|
| 🚨 Breaking Changes (major) | `version: major`, or a `type!:` title |
| ✨ New Features (minor) | `version: minor`, `feature` (a `feat:` title bumps minor too) |
| 🐛 Bug Fixes | `bug`, `fix` |
| ⬆️ Dependency Updates | `dependencies` |
| 📦 Examples | `examples` |
| 🏗️ Infrastructure | `infrastructure` |
| 📚 Documentation | `documentation` |
| 🔧 Internal Changes | everything else |

Anything else is a patch. Left out of the notes, and of version resolution: `skip-changelog` and `internal`, labels a person sets on purpose.

## Using it

`.github/release-drafter.yml`:

```yaml
_extends: Digitale-Maatje/.github:release-drafter.yml@1
```

A repository whose deploy promotes a pre-release also sets `prerelease: true` and `latest: false`.

`.github/workflows/label-pull-request.yml`:

```yaml
name: Pull Request labels
on:
  pull_request:
    types: [opened, reopened, synchronize, edited]
concurrency:
  group: pr-labels-${{ github.event.pull_request.number }}
  cancel-in-progress: true
permissions: {}
jobs:
  labels:
    uses: Digitale-Maatje/.github/.github/workflows/labeler.yml@1
    permissions:
      contents: read
      pull-requests: write
```

`.github/workflows/draft-release.yml`:

```yaml
name: Draft Release
on:
  push:
    branches: [main]
concurrency:
  group: draft-release
  cancel-in-progress: true
permissions: {}
jobs:
  draft:
    uses: Digitale-Maatje/.github/.github/workflows/release-drafter.yml@1
    permissions:
      contents: write
      pull-requests: read
```

The labeler reads its rules from this repository, at the same `1` tag, so a repository never copies `labeler.yml`.

## Changing it

1. A PR here is labelled with the rules it proposes, and checked: the Release Drafter base against its schema, every labeler pattern compiled as RE2 (the labeler ignores a pattern it cannot compile, without a word), and the workflows against theirs.
2. Every merge to `main` updates this repository's own draft release.
3. Publishing the draft tags `X.Y.Z` and moves `1`, `1.Y` and `latest` to its commit. A published pre-release moves nothing.

Repositories pin `@1` and pick up every `1.*` release. A breaking change (a `feat!:` PR) becomes `2`. Its release also moves the `ref` default in [`workflows/labeler.yml`](.github/workflows/labeler.yml) to `refs/tags/2`, and each repository adopts it by moving its `@1` references.

The rules match Luscii's base ([Luscii/.github](https://github.com/Luscii/.github)), comments apart, so a fix found in one carries over to the other.

## Version tags

Releases are tagged `X.Y.Z`, without a `v`: `1.4.2`, not `v1.4.2`. The release *name* keeps the `v`. The one exception is a Go module: Go only resolves `vX.Y.Z` tags, so it sets `tag-template: "v$RESOLVED_VERSION"`.

## Why this repository is public

A workflow's `GITHUB_TOKEN` can read its own repository and public ones. A private base would need an App token in every repository that extends it or calls its workflows. Nothing here is sensitive.
