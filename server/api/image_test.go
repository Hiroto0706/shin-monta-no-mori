package api_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"shin-monta-no-mori/server/api"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestListImages(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	server := setUp(t, config)
	defer tearDown(t, config)

	type args struct {
		page            string
		imageFetchLimit int
	}

	tests := []struct {
		name         string
		arg          args
		want         []model.Illustration
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系（p=0のとき）",
			arg: args{
				page:            "0",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image: db.Image{
						ID:          999991,
						Title:       "test_image_title_999991",
						OriginalSrc: "test_image_original_src_999991.com",
						SimpleSrc: sql.NullString{
							String: "test_image_simple_src_999991.com",
							Valid:  true,
						},
					},
					Character: []db.Character{
						{
							ID:   11001,
							Name: "test_character_name_11001",
							Src:  "test_character_src_11001.com",
						},
					},
					Category: []*model.Category{
						{
							ParentCategory: db.ParentCategory{
								ID:   11001,
								Name: "test_parent_category_name_11001",
								Src:  "test_parent_category_src_11001.com",
							},
							ChildCategory: []db.ChildCategory{
								{
									ID:       11001,
									Name:     "test_child_category_name_11001",
									ParentID: 11001,
								},
							},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（p=1のとき）",
			arg: args{
				page:            "1",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image: db.Image{
						ID:          999990,
						Title:       "test_image_title_999990",
						OriginalSrc: "test_image_original_src_999990.com",
						SimpleSrc: sql.NullString{
							String: "test_image_simple_src_999990.com",
							Valid:  true,
						},
					},
					Character: []db.Character{
						{
							ID:   11001,
							Name: "test_character_name_11001",
							Src:  "test_character_src_11001.com",
						},
					},
					Category: []*model.Category{
						{
							ParentCategory: db.ParentCategory{
								ID:   11001,
								Name: "test_parent_category_name_11001",
								Src:  "test_parent_category_src_11001.com",
							},
							ChildCategory: []db.ChildCategory{
								{
									ID:       11001,
									Name:     "test_child_category_name_11001",
									ParentID: 11001,
								},
							},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（クエリパラメータの値が不正な場合)",
			arg: args{
				page:            "a",
				imageFetchLimit: 1,
			},
			want:         []model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 取得するイメージの数を1にする
			server.Config.ImageFetchLimit = tt.arg.imageFetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/illustrations/list/?p="+tt.arg.page, nil)
			server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got []model.Illustration
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				compareIllustrationsObjects(t, got, tt.want)
			}
		})
	}
}

func compareIllustrationsObjects(t *testing.T, got []model.Illustration, want []model.Illustration) {
	for i, g := range got {
		// イメージ比較
		if d := cmp.Diff(g.Image, want[i].Image, cmpopts.IgnoreFields(g.Image, "CreatedAt", "UpdatedAt")); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}

		// キャラクター比較
		for j, gch := range g.Character {
			if d := cmp.Diff(gch, want[i].Character[j], cmpopts.IgnoreFields(gch, "CreatedAt", "UpdatedAt")); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		}

		// カテゴリー比較
		for j, gca := range g.Category {
			// 親カテゴリー比較
			if d := cmp.Diff(gca.ParentCategory, want[i].Category[j].ParentCategory, cmpopts.IgnoreFields(gca.ParentCategory, "CreatedAt", "UpdatedAt")); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}

			// 子カテゴリー比較
			for k, gcca := range gca.ChildCategory {
				if d := cmp.Diff(gcca, want[i].Category[j].ChildCategory[k], cmpopts.IgnoreFields(gcca, "CreatedAt", "UpdatedAt")); len(d) != 0 {
					t.Errorf("differs: (-got +want)\n%s", d)
				}
			}
		}
	}
}

func createConn(config util.Config) *db.Store {
	conn, err := sql.Open(config.DBDriver, config.TestDBUrl)
	if err != nil {
		log.Fatal("cannot connect to test db: ", err)
	}

	return db.NewStore(conn)
}

func newTestServer(store *db.Store, config util.Config) (*api.Server, error) {
	server := &api.Server{
		Config: config,
		Store:  store,
	}

	router := gin.Default()
	// router.Use(api.CORSMiddleware())
	server.Router = router

	api.SetAdminRouters(server)

	return server, nil
}

func setUp(t *testing.T, config util.Config) *api.Server {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(999990, 'test_image_title_999990', 'test_image_original_src_999990.com', 'test_image_simple_src_999990.com'),
		(999991, 'test_image_title_999991', 'test_image_original_src_999991.com', 'test_image_simple_src_999991.com');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(11001, 'test_character_name_11001', 'test_character_src_11001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_characters_relations (id, image_id, character_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(11001, 'test_parent_category_name_11001', 'test_parent_category_src_11001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_parent_categories_relations (id, image_id, parent_category_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001);
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(11001, 'test_child_category_name_11001', 11001);
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

func tearDown(t *testing.T, config util.Config) {
	store := createConn(config)

	queries := []string{
		"TRUNCATE TABLE images RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE characters RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE operators RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE parent_categories RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE child_categories RESTART IDENTITY CASCADE;",
	}
	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}
}
