package user_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"shin-monta-no-mori/api"
	"shin-monta-no-mori/internal/app"
	db "shin-monta-no-mori/internal/db/sqlc"
	model "shin-monta-no-mori/internal/domains/models"
	"shin-monta-no-mori/pkg/util"
	"testing"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

type illustrationTest struct{}

const AppEnvPath = "../../"

func TestListIllustrations(t *testing.T) {
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	i := illustrationTest{}
	c := i.setUp(t, config)
	defer i.tearDown(t, config)

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
						OriginalFilename: "test_image_original_filename_999991",
					},
					Characters: []*model.Character{
						{
							Character: db.Character{},
						},
					},
					Categories: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
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
						OriginalFilename: "test_image_original_filename_999990",
					},
					Characters: []*model.Character{
						{
							Character: db.Character{},
						},
					},
					Categories: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
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
			c.Server.Config.ImageFetchLimit = tt.arg.imageFetchLimit
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/illustrations/list?p="+tt.arg.page, nil)

			w := httptest.NewRecorder()
			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustrations []model.Illustration `json:"illustrations"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got.Illustrations {
					compareIllustrationsObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func TestGetIllustration(t *testing.T) {
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	i := illustrationTest{}
	c := i.setUp(t, config)
	defer i.tearDown(t, config)

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
				id: "21001",
			},
			want: model.Illustration{
				Image: db.Image{
					ID:          21001,
					Title:       "test_image_title_21001",
					OriginalSrc: "test_image_original_src_21001.com",
					SimpleSrc: sql.NullString{
						String: "test_image_simple_src_21001.com",
						Valid:  true,
					},
					OriginalFilename: "test_image_original_filename_21001",
				},
				Characters: []*model.Character{
					{
						Character: db.Character{
							ID:            21002,
							Name:          "test_character_name_21002",
							Src:           "test_character_src_21002.com",
							PriorityLevel: 2,
						},
					},
				},
				Categories: []*model.Category{
					{
						ParentCategory: db.ParentCategory{
							ID:            21002,
							Name:          "test_parent_category_name_21002",
							Src:           "test_parent_category_src_21002.com",
							PriorityLevel: 2,
						},
						ChildCategory: []db.ChildCategory{
							{
								ID:            21002,
								Name:          "test_child_category_name_21002",
								ParentID:      21002,
								PriorityLevel: 2,
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
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/illustrations/"+tt.arg.id, nil)

			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustration *model.Illustration `json:"illustration"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareIllustrationsObjects(t, *got.Illustration, tt.want, ignoreFields)
			}
		})
	}
}

