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
			categories.GET("/list", s.ListIllustrations)
			categories.GET("/search", s.SearchIllustrations)
			categories.GET("/:id", s.GetIllustration)
			categories.POST("/create", s.CreateIllustration)
			categories.DELETE("/:id", s.DeleteIllustration)
			categories.PUT("/:id", s.EditIllustration)
		}
	}
}
