# Payments API mock kit

A ready-made payment-processor contour (charges, refunds, customers, payment intents) — declarative content, no code. Instantiate to get a Stripe-style API you can point a client at.

    mockarty-cli plugin install mockarty.payments-kit-1.0.0.zip
    mockarty-cli plugin enable mockarty.payments-kit

Then instantiate the kit `payments_api` into any namespace from the Mock Kits catalogue (UI or MCP `mock_kit_instantiate`).
