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

	illustrationGetKey   = "illustration_%d"
	illustrationsListKey = "illustrations_list_offset_%d"
)

func GetIllustrationsListKey(offset int) string {
	return fmt.Sprintf(illustrationsListKey, offset)
}

func GetIllustrationKey(id int) string {
	return fmt.Sprintf(illustrationGetKey, id)
}
