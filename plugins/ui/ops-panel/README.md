# Ops panel

Adds a sidebar item opening a small operations panel — an example of a UI-contribution plugin. The panel ships inside the plugin and renders in a sandboxed iframe (strict CSP, no host-session access), theme-aware with the host.

- **id**: `mockarty.ops-panel`  ·  **version**: 1.0.0  ·  **mechanisms**: ui

## Install

Download the release zip (or run `mockarty-cli plugin pack .` in this folder), then:

```bash
mockarty-cli plugin install ops-panel-1.0.0.zip
mockarty-cli plugin enable mockarty.ops-panel
```

Or in the UI: **Admin -> Plugins -> Install plugin...** (file) / **Install from URL...**.
