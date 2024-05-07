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
	"os"
	"path/filepath"
	"shin-monta-no-mori/server/api"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/internal/domains/service"
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
						OriginalFilename: "test_image_original_filename_999991",
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
						OriginalFilename: "test_image_original_filename_999990",
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
			req, _ := http.NewRequest("GET", "/api/v1/admin/illustrations/list?p="+tt.arg.page, nil)
			server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got []model.Illustration
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got {
					compareIllustrationsObjects(t, g, tt.want[i], ignoreFields)
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
					OriginalFilename: "test_image_original_filename_11001",
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
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareIllustrationsObjects(t, got, tt.want, ignoreFields)
			}
		})
	}
}

func TestSearchIllustrations(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	server := setUp(t, config)
	defer tearDown(t, config)

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
				q:               "test_image_title_12001",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image: db.Image{
						ID:          12001,
						Title:       "test_image_title_12001",
						OriginalSrc: "test_image_original_src_12001.com",
						SimpleSrc: sql.NullString{
							String: "test_image_simple_src_12001.com",
							Valid:  true,
						},
						OriginalFilename: "test_image_original_filename_12001",
					},
					Character: []db.Character{
						{
							ID:   12001,
							Name: "test_character_name_12001",
							Src:  "test_character_src_12001.com",
						},
					},
					Category: []*model.Category{
						{
							ParentCategory: db.ParentCategory{
								ID:   12001,
								Name: "test_parent_category_name_12001",
								Src:  "test_parent_category_src_12001.com",
							},
							ChildCategory: []db.ChildCategory{
								{
									ID:       12001,
									Name:     "test_child_category_name_12001",
									ParentID: 12001,
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
			name: "正常系（存在しないイラストを検索した時）",
			arg: args{
				p:               "0",
				q:               "not exist illustration",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image:     db.Image{},
					Character: []db.Character{},
					Category: []*model.Category{
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
			name: "異常系（クエリの値が不正な時）",
			arg: args{
				p:               "aaa",
				q:               "not exist illustration",
				imageFetchLimit: 1,
			},
			want: []model.Illustration{
				{
					Image:     db.Image{},
					Character: []db.Character{},
					Category: []*model.Category{
						{
							ParentCategory: db.ParentCategory{},
							ChildCategory:  []db.ChildCategory{},
						},
					},
				},
			},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 取得するイメージの数を1にする
			server.Config.ImageFetchLimit = tt.arg.imageFetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/illustrations/search?p="+tt.arg.p+"&q="+tt.arg.q, nil)
			server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got []model.Illustration
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got {
					compareIllustrationsObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func TestCreateIllustration(t *testing.T) {
	os.Setenv("CREDENTIAL_FILE_PATH", "../../credential.json")
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	server := setUp(t, config)
	defer tearDown(t, config)

	tests := []struct {
		name         string
		prepare      func() (*bytes.Buffer, string)
		want         model.Illustration
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				// テキストフィールドを追加
				_ = writer.WriteField("title", "test_illustration_1")
				_ = writer.WriteField("filename", "test_illustration_filename_1")
				_ = writer.WriteField("characters[]", "13001")
				_ = writer.WriteField("parent_categories[]", "13001")
				_ = writer.WriteField("child_categories[]", "13001")

				// ファイルを追加
				filePath := "../tmp/test-image.png"

				// TODO: tmpがgithub上にないので、空のコンテンツをGCSに保存することになってしまっている

				// file, err := os.Open(filePath)
				// require.NoError(t, err)
				// defer file.Close()
				// // fileパートを作成
				// part, err := writer.CreateFormFile("original_image_file", filepath.Base(filePath))
				// require.NoError(t, err)

				// // ファイルの内容を読み込み、書き込む
				// _, err = io.Copy(part, file)
				// require.NoError(t, err)

				file, _ := writer.CreateFormFile("original_image_file", filepath.Base(filePath))
				_, _ = file.Write([]byte("file content"))

				return body, writer.FormDataContentType()
			},
			want: model.Illustration{
				Image: db.Image{
					Title:            "test_illustration_1",
					OriginalSrc:      "https://storage.googleapis.com/shin-monta-no-mori/image/dev/test_illustration_filename_1.png",
					OriginalFilename: "test_illustration_filename_1",
				},
				Character: []db.Character{
					{
						ID:   13001,
						Name: "test_character_name_13001",
						Src:  "test_character_src_13001.com",
					},
				},
				Category: []*model.Category{
					{
						ParentCategory: db.ParentCategory{
							ID:   13001,
							Name: "test_parent_category_name_13001",
							Src:  "test_parent_category_src_13001.com",
						},
						ChildCategory: []db.ChildCategory{
							{
								ID:       13001,
								Name:     "test_child_category_name_13001",
								ParentID: 13001,
							},
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
			body, contentType := tt.prepare()
			req := httptest.NewRequest("POST", "/api/v1/admin/illustrations/create", body)
			req.Header.Set("Content-Type", contentType)

			w := httptest.NewRecorder()
			server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got struct {
					Illustrations model.Illustration `json:"illustrations"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt", "ID"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareIllustrationsObjects(t, got.Illustrations, tt.want, ignoreFields)
				// GCSからテストオブジェクトを削除する
				deleteGCSObject(t, &gin.Context{}, &config, got.Illustrations.Image.OriginalSrc)
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
	for i, gch := range got.Character {
		if d := cmp.Diff(gch, want.Character[i], cmpopts.IgnoreFields(gch, ignoreFieldsMap["Other"]...)); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}
	}

	// カテゴリー比較
	for j, gca := range got.Category {
		// 親カテゴリー比較
		if d := cmp.Diff(gca.ParentCategory, want.Category[j].ParentCategory, cmpopts.IgnoreFields(gca.ParentCategory, ignoreFieldsMap["Other"]...)); len(d) != 0 {
			t.Errorf("differs: (-got +want)\n%s", d)
		}

		// 子カテゴリー比較
		for k, gcca := range gca.ChildCategory {
			if d := cmp.Diff(gcca, want.Category[j].ChildCategory[k], cmpopts.IgnoreFields(gcca, ignoreFieldsMap["Other"]...)); len(d) != 0 {
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
	server.Router = router
	api.SetAdminRouters(server)
	return server, nil
}

func setUp(t *testing.T, config util.Config) *api.Server {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src, original_filename)
		VALUES
		(11001, 'test_image_title_11001', 'test_image_original_src_11001.com', 'test_image_simple_src_11001.com', 'test_image_original_filename_11001'),
		(999990, 'test_image_title_999990', 'test_image_original_src_999990.com', 'test_image_simple_src_999990.com', 'test_image_original_filename_999990'),
		(999991, 'test_image_title_999991', 'test_image_original_src_999991.com', 'test_image_simple_src_999991.com', 'test_image_original_filename_999991'),
		(12001, 'test_image_title_12001', 'test_image_original_src_12001.com', 'test_image_simple_src_12001.com', 'test_image_original_filename_12001');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(11001, 'test_character_name_11001', 'test_character_src_11001.com'),
		(11002, 'test_character_name_11002', 'test_character_src_11002.com'),
		(12001, 'test_character_name_12001', 'test_character_src_12001.com'),
		(13001, 'test_character_name_13001', 'test_character_src_13001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_characters_relations (id, image_id, character_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001),
		(11003, 11001, 11002),
		(11004, 12001, 12001);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(11001, 'test_parent_category_name_11001', 'test_parent_category_src_11001.com'),
		(11002, 'test_parent_category_name_11002', 'test_parent_category_src_11002.com'),
		(12001, 'test_parent_category_name_12001', 'test_parent_category_src_12001.com'),
		(13001, 'test_parent_category_name_13001', 'test_parent_category_src_13001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_parent_categories_relations (id, image_id, parent_category_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001),
		(11003, 11001, 11002),
		(11004, 12001, 12001);
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(11001, 'test_child_category_name_11001', 11001),
		(11002, 'test_child_category_name_11002', 11002),
		(12001, 'test_child_category_name_12001', 12001),
		(13001, 'test_child_category_name_13001', 13001);
		`),
		fmt.Sprintln(`
		INSERT INTO image_child_categories_relations (id, image_id, child_category_id)
		VALUES
		(12001, 12001, 12001);
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

func deleteGCSObject(t require.TestingT, c *gin.Context, config *util.Config, src string) {
	storageService := &service.GCSStorageService{
		Config: *config,
	}
	err := storageService.DeleteFile(c, src)
	require.NoError(t, err)
}
