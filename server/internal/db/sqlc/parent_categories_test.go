package db_test

import (
	"context"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateParentCategory(t *testing.T) {
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.CreateParentCategoryParams
		want    db.ParentCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.CreateParentCategoryParams{
				Name: "test_parent_category_name_00001",
				Src:  "test_parent_category_src_00001",
			},
			want: db.ParentCategory{
				Name: "test_parent_category_name_00001",
				Src:  "test_parent_category_src_00001",
			},
			wantErr: false,
		},
		{
			name: "正常系（Srcが空文字の場合）",
			arg: db.CreateParentCategoryParams{
				Name: "test_parent_category_name_00010",
				Src:  "",
			},
			want: db.ParentCategory{
				Name: "test_parent_category_name_00010",
				Src:  "",
			},
			wantErr: false,
		},
		{
			name: "異常系（nameが空文字の場合）",
			arg: db.CreateParentCategoryParams{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p_category, err := testQueries.CreateParentCategory(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, p_category)
				require.Equal(t, tt.arg.Name, p_category.Name)
				require.Equal(t, tt.arg.Src, p_category.Src)
				require.NotZero(t, p_category.ID)
				require.NotZero(t, p_category.CreatedAt)
			}
		})
	}
}

func TestDeleteParentCategory(t *testing.T) {
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
			err := testQueries.DeleteParentCategory(context.Background(), tt.id)
			if !tt.wantErr {
				require.NoError(t, err)
				_, err := testQueries.GetParentCategory(context.Background(), tt.id)
				require.Error(t, err, "The parent_category should no longer exist.")
			}
		})
	}
}

func TestGetParentCategory(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		arg     args
		want    db.ParentCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: args{
				id: 20001,
			},
			want: db.ParentCategory{
				ID:   20001,
				Name: "test_parent_category_name_20001",
				Src:  "test_parent_category_src_20001",
			},
			wantErr: false,
		},
		{
			name: "正常系（Srcが空文字の場合）",
			arg: args{
				id: 20002,
			},
			want: db.ParentCategory{
				ID:   20002,
				Name: "test_parent_category_name_20002",
				Src:  "",
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
			p_category, err := testQueries.GetParentCategory(context.Background(), tt.arg.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, p_category.ID)
				require.Equal(t, tt.want.Name, p_category.Name)
				require.Equal(t, tt.want.Src, p_category.Src)
				require.NotZero(t, p_category.CreatedAt)
			}
		})
	}
}

func TestListParentCategory(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.ListParentCategoriesParams
		want    []db.ParentCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.ListParentCategoriesParams{
				Limit:  3,
				Offset: 0,
			},
			want: []db.ParentCategory{
				{
					ID:   99992,
					Name: "test_parent_category_name_99992",
					Src:  "test_parent_category_src_99992",
				},
				{
					ID:   99991,
					Name: "test_parent_category_name_99991",
					Src:  "test_parent_category_src_99991",
				},
				{
					ID:   99990,
					Name: "test_parent_category_name_99990",
					Src:  "test_parent_category_src_99990",
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（returnが空の時）",
			arg: db.ListParentCategoriesParams{
				Limit:  3,
				Offset: 1000,
			},
			want:    []db.ParentCategory{{}},
			wantErr: false,
		},
		{
			name: "異常系（argsのLimitの値が不正な場合）",
			arg: db.ListParentCategoriesParams{
				Limit:  -1,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "異常系（argsのOffsetの値が不正な場合）",
			arg: db.ListParentCategoriesParams{
				Limit:  3,
				Offset: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p_categories, err := testQueries.ListParentCategories(context.Background(), tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				for i, pc := range p_categories {
					require.NoError(t, err)
					require.Equal(t, tt.want[i].ID, pc.ID)
					require.Equal(t, tt.want[i].Name, pc.Name)
					require.Equal(t, tt.want[i].Src, pc.Src)
					require.NotZero(t, pc.CreatedAt)
				}
			}
		})
	}
}

func TestUpdateParentCategory(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.UpdateParentCategoryParams
		want    db.ParentCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.UpdateParentCategoryParams{
				ID:   40001,
				Name: "test_parent_category_name_40001_edited",
				Src:  "test_parent_category_src_40001_edited",
			},
			want: db.ParentCategory{
				ID:   40001,
				Name: "test_parent_category_name_40001_edited",
				Src:  "test_parent_category_src_40001_edited",
			},
			wantErr: false,
		},
		{
			name: "正常系（srcが空になる場合）",
			arg: db.UpdateParentCategoryParams{
				ID:   40002,
				Name: "test_parent_category_name_40002_edited",
				Src:  "",
			},
			want: db.ParentCategory{
				ID:   40002,
				Name: "test_parent_category_name_40002_edited",
				Src:  "",
			},
			wantErr: false,
		},
		{
			name: "異常系（存在しないIDを指定している場合）",
			arg: db.UpdateParentCategoryParams{
				ID: 99999,
			},
			wantErr: true,
		},
		{
			name: "異常系（titleが空になる場合）",
			arg: db.UpdateParentCategoryParams{
				ID:   40003,
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p_category, err := testQueries.UpdateParentCategory(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, p_category)
				require.Equal(t, tt.arg.ID, p_category.ID)
				require.Equal(t, tt.arg.Name, p_category.Name)
				require.Equal(t, tt.arg.Src, p_category.Src)
				require.NotZero(t, p_category.CreatedAt)
			}
		})
	}
}
