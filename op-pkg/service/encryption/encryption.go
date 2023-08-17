package encryption

import (
	"context"

	"github.com/grafana/grafana/pkg/infra/usagestats"
	"github.com/grafana/grafana/pkg/services/encryption"
	"github.com/grafana/grafana/pkg/setting"
)

type Service struct {
	Cipher
	Decipher
}

func ProvideEncryptionService(
	provider encryption.Provider,
	usageMetrics usagestats.Service,
	settingsProvider setting.Provider,
) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) EncryptJsonData(ctx context.Context, kv map[string]string, secret string) (map[string][]byte, error) {
	res := make(map[string][]byte)
	for k, v := range kv {
		res[k] = []byte(v)
	}
	return res, nil
}

func (s *Service) DecryptJsonData(ctx context.Context, sjd map[string][]byte, secret string) (map[string]string, error) {
	res := make(map[string]string)
	for k, v := range sjd {
		res[k] = string(v)
	}
	return res, nil
}

func (s *Service) GetDecryptedValue(ctx context.Context, sjd map[string][]byte, key string, fallback string, secret string) string {
	if value, ok := sjd[key]; ok {
		return string(value)
	}
	return fallback
}

func (s *Service) Validate(section setting.Section) error {
	return nil
}

func (s *Service) Reload(_ setting.Section) error {
	return nil
}
