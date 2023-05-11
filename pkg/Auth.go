package pkg

import "context"

type Auth struct {
	Token string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"token": a.Token}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}
