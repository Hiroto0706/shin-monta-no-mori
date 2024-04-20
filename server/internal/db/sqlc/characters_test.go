package db_test

import (
	"context"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCharacter(t *testing.T) {
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.CreateCharacterParams
		want    db.Character
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.CreateCharacterParams{
				Name: "test_character_name_10001",
				Src:  "test_character_src_10001",
			},
			want: db.Character{
				Name: "test_character_name_10001",
				Src:  "test_character_src_10001",
			},
			wantErr: false,
		},
		{
			name: "正常系（Srcが空文字の場合）",
			arg: db.CreateCharacterParams{
				Name: "test_character_name_10001",
				Src:  "",
			},
			want: db.Character{
				Name: "test_character_name_10001",
				Src:  "",
			},
			wantErr: false,
		},
		{
			name: "異常系（nameが空文字の場合）",
			arg: db.CreateCharacterParams{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character, err := testQueries.CreateCharacter(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, character)
				require.Equal(t, tt.arg.Name, character.Name)
				require.Equal(t, tt.arg.Src, character.Src)
				require.NotZero(t, character.ID)
				require.NotZero(t, character.CreatedAt)
			}
		})
	}
}

func TestDeleteCharacter(t *testing.T) {
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
			err := testQueries.DeleteCharacter(context.Background(), tt.id)
			if !tt.wantErr {
				require.NoError(t, err)
				_, err := testQueries.GetCharacter(context.Background(), tt.id)
				require.Error(t, err, "The character should no longer exist.")
			}
		})
	}
}

func TestGetCharacter(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		arg     args
		want    db.Character
		wantErr bool
	}{
		{
			name: "正常系",
			arg: args{
				id: 20001,
			},
			want: db.Character{
				ID:   20001,
				Name: "test_character_name_20001",
				Src:  "test_character_src_20001",
			},
			wantErr: false,
		},
		{
			name: "正常系（Srcが空文字の場合）",
			arg: args{
				id: 20002,
			},
			want: db.Character{
				ID:   20002,
				Name: "test_character_name_20002",
				Src:  "",
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
			character, err := testQueries.GetCharacter(context.Background(), tt.arg.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, character.ID)
				require.Equal(t, tt.want.Name, character.Name)
				require.Equal(t, tt.want.Src, character.Src)
				require.NotZero(t, character.CreatedAt)
			}
		})
	}
}

func TestListCharacter(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.ListCharactersParams
		want    []db.Character
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.ListCharactersParams{
				Limit:  3,
				Offset: 0,
			},
			want: []db.Character{
				{
					ID:   99992,
					Name: "test_character_name_99992",
					Src:  "",
				},
				{
					ID:   99991,
					Name: "test_character_name_99991",
					Src:  "test_character_src_99991",
				},
				{
					ID:   99990,
					Name: "test_character_name_99990",
					Src:  "test_character_src_99990",
				},
			},
			wantErr: false,
		},
		{
			name: "正常系（returnが空の時）",
			arg: db.ListCharactersParams{
				Limit:  3,
				Offset: 1000,
			},
			want:    []db.Character{{}},
			wantErr: false,
		},
		{
			name: "異常系（argsの値が不正な場合）",
			arg: db.ListCharactersParams{
				Limit:  -1,
				Offset: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			characters, err := testQueries.ListCharacters(context.Background(), tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				for i, character := range characters {
					require.NoError(t, err)
					require.Equal(t, tt.want[i].ID, character.ID)
					require.Equal(t, tt.want[i].Name, character.Name)
					require.Equal(t, tt.want[i].Src, character.Src)
					require.NotZero(t, character.CreatedAt)
				}
			}
		})
	}
}

func TestUpdateCharacter(t *testing.T) {
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.UpdateCharacterParams
		want    db.Character
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.UpdateCharacterParams{
				ID:   40001,
				Name: "test_character_name_40001_edited",
				Src:  "test_character_src_40001_edited",
			},
			want: db.Character{
				ID:   40001,
				Name: "test_character_name_40001_edited",
				Src:  "test_character_src_40001_edited",
			},
			wantErr: false,
		},
		{
			name: "異常系（存在しないIDを指定している場合）",
			arg: db.UpdateCharacterParams{
				ID: 99999,
			},
			wantErr: true,
		},
		{
			name: "異常系（titleが空になる場合）",
			arg: db.UpdateCharacterParams{
				ID:   40002,
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character, err := testQueries.UpdateCharacter(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, character)
				require.Equal(t, tt.arg.ID, character.ID)
				require.Equal(t, tt.arg.Name, character.Name)
				require.Equal(t, tt.arg.Src, character.Src)
				require.NotZero(t, character.CreatedAt)
			}
		})
	}
}
