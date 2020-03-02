package infinote

import (
	"boilerplate/store"
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// BlacklistProvider methods for token blacklist management
type BlacklistProvider interface {
	OnList(tokenID string) bool
	CleanIssuedTokens() error
	RefreshBlacklist() error
	StartTicker(context.Context)
	BlacklistAll(userID string) error
	BlacklistOne(tokenID string) error
}

// Blacklister implements BlacklistProvider methods
type Blacklister struct {
	blacklist             store.Blacklist // blacklisted JWT IDs
	mutex                 sync.Mutex
	log                   *zap.SugaredLogger
	TokenStore            TokenStorer
	blacklistRefreshHours int
}

// NewBlacklister returns a token blacklist provider
func NewBlacklister(log *zap.SugaredLogger, tokenStorer TokenStorer, blacklistRefreshHours int) *Blacklister {
	blacklist, err := tokenStorer.Blacklist()
	if err != nil {
		panic(err)
	}
	b := &Blacklister{
		blacklist:             blacklist,
		log:                   log,
		TokenStore:            tokenStorer,
		blacklistRefreshHours: blacklistRefreshHours,
	}
	return b
}

// OnList checks if token id is on the blacklist
func (b *Blacklister) OnList(tokenID string) bool {
	_, found := b.blacklist[tokenID]
	if found {
		return true
	}
	return false
}

// CleanIssuedTokens will delete rows from the issued_tokens list that have passed token expiry and then reload list
// into memory.
func (b *Blacklister) CleanIssuedTokens() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.log.Info("cleaning blacklist...")

	expiredTokens, err := b.TokenStore.GetAllExpired()
	if err != nil {
		return fmt.Errorf("get expired: %w", err)
	}

	for _, token := range expiredTokens {
		err := b.TokenStore.Delete(token)
		if err != nil {
			return fmt.Errorf("delete expired: %w", err)
		}
	}

	err = b.RefreshBlacklist()
	if err != nil {
		b.log.Error(err)
	}

	return nil
}

// RefreshBlacklist reloads the blacklist into memory
func (b *Blacklister) RefreshBlacklist() error {
	b.log.Info("refreshing blacklist...")

	blacklist, err := b.TokenStore.Blacklist()
	if err != nil {
		return fmt.Errorf("get blacklist: %w", err)
	}

	b.blacklist = blacklist
	return nil
}

// StartTicker will start the ticker to do issued tokens maintenance to clear expired tokens etc.
func (b *Blacklister) StartTicker(ctx context.Context) {

	dur := time.Duration(b.blacklistRefreshHours) * time.Hour

	ticker := time.NewTicker(dur)
	stop := make(chan bool, 1)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := b.CleanIssuedTokens()
			if err != nil {
				b.log.Error(err)
			}
		case <-stop:
			return
		}
	}
}

// BlacklistAll will mark all tokens from that userID as blacklisted. Used when changing password.
func (b *Blacklister) BlacklistAll(userID string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	tokens, err := b.TokenStore.GetAllByUser(userID)
	if err != nil {
		return fmt.Errorf("get user blacklist: %w", err)
	}

	for _, token := range tokens {
		token.Blacklisted = true
		_, err = b.TokenStore.Update(token)
		if err != nil {
			return fmt.Errorf("update token: %w", err)
		}
	}

	err = b.RefreshBlacklist()
	if err != nil {
		return fmt.Errorf("refresh blacklist: %w", err)
	}

	return nil
}

// BlacklistOne will mark a single token as blacklisted. Could be used to log out a single device
func (b *Blacklister) BlacklistOne(tokenID string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	token, err := b.TokenStore.Get(tokenID)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	token.Blacklisted = true
	_, err = b.TokenStore.Update(token)
	if err != nil {
		return fmt.Errorf("update token: %w", err)
	}

	err = b.RefreshBlacklist()
	if err != nil {
		return fmt.Errorf("refresh blacklist: %w", err)
	}

	return nil
}
