package admin

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

const (
	TokenSize           = 32
	TokenExpiration     = 10 * time.Minute
	TokenRotateInterval = 5 * time.Minute
)

var (
	tokenStore = NewTokenStore(TokenExpiration, TokenRotateInterval)
)

func init() {
	go tokenStore.Start()
}

type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).Unix())), nil
}

type Token struct {
	Token     string   `json:"token"`
	ExpiresAt UnixTime `json:"expiresAt"`
}

func NewToken(expiration time.Duration) *Token {
	b := make([]byte, TokenSize)
	rand.Read(b)
	return &Token{
		Token:     hex.EncodeToString(b),
		ExpiresAt: UnixTime(time.Now().Add(expiration)),
	}
}

func (t *Token) IsValid(token string) bool {
	if t == nil {
		return false
	}
	return token == t.Token && time.Now().Before(time.Time(t.ExpiresAt))
}

type TokenStore struct {
	tokenExpiration time.Duration
	rotateInterval  time.Duration

	lock    sync.RWMutex
	tokens  []*Token
	current int
}

func NewTokenStore(tokenExpiration, rotateInterval time.Duration) *TokenStore {
	tokens := make([]*Token, tokenExpiration/rotateInterval+1)
	tokens[0] = NewToken(tokenExpiration)
	return &TokenStore{
		tokenExpiration: tokenExpiration,
		rotateInterval:  rotateInterval,
		tokens:          tokens,
		current:         0,
	}
}

func (tr *TokenStore) Start() {
	for {
		time.Sleep(tr.rotateInterval)
		tr.lock.Lock()
		tr.current = (tr.current + 1) % len(tr.tokens)
		tr.tokens[tr.current] = NewToken(tr.tokenExpiration)
		tr.lock.Unlock()
	}
}

func (tr *TokenStore) GetCurrentToken() *Token {
	tr.lock.RLock()
	defer tr.lock.RUnlock()
	return tr.tokens[tr.current]
}

func (tr *TokenStore) IsValid(token string) bool {
	tr.lock.RLock()
	defer tr.lock.RUnlock()
	for i := range tr.tokens {
		if tr.tokens[i] != nil && tr.tokens[i].IsValid(token) {
			return true
		}
	}
	return false
}
