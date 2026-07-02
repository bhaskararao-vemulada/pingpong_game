package server




func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/ws", s.HandleWebSocket)
}