# GitLab connector

A ready-to-use GitLab connector — reuses the built-in GitLab adapter, no code. Links issues, merge requests and pipelines to Mockarty.

    mockarty-cli plugin install mockarty.gitlab-connector-1.0.0.zip
    mockarty-cli plugin enable mockarty.gitlab-connector

Then create an integration from the connector `gitlab_self` (Integrations → New → from connector) and add your credentials per-namespace.
