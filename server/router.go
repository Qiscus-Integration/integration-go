package server

import (
	httpDelivery "integration-go/delivery/http"
	"net/http"
)

// Routers adds the routes for the server's endpoints to the router instance of the server.
func (s *Server) routers() {
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	httpDelivery.NewRoom(s.roomUC).HandleRoute(s.Router)
}
