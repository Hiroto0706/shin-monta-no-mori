package db

import (
	"context"
	"database/sql"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateImage(t *testing.T) {
	arg := CreateImageParams{
		Title:       util.RandomTitle(),
		OriginalSrc: util.RandomTitle(),
		SimpleSrc: sql.NullString{
			String: util.RandomTitle(),
			Valid:  true,
		},
	}

	image, err := testQueries.CreateImage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, image)

	require.Equal(t, arg.Title, image.Title)
	require.Equal(t, arg.OriginalSrc, image.OriginalSrc)
	require.Equal(t, arg.SimpleSrc, image.SimpleSrc)

	require.NotZero(t, image.ID)
	require.NotZero(t, image.CreatedAt)
}
