package coalmer

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ECDSASignature struct {
	R, S *big.Int
}

type payload struct {
	Iat int64  `json:"iat"`
	Jti string `json:"jti"`
	Htu string `json:"htu"`
	Htm string `json:"htm"`
	Uid string `json:"uuid"`
}

type pkeyJwk struct {
	Crv string `json:"crv"`
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type pkeyHeader struct {
	Typ string  `json:"typ"`
	Alg string  `json:"alg"`
	Jwk pkeyJwk `json:"jwk"`
}

func byteToBase64URL(target []byte) string {
	return base64.RawURLEncoding.EncodeToString(target)
}

// dPoPGenerator is generate the dPoP token
func dPoPGenerator(uuid_ string, method string, url_ string) (string, error) { //因为有 url和uuid 包了
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", errors.New("error at dPoPGenerator/ecdsa.GenerateKey(): " + err.Error())
	}

	//pl := payload{time.Now().Unix(), uuid_, url_, strings.ToUpper(method), CoalMiner.ClientID}
	clientId := uuid.NewString()
	pl := payload{time.Now().Unix(), uuid_, url_, strings.ToUpper(method), clientId}
	pkjwk := pkeyJwk{"P-256", "EC", byteToBase64URL(privateKey.PublicKey.X.Bytes()), byteToBase64URL(privateKey.PublicKey.Y.Bytes())}
	pkh := pkeyHeader{"dpop+jwt", "ES256", pkjwk}

	headerString, err := json.Marshal(pkh)
	if err != nil {
		return "", errors.New("error at dPoPGenerator/json.Marshal(pkh): " + err.Error())
	}
	payloadString, err := json.Marshal(pl)
	if err != nil {
		return "", errors.New("error at dPoPGenerator/json.Marshal(pl): " + err.Error())
	}

	data_unsigned := fmt.Sprintf("%s.%s", byteToBase64URL(headerString), byteToBase64URL(payloadString))
	hval := sha256.Sum256([]byte(data_unsigned))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hval[:])
	if err != nil {
		return "", errors.New("error at dPoPGenerator/ecdsa.Sign(): " + err.Error())
	}

	signatured := append(r.Bytes(), s.Bytes()...)
	signaturedString := byteToBase64URL(signatured)
	result := fmt.Sprintf("%s.%s", data_unsigned, signaturedString)
	return result, nil
}

func generateSearchSessionId(length int) string {
	buflen := length
	if buflen%2 != 0 {
		buflen += 1
	}
	buf := make([]byte, buflen/2)
	rand.Read(buf)
	return hex.EncodeToString(buf)[:length]
}

func generateHeader(link, method string) (headers map[string]string, err error) {
	u, _ := url.Parse(link)
	uri := u.Scheme + "://" + u.Host + u.EscapedPath()
	dPoP, err := dPoPGenerator(uuid.NewString(), strings.ToUpper(method), uri)

	if err != nil {
		return nil, err
	}

	headers = map[string]string{
		"DPoP":            dPoP,
		"x-platform":      "web",
		"authority":       ApiURL,
		"accept-language": "ja",
		//"accept-encoding": "gzip, deflate, br",
		//"accept-encoding": "gzip, deflate",
		"accept":       "application/json, text/plain, */*",
		"Content-Type": "application/json",
		"origin":       RootURL,
		"referer":      RootURL,
	}
	return headers, nil
}
