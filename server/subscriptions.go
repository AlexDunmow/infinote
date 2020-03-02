package infinote

import (
	"context"
	"sync"
)

type CursorPlacement struct {
	LineNumber int    `json:"lineNumber"`
	Column     int    `json:"column"`
	UserID     string `json:"userID"`
	UserName   string `json:"userName"`
}

// BlacklistProvider methods for token blacklist management
type SubscriptionProvider interface {
	NoteEvent(ctx context.Context, noteID string) (<-chan *CursorPlacement, error)
}

// Blacklister implements BlacklistProvider methods
type SubscriptionHub struct {
	*sync.RWMutex
}

// NewBlacklister returns a token blacklist provider
func NewSubHub() *SubscriptionHub {
	subhub := &SubscriptionHub{
		&sync.RWMutex{},
	}
	return subhub
}

// OnList checks if token id is on the blacklist
func (b *Blacklister) CursorChanged(tokenID string) bool {
	_, found := b.blacklist[tokenID]
	if found {
		return true
	}
	return false
}
