# OIDC auth mock kit

An OpenID-Connect provider contour (token, userinfo, jwks, discovery) — mock an identity provider for local auth flows without a real IdP.

    mockarty-cli plugin install mockarty.auth-oidc-kit-1.0.0.zip
    mockarty-cli plugin enable mockarty.auth-oidc-kit

Then instantiate the kit `oidc_provider` into any namespace from the Mock Kits catalogue (UI or MCP `mock_kit_instantiate`).