func TestSearchIllustrations(t *testing.T) {
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	i := illustrationTest{}
	c := i.setUp(t, config)
	defer i.tearDown(t, config)

	type args struct {
		p               string
		q               string
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
			name: "正常系",
			arg: args{
				p:               "0",
				q:               "test_image_title_22001",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image: db.Image{
						ID:          22001,
						Title:       "test_image_title_22001",
						OriginalSrc: "test_image_original_src_22001.com",
						SimpleSrc: sql.NullString{
							String: "test_image_simple_src_22001.com",
							Valid:  true,
						},
						OriginalFilename: "test_image_original_filename_22001",
					},
					Characters: []*model.Character{
						{
							Character: db.Character{},
						},
					},
					Categories: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（存在しないイラストを検索した時）",
			arg: args{
				p:               "0",
				q:               "not exist illustration",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image:      db.Image{},
					Characters: []*model.Character{},
					Categories: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（pageの値が不正な時）",
			arg: args{
				p:               "-1",
				q:               "test",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image:      db.Image{},
					Characters: []*model.Character{},
					Categories: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
						},
					},
				},
			},
			wantErr:      true,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 取得するイメージの数を1にする
			c.Server.Config.ImageFetchLimit = tt.arg.imageFetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/illustrations/search?p="+tt.arg.p+"&q="+tt.arg.q, nil)

			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustrations []model.Illustration `json:"illustrations"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got.Illustrations {
					compareIllustrationsObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func TestListIllustrationsByCharacterID(t *testing.T) {
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	i := illustrationTest{}
	c := i.setUp(t, config)
	defer i.tearDown(t, config)

	type args struct {
		p               string
		id              string
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
			name: "正常系",
			arg: args{
				p:               "0",
				id:              "24001",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image: db.Image{
						ID:          24001,
						Title:       "test_image_title_24001",
						OriginalSrc: "test_image_original_src_24001.com",
						SimpleSrc: sql.NullString{
							String: "test_image_simple_src_24001.com",
							Valid:  true,
						},
						OriginalFilename: "test_image_original_filename_24001",
					},
					Characters: []*model.Character{
						{
							Character: db.Character{},
						},
					},
					Categories: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（存在しないcharacter_idを指定している場合）",
			arg: args{
				p:               "0",
				id:              "999999",
				imageFetchLimit: 1,
			},
			want:         []model.Illustration{},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idの値が不正な場合）",
			arg: args{
				p:               "0",
				id:              "aaa",
				imageFetchLimit: 0,
			},
			want:         []model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（pagerの値が不正な場合）",
			arg: args{
				p:               "aaa",
				id:              "999999",
				imageFetchLimit: 0,
			},
			want:         []model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 取得するイメージの数を1にする
			c.Server.Config.ImageFetchLimit = tt.arg.imageFetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/illustrations/character/"+tt.arg.id+"?p="+tt.arg.p, nil)

			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustrations []model.Illustration `json:"illustrations"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got.Illustrations {
					compareIllustrationsObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}
func TestListIllustrationsByChildCategoryID(t *testing.T) {
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	i := illustrationTest{}
	c := i.setUp(t, config)
	defer i.tearDown(t, config)

	type args struct {
		p               string
		id              string
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
			name: "正常系",
			arg: args{
				p:               "0",
				id:              "25001",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image: db.Image{
						ID:          25001,
						Title:       "test_image_title_25001",
						OriginalSrc: "test_image_original_src_25001.com",
						SimpleSrc: sql.NullString{
							String: "test_image_simple_src_25001.com",
							Valid:  true,
						},
						OriginalFilename: "test_image_original_filename_25001",
					},
					Characters: []*model.Character{
						{
							Character: db.Character{},
						},
					},
					Categories: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（存在しないchild_idを指定している場合）",
			arg: args{
				p:               "0",
				id:              "999999",
				imageFetchLimit: 1,
			},
			want:         []model.Illustration{},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idの値が不正な場合）",
			arg: args{
				p:               "0",
				id:              "aaa",
				imageFetchLimit: 0,
			},
			want:         []model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（pagerの値が不正な場合）",
			arg: args{
				p:               "aaa",
				id:              "999999",
				imageFetchLimit: 0,
			},
			want:         []model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 取得するイメージの数を1にする
			c.Server.Config.ImageFetchLimit = tt.arg.imageFetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/illustrations/category/child/"+tt.arg.id+"?p="+tt.arg.p, nil)

			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustrations []model.Illustration `json:"illustrations"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got.Illustrations {
					compareIllustrationsObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func compareIllustrationsObjects(t *testing.T, got model.Illustration, want model.Illustration, ignoreFieldsMap map[string][]string) {
	// イメージ比較
	if d := cmp.Diff(got.Image, want.Image, cmpopts.IgnoreFields(got.Image, ignoreFieldsMap["Image"]...)); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}

	// キャラクター比較
	for i, gch := range got.Characters {
		if d := cmp.Diff(gch.Character, want.Characters[i].Character, cmpopts.IgnoreFields(gch.Character, ignoreFieldsMap["Other"]...)); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}
	}

	// カテゴリー比較
	for j, gca := range got.Categories {
		// 親カテゴリー比較
		if d := cmp.Diff(gca.ParentCategory, want.Categories[j].ParentCategory, cmpopts.IgnoreFields(gca.ParentCategory, ignoreFieldsMap["Other"]...)); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}

		// 子カテゴリー比較
		for k, gcca := range gca.ChildCategory {
			if d := cmp.Diff(gcca, want.Categories[j].ChildCategory[k], cmpopts.IgnoreFields(gcca, ignoreFieldsMap["Other"]...)); len(d) != 0 {
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

func newTestServer(store *db.Store, config util.Config) (*app.AppContext, error) {
	s := &app.Server{
		Config: config,
		Store:  store,
	}
	router := gin.Default()
	s.Router = router
	api.SetUserRouters(s)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx := app.NewAppContext(c, s)
	return ctx, nil
}

func (i illustrationTest) setUp(t *testing.T, config util.Config) *app.AppContext {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src, original_filename)
		VALUES
		(21001, 'test_image_title_21001', 'test_image_original_src_21001.com', 'test_image_simple_src_21001.com', 'test_image_original_filename_21001'),
		(999990, 'test_image_title_999990', 'test_image_original_src_999990.com', 'test_image_simple_src_999990.com', 'test_image_original_filename_999990'),
		(999991, 'test_image_title_999991', 'test_image_original_src_999991.com', 'test_image_simple_src_999991.com', 'test_image_original_filename_999991'),
		(22001, 'test_image_title_22001', 'test_image_original_src_22001.com', 'test_image_simple_src_22001.com', 'test_image_original_filename_22001'),
		(23001, 'test_image_title_23001', 'test_image_original_src_23001.com', 'test_image_simple_src_23001.com', 'test_image_original_filename_23001'),
		(24001, 'test_image_title_24001', 'test_image_original_src_24001.com', 'test_image_simple_src_24001.com', 'test_image_original_filename_24001'),
		(25001, 'test_image_title_25001', 'test_image_original_src_25001.com', 'test_image_simple_src_25001.com', 'test_image_original_filename_25001');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(21001, 'test_character_name_21001', 'test_character_src_21001.com'),
		(21002, 'test_character_name_21002', 'test_character_src_21002.com'),
		(22001, 'test_character_name_22001', 'test_character_src_22001.com'),
		(23001, 'test_character_name_23001', 'test_character_src_23001.com'),
		(24001, 'test_character_name_24001', 'test_character_src_24001.com'),
		(25001, 'test_character_name_25001', 'test_character_src_25001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_characters_relations (id, image_id, character_id)
		VALUES
		(21001, 999990, 21001),
		(21002, 999991, 21001),
		(21003, 21001, 21002),
		(22001, 22001, 22001),
		(23001, 23001, 23001),
		(24001, 24001, 24001),
		(25001, 25001, 25001);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(21001, 'test_parent_category_name_21001', 'test_parent_category_src_21001.com'),
		(21002, 'test_parent_category_name_21002', 'test_parent_category_src_21002.com'),
		(22001, 'test_parent_category_name_22001', 'test_parent_category_src_22001.com'),
		(23001, 'test_parent_category_name_23001', 'test_parent_category_src_23001.com'),
		(24001, 'test_parent_category_name_24001', 'test_parent_category_src_24001.com'),
		(25001, 'test_parent_category_name_25001', 'test_parent_category_src_25001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_parent_categories_relations (id, image_id, parent_category_id)
		VALUES
		(21001, 999990, 21001),
		(21002, 999991, 21001),
		(21003, 21001, 21002),
		(22001, 22001, 22001),
		(23001, 23001, 23001),
		(24001, 24001, 24001),
		(25001, 25001, 25001);
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(21001, 'test_child_category_name_21001', 21001),
		(21002, 'test_child_category_name_21002', 21002),
		(22001, 'test_child_category_name_22001', 22001),
		(23001, 'test_child_category_name_23001', 23001),
		(24001, 'test_child_category_name_24001', 24001),
		(25001, 'test_child_category_name_25001', 25001);
		`),
		fmt.Sprintln(`
		INSERT INTO image_child_categories_relations (id, image_id, child_category_id)
		VALUES
		(22001, 22001, 22001),
		(23001, 23001, 23001),
		(24001, 24001, 24001),
		(25001, 25001, 25001);
		`),
	}

	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to exec query: %v", err)
		}
	}

	s, err := newTestServer(store, config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	return s
}

func (i illustrationTest) tearDown(t *testing.T, config util.Config) {
	store := createConn(config)

	queries := []string{
		"TRUNCATE TABLE image_child_categories_relations RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE image_parent_categories_relations RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE image_characters_relations RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE child_categories RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE parent_categories RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE characters RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE images RESTART IDENTITY CASCADE;",
	}
	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}
}

// func deleteGCSObject(t require.TestingT, c *gin.Context, config *util.Config, src string) {
// 	storageService := &service.GCSStorageService{
// 		Config: *config,
// 	}
// 	err := storageService.DeleteFile(c, src)
// 	require.NoError(t, err)
// }
