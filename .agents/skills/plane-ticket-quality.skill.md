# Skill: High-Quality Plane Ticket Writing

## Purpose

Create high-signal Plane tickets from commits/PRs with clear technical and
business impact.

## Mandatory Rules

1. Never write dates in the ticket title or description.
2. Never write vague summaries like "fix bug" without explaining:
   - what failed
   - who/what was impacted
   - why it happened (if known)
   - how it was fixed
3. If context is incomplete, investigate first. Do not publish low-context
   tickets.
4. Always include business impact, not only technical details.

## Required Investigation Before Writing

Collect context in this order when available:

1. PR title and description
2. PR/commit diff (files changed, meaningful code paths)
3. Commit titles and commit messages
4. Related docs/design notes
5. Validation evidence (tests/checks/runbooks)

If a source is missing, say that explicitly and continue with remaining
evidence.

## Ticket Structure (Use Exactly)

### Title

- Specific outcome, not process wording.
- Include domain and affected area.
- Avoid generic prefixes unless the team standard requires them.

### Original Problem

- Describe the real failure mode.
- Include user/system symptom and scope.

### Previous State

- Explain how the system behaved before the change.
- Mention constraints, gaps, or risks in prior implementation.

### Work Done

- Summarize the implemented changes by component.
- Mention key files/services touched when useful.
- Link behavior changes to the problem they resolve.

### Result

- Describe post-change behavior in production terms.
- Include business/user outcome and risk reduction.

### Business Impact

- Explicitly state impact on reliability, support load, developer velocity,
  cost, security, or conversion/retention.

### Evidence

- List commit subjects and PR references used to derive the ticket.

## Quality Bar (Must Pass)

1. A reader can understand the issue without opening code.
2. Problem statement is concrete and falsifiable.
3. Result section states measurable or observable effect.
4. Business impact is present and non-generic.
5. No dates appear anywhere in the ticket body/title.
6. No placeholder text remains.

## Guardrails

- Do not invent root causes. If unknown, write: "Root cause under investigation;
  observed behavior and mitigation documented."
- Do not claim business metrics you cannot support.
- Prefer concise, dense writing over long narrative.

## Output Template

Title: <specific outcome>

Original Problem:
<concrete failure and symptoms>

Previous State:
<how things worked before and why that was insufficient>

Work Done:

- <change 1>
- <change 2>
- <change 3>

Result:
<post-change system behavior>

Business Impact:

- <impact 1>
- <impact 2>

Evidence:

- PR: <link or id>
- Commits:
  - <sha> <subject>
  - <sha> <subject>

## Suggested Invocation

"Use skill plane-ticket-quality to draft Plane issues from these commits/PRs. If
context is weak, inspect diffs and PR descriptions first, then write tickets
with business impact and no dates."
