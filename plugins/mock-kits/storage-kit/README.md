# Object storage mock kit

An S3-style object-storage contour (buckets, objects, presigned URLs) — mock a blob store for upload/download flows.

    mockarty-cli plugin install mockarty.storage-kit-1.0.0.zip
    mockarty-cli plugin enable mockarty.storage-kit

Then instantiate the kit `object_storage` into any namespace from the Mock Kits catalogue (UI or MCP `mock_kit_instantiate`).
