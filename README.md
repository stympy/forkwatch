# Forkwatch

Discover meaningful patches hiding in GitHub forks.

Forkwatch analyzes forks of a repository to find changes that haven't been submitted as pull requests. It groups forks by the files they modify and highlights **convergence** — when multiple independent forks touch the same code, that's a strong signal something needs fixing upstream.

Inspired by the idea of [Respectful Open Source](https://radicle.xyz/2024/09/24/respectful-open-source.html): maintainers and users should be able to find meaningful work happening in forks without anyone having to push a PR.

## Install

```
go install github.com/stympy/forkwatch@latest
```

Or build from source:

```
git clone https://github.com/stympy/forkwatch.git
cd forkwatch
go build -o forkwatch .
```

## Prerequisites

Forkwatch uses the [GitHub CLI](https://cli.github.com/) for authentication. Install it and run `gh auth login` before using forkwatch.

## Usage

```
forkwatch analyze owner/repo [flags]
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--min-ahead` | 1 | Minimum commits ahead to consider |
| `--limit` | 100 | Max forks to analyze (sorted by most recently pushed) |
| `--json` | false | Output as JSON |

### Examples

```
# Analyze forks of a repository
forkwatch analyze expressjs/express

# Only show forks with 3+ commits ahead
forkwatch analyze expressjs/express --min-ahead 3

# Get JSON output for scripting
forkwatch analyze expressjs/express --json

# Analyze more forks (slower, uses more API calls)
forkwatch analyze expressjs/express --limit 500
```

## Example output

```
$ forkwatch analyze maximadeka/convertkit-ruby

maximadeka/convertkit-ruby
Forks: 46 total, 17 analyzed, 17 with meaningful changes

convertkit-ruby.gemspec (11 forks converge here)
  WebinarGeek          +1 -2 (2 commits ahead)
    Change gitspec faraday version
  roelbondoc           +3 -3 (1 commits ahead)
    Update convertkit-ruby.gemspec
  chaiandconversation  +1 -0 (4 commits ahead)
    Updating gem to use newer faraday to work with Rails 6
  excid3               +1 -2 (5 commits ahead)
    Update faraday_middleware dependency to the latest
  mikefogg             +13 -13 (1 commits ahead)
    Upgrading faraday
  ...and 6 more

lib/convertkit/connection.rb (4 forks converge here)
  WebinarGeek          +0 -1 (2 commits ahead)
    Change gitspec faraday version
  chaiandconversation  +0 -1 (4 commits ahead)
    Updating gem to use newer faraday to work with Rails 6
  alexbndk             +0 -1 (2 commits ahead)
    Update convertkit-ruby.gemspec
  excid3               +0 -1 (5 commits ahead)
    Update faraday_middleware dependency to the latest

lib/convertkit/client/tags.rb (3 forks converge here)
  ericalli             +4 -0 (3 commits ahead)
    Add form subscriptions endpoint
  jaswinder97          +6 -0 (2 commits ahead)
    Allow gemspec to use latest versions of dependencies
  jamesknelson         +11 -2 (3 commits ahead)
    add webhoks and remove tag support
...
```

11 independent forks all updating the gemspec to fix the faraday dependency — that's a clear signal the maintainer should bump it upstream.

## How it works

1. Fetches forks sorted by most recently pushed
2. Compares each fork's default branch to upstream
3. Filters out noise: bot commits (dependabot, renovate), lock file changes, CI config tweaks
4. Groups forks by the files they modify
5. Highlights convergence — files modified by multiple independent forks

## Rate limits

Forkwatch uses one GitHub API call per fork analyzed plus a few for setup. It monitors rate limits and stops gracefully before hitting 403s. With the default `--limit 100`, a typical run uses ~100 API calls out of GitHub's 5,000/hour allowance.
