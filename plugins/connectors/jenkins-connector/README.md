# Jenkins connector

A ready-to-use Jenkins connector — reuses the built-in Jenkins adapter to trigger and read CI jobs from Mockarty.

    mockarty-cli plugin install mockarty.jenkins-connector-1.0.0.zip
    mockarty-cli plugin enable mockarty.jenkins-connector

Then create an integration from the connector `jenkins_ci` (Integrations → New → from connector) and add your credentials per-namespace.
