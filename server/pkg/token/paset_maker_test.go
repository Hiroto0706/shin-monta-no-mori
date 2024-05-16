package token_test

import (
	"errors"
	"log"
	"shin-monta-no-mori/server/pkg/token"
	"shin-monta-no-mori/server/pkg/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

func TestPasetoMaker(t *testing.T) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		setupToken    func() (string, *token.Payload, error)
		checkResponse func(t *testing.T, payload *token.Payload, err error)
	}{
		{
			name: "正常系",
			setupToken: func() (string, *token.Payload, error) {
				return maker.CreateToken("testuser", time.Minute)
			},
			checkResponse: func(t *testing.T, payload *token.Payload, err error) {
				require.NoError(t, err)
				require.NotNil(t, payload)
				require.NotZero(t, payload.ID)
				require.Equal(t, "testuser", payload.Username)
				require.WithinDuration(t, time.Now().Add(time.Minute), payload.ExpiredAt, time.Second)
			},
		},
		{
			name: "異常系（InvalidToken）",
			setupToken: func() (string, *token.Payload, error) {
				return "invalidToken", nil, ErrInvalidToken
			},
			checkResponse: func(t *testing.T, payload *token.Payload, err error) {
				require.Error(t, err)
				require.EqualError(t, err, ErrInvalidToken.Error())
				require.Nil(t, payload)
			},
		},
		{
			name: "異常系（ExpiredToken）",
			setupToken: func() (string, *token.Payload, error) {
				return maker.CreateToken("testuser", -time.Minute)
			},
			checkResponse: func(t *testing.T, payload *token.Payload, err error) {
				require.Error(t, err)
				require.EqualError(t, err, ErrExpiredToken.Error())
				require.Nil(t, payload)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenStr, _, err := tc.setupToken()
			if err == nil {
				payload, err := maker.VerifyToken(tokenStr)
				tc.checkResponse(t, payload, err)
			} else {
				tc.checkResponse(t, nil, err)
			}
		})
	}
}
