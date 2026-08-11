# condition-matcher — luhn / divisible_by (WASM condition-matcher example)

Adds two mock condition operators via a sandboxed WASM module:

- `plugin:luhn` — the value extracted by the condition's path passes the Luhn
  checksum (credit-card / IMEI numbers).
- `plugin:divisible_by` — the value is an integer divisible by the operand
  (`value` field of the condition).

    ./build.sh                 # needs TinyGo
    mockarty-cli plugin pack .
    mockarty-cli plugin install mockarty.matchers-1.0.0.zip
    mockarty-cli plugin enable mockarty.matchers

Then set a condition `assertAction: "plugin:luhn"` (path → the field to check) on
any mock.
