# Board templates: retro + C4 sketch

A `content_packs` plugin with `kind: "board"` — whiteboard TEMPLATES delivered
through the plugin system, so a template gallery can grow without a product
release.

Install + enable the plugin, then instantiate the pack (Content packs picker,
or `content_pack_instantiate` over MCP) — two boards appear in the namespace:

- **Retrospective — Start / Stop / Continue** — three-column retro frame;
- **C4 container sketch** — user → app → database / external API with bound
  arrows (drag a box, the arrows follow).

Each item is `{name, description, tags, scene}` where `scene` is plain
Excalidraw JSON — copy this plugin and swap in your own scenes to ship a
company template pack.
