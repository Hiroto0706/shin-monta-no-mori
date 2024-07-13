package admin_test

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
	"shin-monta-no-mori/server/api"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/pkg/lib/password"
	"shin-monta-no-mori/server/pkg/token"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lib/pq"
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

	// 認証用トークンの生成
	accessToken := setAuthUser(t, c)

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
							Character: db.Character{
								ID:   11001,
								Name: "test_character_name_11001",
								Src:  "test_character_src_11001.com",
							},
						},
					},
					Categories: []*model.Category{
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
					Characters: []*model.Character{
						{
							Character: db.Character{
								ID:   11001,
								Name: "test_character_name_11001",
								Src:  "test_character_src_11001.com",
							},
						},
					},
					Categories: []*model.Category{
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
			c.Server.Config.ImageFetchLimit = tt.arg.imageFetchLimit
			req, _ := http.NewRequest("GET", "/api/v1/admin/illustrations/list?p="+tt.arg.page, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			w := httptest.NewRecorder()
			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustration []model.Illustration `json:"illustrations"`
					TotalPages   int64                `json:"total_pages"`
					TotalCount   int64                `json:"total_count"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got.Illustration {
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

	// 認証用トークンの生成
	accessToken := setAuthUser(t, c)

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
				Characters: []*model.Character{
					{
						Character: db.Character{
							ID:   11002,
							Name: "test_character_name_11002",
							Src:  "test_character_src_11002.com",
						},
					},
				},
				Categories: []*model.Category{
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
			req.Header.Set("Authorization", "Bearer "+accessToken)

			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustration model.Illustration `json:"illustration"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareIllustrationsObjects(t, got.Illustration, tt.want, ignoreFields)
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

	// 認証用トークンの生成
	accessToken := setAuthUser(t, c)

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
					Characters: []*model.Character{
						{
							Character: db.Character{
								ID:   12001,
								Name: "test_character_name_12001",
								Src:  "test_character_src_12001.com",
							},
						},
					},
					Categories: []*model.Category{
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
			req, _ := http.NewRequest("GET", "/api/v1/admin/illustrations/search?p="+tt.arg.p+"&q="+tt.arg.q, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Illustration []model.Illustration `json:"illustrations"`
					TotalPages   int64                `json:"total_pages"`
					TotalCount   int64                `json:"total_count"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got.Illustration {
					compareIllustrationsObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

// TODO: credential.jsonがgithub上にないためアップロード処理が実行できない
// func TestCreateIllustration(t *testing.T) {
// 	os.Setenv("CREDENTIAL_FILE_PATH", "../../credential.json")
// 	config, err := util.LoadConfig(AppEnvPath)
// 	if err != nil {
// 		log.Fatal("cannot load config :", err)
// 	}
// 	server := setUp(t, config)
// 	defer tearDown(t, config)

// 	tests := []struct {
// 		name         string
// 		prepare      func() (*bytes.Buffer, string)
// 		want         model.Illustration
// 		wantErr      bool
// 		expectedCode int
// 	}{
// 		{
// 			name: "正常系",
// 			prepare: func() (*bytes.Buffer, string) {
// 				body := &bytes.Buffer{}
// 				writer := multipart.NewWriter(body)
// 				defer writer.Close()

// 				// テキストフィールドを追加
// 				_ = writer.WriteField("title", "test_illustration_1")
// 				_ = writer.WriteField("filename", "test_illustration_filename_1")
// 				_ = writer.WriteField("characters[]", "13001")
// 				_ = writer.WriteField("parent_categories[]", "13001")
// 				_ = writer.WriteField("child_categories[]", "13001")

// 				// ファイルを追加
// 				filePath := "../tmp/test-image.png"

// 				// TODO: tmpがgithub上にないので、空のコンテンツをGCSに保存することになってしまっている

// 				// file, err := os.Open(filePath)
// 				// require.NoError(t, err)
// 				// defer file.Close()
// 				// // fileパートを作成
// 				// part, err := writer.CreateFormFile("original_image_file", filepath.Base(filePath))
// 				// require.NoError(t, err)

// 				// // ファイルの内容を読み込み、書き込む
// 				// _, err = io.Copy(part, file)
// 				// require.NoError(t, err)

// 				file1, _ := writer.CreateFormFile("original_image_file", filepath.Base(filePath))
// 				_, _ = file1.Write([]byte("file content"))
// 				file2, _ := writer.CreateFormFile("simple_image_file", filepath.Base(filePath))
// 				_, _ = file2.Write([]byte("file content"))

// 				return body, writer.FormDataContentType()
// 			},
// 			want: model.Illustration{
// 				Image: db.Image{
// 					Title:            "test_illustration_1",
// 					OriginalSrc:      "https://storage.googleapis.com/shin-monta-no-mori/image/dev/test_illustration_filename_1.png",
// 					OriginalFilename: "test_illustration_filename_1",
// 					SimpleSrc: sql.NullString{
// 						String: "https://storage.googleapis.com/shin-monta-no-mori/image/dev/test_illustration_filename_1_s.png",
// 						Valid:  true,
// 					},
// 					SimpleFilename: sql.NullString{
// 						String: "test_illustration_filename_1_s",
// 						Valid:  true,
// 					},
// 				},
// 				Character: []db.Character{
// 					{
// 						ID:   13001,
// 						Name: "test_character_name_13001",
// 						Src:  "test_character_src_13001.com",
// 					},
// 				},
// 				Categories: []*model.Category{
// 					{
// 						ParentCategory: db.ParentCategory{
// 							ID:   13001,
// 							Name: "test_parent_category_name_13001",
// 							Src:  "test_parent_category_src_13001.com",
// 						},
// 						ChildCategory: []db.ChildCategory{
// 							{
// 								ID:       13001,
// 								Name:     "test_child_category_name_13001",
// 								ParentID: 13001,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			wantErr:      false,
// 			expectedCode: http.StatusOK,
// 		},
// 		{
// 			name: "異常系（requestの型が不正な場合）",
// 			prepare: func() (*bytes.Buffer, string) {
// 				body := &bytes.Buffer{}
// 				writer := multipart.NewWriter(body)
// 				defer writer.Close()

// 				// テキストフィールドを追加
// 				_ = writer.WriteField("aaa", "aaa")

// 				return body, writer.FormDataContentType()
// 			},
// 			want: model.Illustration{
// 				Image:     db.Image{},
// 				Character: []db.Character{},
// 				Category:  []*model.Category{},
// 			},
// 			wantErr:      true,
// 			expectedCode: http.StatusBadRequest,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			body, contentType := tt.prepare()
// 			req := httptest.NewRequest("POST", "/api/v1/admin/illustrations/create", body)
// 			req.Header.Set("Content-Type", contentType)

// 			w := httptest.NewRecorder()
// 			s.Router.ServeHTTP(w, req)

// 			require.Equal(t, tt.expectedCode, w.Code)

// 			if tt.wantErr {
// 				require.NotEmpty(t, w.Body.String())
// 			} else {
// 				var got struct {
// 					Illustrations model.Illustration `json:"illustrations"`
// 				}
// 				err := json.Unmarshal(w.Body.Bytes(), &got)
// 				require.NoError(t, err)
// 				ignoreFields := map[string][]string{
// 					"Image": {"CreatedAt", "UpdatedAt", "ID"},
// 					"Other": {"CreatedAt", "UpdatedAt"},
// 				}
// 				compareIllustrationsObjects(t, got.Illustrations, tt.want, ignoreFields)
// 				// GCSからテストオブジェクトを削除する
// 				deleteGCSObject(t, &gin.Context{}, &config, got.Illustrations.Image.OriginalSrc)
// 			}
// 		})
// 	}
// }

func TestEditIllustration(t *testing.T) {
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	i := illustrationTest{}
	c := i.setUp(t, config)
	defer i.tearDown(t, config)

	// 認証用トークンの生成
	accessToken := setAuthUser(t, c)

	tests := []struct {
		name         string
		arg          string
		prepare      func() (*bytes.Buffer, string)
		want         model.Illustration
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系",
			arg:  "14001",
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				// テキストフィールドを追加
				_ = writer.WriteField("title", "test_image_title_14001_edited")
				_ = writer.WriteField("filename", "test_image_original_filename_14001")
				_ = writer.WriteField("characters[]", "14001")
				_ = writer.WriteField("characters[]", "14002")
				_ = writer.WriteField("parent_categories[]", "14001")
				_ = writer.WriteField("parent_categories[]", "14002")
				_ = writer.WriteField("child_categories[]", "14001")
				_ = writer.WriteField("child_categories[]", "14002")

				return body, writer.FormDataContentType()
			},
			want: model.Illustration{
				Image: db.Image{
					Title:            "test_image_title_14001_edited",
					OriginalFilename: "test_image_original_filename_14001",
				},
				Characters: []*model.Character{
					{
						Character: db.Character{
							ID:   14001,
							Name: "test_character_name_14001",
							Src:  "test_character_src_14001.com",
						},
					},
					{
						Character: db.Character{
							ID:   14002,
							Name: "test_character_name_14002",
							Src:  "test_character_src_14002.com",
						},
					},
				},
				Categories: []*model.Category{
					{
						ParentCategory: db.ParentCategory{
							ID:   14001,
							Name: "test_parent_category_name_14001",
							Src:  "test_parent_category_src_14001.com",
						},
						ChildCategory: []db.ChildCategory{
							{
								ID:       14001,
								Name:     "test_child_category_name_14001",
								ParentID: 14001,
							},
						},
					},
					{
						ParentCategory: db.ParentCategory{
							ID:   14002,
							Name: "test_parent_category_name_14002",
							Src:  "test_parent_category_src_14002.com",
						},
						ChildCategory: []db.ChildCategory{
							{
								ID:       14002,
								Name:     "test_child_category_name_14002",
								ParentID: 14002,
							},
						},
					},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（idの値が不正な場合）",
			arg:  "aaa",
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				return body, writer.FormDataContentType()
			},
			want:         model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "異常系（存在しないillustrationを編集しようとした場合）",
			arg:  "999999",
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				return body, writer.FormDataContentType()
			},
			want:         model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "異常系（存在しないcharacterのIDを指定している場合）",
			arg:  "14002",
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				// テキストフィールドを追加
				_ = writer.WriteField("title", "test_image_title_14002_edited")
				_ = writer.WriteField("filename", "test_image_original_filename_14002")
				_ = writer.WriteField("characters[]", "999999")

				return body, writer.FormDataContentType()
			},
			want:         model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "異常系（存在しないparent_categoryのIDを指定している場合）",
			arg:  "14003",
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				// テキストフィールドを追加
				_ = writer.WriteField("title", "test_image_title_14003_edited")
				_ = writer.WriteField("filename", "test_image_original_filename_14003")
				_ = writer.WriteField("parent_categories[]", "999999")

				return body, writer.FormDataContentType()
			},
			want:         model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "異常系（存在しないchild_categoryのIDを指定している場合）",
			arg:  "14004",
			prepare: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				defer writer.Close()

				// テキストフィールドを追加
				_ = writer.WriteField("title", "test_image_title_14004_edited")
				_ = writer.WriteField("filename", "test_image_original_filename_14004")
				_ = writer.WriteField("child_categories[]", "999999")

				return body, writer.FormDataContentType()
			},
			want:         model.Illustration{},
			wantErr:      true,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, contentType := tt.prepare()
			req := httptest.NewRequest("PUT", "/api/v1/admin/illustrations/"+tt.arg, body)
			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			w := httptest.NewRecorder()
			c.Server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				var got struct {
					Illustration model.Illustration `json:"illustration"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Image": {"CreatedAt", "UpdatedAt", "ID", "OriginalSrc", "SimpleSrc", "SimpleFilename"},
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				compareIllustrationsObjects(t, got.Illustration, tt.want, ignoreFields)
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
	token, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}
	s := &app.Server{
		Config:     config,
		Store:      store,
		TokenMaker: token,
	}
	router := gin.Default()
	s.Router = router
	api.SetAdminRouters(s)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx := app.NewAppContext(c, s)
	return ctx, nil
}

func newTestUserCreation(ctx *app.AppContext, name, pw, email string) (db.Operator, error) {
	hashedPassword, err := password.HashPassword(pw)
	if err != nil {
		return db.Operator{}, err
	}
	arg := db.CreateOperatorParams{
		Name:           name,
		HashedPassword: hashedPassword,
		Email:          email,
	}

	user, err := ctx.Server.Store.CreateOperator(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return db.Operator{}, err
			}
		}
		return db.Operator{}, err
	}

	return user, nil
}

// 認証用トークンの生成
func setAuthUser(t *testing.T, c *app.AppContext) string {
	user, err := newTestUserCreation(c, "testuser", "testtest", "test@test.com")
	require.NoError(t, err)
	accessToken, _, err := c.Server.TokenMaker.CreateToken(
		user.Name,
		c.Server.Config.AccessTokenDuration,
	)
	require.NoError(t, err)

	refreshToken, refreshPayload, err := c.Server.TokenMaker.CreateToken(
		user.Name,
		c.Server.Config.RefreshTokenDuration,
	)
	require.NoError(t, err)

	_, err = c.Server.Store.CreateSession(c, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Name:         user.Name,
		Email:        sql.NullString{String: user.Email, Valid: true},
		RefreshToken: refreshToken,
		ExpiresAt:    refreshPayload.ExpiredAt,
		UserAgent:    "test",
		ClientIp:     "test",
	})
	require.NoError(t, err)

	return accessToken
}

