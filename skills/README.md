# Mockarty agent skills

Ready-to-use **agent skills** for driving and extending Mockarty from an AI
coding assistant. Each skill is a self-contained `SKILL.md` (YAML front-matter +
Markdown body) — the format Claude Code loads natively, and plain enough that
Cursor and OpenCode read it as an instruction file too.

| Skill | Use it when you want to… |
|-------|--------------------------|
| [`mockarty-drive`](mockarty-drive/SKILL.md) | Drive Mockarty over MCP/REST — create mocks, run functional/load/fuzz suites, manage test cases, read reports — as an autonomous tester. |
| [`mockarty-plugin-author`](mockarty-plugin-author/SKILL.md) | Author, build, test and install a Mockarty **plugin** (mock kit, WASM faker/transformer/matcher, connector, content pack, UI panel, links, event/task types). |

## Install

### Claude Code
Copy a skill folder into your project's (or user's) skills directory — Claude
Code discovers it automatically:

```bash
mkdir -p .claude/skills
cp -r mockarty-drive .claude/skills/
# or user-wide:  cp -r mockarty-drive ~/.claude/skills/
```

Invoke it by name (`/mockarty-drive`) or just describe the task — the
`description:` front-matter lets the model pick it up on its own.

### Cursor
Cursor reads project rules from `.cursor/rules/`. Drop the skill body in as a
rule file:

```bash
mkdir -p .cursor/rules
cp mockarty-drive/SKILL.md .cursor/rules/mockarty-drive.mdc
```

(You can delete the YAML front-matter or keep it — Cursor ignores it.)

### OpenCode / other assistants
Any assistant that accepts a Markdown instruction/context file can use these:
point it at the `SKILL.md` (e.g. add it to your `AGENTS.md`, a system prompt, or
a context include). The body is written to stand alone.

## Connect the assistant to Mockarty

The `mockarty-drive` skill assumes an MCP connection to a running Mockarty. Add
the server to your assistant's MCP config (Claude Code example — `.mcp.json`):

```json
{
  "mcpServers": {
    "mockarty": {
      "type": "http",
      "url": "https://your-mockarty-host/mcp",
      "headers": { "Authorization": "Bearer mk_your_api_token" }
    }
  }
}
```

Get an API token from **Admin → API tokens** in the Mockarty UI. For
plain-REST-only assistants, the skills also show the equivalent `curl` calls.

## License

These skills are MIT — copy, adapt, and ship them with your own Mockarty
extensions.
