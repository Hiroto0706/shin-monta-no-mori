package db_test

import (
	"context"
	"database/sql"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateImage(t *testing.T) {
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.CreateImageParams
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.CreateImageParams{
				Title:       util.RandomTitle(),
				OriginalSrc: util.RandomTitle(),
				SimpleSrc: sql.NullString{
					String: util.RandomTitle(),
					Valid:  true,
				},
				OriginalFilename: util.RandomTitle(),
				SimpleFilename: sql.NullString{
					String: util.RandomTitle(),
					Valid:  true,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（SimpleSrcが空文字の場合）",
			arg: db.CreateImageParams{
				Title:            util.RandomTitle(),
				OriginalSrc:      util.RandomTitle(),
				OriginalFilename: util.RandomTitle(),
			},
			wantErr: false,
		},
		{
			name: "異常系（titleが空文字の場合）",
			arg: db.CreateImageParams{
				Title: "",
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
				require.Equal(t, tt.arg.OriginalFilename, image.OriginalFilename)
				require.Equal(t, tt.arg.SimpleFilename, image.SimpleFilename)
				require.NotZero(t, image.ID)
				require.NotZero(t, image.CreatedAt)
			}
		})
	}
}

func TestDeleteImage(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

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
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

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

func TestListImage(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	type wants struct {
		ID          int64
		Title       string
		OriginalSrc string
		SimpleSrc   sql.NullString
	}
	tests := []struct {
		name    string
		arg     db.ListImageParams
		want    []wants
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.ListImageParams{
				Limit:  3,
				Offset: 0,
			},
			want: []wants{
				{
					ID:          99992,
					Title:       "test_image_title_99992",
					OriginalSrc: "test_image_original_src_99992",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_99992",
						Valid:  true,
					},
				},
				{
					ID:          99991,
					Title:       "test_image_title_99991",
					OriginalSrc: "test_image_original_src_99991",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_99991",
						Valid:  true,
					},
				},
				{
					ID:          99990,
					Title:       "test_image_title_99990",
					OriginalSrc: "test_image_original_src_99990",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_99990",
						Valid:  true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（returnが空の時）",
			arg: db.ListImageParams{
				Limit:  3,
				Offset: 1000,
			},
			want:    []wants{{}},
			wantErr: false,
		},
		{
			name: "異常系（argsの値が不正な場合）",
			arg: db.ListImageParams{
				Limit:  -1,
				Offset: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images, err := testQueries.ListImage(context.Background(), tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				for i, image := range images {
					require.NoError(t, err)
					require.Equal(t, tt.want[i].ID, image.ID)
					require.Equal(t, tt.want[i].Title, image.Title)
					require.Equal(t, tt.want[i].OriginalSrc, image.OriginalSrc)
					require.Equal(t, tt.want[i].SimpleSrc.String, image.SimpleSrc.String)
					require.NotZero(t, image.CreatedAt)
				}
			}
		})
	}
}

func TestUpdateImage(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.UpdateImageParams
		want    db.Image
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.UpdateImageParams{
				ID:          40001,
				Title:       "test_image_title_40001_edited",
				OriginalSrc: "test_image_original_src_40001_edited",
				SimpleSrc: sql.NullString{
					String: "test_image_simple_src_40001_edited",
					Valid:  true,
				},
			},
			want: db.Image{
				ID:          40001,
				Title:       "test_image_title_40001_edited",
				OriginalSrc: "test_image_original_src_40001_edited",
				SimpleSrc: sql.NullString{
					String: "test_image_simple_src_40001_edited",
					Valid:  true,
				},
			},
			wantErr: false,
		},
		{
			name: "異常系（存在しないIDを指定している場合）",
			arg: db.UpdateImageParams{
				ID: 99999,
			},
			wantErr: true,
		},
		{
			name: "異常系（titleが空になる場合）",
			arg: db.UpdateImageParams{
				ID:    40002,
				Title: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, err := testQueries.UpdateImage(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, image)
				require.Equal(t, tt.arg.ID, image.ID)
				require.Equal(t, tt.arg.Title, image.Title)
				require.Equal(t, tt.arg.OriginalSrc, image.OriginalSrc)
				require.Equal(t, tt.arg.SimpleSrc.String, image.SimpleSrc.String)
				require.NotZero(t, image.CreatedAt)
			}
		})
	}
}
