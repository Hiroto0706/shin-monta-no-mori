package db_test

import (
	"context"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateChildCategory(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.CreateChildCategoryParams
		want    db.ChildCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.CreateChildCategoryParams{
				Name:     "test_child_category_name_00001",
				ParentID: 50001,
			},
			want: db.ChildCategory{
				Name:     "test_child_category_name_00001",
				ParentID: 50001,
			},
			wantErr: false,
		},
		{
			name: "異常系（ParentIDがnullの場合）",
			arg: db.CreateChildCategoryParams{
				Name: "test_child_category_name_00010",
			},
			want: db.ChildCategory{
				Name: "test_child_category_name_00010",
			},
			wantErr: true,
		},
		{
			name: "異常系（nameが空文字の場合）",
			arg: db.CreateChildCategoryParams{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c_category, err := testQueries.CreateChildCategory(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, c_category)
				require.Equal(t, tt.arg.Name, c_category.Name)
				require.Equal(t, tt.arg.ParentID, c_category.ParentID)
				require.NotZero(t, c_category.ID)
				require.NotZero(t, c_category.CreatedAt)
			}
		})
	}
}

func TestDeleteChildCategory(t *testing.T) {
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
			err := testQueries.DeleteChildCategory(context.Background(), tt.id)
			if !tt.wantErr {
				require.NoError(t, err)
				_, err := testQueries.GetChildCategory(context.Background(), tt.id)
				require.Error(t, err, "The child_category should no longer exist.")
			}
		})
	}
}

func TestGetChildCategory(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		arg     args
		want    db.ChildCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: args{
				id: 20001,
			},
			want: db.ChildCategory{
				ID:       20001,
				Name:     "test_child_category_name_20001",
				ParentID: 70001,
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
			c_category, err := testQueries.GetChildCategory(context.Background(), tt.arg.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, c_category.ID)
				require.Equal(t, tt.want.Name, c_category.Name)
				require.Equal(t, tt.want.ParentID, c_category.ParentID)
				require.NotZero(t, c_category.CreatedAt)
			}
		})
	}
}

func TestListChildCategory(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.ListChildCategoriesParams
		want    []db.ChildCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.ListChildCategoriesParams{
				Limit:  3,
				Offset: 0,
			},
			want: []db.ChildCategory{
				{
					ID:       99993,
					Name:     "test_child_category_name_99993",
					ParentID: 80003,
				},
				{
					ID:       99992,
					Name:     "test_child_category_name_99992",
					ParentID: 80002,
				},
				{
					ID:       99991,
					Name:     "test_child_category_name_99991",
					ParentID: 80001,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（returnが空の時）",
			arg: db.ListChildCategoriesParams{
				Limit:  3,
				Offset: 1000,
			},
			want:    []db.ChildCategory{{}},
			wantErr: false,
		},
		{
			name: "異常系（argsのLimitの値が不正な場合）",
			arg: db.ListChildCategoriesParams{
				Limit:  -1,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "異常系（argsのOffsetの値が不正な場合）",
			arg: db.ListChildCategoriesParams{
				Limit:  3,
				Offset: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c_categories, err := testQueries.ListChildCategories(context.Background(), tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				for i, pc := range c_categories {
					require.NoError(t, err)
					require.Equal(t, tt.want[i].ID, pc.ID)
					require.Equal(t, tt.want[i].Name, pc.Name)
					require.Equal(t, tt.want[i].ParentID, pc.ParentID)
					require.NotZero(t, pc.CreatedAt)
				}
			}
		})
	}
}

func TestUpdateChildCategory(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.UpdateChildCategoryParams
		want    db.ChildCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.UpdateChildCategoryParams{
				ID:       30001,
				Name:     "test_child_category_name_30001_edited",
				ParentID: 90001,
			},
			want: db.ChildCategory{
				ID:       30001,
				Name:     "test_child_category_name_30001_edited",
				ParentID: 90001,
			},
			wantErr: false,
		},
		{
			name: "異常系（存在しないIDを指定している場合）",
			arg: db.UpdateChildCategoryParams{
				ID: 99999,
			},
			wantErr: true,
		},
		{
			name: "異常系（titleが空になる場合）",
			arg: db.UpdateChildCategoryParams{
				ID:   30002,
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "異常系（ParentIDが不正な場合）",
			arg: db.UpdateChildCategoryParams{
				ID:       30003,
				ParentID: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p_category, err := testQueries.UpdateChildCategory(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, p_category)
				require.Equal(t, tt.arg.ID, p_category.ID)
				require.Equal(t, tt.arg.Name, p_category.Name)
				require.Equal(t, tt.arg.ParentID, p_category.ParentID)
				require.NotZero(t, p_category.CreatedAt)
			}
		})
	}
}

func TestGetChildCategoriesByParentID(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	type args struct {
		parentID int64
	}
	tests := []struct {
		name    string
		arg     args
		want    []db.ChildCategory
		wantErr bool
	}{
		{
			name: "正常系",
			arg: args{
				parentID: 90003,
			},
			want: []db.ChildCategory{
				{
					ID:       40001,
					Name:     "test_child_category_name_40001",
					ParentID: 90003,
				},
				{
					ID:       40002,
					Name:     "test_child_category_name_40002",
					ParentID: 90003,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（存在しないIDの場合）",
			arg: args{
				parentID: 99999,
			},
			want:    []db.ChildCategory{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c_categories, err := testQueries.GetChildCategoriesByParentID(context.Background(), tt.arg.parentID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				for i, pc := range c_categories {
					require.NoError(t, err)
					require.Equal(t, tt.want[i].ID, pc.ID)
					require.Equal(t, tt.want[i].Name, pc.Name)
					require.Equal(t, tt.want[i].ParentID, pc.ParentID)
					require.NotZero(t, pc.CreatedAt)
				}
			}
		})
	}
}
