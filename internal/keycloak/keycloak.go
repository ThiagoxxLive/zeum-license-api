// Package keycloak busca e mantém em cache a chave pública do realm do
// Keycloak, usada para validar a assinatura RS256 dos JWTs emitidos por ele.
// Equivalente ao KeycloakJWSProvider do zeum-admin-api.
package keycloak

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const publicKeyTTL = 1 * time.Hour

type PublicKeyProvider struct {
	baseURL    string
	realm      string
	httpClient *http.Client

	mu        sync.Mutex
	cachedKey *rsa.PublicKey
	expiresAt time.Time
}

func NewPublicKeyProvider(baseURL, realm string) *PublicKeyProvider {

	return &PublicKeyProvider{
		baseURL:    strings.TrimRight(baseURL, "/"),
		realm:      realm,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// PublicKey retorna a chave pública do realm, buscando-a do Keycloak apenas
// quando o cache expirou ou quando forceRefresh é true (usado para uma nova
// tentativa após falha de verificação, caso a chave tenha rotacionado).
func (p *PublicKeyProvider) PublicKey(forceRefresh bool) (*rsa.PublicKey, error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	if !forceRefresh && p.cachedKey != nil && time.Now().Before(p.expiresAt) {
		return p.cachedKey, nil
	}

	key, err := p.fetchPublicKey()

	if err != nil {
		return nil, err
	}

	p.cachedKey = key
	p.expiresAt = time.Now().Add(publicKeyTTL)

	return key, nil
}

func (p *PublicKeyProvider) fetchPublicKey() (*rsa.PublicKey, error) {

	url := fmt.Sprintf("%s/realms/%s", p.baseURL, p.realm)

	resp, err := p.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak retornou status %d ao buscar a chave pública do realm", resp.StatusCode)
	}

	var body struct {
		PublicKey string `json:"public_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	if body.PublicKey == "" {
		return nil, errors.New("resposta do keycloak não contém public_key")
	}

	der, err := base64.StdEncoding.DecodeString(body.PublicKey)

	if err != nil {
		return nil, err
	}

	parsed, err := x509.ParsePKIXPublicKey(der)

	if err != nil {
		return nil, err
	}

	rsaKey, ok := parsed.(*rsa.PublicKey)

	if !ok {
		return nil, errors.New("chave pública do keycloak não é RSA")
	}

	return rsaKey, nil
}
