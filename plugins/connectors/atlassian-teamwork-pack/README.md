# Atlassian teamwork pack

Everything a team moving off Jira and Confluence needs on the first day, in one
package. It is entirely declarative — no code from this package runs on your
server.

## What you get

| Surface | What appears |
|---------|--------------|
| Integrations | Jira Cloud and Confluence Cloud connection presets |
| Issue tracker | a starter project with an epic, three stories and a defect |
| Wiki | a team home page, a decision record and a runbook |
| Test cases | three checks that tell you the migration actually worked |
| Links | deep links from a migrated issue or page back to the original |
| Wiki editor | a `plugin-jira-issue` macro that shows an issue's key and state |
| Notifications | events for "import finished" and "import conflict" |

## Installing

```bash
mockarty-cli plugin install mockarty.atlassian-teamwork-pack-1.0.0.zip
mockarty-cli plugin enable mockarty.atlassian-teamwork-pack
```

Then open the extension's settings and set **Atlassian site URL** to your site,
so the deep links point at the right place. The API token field holds a
*reference* to a secret — the token itself never enters the extension.

## Using the macro

On any wiki page:

```
!plugin-jira-issue:key=MIGR-12 state="In review"
```

Arguments are separated by **spaces**, and a value containing a space must be
quoted — `state=In review` would take only the first word.

It renders the key, the state and a link back to the original issue.

## What it does not do

It does not import your data — it prepares the shape your data will land in.
Use the tracker's own import for the issues themselves, then instantiate the
content packs for the parts an import does not carry.
