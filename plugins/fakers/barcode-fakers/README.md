# barcode-fakers — EAN-13 / UPC-A / ISBN-13 (WASM faker example)

Adds `$.fake.ean13`, `$.fake.upc`, `$.fake.isbn13` — valid retail barcodes with
correct GS1 check digits.

    ./build.sh                 # needs TinyGo
    mockarty-cli plugin install mockarty.barcode-fakers-1.0.0.zip
    mockarty-cli plugin enable mockarty.barcode-fakers
