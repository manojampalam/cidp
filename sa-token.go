package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func getSAToken(namespace, name string) (string, error) {
	k := getSwSKey()
	iat := time.Now().Unix()
	//valid for 30 days
	exp := iat + (30 * 24 * 60 * 60)
	header := fmt.Sprintf(`{"alg":"%s","kid":"%s"}`, k.Alg(), k.KeyID())
	headerB64 := base64.RawURLEncoding.EncodeToString([]byte(header))
	aud := "https://login.microsoftonline.com/"
	iss := "arc.azure.com"
	sub := "system:serviceaccount:" + namespace + ":" + name
	payload := fmt.Sprintf(`{"aud":["%s"],"iat":%d,"nbf":%d,"exp":%d,"iss":"%s","sub":"%s"}`, aud, iat, iat, exp, iss, sub)
	payloadB64 := base64.RawURLEncoding.EncodeToString([]byte(payload))
	h256 := sha256.Sum256([]byte(headerB64 + "." + payloadB64))
	signature, err := k.Sign(h256[:])
	if err != nil {
		return "", err
	}
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	return headerB64 + "." + payloadB64 + "." + signatureB64, nil
}