func (i illustrationTest) setUp(t *testing.T, config util.Config) *app.AppContext {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src, original_filename)
		VALUES
		(11001, 'test_image_title_11001', 'test_image_original_src_11001.com', 'test_image_simple_src_11001.com', 'test_image_original_filename_11001'),
		(999990, 'test_image_title_999990', 'test_image_original_src_999990.com', 'test_image_simple_src_999990.com', 'test_image_original_filename_999990'),
		(999991, 'test_image_title_999991', 'test_image_original_src_999991.com', 'test_image_simple_src_999991.com', 'test_image_original_filename_999991'),
		(12001, 'test_image_title_12001', 'test_image_original_src_12001.com', 'test_image_simple_src_12001.com', 'test_image_original_filename_12001'),
		(14001, 'test_image_title_14001', 'test_image_original_src_14001.com', 'test_image_simple_src_14001.com', 'test_image_original_filename_14001'),
		(14002, 'test_image_title_14002', 'test_image_original_src_14002.com', 'test_image_simple_src_14002.com', 'test_image_original_filename_14002'),
		(14003, 'test_image_title_14003', 'test_image_original_src_14003.com', 'test_image_simple_src_14003.com', 'test_image_original_filename_14003'),
		(14004, 'test_image_title_14004', 'test_image_original_src_14004.com', 'test_image_simple_src_14004.com', 'test_image_original_filename_14004');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(11001, 'test_character_name_11001', 'test_character_src_11001.com'),
		(11002, 'test_character_name_11002', 'test_character_src_11002.com'),
		(12001, 'test_character_name_12001', 'test_character_src_12001.com'),
		(13001, 'test_character_name_13001', 'test_character_src_13001.com'),
		(14001, 'test_character_name_14001', 'test_character_src_14001.com'),
		(14002, 'test_character_name_14002', 'test_character_src_14002.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_characters_relations (id, image_id, character_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001),
		(11003, 11001, 11002),
		(12001, 12001, 12001),
		(14001, 14001, 14001);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(11001, 'test_parent_category_name_11001', 'test_parent_category_src_11001.com'),
		(11002, 'test_parent_category_name_11002', 'test_parent_category_src_11002.com'),
		(12001, 'test_parent_category_name_12001', 'test_parent_category_src_12001.com'),
		(13001, 'test_parent_category_name_13001', 'test_parent_category_src_13001.com'),
		(14001, 'test_parent_category_name_14001', 'test_parent_category_src_14001.com'),
		(14002, 'test_parent_category_name_14002', 'test_parent_category_src_14002.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_parent_categories_relations (id, image_id, parent_category_id)
		VALUES
		(11001, 999990, 11001),
		(11002, 999991, 11001),
		(11003, 11001, 11002),
		(12001, 12001, 12001),
		(14001, 14001, 14001);
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(11001, 'test_child_category_name_11001', 11001),
		(11002, 'test_child_category_name_11002', 11002),
		(12001, 'test_child_category_name_12001', 12001),
		(13001, 'test_child_category_name_13001', 13001),
		(14001, 'test_child_category_name_14001', 14001),
		(14002, 'test_child_category_name_14002', 14002);
		`),
		fmt.Sprintln(`
		INSERT INTO image_child_categories_relations (id, image_id, child_category_id)
		VALUES
		(12001, 12001, 12001),
		(14001, 14001, 14001);
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
		"TRUNCATE TABLE operators RESTART IDENTITY CASCADE;",
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
