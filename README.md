# Forkwatch

Discover meaningful patches hiding in GitHub forks.

Forkwatch analyzes forks of a repository to find changes that haven't been submitted as pull requests. It groups forks by the files they modify and highlights **convergence** — when multiple independent forks touch the same code, that's a strong signal something needs fixing upstream.

Inspired by the idea of [Respectful Open Source](https://nesbitt.io/2026/02/13/respectful-open-source.html): maintainers and users should be able to find meaningful work happening in forks without anyone having to push a PR.

## Install

### Homebrew

```
brew install stympy/tap/forkwatch
```

### Go

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

  WebinarGeek +1 -2 — Change gitspec faraday version
    -  spec.add_runtime_dependency "faraday", "~> 1.0"
    -  spec.add_runtime_dependency "faraday_middleware", "~> 1.0"
    +  spec.add_runtime_dependency "faraday", '>= 2.0'

  alexbndk +1 -2 — Update convertkit-ruby.gemspec
    -  spec.add_runtime_dependency "faraday", "~> 1.0"
    -  spec.add_runtime_dependency "faraday_middleware", "~> 1.0"
    +  spec.add_runtime_dependency "faraday", "~> 2.7.4"

  excid3 +1 -2 — Update faraday_middleware dependency to the latest
    -  spec.add_runtime_dependency "faraday", "~> 1.0"
    -  spec.add_runtime_dependency "faraday_middleware", "~> 1.0"
    +  spec.add_runtime_dependency "faraday", ">= 1.0", "< 3.0"
  ...

lib/convertkit/connection.rb (4 forks converge here)

  Most common change pattern:
     require "faraday"
    -require "faraday_middleware"
     require "json"
  WebinarGeek, chaiandconversation, alexbndk, excid3
...
```

11 independent forks all updating the gemspec to fix the faraday dependency — and now you can see exactly what each one changed. When forks make identical changes (like removing the `faraday_middleware` require), they're grouped together automatically.

## How it works

1. Fetches forks sorted by most recently pushed
2. Compares each fork's default branch to upstream
3. Filters out noise: bot commits (dependabot, renovate), lock file changes, CI config tweaks
4. Groups forks by the files they modify
5. Highlights convergence — files modified by multiple independent forks
6. Shows the actual patches — when multiple forks make identical changes, they're grouped together; unique changes are shown inline with their diffs

## Rate limits

Forkwatch uses one GitHub API call per fork analyzed plus a few for setup. It monitors rate limits and stops gracefully before hitting 403s. With the default `--limit 100`, a typical run uses ~100 API calls out of GitHub's 5,000/hour allowance.
