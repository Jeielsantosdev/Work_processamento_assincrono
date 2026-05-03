package routes

import (
	"encoding/json"
	"net/http"
)

// Context é o contexto de requisição
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string
}

// JSON escreve uma resposta JSON
func (c *Context) JSON(statusCode int, data interface{}) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	return json.NewEncoder(c.Writer).Encode(data)
}

// BindJSON vincula o corpo JSON a um struct
func (c *Context) BindJSON(v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// Handler é um manipulador de rota
type Handler func(c *Context) error

// AppRouter é um roteador simples HTTP
type AppRouter struct {
	mux    *http.ServeMux
	routes map[string]Handler
}

// NewAppRouter cria um novo roteador
func NewAppRouter() *AppRouter {
	return &AppRouter{
		mux:    http.NewServeMux(),
		routes: make(map[string]Handler),
	}
}

// POST registra uma rota POST
func (r *AppRouter) POST(path string, handler Handler) {
	r.routes["POST:"+path] = handler
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx := &Context{Writer: w, Request: req, Params: make(map[string]string)}
		if err := handler(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// GET registra uma rota GET
func (r *AppRouter) GET(path string, handler Handler) {
	r.routes["GET:"+path] = handler
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx := &Context{Writer: w, Request: req, Params: make(map[string]string)}
		if err := handler(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// ServeHTTP implementa http.Handler
func (r *AppRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
