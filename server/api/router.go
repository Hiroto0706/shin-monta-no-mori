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
			illustrations.GET("/random", app.HandlerFuncWrapper(s, user.FetchRandomIllustrations))
			illustrations.GET("/character/:id", app.HandlerFuncWrapper(s, user.ListIllustrationsByCharacterID))
			illustrations.GET("/category/parent/:id", app.HandlerFuncWrapper(s, user.ListIllustrationsByParentCategoryID))
			illustrations.GET("/category/child/:id", app.HandlerFuncWrapper(s, user.ListIllustrationsByChildCategoryID))
		}
		characters := v1.Group("/characters")
		{
			characters.GET("/:id", app.HandlerFuncWrapper(s, user.GetCharacter))
			characters.GET("/list", app.HandlerFuncWrapper(s, user.ListCharacters))
			characters.GET("/list/all", app.HandlerFuncWrapper(s, user.ListAllCharacters))
		}
		categories := v1.Group("/categories")
		{
			categories.GET("/list", app.HandlerFuncWrapper(s, user.ListCategories))
			categories.GET("/list/all", app.HandlerFuncWrapper(s, user.ListAllCategories))
			child_categories := categories.Group("/child")
			{
				child_categories.GET("/:id", app.HandlerFuncWrapper(s, user.GetChildCategory))
				child_categories.GET("/list", app.HandlerFuncWrapper(s, user.ListChildCategories))
			}
		}
	}
}

func SetAdminRouters(s *app.Server) {
	v1 := s.Router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/verify", app.HandlerFuncWrapper(s, admin.VerifyAccessToken))
		auth.POST("/login", app.HandlerFuncWrapper(s, admin.Login))
	}

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
			characters.GET("/list/all", app.HandlerFuncWrapper(s, admin.ListAllCharacters))
			characters.GET("/search", app.HandlerFuncWrapper(s, admin.SearchCharacters))
			characters.GET("/:id", app.HandlerFuncWrapper(s, admin.GetCharacter))
			characters.POST("/create", app.HandlerFuncWrapper(s, admin.CreateCharacter))
			characters.DELETE("/:id", app.HandlerFuncWrapper(s, admin.DeleteCharacter))
			characters.PUT("/:id", app.HandlerFuncWrapper(s, admin.EditCharacter))
		}
		categories := adminGroup.Group("/categories")
		{
			categories.GET("/list", app.HandlerFuncWrapper(s, admin.ListCategories))
			categories.GET("/list/all", app.HandlerFuncWrapper(s, admin.ListAllCategories))
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
				child_categories.GET("/:id", app.HandlerFuncWrapper(s, admin.GetChildCategory))
				child_categories.POST("/create", app.HandlerFuncWrapper(s, admin.CreateChildCategory))
				child_categories.PUT("/:id", app.HandlerFuncWrapper(s, admin.EditChildCategory))
				child_categories.DELETE("/:id", app.HandlerFuncWrapper(s, admin.DeleteChildCategory))
			}
		}
	}
}
