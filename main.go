package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func keyFunc(secret string, isBase64 bool) func(*jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Expected HS256 signing algorithm")
		}

		if isBase64 {
			return base64.StdEncoding.DecodeString(secret)
		}

		return []byte(secret), nil
	}
}

func validateJWTAndShowClaims(jwtString string, secret string, b64 bool) {
	claims, err := jwt.Parse(jwtString, keyFunc(secret, b64))
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				println("Provided value was not a recognizable JWT")
			}

			if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				println("JWT is expired")
			}

			if ve.Errors&(jwt.ValidationErrorNotValidYet) != 0 {
				println("JWT is not valid yet")
			}

			// Ok lets try some stuff - did they not base64 encode?
			_, err = jwt.Parse(jwtString, keyFunc(secret, !b64))
			if err == nil {
				if b64 {
					println("Provided secret was base64 encoded, but the JWT was signed with the encoded secret instead of the decoded bytes")
				} else {
					println("Provided secret was not encoded, but the JWT was signed with the base64 decoded bytes instead.")
				}
			}

			os.Exit(-1)
		}
	} else {
		b, _ := json.MarshalIndent(claims.Claims, "", "  ")
		println(string(b))
	}
}

func main() {
	// Expects arguments in the form echo 'whatever' | jwt-sign --secret=<secret>  --expires=30d
	// or jwt-sign --secret=<secret> --claims-file='whatever-file' --expires=30d
	// Only supports HS256 signing for now.

	var file = flag.String("claims-file", "-", "A json file that represents the claims the JWT will have. Set to - to read from stdin")
	var secret = flag.String("secret", "", "The secret used to sign the JWT. Must be non-empty.")
	var claimBlob = flag.String("claims", "", "JSON blob to set claims. Cannot be used in combination with -claims-file")
	var b64 = flag.Bool("base64", true, "Set to false if the secret is not base64 encoded")
	var duration = flag.Duration("exp", time.Minute*20, "Expiry time. Can be specified in 's', 'm' or 'h'.")
	var jwtString = flag.String("jwt", "", "If provided, the jwt will be validated and the claims shown.")

	flag.Parse()

	if len(*secret) == 0 {
		flag.PrintDefaults()
		return
	}

	if *jwtString != "" {
		validateJWTAndShowClaims(*jwtString, *secret, *b64)
		return
	}

	var err error

	// Read claims
	var in io.Reader = os.Stdin
	if *file != "-" {
		in, err = os.Open(*file)
		if err != nil {
			log.Panic(err)
		}
	} else if *claimBlob != "" {
		in = bytes.NewReader([]byte(*claimBlob))
	}

	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		log.Panic(err)
	}

	claims := map[string]interface{}{}
	err = json.Unmarshal(bytes, &claims)
	if err != nil {
		log.Panic(err)
	}

	// Setup iat and expires fields
	claims["iat"] = time.Now().Unix()

	claims["exp"] = time.Now().Add(*duration).Unix()

	// Setup secret, decoding from base64 if required
	secretBytes := []byte(*secret)
	if *b64 {
		secretBytes, err = base64.StdEncoding.DecodeString(*secret)
		if err != nil {
			log.Panic(err)
		}
	}

	// Sign and print the JWT.
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	signed, err := tk.SignedString(secretBytes)
	if err != nil {
		log.Panic(err)
	}

	println(signed)
	return
}
