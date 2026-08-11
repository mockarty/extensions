# response-transformer — envelope (WASM response-transformer example)

A `response-transformer` plugin: tag any mock with `transform:envelope` and its
rendered body is wrapped as `{"data":<body>,"transformedBy":"plugin"}` by a
sandboxed WASM module that runs inline on the response path.

    ./build.sh                 # needs TinyGo
    mockarty-cli plugin pack .
    mockarty-cli plugin install mockarty.response-envelope-1.0.0.zip
    mockarty-cli plugin enable mockarty.response-envelope

Then add the tag `transform:envelope` to any mock.
