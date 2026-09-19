# Wiki status lozenges

Confluence-parity **wiki macros**, contributed declaratively — no code.

| Macro | Usage | Renders |
|---|---|---|
| Status lozenge | `!plugin-status:colour=green SHIPPED` | a coloured pill (grey/green/red/yellow/blue) |
| Decision panel | `!plugin-decision:owner=alice We ship on Friday` | a success panel titled **DECISION — alice** |

After install + enable, both appear in the wiki editor's `/` slash-menu. The
Confluence importer's `{status}` output pairs naturally with the lozenge.

Templates are Go `text/template` producing markdown/pinned-attribute HTML; the
output passes the same sanitizer as hand-typed content, so macros cannot smuggle
scripts or styles.

```bash
mockarty-cli plugin install wiki-status-macros-1.0.0.zip
mockarty-cli plugin enable mockarty.wiki-status-macros
```
