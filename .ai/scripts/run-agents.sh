#!/usr/bin/env bash
# Run one Claude agent per task file in parallel, each in an isolated git worktree.
# Each agent creates a branch, makes the fix, commits, pushes, and opens a PR.
#
# Usage:
#   ./run-agents.sh                    # run all tasks
#   ./run-agents.sh fix-some-task      # run a single task by id (no .json)
#   MAX_AGENTS=5 ./run-agents.sh       # override parallelism
#
# Prerequisites: claude CLI, gh CLI, jq, git

set -euo pipefail

REPO_ROOT="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"
TASK_DIR="$REPO_ROOT/.ai/tasks"
PROMPT_FILE="$REPO_ROOT/.ai/prompts/worker.md"
WORKTREE_DIR="$REPO_ROOT/.git/worktrees-ai"
MAX_AGENTS="${MAX_AGENTS:-3}"

# ── Sanity checks ─────────────────────────────────────────────────────────────

for cmd in claude gh jq git; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "ERROR: '$cmd' is required but not found in PATH." >&2
    exit 1
  fi
done

if [[ ! -f "$PROMPT_FILE" ]]; then
  echo "ERROR: Worker prompt not found at $PROMPT_FILE" >&2
  exit 1
fi

mkdir -p "$WORKTREE_DIR"

# ── Agent runner ──────────────────────────────────────────────────────────────

run_agent() {
  local task_file="$1"
  local name branch worktree_path prompt exit_code

  name="$(jq -r '.id' "$task_file")"
  branch="ai/${name}"
  worktree_path="$WORKTREE_DIR/$name"

  echo "[${name}] Starting..."

  # Skip if PR already exists
  if gh pr view "$branch" &>/dev/null 2>&1; then
    echo "[${name}] PR already exists, skipping."
    return
  fi

  # Clean up stale worktree if the directory exists but isn't registered
  if [[ -d "$worktree_path" ]]; then
    git -C "$REPO_ROOT" worktree remove --force "$worktree_path" 2>/dev/null || rm -rf "$worktree_path"
  fi

  # Remove stale branch if it exists
  git -C "$REPO_ROOT" branch -D "$branch" 2>/dev/null || true

  git -C "$REPO_ROOT" worktree add -b "$branch" "$worktree_path"

  # Build the prompt: worker instructions + task JSON
  prompt="$(cat "$PROMPT_FILE")

---

## Your Task

$(cat "$task_file")"

  # Run the agent inside the worktree
  pushd "$worktree_path" >/dev/null
  exit_code=0
  claude --dangerously-skip-permissions -p "$prompt" || exit_code=$?
  popd >/dev/null

  if [[ $exit_code -ne 0 ]]; then
    echo "[${name}] Claude exited with code $exit_code, skipping commit." >&2
    git -C "$REPO_ROOT" worktree remove --force "$worktree_path" 2>/dev/null || true
    return
  fi

  # Commit if there are changes
  pushd "$worktree_path" >/dev/null
  git add -A
  if git diff --cached --quiet; then
    echo "[${name}] No changes made."
    popd >/dev/null
    git -C "$REPO_ROOT" worktree remove --force "$worktree_path" 2>/dev/null || true
    return
  fi

  git commit -m "refactor: $(jq -r '.title' "$task_file")

Pain score: $(jq -r '.pain_score' "$task_file")
Task: ${name}

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"

  git push -u origin "$branch"
  popd >/dev/null

  # Open PR
  gh pr create \
    --repo "$(gh repo view --json nameWithOwner -q .nameWithOwner)" \
    --head "$branch" \
    --base main \
    --title "$(jq -r '.title' "$task_file")" \
    --body "$(jq -r '"## Summary\n\n" + .description + "\n\n**Pain score:** " + (.pain_score | tostring) + "/10\n**Task id:** " + .id + "\n\n## Files\n\n" + (.files | map("- `" + . + "`") | join("\n")) + "\n\n---\n🤖 Generated with [Claude Code](https://claude.com/claude-code)"' "$task_file")" \
    --label "ai-refactor" 2>/dev/null || \
  gh pr create \
    --repo "$(gh repo view --json nameWithOwner -q .nameWithOwner)" \
    --head "$branch" \
    --base main \
    --title "$(jq -r '.title' "$task_file")" \
    --body "$(jq -r '.description' "$task_file")" 2>/dev/null || \
  echo "[${name}] WARNING: PR creation failed — branch $branch is pushed."

  echo "[${name}] Done."
}

# ── Task selection ────────────────────────────────────────────────────────────

if [[ $# -ge 1 ]]; then
  # Run a single named task
  task_file="$TASK_DIR/${1}.json"
  if [[ ! -f "$task_file" ]]; then
    echo "ERROR: Task file not found: $task_file" >&2
    exit 1
  fi
  run_agent "$task_file"
  exit 0
fi

# Run all tasks, capped at MAX_AGENTS in parallel
tasks=("$TASK_DIR"/*.json)
total="${#tasks[@]}"
echo "Found $total tasks. Running up to $MAX_AGENTS agents in parallel."

count=0
for task_file in "${tasks[@]}"; do
  run_agent "$task_file" &
  count=$((count + 1))
  if [[ $count -ge $MAX_AGENTS ]]; then
    wait
    count=0
  fi
done

wait
echo "All agents finished."
