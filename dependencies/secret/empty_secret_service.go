package secret

import (
	"context"
	"github.com/wolffcm/flux/codes"
	"github.com/wolffcm/flux/internal/errors"
)

func (ess EmptySecretService) LoadSecret(ctx context.Context, k string) (string, error) {
	return "", errors.Newf(codes.NotFound, "secret key %q not found", k)
}

// Secret service that always reports no secrets exist
type EmptySecretService struct {
}
