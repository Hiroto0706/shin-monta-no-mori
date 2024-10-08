package cache

import (
	"fmt"
	"time"
)

const (
	// キャッシュの有効期間
	CacheDurationWeek    = 24 * 7 * time.Hour
	CacheDurationDay     = 24 * time.Hour
	CacheDurationHalfDay = 12 * time.Hour
	CacheDurationHour    = 1 * time.Hour

	// イラスト
	illustrationGetKey = "illustration_%d"

	IllustrationsPrefix               = "illustrations_list"
	illustrationsListKey              = IllustrationsPrefix + "_offset_%d"
	illustrationsListByCharacterIDKey = IllustrationsPrefix + "_by_character_%d_%d"
	illustrationsListByCategoryIDKey  = IllustrationsPrefix + "_by_category_%d_%d"

	// カテゴリ
	CategoriesPrefix     = "categories_list"
	categoriesListAllKey = CategoriesPrefix + "_all"

	// キャラクター
	CharactersPrefix     = "characters_list"
	charactersListAllKey = CharactersPrefix + "_all"
)

func GetIllustrationsListKey(offset int) string {
	return fmt.Sprintf(illustrationsListKey, offset)
}

func GetIllustrationKey(id int) string {
	return fmt.Sprintf(illustrationGetKey, id)
}

func GetIllustrationsListByCharacterKey(id, offset int) string {
	return fmt.Sprintf(illustrationsListByCharacterIDKey, id, offset)
}

func GetIllustrationsListByCategoryKey(id, offset int) string {
	return fmt.Sprintf(illustrationsListByCategoryIDKey, id, offset)
}

func GetCategoriesAllKey() string {
	return categoriesListAllKey
}

func GetCharactersAllKey() string {
	return charactersListAllKey
}
