# Community connectors (n8n import)

Nine integration presets — Slack, Telegram, Discord, Notion, GitHub, Jira,
Trello, Asana, Mattermost — produced by converting real n8n community nodes:

```bash
mockarty-cli plugin import ./nodes --from n8n
```

What the converter does:

- a node matching a built-in Mockarty adapter (`github`, `jira`, …) reuses it;
- every other node becomes a `webhook_generic` preset, seeded with the API base
  URL when the node declares one (see `asana`'s `url_template` here);
- name, description, category and docs link carry over; `*Trigger` variants
  fold into their base node.

Install the pack and each service shows up as a one-click preset in
Integrations. This directory is itself the converter's showcase output — study
`plugin.json` to see exactly what your own nodes will produce.
