package api

func SetUserRouters(s *Server) {
	v1 := s.Router.Group("/api/v1")
	{
		v1.GET("/", s.Greet)
	}
}

func SetAdminRouters(s *Server) {
	v1 := s.Router.Group("/api/v1")
	admin := v1.Group("/admin")
	// ログイン認証
	// admin.Use(authMiddleware(s.tokenMaker))
	{
		illustrations := admin.Group("/illustrations")
		{
			illustrations.GET("/list", s.ListIllustrations)
			illustrations.GET("/search", s.SearchIllustrations)
			illustrations.GET("/:id", s.GetIllustration)
			illustrations.POST("/create", s.CreateIllustration)
			illustrations.DELETE("/:id", s.DeleteIllustration)
			illustrations.PUT("/:id", s.EditIllustration)
		}
		categories := admin.Group("/categories")
		{
			categories.GET("/list", s.ListCategories)
			categories.GET("/search", s.SearchCategories)
			categories.GET("/:id", s.GetCategory)
			parent_categories := categories.Group("/parent")
			{
				parent_categories.POST("/create", s.CreateParentCategory)
				parent_categories.PUT("/:id", s.EditParentCategory)
				parent_categories.DELETE("/:id", s.DeleteParentCategory)
			}
			child_categories := categories.Group("/child")
			{
				child_categories.POST("/create", s.CreateChildCategory)
				child_categories.PUT("/:id", s.EditChildCategory)
				child_categories.DELETE("/:id", s.DeleteChildCategory)
			}
		}
	}
}
