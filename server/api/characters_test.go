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
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

type charactersTest struct{}

func TestListCharacters(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := charactersTest{}
	server := c.setUp(t, config)
	defer c.tearDown(t, config)

	type args struct {
		page         string
		fetchLimit   int
		compareLimit int
	}

	tests := []struct {
		name         string
		arg          args
		want         []db.Character
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系（p=0）",
			arg: args{
				page:         "0",
				fetchLimit:   1,
				compareLimit: 1,
			},
			want: []db.Character{
				{
					ID:       29001,
					Name:     "test_character_name_29001",
					Src:      "test_character_src_29001.com",
					Filename: sql.NullString{String: "test_character_filename_29001", Valid: true},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（存在しないページを指定している場合）",
			arg: args{
				page:         "9999",
				fetchLimit:   1,
				compareLimit: 0,
			},
			want: []db.Character{
				{},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（pageの値が負の場合）",
			arg: args{
				page:         "-1",
				fetchLimit:   1,
				compareLimit: 1,
			},
			want: []db.Character{
				{},
			},
			wantErr:      true,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 取得するイメージの数を1にする
			server.Config.CharacterFetchLimit = tt.arg.fetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/characters/list?p="+tt.arg.page, nil)
			server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Characters []db.Character `json:"characters"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				log.Println(got.Characters)
				for i, g := range got.Characters[:tt.arg.compareLimit] {
					compareCharactersObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func TestSearchCharacters(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := charactersTest{}
	server := c.setUp(t, config)
	defer c.tearDown(t, config)

	type args struct {
		page         string
		query        string
		fetchLimit   int
		compareLimit int
	}

	tests := []struct {
		name         string
		arg          args
		want         []db.Character
		wantErr      bool
		expectedCode int
	}{
		{
			name: "正常系（p=0, q=20001）",
			arg: args{
				page:         "0",
				query:        "20001",
				fetchLimit:   1,
				compareLimit: 1,
			},
			want: []db.Character{
				{
					ID:       20001,
					Name:     "test_character_name_20001",
					Src:      "test_character_src_20001.com",
					Filename: sql.NullString{String: "test_character_filename_20001", Valid: true},
				},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（p=9999, q=なし）",
			arg: args{
				page:         "9999",
				query:        "",
				fetchLimit:   1,
				compareLimit: 1,
			},
			want: []db.Character{
				{},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "正常系（p=0, q=存在しないquery）",
			arg: args{
				page:         "0",
				query:        "not exist character name",
				fetchLimit:   1,
				compareLimit: 1,
			},
			want: []db.Character{
				{},
			},
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name: "異常系（pageの値が負になる場合）",
			arg: args{
				page:         "-1",
				query:        "",
				fetchLimit:   1,
				compareLimit: 1,
			},
			want: []db.Character{
				{},
			},
			wantErr:      true,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 取得するイメージの数を1にする
			server.Config.CharacterFetchLimit = tt.arg.fetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/characters/search?p="+tt.arg.page+"&q="+tt.arg.query, nil)
			server.Router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.wantErr {
				require.NotEmpty(t, w.Body.String())
			} else {
				type wantType struct {
					Characters []db.Character `json:"characters"`
				}
				var got wantType
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				ignoreFields := map[string][]string{
					"Other": {"CreatedAt", "UpdatedAt"},
				}
				for i, g := range got.Characters {
					compareCharactersObjects(t, g, tt.want[i], ignoreFields)
				}
			}
		})
	}
}

func compareCharactersObjects(t *testing.T, got db.Character, want db.Character, ignoreFieldsMap map[string][]string) {
	if d := cmp.Diff(got, want, cmpopts.IgnoreFields(got, ignoreFieldsMap["Other"]...)); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func (c charactersTest) setUp(t *testing.T, config util.Config) *api.Server {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src, filename)
		VALUES
		(29001, 'test_character_name_29001', 'test_character_src_29001.com', 'test_character_filename_29001'),
		(20001, 'test_character_name_20001', 'test_character_src_20001.com', 'test_character_filename_20001');
		`),
		// fmt.Sprintln(`
		// INSERT INTO image_characters_relations (id, image_id, character_id)
		// VALUES
		// (11001, 999990, 11001);
		// `),
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

func (c charactersTest) tearDown(t *testing.T, config util.Config) {
	store := createConn(config)

	queries := []string{
		"TRUNCATE TABLE characters RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE image_characters_relations RESTART IDENTITY CASCADE;",
	}
	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}
}
