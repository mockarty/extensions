# Generic webhook connector

A ready-to-use generic webhook connector — reuses the built-in webhook adapter to POST Mockarty events to any URL.

    mockarty-cli plugin install mockarty.webhook-connector-1.0.0.zip
    mockarty-cli plugin enable mockarty.webhook-connector

Then create an integration from the connector `generic_webhook` (Integrations → New → from connector) and add your credentials per-namespace.
