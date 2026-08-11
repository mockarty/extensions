# tracker-links

Declarative **link types** (D6): issue↔Jira ticket and mock↔GitHub issue with deep-link URL templates.

    mockarty-cli plugin install mockarty.tracker-links-1.0.0.zip
    mockarty-cli plugin enable mockarty.tracker-links

Catalogue: `GET /api/v1/entity-link-types?entity=issue`.

Stored links: issue cards show an **External links** section (add/remove; each
renders as a deep link). REST: `POST/GET/DELETE /api/v1/entity-links`; MCP:
`entity_link_add` / `entity_link_list` / `entity_link_remove`.
