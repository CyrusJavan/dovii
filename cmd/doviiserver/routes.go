package main

func (s *server) routes() {
	s.router.HandleFunc("/{key}", s.getHandler()).
		Methods("GET")

	s.router.HandleFunc("/{key}/{value}", s.setHandler()).
		Methods("POST")
}
