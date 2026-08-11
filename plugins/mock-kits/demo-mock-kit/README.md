# demo-mock-kit — starter Users API (mock-kit plugin example)

The simplest kind of Mockarty plugin: **declarative content, no code**. It
contributes a mock kit — a ready-made contour of `/users` endpoints with dynamic
fakers — to the mock-kits catalogue.

## Install & use

```bash
mockarty-cli plugin pack .
mockarty-cli plugin install mockarty.demo-mock-kit-1.0.0.zip
mockarty-cli plugin enable mockarty.demo-mock-kit
```

Once enabled, the "Users API (demo)" kit shows up in the mock-kits catalogue
(UI, and the `mock_kit_list` / `mock_kit_instantiate` MCP tools). Instantiate it
into a namespace and you have working `GET /users`, `GET /users/:id` and
`POST /users` mocks immediately.

This is the fastest way to seed a shareable contour — no build step, no
dependencies, pure JSON.
