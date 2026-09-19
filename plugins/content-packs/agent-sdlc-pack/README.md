# Agent SDLC pack

Turns an empty namespace into one an autonomous tester can actually work in —
and, because every surface here is a normal Mockarty surface, a person works in
exactly the same place.

## What you get

| Surface | What appears |
|---------|--------------|
| Wiki | the knowledge an agent reads first: how testing works here, what it must never do alone, how to read a failed run |
| Issue tracker | a project shaped like the loop the knowledge describes |
| Test cases | the smoke suite an agent runs before anything else |
| Task types | `plugin-agent-triage` and `plugin-agent-cover`, runnable by a runner |
| Notifications | events for "asked for review", "filed a defect", "loop finished" |
| Links | defect → evidence, case → requirement |
| Wiki editor | a `plugin-agent-rule` macro for rules that must stand out |

## Why the knowledge pages matter

An agent is only as good as what it reads. The pages in this pack are written to
be *quoted*: short rules, a decision table for triage, and an explicit list of
things that require a person. If you change how your team works, edit those
pages — the agent follows them, not this README.

## Installing

```bash
mockarty-cli plugin install mockarty.agent-sdlc-pack-1.0.0.zip
mockarty-cli plugin enable mockarty.agent-sdlc-pack
```

Then instantiate the content packs into the namespace you want them in.

## Settings

- **Require human review before an agent closes an issue** — keeps a review step
  the agent cannot skip. On by default.
- **Label for agent-authored work** — applied to everything the agent creates,
  so a person can always tell where something came from.

## What it does not do

It ships no code and runs nothing by itself. The task types it declares are
executed by a runner you already have; without one they simply sit in the
catalogue, visible and unused.
