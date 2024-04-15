package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCharacter(t *testing.T) {
	defer TearDown(t, testQueries.db)

	tests := []struct {
		name    string
		arg     CreateCharacterParams
		want    Character
		wantErr bool
	}{
		{
			name: "正常系",
			arg: CreateCharacterParams{
				Name: "test_character_name_10001",
				Src:  "test_character_src_10001",
			},
			want: Character{
				Name: "test_character_name_10001",
				Src:  "test_character_src_10001",
			},
			wantErr: false,
		},
		{
			name: "正常系（Srcが空文字の場合）",
			arg: CreateCharacterParams{
				Name: "test_character_name_10001",
				Src:  "",
			},
			want: Character{
				Name: "test_character_name_10001",
				Src:  "",
			},
			wantErr: false,
		},
		{
			name: "異常系（nameが空文字の場合）",
			arg: CreateCharacterParams{
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
	SetUp(t, testQueries.db)
	defer TearDown(t, testQueries.db)

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
	SetUp(t, testQueries.db)
	defer TearDown(t, testQueries.db)

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		arg     args
		want    Character
		wantErr bool
	}{
		{
			name: "正常系",
			arg: args{
				id: 20001,
			},
			want: Character{
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
			want: Character{
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
