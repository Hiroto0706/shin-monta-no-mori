package db

import (
	"context"
	"database/sql"
	"fmt"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, db DBTX) {
	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(10001, 'test_image_title_10001', 'test_image_original_src_10001', 'test_image_simple_src_10001');
		`),
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(20001, 'test_image_title_20001', 'test_image_original_src_20001', 'test_image_simple_src_20001'),
		(20002, 'test_image_title_20002', 'test_image_original_src_20002', '');
		`),
	}
	for _, query := range queries {
		if _, err := db.ExecContext(context.Background(), query); err != nil {
			t.Fatalf("Failed to exec query: %v", err)
		}
	}
}

func tearDown(t *testing.T, db DBTX) {
	query := "TRUNCATE TABLE images RESTART IDENTITY CASCADE;"
	if _, err := db.ExecContext(context.Background(), query); err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}
}

func TestCreateImage(t *testing.T) {
	defer tearDown(t, testQueries.db)

	tests := []struct {
		name    string
		arg     CreateImageParams
		wantErr bool
	}{
		{
			name: "正常系",
			arg: CreateImageParams{
				Title:       util.RandomTitle(),
				OriginalSrc: util.RandomTitle(),
				SimpleSrc: sql.NullString{
					String: util.RandomTitle(),
					Valid:  true,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（SimpleSrcが空文字の場合）",
			arg: CreateImageParams{
				Title:       util.RandomTitle(),
				OriginalSrc: util.RandomTitle(),
				SimpleSrc: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			wantErr: false,
		},
		{
			name: "異常系（titleが空文字の場合）",
			arg: CreateImageParams{
				Title:       "",
				OriginalSrc: util.RandomTitle(),
				SimpleSrc: sql.NullString{
					String: util.RandomTitle(),
					Valid:  true,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, err := testQueries.CreateImage(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, image)
				require.Equal(t, tt.arg.Title, image.Title)
				require.Equal(t, tt.arg.OriginalSrc, image.OriginalSrc)
				require.Equal(t, tt.arg.SimpleSrc.String, image.SimpleSrc.String)
				require.NotZero(t, image.ID)
				require.NotZero(t, image.CreatedAt)
			}
		})
	}
}

func TestDeleteImage(t *testing.T) {
	setup(t, testQueries.db)
	defer tearDown(t, testQueries.db)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "正常系",
			id:      10001,
			wantErr: false,
		},
		{
			name:    "正常系（存在しないIDを消そうとする時）",
			id:      99999,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testQueries.DeleteImage(context.Background(), tt.id)
			if !tt.wantErr {
				require.NoError(t, err)
				_, err := testQueries.GetImage(context.Background(), tt.id)
				require.Error(t, err, "The image should no longer exist.")
			}
		})
	}
}

func TestGetImage(t *testing.T) {
	setup(t, testQueries.db)
	defer tearDown(t, testQueries.db)

	type args struct {
		id int64
	}
	type wants struct {
		ID          int64
		Title       string
		OriginalSrc string
		SimpleSrc   sql.NullString
	}
	tests := []struct {
		name    string
		arg     args
		want    wants
		wantErr bool
	}{
		{
			name: "正常系",
			arg: args{
				id: 20001,
			},
			want: wants{
				ID:          20001,
				Title:       "test_image_title_20001",
				OriginalSrc: "test_image_original_src_20001",
				SimpleSrc: sql.NullString{
					String: "test_image_simple_src_20001",
					Valid:  true,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（SimpleSrcが空文字の場合）",
			arg: args{
				id: 20002,
			},
			want: wants{
				ID:          20002,
				Title:       "test_image_title_20002",
				OriginalSrc: "test_image_original_src_20002",
				SimpleSrc: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			wantErr: false,
		},
		{
			name: "異常系（存在しないIDの場合）",
			arg: args{
				id: 99999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, err := testQueries.GetImage(context.Background(), tt.arg.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, image.ID)
				require.Equal(t, tt.want.Title, image.Title)
				require.Equal(t, tt.want.OriginalSrc, image.OriginalSrc)
				require.Equal(t, tt.want.SimpleSrc.String, image.SimpleSrc.String)
				require.NotZero(t, image.CreatedAt)
			}
		})
	}
}
