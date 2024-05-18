package api

import (
	"shin-monta-no-mori/server/api/admin"
	"shin-monta-no-mori/server/api/middleware"
	"shin-monta-no-mori/server/api/user"
	"shin-monta-no-mori/server/internal/app"
)

func SetUserRouters(s *app.Server) {
	v1 := s.Router.Group("/api/v1")
	{
		illustrations := v1.Group("/illustrations")
		{
			illustrations.GET("/:id", app.HandlerFuncWrapper(s, user.GetIllustration))
			illustrations.GET("/list", app.HandlerFuncWrapper(s, user.ListIllustrations))
			illustrations.GET("/search", app.HandlerFuncWrapper(s, user.SearchIllustrations))
		}
	}
}

func SetAdminRouters(s *app.Server) {
	v1 := s.Router.Group("/api/v1")
	adminGroup := v1.Group("/admin")
	// ログイン認証
	adminGroup.Use(middleware.AuthMiddleware(s.TokenMaker))
	{
		illustrations := adminGroup.Group("/illustrations")
		{
			illustrations.GET("/:id", app.HandlerFuncWrapper(s, admin.GetIllustration))
			illustrations.GET("/list", app.HandlerFuncWrapper(s, admin.ListIllustrations))
			illustrations.GET("/search", app.HandlerFuncWrapper(s, admin.SearchIllustrations))
			illustrations.POST("/create", app.HandlerFuncWrapper(s, admin.CreateIllustration))
			illustrations.DELETE("/:id", app.HandlerFuncWrapper(s, admin.DeleteIllustration))
			illustrations.PUT("/:id", app.HandlerFuncWrapper(s, admin.EditIllustration))
		}
		characters := adminGroup.Group("/characters")
		{
			characters.GET("/list", app.HandlerFuncWrapper(s, admin.ListCharacters))
			characters.GET("/search", app.HandlerFuncWrapper(s, admin.SearchCharacters))
			characters.GET("/:id", app.HandlerFuncWrapper(s, admin.GetCharacter))
			characters.POST("/create", app.HandlerFuncWrapper(s, admin.CreateCharacter))
			characters.DELETE("/:id", app.HandlerFuncWrapper(s, admin.DeleteCharacter))
			characters.PUT("/:id", app.HandlerFuncWrapper(s, admin.EditCharacter))
		}
		categories := adminGroup.Group("/categories")
		{
			categories.GET("/list", app.HandlerFuncWrapper(s, admin.ListCategories))
			categories.GET("/search", app.HandlerFuncWrapper(s, admin.SearchCategories))
			categories.GET("/:id", app.HandlerFuncWrapper(s, admin.GetCategory))
			parent_categories := categories.Group("/parent")
			{
				parent_categories.POST("/create", app.HandlerFuncWrapper(s, admin.CreateParentCategory))
				parent_categories.PUT("/:id", app.HandlerFuncWrapper(s, admin.EditParentCategory))
				parent_categories.DELETE("/:id", app.HandlerFuncWrapper(s, admin.DeleteParentCategory))
			}
			child_categories := categories.Group("/child")
			{
				child_categories.POST("/create", app.HandlerFuncWrapper(s, admin.CreateChildCategory))
				child_categories.PUT("/:id", app.HandlerFuncWrapper(s, admin.EditChildCategory))
				child_categories.DELETE("/:id", app.HandlerFuncWrapper(s, admin.DeleteChildCategory))
			}
		}
	}
}
