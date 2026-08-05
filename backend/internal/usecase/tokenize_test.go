package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Fabryny/token-vault/backend/internal/domain"
	"github.com/Fabryny/token-vault/backend/internal/infra/crypto"
	"github.com/Fabryny/token-vault/backend/internal/usecase"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) Save(ctx context.Context, t *domain.Token) error {
	return m.Called(ctx, t).Error(0)
}

func (m *MockTokenRepository) FindByValue(ctx context.Context, value string) (*domain.Token, error) {
	args := m.Called(ctx, value)
	if v := args.Get(0); v != nil {
		return v.(*domain.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

// newEncryptor: o Encryptor REAL. Ele é código nosso, puro, sem I/O —
// mocka-se a FRONTEIRA (banco, rede), não a lógica do próprio projeto.
func newEncryptor(t *testing.T) domain.Encryptor {
	t.Helper()
	enc, err := crypto.NewEncryptor(make([]byte, 32))
	require.NoError(t, err)
	return enc
}

func TestTokenizeSuccess(t *testing.T) {
	const pan = "4111111111111111"

	repo := new(MockTokenRepository)
	repo.On("Save", mock.Anything, mock.Anything).Return(nil)

	uc := usecase.NewTokenize(repo, newEncryptor(t))
	got, err := uc.Execute(context.Background(), pan, uuid.New())

	require.NoError(t, err)
	assert.NotEmpty(t, got.Value)
	assert.NotEqual(t, pan, got.Value) // a premissa do sistema
	assert.NotContains(t, string(got.Ciphertext), pan)
	assert.Equal(t, "1111", got.Last4)
	assert.Len(t, got.Nonce, 12)
	repo.AssertExpectations(t)
}

func TestTokenizeRejectsInvalidPAN(t *testing.T) {
	repo := new(MockTokenRepository)

	uc := usecase.NewTokenize(repo, newEncryptor(t))
	_, err := uc.Execute(context.Background(), "4111111111111112", uuid.New())

	assert.ErrorIs(t, err, domain.ErrInvalidPAN)
	repo.AssertNotCalled(t, "Save") // não pode chegar no banco
}

func TestDetokenizeRoundTrip(t *testing.T) {
	const pan = "4111111111111111"

	repo := new(MockTokenRepository)
	var saved *domain.Token
	repo.On("Save", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { saved = args.Get(1).(*domain.Token) }).
		Return(nil)

	enc := newEncryptor(t)

	created, err := usecase.NewTokenize(repo, enc).Execute(context.Background(), pan, uuid.New())
	require.NoError(t, err)

	repo.On("FindByValue", mock.Anything, created.Value).Return(saved, nil)

	got, err := usecase.NewDetokenize(repo, enc).Execute(context.Background(), created.Value)
	require.NoError(t, err)
	assert.Equal(t, pan, got)
}

func TestDetokenizeNotFound(t *testing.T) {
	repo := new(MockTokenRepository)
	repo.On("FindByValue", mock.Anything, "INEXISTENTE").
		Return(nil, domain.ErrTokenNotFound)

	_, err := usecase.NewDetokenize(repo, newEncryptor(t)).
		Execute(context.Background(), "INEXISTENTE")

	assert.ErrorIs(t, err, domain.ErrTokenNotFound)
}
