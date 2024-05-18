package user_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"shin-monta-no-mori/server/internal/app"
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
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := categoriesTest{}
	ctx := c.setUp(t, config)
	defer c.tearDown(t, config)

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
			req, _ := http.NewRequest("GET", "/api/v1/categories/list", nil)

			ctx.Server.Router.ServeHTTP(w, req)

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

func (c categoriesTest) setUp(t *testing.T, config util.Config) *app.AppContext {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src, filename)
		VALUES
		(99999, 'test_parent_category_name_99999', 'test_parent_category_src_99999.com', 'test_parent_category_filename_99999');
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(99999, 'test_child_category_name_99999', 99999);
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
