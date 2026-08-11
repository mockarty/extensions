# team-wiki-pack

A Confluence-style content pack for the **Wiki** module — instantiate it to seed a namespace with runbook/ADR/meeting-notes/onboarding/PRD pages.

    mockarty-cli plugin install mockarty.team-wiki-pack-1.0.0.zip
    mockarty-cli plugin enable mockarty.team-wiki-pack

Then instantiate the pack `team_wiki` from the Content Packs catalogue (`POST /api/v1/content-packs/team_wiki/instantiate`) into any namespace.
