# Contract starter pack

A **content pack** for the Contract testing module. It publishes ready-to-use API
contracts into the registry so a team can start verifying consumers/providers
immediately instead of authoring specs from scratch.

## What it seeds

| Service | Spec type | What |
|---------|-----------|------|
| `orders-api` | OpenAPI 3.0 | List + create orders, with an `Order` schema |
| `users-api` | GraphQL | `Query.me`/`user`, `Mutation.updateUser`, `User` type |

Each is published as an **active** registry entry in the target namespace.

## Requirements

The **contract testing** feature must be licensed for the namespace — the pack is
gated server-side, same as the API-tester and TCM content packs.

## Install & instantiate

```bash
# install the bundle (zip of this folder)
curl -sS -X POST "$MOCKARTY/api/v1/plugins?namespace=$NS" -F "bundle=@contract-starter-pack.zip"
# enable it for the namespace
curl -sS -X POST "$MOCKARTY/api/v1/plugins/mockarty.contract-starter-pack/enable?namespace=$NS"
# seed the contracts
curl -sS -X POST "$MOCKARTY/api/v1/content-packs/contract_starter/instantiate?namespace=$NS"
```

Or, for an agent, over MCP: `content_pack_list` → `content_pack_instantiate`.

Re-instantiating seeds another copy — the pack is a starter, not a sync source.

No code: this is a pure declarative content pack.
