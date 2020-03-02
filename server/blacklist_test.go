package infinote_test

import (
	infinote "boilerplate"
	"boilerplate/db"
	"boilerplate/mocks"
	"boilerplate/store"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestNewBlacklister(t *testing.T) {
	l := zap.NewNop().Sugar()
	tokenStore := &mocks.TokenStorer{}
	tokenStore.On("Blacklist").Return(store.Blacklist{}, nil)
	userAuth := &infinote.UserAuth{}
	type args struct {
		log         *zap.SugaredLogger
		tokenStorer *mocks.TokenStorer
		authConfig  *infinote.UserAuth
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{"happy path", args{l, tokenStore, userAuth}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := infinote.NewBlacklister(tt.args.log, tt.args.tokenStorer, 24)
			if !tt.wantNil && bl == nil {
				t.Errorf("blacklister: got %v, expected not nil", bl)
			}
		})
	}
}

func TestBlacklister_OnList(t *testing.T) {
	l := zap.NewNop().Sugar()
	tokenStore := &mocks.TokenStorer{}
	tokenStore.On("Blacklist").Return(store.Blacklist{}, nil)
	userAuth := &infinote.UserAuth{}
	blacklist := &mocks.BlacklistProvider{}
	mutex := &sync.Mutex{}

	type fields struct {
		blacklist  *mocks.BlacklistProvider
		mutex      *sync.Mutex
		log        *zap.SugaredLogger
		tokenStore *mocks.TokenStorer
		authConfig *infinote.UserAuth
	}
	type args struct {
		tokenID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"happy path", fields{blacklist, mutex, l, tokenStore, userAuth}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := infinote.NewBlacklister(tt.fields.log, tt.fields.tokenStore, 24)
			if got := bl.OnList(tt.args.tokenID); got != tt.want {
				t.Errorf("Blacklister.OnList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlacklister_CleanIssuedTokens(t *testing.T) {
	l := zap.NewNop().Sugar()
	tokenStore := &mocks.TokenStorer{}
	tokenStore.On("Blacklist").Return(store.Blacklist{}, nil)
	tokenStore.On("GetAllExpired").Return([]*db.IssuedToken{}, nil)
	userAuth := &infinote.UserAuth{}
	blacklist := &mocks.BlacklistProvider{}
	mutex := &sync.Mutex{}

	type fields struct {
		blacklist  *mocks.BlacklistProvider
		mutex      *sync.Mutex
		log        *zap.SugaredLogger
		tokenStore *mocks.TokenStorer
		authConfig *infinote.UserAuth
	}
	type args struct {
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy path", fields{blacklist, mutex, l, tokenStore, userAuth}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := infinote.NewBlacklister(tt.fields.log, tt.fields.tokenStore, 24)
			err := bl.CleanIssuedTokens()
			if !tt.wantErr && err != nil {
				t.Errorf("got %v, expected nil", err)
			}
		})
	}
}

func TestBlacklister_RefreshBlacklist(t *testing.T) {
	l := zap.NewNop().Sugar()
	tokenStore := &mocks.TokenStorer{}
	tokenStore.On("Blacklist").Return(store.Blacklist{}, nil)
	userAuth := &infinote.UserAuth{}
	blacklist := &mocks.BlacklistProvider{}
	mutex := &sync.Mutex{}

	type fields struct {
		blacklist  *mocks.BlacklistProvider
		mutex      *sync.Mutex
		log        *zap.SugaredLogger
		tokenStore *mocks.TokenStorer
		authConfig *infinote.UserAuth
	}
	type args struct {
		tokenID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy path", fields{blacklist, mutex, l, tokenStore, userAuth}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := infinote.NewBlacklister(tt.fields.log, tt.fields.tokenStore, 24)
			err := bl.RefreshBlacklist()
			if !tt.wantErr && err != nil {
				t.Errorf("got %v, expected nil", err)
			}
		})
	}
}

func TestBlacklister_StartTicker(t *testing.T) {
	l := zap.NewNop().Sugar()
	tokenStore := &mocks.TokenStorer{}
	tokenStore.On("Blacklist").Return(store.Blacklist{}, nil)
	userAuth := &infinote.UserAuth{}
	userAuth.BlacklistRefreshHours = 1
	userAuth.TokenExpiryDays = 1
	blacklist := &mocks.BlacklistProvider{}
	mutex := &sync.Mutex{}

	type fields struct {
		blacklist  *mocks.BlacklistProvider
		mutex      *sync.Mutex
		log        *zap.SugaredLogger
		tokenStore *mocks.TokenStorer
		authConfig *infinote.UserAuth
	}
	type args struct {
		tokenID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy path", fields{blacklist, mutex, l, tokenStore, userAuth}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := infinote.NewBlacklister(tt.fields.log, tt.fields.tokenStore, 24)
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
			defer cancel()
			bl.StartTicker(ctx)
		})
	}
}

func TestBlacklister_BlacklistAll(t *testing.T) {
	l := zap.NewNop().Sugar()
	tokenStore := &mocks.TokenStorer{}
	tokenStore.On("Blacklist").Return(store.Blacklist{}, nil)
	tokenStore.On("GetAllByUser", mock.AnythingOfType("string")).Return([]*db.IssuedToken{}, nil)
	userAuth := &infinote.UserAuth{}
	blacklist := &mocks.BlacklistProvider{}
	mutex := &sync.Mutex{}

	type fields struct {
		blacklist  *mocks.BlacklistProvider
		mutex      *sync.Mutex
		log        *zap.SugaredLogger
		tokenStore *mocks.TokenStorer
		authConfig *infinote.UserAuth
	}
	type args struct {
		tokenID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy path", fields{blacklist, mutex, l, tokenStore, userAuth}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := infinote.NewBlacklister(tt.fields.log, tt.fields.tokenStore, 24)

			err := bl.BlacklistAll(uuid.Must(uuid.NewV4()).String())
			if !tt.wantErr && err != nil {
				t.Errorf("got %v, expected nil", err)
			}
		})
	}
}

func TestBlacklister_BlacklistOne(t *testing.T) {
	l := zap.NewNop().Sugar()
	tokenStore := &mocks.TokenStorer{}
	tokenStore.On("Blacklist").Return(store.Blacklist{}, nil)
	tokenStore.On("Get", mock.AnythingOfType("string")).Return(&db.IssuedToken{}, nil)
	tokenStore.On("Update", mock.AnythingOfType("*db.IssuedToken")).Return(&db.IssuedToken{}, nil)
	userAuth := &infinote.UserAuth{}
	blacklist := &mocks.BlacklistProvider{}
	mutex := &sync.Mutex{}

	type fields struct {
		blacklist  *mocks.BlacklistProvider
		mutex      *sync.Mutex
		log        *zap.SugaredLogger
		tokenStore *mocks.TokenStorer
		authConfig *infinote.UserAuth
	}
	type args struct {
		tokenID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy path", fields{blacklist, mutex, l, tokenStore, userAuth}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := infinote.NewBlacklister(tt.fields.log, tt.fields.tokenStore, 24)
			err := bl.BlacklistOne(uuid.Must(uuid.NewV4()).String())
			if !tt.wantErr && err != nil {
				t.Errorf("got %v, expected nil", err)
			}
		})
	}
}
