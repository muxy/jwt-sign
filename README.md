jwt-sign
============
`jwt-sign` is a basic HS256 jwt signing and validation command line application.

Signing
-----
To sign a JWT, provide the `--secret` parameter and a claims object through either `--claims` or `--claims-file`.

    jwt-sign --secret "this-is-a-secret" --base64=false \
       --claims '{ "falcons": 42, "role": "hawker" }' --exp 60m

or

    echo '{ "falcons": 42, "role": "hawker" }' | jwt-sign --secret "this-is-a-secret" --base64=false --exp 60m

or

    echo '{ "falcons": 42, "role": "hawker" }' > claims.json
    jwt-sign --secret "this-is-a-secret" --base64=false --exp 60m --claims-file claims.json

Validating
-------------
To validate a JWT and show the claims object, simply pass in `--jwt` instead of a claims object.

    jwt-sign --jwt "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzkyMTUxNDYsImZhbGNvbnMiOjQyLCJpYXQiOjE1MzkyMTE1NDYsInJvbGUiOiJoYXdrZXIifQ.RFmyDuz8MTqMYjDzD4o3S1Kb_cNr48_RBacHZJes7d8" --secret "this-is-a-secret" --base64=false

Invalid JWTs will have an error printed instead.