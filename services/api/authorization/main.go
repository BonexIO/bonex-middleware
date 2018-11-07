package authorization

import (
	"crypto/sha256"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"net/http"
)

func Authorize(r *http.Request, path string) (bool, string) {

	pubkey := []byte(r.Header.Get("XPubKey"))
	signature := []byte(r.Header.Get("XSignature"))
	tonce := []byte(r.Header.Get("XTonce"))

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return false, "body"
	}

	payload := append([]byte(path), append(body, tonce...)...)

	return Verify(pubkey, signature, payload), "ok"

}

func Verify(pubkey ed25519.PublicKey, signature, payload []byte) bool {

	hash := sha256.Sum256(payload)

	message := hash[:]

	return ed25519.Verify(pubkey, message, signature)
}
