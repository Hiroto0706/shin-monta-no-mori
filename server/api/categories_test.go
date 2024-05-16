package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"shin-monta-no-mori/server/api"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

type categoriesTest struct{}

func TestListCategories(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := categoriesTest{}
	s := c.setUp(t, config)
	defer c.tearDown(t, config)

	// 認証用トークンの生成
	accessToken := setAuthUser(t, s)

	type args struct {
		compareLimit int
	}

	tests := []struct {
		name         string
		arg          args
		want         []model.Category
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg: args{
				compareLimit: 1,
			},
			want: []model.Category{
				{
					ParentCategory: db.ParentCategory{
						ID:   99999,
						Name: "test_parent_category_name_99999",
						Src:  "test_parent_category_src_99999.com",
						Filename: sql.NullString{
							String: "test_parent_category_filename_99999",
							Valid:  true,
						},
					},
					ChildCategory: []db.ChildCategory{
						{
							ID:       99999,
							Name:     "test_child_category_name_99999",
							ParentID: 99999,
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/categories/list", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			s.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got []model.Category
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got[:tt.arg.compareLimit] {
					compareCategoriesObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := categoriesTest{}
	s := c.setUp(t, config)
	defer c.tearDown(t, config)

	// 認証用トークンの生成
	accessToken := setAuthUser(t, s)

	type args struct {
		id string
	}
	tests := []struct {
		name         string
		arg          args
		want         model.Category
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg: args{
				id: "10001",
			},
			want: model.Category{
				ParentCategory: db.ParentCategory{
					ID:   10001,
					Name: "test_parent_category_name_10001",
					Src:  "test_parent_category_src_10001.com",
					Filename: sql.NullString{
						String: "test_parent_category_filename_10001",
						Valid:  true,
					},
				},
				ChildCategory: []db.ChildCategory{
					{
						ID:       10001,
						Name:     "test_child_category_name_10001",
						ParentID: 10001,
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（ParentCategoryは存在するが、ChildCategoryが存在しない場合）",
			arg: args{
				id: "10002",
			},
			want: model.Category{
				ParentCategory: db.ParentCategory{
					ID:   10002,
					Name: "test_parent_category_name_10002",
					Src:  "test_parent_category_src_10002.com",
					Filename: sql.NullString{
						String: "test_parent_category_filename_10002",
						Valid:  true,
					},
				},
				ChildCategory: []db.ChildCategory{},
			},
			wantErr:      true,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idの値が不正な場合）",
			arg: args{
				id: "aaa",
			},
			want: model.Category{
				ParentCategory: db.ParentCategory{},
				ChildCategory:  []db.ChildCategory{},
			},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（存在しないParentCategoryを取得しようとした場合）",
			arg: args{
				id: "999999",
			},
			want: model.Category{
				ParentCategory: db.ParentCategory{},
				ChildCategory:  []db.ChildCategory{},
			},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/categories/"+tt.arg.id, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			s.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got model.Category
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareCategoriesObjects(t, got, tt.want, ignoreFields)
			}
		})
	}
}

func TestSearchCategories(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := categoriesTest{}
	s := c.setUp(t, config)
	defer c.tearDown(t, config)

	// 認証用トークンの生成
	accessToken := setAuthUser(t, s)

	type args struct {
		q            string
		compareLimit int
	}
	tests := []struct {
		name         string
		arg          args
		want         []model.Category
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg: args{
				q:            "test_parent_category_name_10003",
				compareLimit: 1,
			},
			want: []model.Category{
				{
					ParentCategory: db.ParentCategory{
						ID:   10003,
						Name: "test_parent_category_name_10003",
						Src:  "test_parent_category_src_10003.com",
						Filename: sql.NullString{
							String: "test_parent_category_filename_10003",
							Valid:  true,
						},
					},
					ChildCategory: []db.ChildCategory{
						{
							ID:       10003,
							Name:     "test_child_category_name_10003",
							ParentID: 10003,
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（queryを何も指定しない場合）",
			arg: args{
				q:            "",
				compareLimit: 1,
			},
			want: []model.Category{
				{
					ParentCategory: db.ParentCategory{
						ID:   99999,
						Name: "test_parent_category_name_99999",
						Src:  "test_parent_category_src_99999.com",
						Filename: sql.NullString{
							String: "test_parent_category_filename_99999",
							Valid:  true,
						},
					},
					ChildCategory: []db.ChildCategory{
						{
							ID:       99999,
							Name:     "test_child_category_name_99999",
							ParentID: 99999,
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（存在しないparent_categoriesを検索した場合）",
			arg: args{
				q:            "not exist category",
				compareLimit: 0,
			},
			want: []model.Category{
				{
					ParentCategory: db.ParentCategory{},
					ChildCategory:  []db.ChildCategory{},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（child_categoryがないparent_categoryを検索した場合）",
			arg: args{
				q:            "test_parent_category_name_10004",
				compareLimit: 1,
			},
			want: []model.Category{
				{
					ParentCategory: db.ParentCategory{
						ID:   10004,
						Name: "test_parent_category_name_10004",
						Src:  "test_parent_category_src_10004.com",
						Filename: sql.NullString{
							String: "test_parent_category_filename_10004",
							Valid:  true,
						},
					},
					ChildCategory: []db.ChildCategory{},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/categories/search?q="+tt.arg.q, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			s.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got []model.Category
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got[:tt.arg.compareLimit] {
					compareCategoriesObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func TestEditParentCategory(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := categoriesTest{}
	s := c.setUp(t, config)
	defer c.tearDown(t, config)

	// 認証用トークンの生成
	accessToken := setAuthUser(t, s)

	type args struct {
		ID       string
		Name     string
		Filename string
	}
	tests := []struct {
		name         string
		arg          args
		prepare      func() (*bytes.Buffer, string)
		want         db.ParentCategory
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg: args{
				ID: "11001",
			},
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				// テキストフィールドを追加
				_ = writer.WriteField("name", "test_parent_category_name_11001_edited")
				_ = writer.WriteField("filename", "test_parent_category_filename_11001")

				require.NoError(t, err)

				return body, writer.FormDataContentType()
			},
			want: db.ParentCategory{
				ID:   11001,
				Name: "test_parent_category_name_11001_edited",
				Filename: sql.NullString{
					String: "test_parent_category_filename_11001",
					Valid:  true,
				},
				Src: "test_parent_category_src_11001.com",
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idの値が不正な場合）",
			arg: args{
				ID: "aaa",
			},
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()
				return body, writer.FormDataContentType()
			},
			want:         db.ParentCategory{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（存在しないparent_categoryのIDを指定した場合）",
			arg: args{
				ID: "999999",
			},
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				return body, writer.FormDataContentType()
			},
			want:         db.ParentCategory{},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, contentType := tt.prepare()
			req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/parent/"+tt.arg.ID, body)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			s.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					ParentCategory db.ParentCategory `json:"parent_category"`
					Message        string            `json:"message"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareParentCategoryObjects(t, got.ParentCategory, tt.want, ignoreFields)
			}
		})
	}
}

func TestEditChildCategory(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := categoriesTest{}
	s := c.setUp(t, config)
	defer c.tearDown(t, config)

	// 認証用トークンの生成
	accessToken := setAuthUser(t, s)

	type args struct {
		ID       string
		Name     string
		Filename string
	}
	tests := []struct {
		name         string
		arg          args
		prepare      func() (*bytes.Buffer, string)
		want         db.ChildCategory
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg: args{
				ID: "12001",
			},
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				// テキストフィールドを追加
				_ = writer.WriteField("name", "test_child_category_name_12001_edited")
				_ = writer.WriteField("parent_id", "12001")

				require.NoError(t, err)

				return body, writer.FormDataContentType()
			},
			want: db.ChildCategory{
				ID:       12001,
				Name:     "test_child_category_name_12001_edited",
				ParentID: 12001,
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idの値が不正な場合）",
			arg: args{
				ID: "aaa",
			},
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				require.NoError(t, err)

				return body, writer.FormDataContentType()
			},
			want:         db.ChildCategory{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（存在しないchild_categoryを編集しようとしている場合）",
			arg: args{
				ID: "999999",
			},
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				require.NoError(t, err)

				return body, writer.FormDataContentType()
			},
			want:         db.ChildCategory{},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, contentType := tt.prepare()
			req, _ := http.NewRequest("PUT", "/api/v1/admin/categories/child/"+tt.arg.ID, body)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			s.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					ChildCategory db.ChildCategory `json:"child_category"`
					Message       string           `json:"message"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareChildCategoryObjects(t, got.ChildCategory, tt.want, ignoreFields)
			}
		})
	}
}
func TestDeleteChildCategory(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := categoriesTest{}
	s := c.setUp(t, config)
	defer c.tearDown(t, config)

	// 認証用トークンの生成
	accessToken := setAuthUser(t, s)

	type args struct {
		ID string
	}
	tests := []struct {
		name         string
		arg          args
		want         string
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg: args{
				ID: "13001",
			},
			want:         "child_categoryの削除に成功しました",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idの値が不正な場合）",
			arg: args{
				ID: "aaa",
			},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（存在しないchild_categoryを削除しようとした場合）",
			arg: args{
				ID: "999999",
			},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "異常系（idの値が0の場合）",
			arg: args{
				ID: "0",
			},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "異常系（idの値が負の場合）",
			arg: args{
				ID: "0",
			},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/v1/admin/categories/child/"+tt.arg.ID, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			s.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Message string `json:"message"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Message)
			}
		})
	}
}

func compareCategoriesObjects(t *testing.T, got model.Category, want model.Category, ignoreFieldsMap map[string][]string) {
	// 親カテゴリー比較
	compareParentCategoryObjects(t, got.ParentCategory, want.ParentCategory, ignoreFieldsMap)

	// 子カテゴリー比較
	for k, gcc := range got.ChildCategory {
		compareChildCategoryObjects(t, gcc, want.ChildCategory[k], ignoreFieldsMap)
	}
}

func compareParentCategoryObjects(t *testing.T, got db.ParentCategory, want db.ParentCategory, ignoreFieldsMap map[string][]string) {
	// 親カテゴリー比較
	if d := cmp.Diff(got, want, cmpopts.IgnoreFields(got, ignoreFieldsMap["Other"]...)); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}
func compareChildCategoryObjects(t *testing.T, got db.ChildCategory, want db.ChildCategory, ignoreFieldsMap map[string][]string) {
	if d := cmp.Diff(got, want, cmpopts.IgnoreFields(got, ignoreFieldsMap["Other"]...)); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func (c categoriesTest) setUp(t *testing.T, config util.Config) *api.Server {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src, filename)
		VALUES
		(99999, 'test_parent_category_name_99999', 'test_parent_category_src_99999.com', 'test_parent_category_filename_99999'),
		(10001, 'test_parent_category_name_10001', 'test_parent_category_src_10001.com', 'test_parent_category_filename_10001'),
		(10002, 'test_parent_category_name_10002', 'test_parent_category_src_10002.com', 'test_parent_category_filename_10002'),
		(10003, 'test_parent_category_name_10003', 'test_parent_category_src_10003.com', 'test_parent_category_filename_10003'),
		(10004, 'test_parent_category_name_10004', 'test_parent_category_src_10004.com', 'test_parent_category_filename_10004'),
		(11001, 'test_parent_category_name_11001', 'test_parent_category_src_11001.com', 'test_parent_category_filename_11001'),
		(12001, 'test_parent_category_name_12001', 'test_parent_category_src_12001.com', 'test_parent_category_filename_12001'),
		(13001, 'test_parent_category_name_13001', 'test_parent_category_src_13001.com', 'test_parent_category_filename_13001');
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(99999, 'test_child_category_name_99999', 99999),
		(10001, 'test_child_category_name_10001', 10001),
		(10003, 'test_child_category_name_10003', 10003),
		(12001, 'test_child_category_name_12001', 12001),
		(13001, 'test_child_category_name_13001', 13001);
		`),
	}

	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to exec query: %v", err)
		}
	}

	server, err := newTestServer(store, config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	return server
}

func (c categoriesTest) tearDown(t *testing.T, config util.Config) {
	store := createConn(config)

	queries := []string{
		"TRUNCATE TABLE child_categories RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE parent_categories RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE operators RESTART IDENTITY CASCADE;",
	}
	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}
}
