# Healthcare (FHIR-lite) mock kit

A FHIR-style healthcare contour (Patient, Observation) — mock a health-records API for integration testing. Structure mirrors FHIR resources; values are synthetic.

    mockarty-cli plugin install mockarty.healthcare-fhir-kit-1.0.0.zip
    mockarty-cli plugin enable mockarty.healthcare-fhir-kit

Then instantiate the kit `fhir_api` into any namespace from the Mock Kits catalogue (UI or MCP `mock_kit_instantiate`).
