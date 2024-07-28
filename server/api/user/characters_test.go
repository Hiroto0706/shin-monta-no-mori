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
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

type charactersTest struct{}

func TestListCharacters(t *testing.T) {
	config, err := util.LoadConfig(AppEnvPath)
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	c := charactersTest{}
	ctx := c.setUp(t, config)
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
					ID:       29999,
					Name:     "test_character_name_29999",
					Src:      "test_character_src_29999.com",
					Filename: sql.NullString{String: "test_character_filename_29999", Valid: true},
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
			ctx.Server.Config.CharacterFetchLimit = tt.arg.fetchLimit
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/characters/list?p="+tt.arg.page, nil)

			ctx.Server.Router.ServeHTTP(w, req)

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
				for i, g := range got.Characters[:tt.arg.compareLimit] {
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

func (c charactersTest) setUp(t *testing.T, config util.Config) *app.AppContext {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src, filename)
		VALUES
		(29999, 'test_character_name_29999', 'test_character_src_29999.com', 'test_character_filename_29999');
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

func (c charactersTest) tearDown(t *testing.T, config util.Config) {
	store := createConn(config)

	queries := []string{
		"TRUNCATE TABLE characters RESTART IDENTITY CASCADE;",
	}
	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}
}
