package api_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestListIllustrations(t *testing.T) {
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
			name: "正常系（データが存在しない場合)",
			arg: args{
				page:            "100",
				imageFetchLimit: 100,
			},
			want:         []model.Illustration{},
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
				for i, g := range got {
					compareIllustrationsObjects(t, g, tt.want[i])
				}
			}
		})
	}
}

func TestGetIllustration(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	server := setUp(t, config)
	defer tearDown(t, config)

	type args struct {
		id string
	}

	tests := []struct {
		name         string
		arg          args
		want         model.Illustration
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg: args{
				id: "11001",
			},
			want: model.Illustration{
				Image: db.Image{
					ID:          11001,
					Title:       "test_image_title_11001",
					OriginalSrc: "test_image_original_src_11001.com",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_11001.com",
						Valid:  true,
					},
				},
				Character: []db.Character{
					{
						ID:   11002,
						Name: "test_character_name_11002",
						Src:  "test_character_src_11002.com",
					},
				},
				Category: []*model.Category{
					{
						ParentCategory: db.ParentCategory{
							ID:   11002,
							Name: "test_parent_category_name_11002",
							Src:  "test_parent_category_src_11002.com",
						},
						ChildCategory: []db.ChildCategory{
							{
								ID:       11002,
								Name:     "test_child_category_name_11002",
								ParentID: 11002,
							},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idが不正な値の場合）",
			arg: args{
				id: "aaa",
			},
			want:         model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（存在しないidを指定している場合）",
			arg: args{
				id: "999999",
			},
			want:         model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/illustrations/"+tt.arg.id, nil)
			server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got model.Illustration
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				compareIllustrationsObjects(t, got, tt.want)
			}
		})
	}
}

func compareIllustrationsObjects(t *testing.T, got model.Illustration, want model.Illustration) {
	// イメージ比較
	if d := cmp.Diff(got.Image, want.Image, cmpopts.IgnoreFields(got.Image, "CreatedAt", "UpdatedAt")); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}

	// キャラクター比較
	for i, gch := range got.Character {
		if d := cmp.Diff(gch, want.Character[i], cmpopts.IgnoreFields(gch, "CreatedAt", "UpdatedAt")); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}
	}

	// カテゴリー比較
	for j, gca := range got.Category {
		// 親カテゴリー比較
		if d := cmp.Diff(gca.ParentCategory, want.Category[j].ParentCategory, cmpopts.IgnoreFields(gca.ParentCategory, "CreatedAt", "UpdatedAt")); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}

		// 子カテゴリー比較
		for k, gcca := range gca.ChildCategory {
			if d := cmp.Diff(gcca, want.Category[j].ChildCategory[k], cmpopts.IgnoreFields(gcca, "CreatedAt", "UpdatedAt")); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
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

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	os.Exit(m.Run())
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
		(11001, 'test_image_title_11001', 'test_image_original_src_11001.com', 'test_image_simple_src_11001.com'),
		(999990, 'test_image_title_999990', 'test_image_original_src_999990.com', 'test_image_simple_src_999990.com'),
		(999991, 'test_image_title_999991', 'test_image_original_src_999991.com', 'test_image_simple_src_999991.com');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(11001, 'test_character_name_11001', 'test_character_src_11001.com'),
		(11002, 'test_character_name_11002', 'test_character_src_11002.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_characters_relations (id, image_id, character_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001),
		(11003, 11001, 11002);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(11001, 'test_parent_category_name_11001', 'test_parent_category_src_11001.com'),
		(11002, 'test_parent_category_name_11002', 'test_parent_category_src_11002.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_parent_categories_relations (id, image_id, parent_category_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001),
		(11003, 11001, 11002);
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(11001, 'test_child_category_name_11001', 11001),
		(11002, 'test_child_category_name_11002', 11002);
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
		"TRUNCATE TABLE image_parent_categories_relations RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE image_characters_relations RESTART IDENTITY CASCADE;",
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
