# finance-fakers — card/IBAN (WASM faker example)

A second WASM faker (alongside ru-fakers): `$.fake.credit_card` (16 digits with a
valid Luhn check — a **test** PAN, never a real card) and `$.fake.iban`.

    ./build.sh                 # needs TinyGo
    mockarty-cli plugin pack .
    mockarty-cli plugin install mockarty.finance-fakers-1.0.0.zip
    mockarty-cli plugin enable mockarty.finance-fakers
