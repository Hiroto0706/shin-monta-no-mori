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
		// admin.GET("/", s.Greet)
		admin.GET("/illustrations/list/", s.ListIllustrations)
		admin.GET("/illustrations/search/", s.SearchIllustration)
		admin.GET("/illustrations/:id", s.GetIllustration)
	}
}
