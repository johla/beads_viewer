// Package agents provides AGENTS.md integration for AI coding agents.
// It handles detection, content injection, and preference storage for
// automatically adding beads_viewer usage instructions to agent configuration files.
package agents

import (
	"regexp"
	"strings"
)

// BlurbVersion is the current version of the agent instructions blurb.
// Increment this when making breaking changes to the blurb format.
const BlurbVersion = 1

// BlurbStartMarker marks the beginning of injected agent instructions.
const BlurbStartMarker = "<!-- bv-agent-instructions-v1 -->"

// BlurbEndMarker marks the end of injected agent instructions.
const BlurbEndMarker = "<!-- end-bv-agent-instructions -->"

// AgentBlurb contains the instructions to be appended to AGENTS.md files.
// This content helps AI coding agents understand how to use beads_viewer
// for issue tracking and project management.
const AgentBlurb = `<!-- bv-agent-instructions-v1 -->

---

## Beads Workflow Integration

This project uses [beads_viewer](https://github.com/Dicklesworthstone/beads_viewer) for issue tracking. Issues are stored in ` + "`" + `.beads/` + "`" + ` and tracked in git.

### Essential Commands

` + "```" + `bash
# View issues (launches TUI - avoid in automated sessions)
bv

# CLI commands for agents (use these instead)
bd ready              # Show issues ready to work (no blockers)
bd list --status=open # All open issues
bd show <id>          # Full issue details with dependencies
bd create --title="..." --type=task --priority=2
bd update <id> --status=in_progress
bd close <id> --reason="Completed"
bd close <id1> <id2>  # Close multiple issues at once
bd sync               # Commit and push changes
` + "```" + `

### Workflow Pattern

1. **Start**: Run ` + "`" + `bd ready` + "`" + ` to find actionable work
2. **Claim**: Use ` + "`" + `bd update <id> --status=in_progress` + "`" + `
3. **Work**: Implement the task
4. **Complete**: Use ` + "`" + `bd close <id>` + "`" + `
5. **Sync**: Always run ` + "`" + `bd sync` + "`" + ` at session end

### Key Concepts

- **Dependencies**: Issues can block other issues. ` + "`" + `bd ready` + "`" + ` shows only unblocked work.
- **Priority**: P0=critical, P1=high, P2=medium, P3=low, P4=backlog (use numbers, not words)
- **Types**: task, bug, feature, epic, question, docs
- **Blocking**: ` + "`" + `bd dep add <issue> <depends-on>` + "`" + ` to add dependencies

### Session Protocol

**Before ending any session, run this checklist:**

` + "```" + `bash
git status              # Check what changed
git add <files>         # Stage code changes
bd sync                 # Commit beads changes
git commit -m "..."     # Commit code
bd sync                 # Commit any new beads changes
git push                # Push to remote
` + "```" + `

### Best Practices

- Check ` + "`" + `bd ready` + "`" + ` at session start to find available work
- Update status as you work (in_progress → closed)
- Create new issues with ` + "`" + `bd create` + "`" + ` when you discover tasks
- Use descriptive titles and set appropriate priority/type
- Always ` + "`" + `bd sync` + "`" + ` before ending session

<!-- end-bv-agent-instructions -->`

// SupportedAgentFiles lists the filenames that can contain agent instructions.
var SupportedAgentFiles = []string{
	"AGENTS.md",
	"CLAUDE.md",
	"agents.md",
	"claude.md",
}

// blurbVersionRegex extracts the version number from a blurb marker.
var blurbVersionRegex = regexp.MustCompile(`<!-- bv-agent-instructions-v(\d+) -->`)

// ContainsBlurb checks if the content already contains a beads_viewer agent blurb.
// Returns true if any version of the blurb marker is found.
func ContainsBlurb(content string) bool {
	return strings.Contains(content, "<!-- bv-agent-instructions-v")
}

// GetBlurbVersion extracts the version number from existing blurb content.
// Returns 0 if no blurb is found.
func GetBlurbVersion(content string) int {
	matches := blurbVersionRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return 0
	}
	// Parse version number
	var version int
	_, _ = strings.NewReader(matches[1]).Read(make([]byte, 1))
	if matches[1] == "1" {
		version = 1
	}
	// For future versions, add more cases or use strconv
	return version
}

// NeedsUpdate checks if the content has an older version of the blurb
// that should be updated to the current version.
func NeedsUpdate(content string) bool {
	if !ContainsBlurb(content) {
		return false
	}
	return GetBlurbVersion(content) < BlurbVersion
}

// AppendBlurb appends the agent blurb to the given content.
// It adds proper spacing before the blurb.
func AppendBlurb(content string) string {
	// Ensure content ends with newline
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	// Add extra newline for spacing
	content += "\n"
	content += AgentBlurb
	content += "\n"
	return content
}

// RemoveBlurb removes an existing blurb from the content.
// This is useful for updating to a new version.
func RemoveBlurb(content string) string {
	// Find start marker
	startIdx := strings.Index(content, "<!-- bv-agent-instructions-v")
	if startIdx == -1 {
		return content
	}

	// Find end marker
	endIdx := strings.Index(content, BlurbEndMarker)
	if endIdx == -1 {
		// Malformed blurb - just return as-is
		return content
	}
	endIdx += len(BlurbEndMarker)

	// Remove any trailing newlines after the end marker
	for endIdx < len(content) && (content[endIdx] == '\n' || content[endIdx] == '\r') {
		endIdx++
	}

	// Remove any leading newlines before the start marker
	for startIdx > 0 && (content[startIdx-1] == '\n' || content[startIdx-1] == '\r') {
		startIdx--
	}

	return content[:startIdx] + content[endIdx:]
}

// UpdateBlurb replaces an existing blurb with the current version.
func UpdateBlurb(content string) string {
	content = RemoveBlurb(content)
	return AppendBlurb(content)
}
