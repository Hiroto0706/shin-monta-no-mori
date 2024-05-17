package api

import (
	"shin-monta-no-mori/server/api/admin"
	"shin-monta-no-mori/server/api/middleware"
	"shin-monta-no-mori/server/internal/app"
)

// func SetUserRouters(s *Server) {
// v1 := s.Router.Group("/api/v1")
// {
// v1.GET("/", s.Greet)
// }
// }

func SetAdminRouters(s *app.Server) {
	v1 := s.Router.Group("/api/v1")
	adminGroup := v1.Group("/admin")
	// ログイン認証
	adminGroup.Use(middleware.AuthMiddleware(s.TokenMaker))
	{
		illustrations := adminGroup.Group("/illustrations")
		{
			illustrations.GET("/list", app.HandlerFuncWrapper(s, admin.ListIllustrations))
			// illustrations.GET("/search", s.SearchIllustrations)
			// illustrations.GET("/:id", s.GetIllustration)
			// illustrations.POST("/create", s.CreateIllustration)
			// illustrations.DELETE("/:id", s.DeleteIllustration)
			// illustrations.PUT("/:id", s.EditIllustration)
		}
		// characters := admin.Group("/characters")
		// {
		// 	characters.GET("/list", s.ListCharacters)
		// 	characters.GET("/search", s.SearchCharacters)
		// 	characters.GET("/:id", s.GetCharacter)
		// 	characters.POST("/create", s.CreateCharacter)
		// 	characters.DELETE("/:id", s.DeleteCharacter)
		// 	characters.PUT("/:id", s.EditCharacter)
		// }
		// categories := admin.Group("/categories")
		// {
		// 	categories.GET("/list", s.ListCategories)
		// 	categories.GET("/search", s.SearchCategories)
		// 	categories.GET("/:id", s.GetCategory)
		// 	parent_categories := categories.Group("/parent")
		// 	{
		// 		parent_categories.POST("/create", s.CreateParentCategory)
		// 		parent_categories.PUT("/:id", s.EditParentCategory)
		// 		parent_categories.DELETE("/:id", s.DeleteParentCategory)
		// 	}
		// 	child_categories := categories.Group("/child")
		// 	{
		// 		child_categories.POST("/create", s.CreateChildCategory)
		// 		child_categories.PUT("/:id", s.EditChildCategory)
		// 		child_categories.DELETE("/:id", s.DeleteChildCategory)
		// 	}
		// }
	}
}
