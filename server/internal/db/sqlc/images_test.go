package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(30001, 'test_image_title_30001', 'test_image_original_src_30001', 'test_image_simple_src_30001'),
		(30002, 'test_image_title_30002', 'test_image_original_src_30002', 'test_image_simple_src_30002'),
		(30003, 'test_image_title_30003', 'test_image_original_src_30003', 'test_image_simple_src_30003');
		`),
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(40001, 'test_image_title_40001', 'test_image_original_src_40001', 'test_image_simple_src_40001'),
		(40002, 'test_image_title_40002', 'test_image_original_src_40002', 'test_image_simple_src_40002'),
		(40003, 'test_image_title_40003', 'test_image_original_src_40003', 'test_image_simple_src_40003');
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

func TestListImage(t *testing.T) {
	setup(t, testQueries.db)
	defer tearDown(t, testQueries.db)

	type wants struct {
		ID          int64
		Title       string
		OriginalSrc string
		SimpleSrc   sql.NullString
	}
	tests := []struct {
		name    string
		arg     ListImageParams
		want    []wants
		wantErr bool
	}{
		{
			name: "正常系",
			arg: ListImageParams{
				Limit:  3,
				Offset: 0,
			},
			want: []wants{
				{
					ID:          30003,
					Title:       "test_image_title_30003",
					OriginalSrc: "test_image_original_src_30003",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_30003",
						Valid:  true,
					},
				},
				{
					ID:          30002,
					Title:       "test_image_title_30002",
					OriginalSrc: "test_image_original_src_30002",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_30002",
						Valid:  true,
					},
				},
				{
					ID:          30001,
					Title:       "test_image_title_30001",
					OriginalSrc: "test_image_original_src_30001",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_30001",
						Valid:  true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（returnが空の時）",
			arg: ListImageParams{
				Limit:  3,
				Offset: 1000,
			},
			want:    []wants{{}},
			wantErr: false,
		},
		{
			name: "異常系（argsの値が不正な場合）",
			arg: ListImageParams{
				Limit:  -1,
				Offset: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images, err := testQueries.ListImage(context.Background(), tt.arg)
			log.Println(images)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				for i, image := range images {
					log.Println(i, image)
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
	setup(t, testQueries.db)
	defer tearDown(t, testQueries.db)

	tests := []struct {
		name    string
		arg     UpdateImageParams
		want    UpdateImageParams
		wantErr bool
	}{
		{
			name: "正常系",
			arg: UpdateImageParams{
				ID:          40001,
				Title:       "test_image_title_40001_edited",
				OriginalSrc: "test_image_original_src_40001_edited",
				SimpleSrc: sql.NullString{
					String: "test_image_simple_src_40001_edited",
					Valid:  true,
				},
			},
			want: UpdateImageParams{
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
			arg: UpdateImageParams{
				ID: 99999,
			},
			wantErr: true,
		},
		{
			name: "異常系（titleが空になる場合）",
			arg: UpdateImageParams{
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
				require.Equal(t, tt.arg.Title, image.Title)
				require.Equal(t, tt.arg.OriginalSrc, image.OriginalSrc)
				require.Equal(t, tt.arg.SimpleSrc.String, image.SimpleSrc.String)
				require.NotZero(t, image.ID)
				require.NotZero(t, image.CreatedAt)
			}
		})
	}
}
