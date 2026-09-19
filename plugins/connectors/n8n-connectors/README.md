# Community connectors (n8n import)

Nine ready-made connections — Slack, Telegram, Discord, Notion, GitHub, Jira,
Trello, Asana, Mattermost — converted from n8n community nodes.

## What this pack does

It gives Mockarty **pre-filled connection presets**: the endpoint shape, the
documentation link and the fields each service needs. Install it, point a preset
at your workspace, and the integration is configured in a couple of clicks
instead of by hand.

## What it does not do

It does **not** run n8n workflows, and installing it does not make Mockarty an
n8n runtime. A preset is a connection, not an execution engine — if you need a
service's behaviour reproduced, mock it or write a test against it.

## Credentials

Every credential field in this pack's settings holds a **reference** to an entry
in your secret store, never the secret itself. The token stays where your
security team put it; the extension only knows where to ask for it.

Configure them under the extension's settings after installing:

| Field | Holds |
|-------|-------|
| Slack bot token | a reference to the Slack token |
| Telegram bot token | a reference to the Telegram token |
| GitHub token | a reference to a GitHub personal access token |
| Jira API token | a reference to the Jira token |
| Notion integration token | a reference to the Notion token |
| Default connection | an existing connection to reuse |

## Installing

```bash
mockarty-cli plugin install community.n8n-connectors-1.0.0.zip
mockarty-cli plugin enable community.n8n-connectors
```
