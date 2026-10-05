package web

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"base-go/internal/models"
	"base-go/internal/storage"

	"go.uber.org/zap"
)

//go:embed static
var staticFiles embed.FS

// Server обслуживает веб-интерфейс для хранилища резюме
type Server struct {
	storage storage.Storage
	log     *zap.Logger
	mux     *http.ServeMux
}

// NewServer создаёт новый веб-сервер поверх переданного хранилища
func NewServer(str storage.Storage, log *zap.Logger) *Server {
	s := &Server{storage: str, log: log, mux: http.NewServeMux()}
	s.routes()
	return s
}

// ListenAndServe запускает веб-сервер на указанном адресе
func (s *Server) ListenAndServe(addr string) error {
	s.log.Info("запуск веб-интерфейса", zap.String("addr", addr))
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	static, err := fs.Sub(staticFiles, "static")
	if err != nil {
		s.log.Error("ошибка инициализации статических файлов веб-интерфейса", zap.Error(err))
	} else {
		s.mux.Handle("/", http.FileServer(http.FS(static)))
	}

	s.mux.HandleFunc("GET /api/resumes", s.handleList)
	s.mux.HandleFunc("POST /api/resumes", s.handleSave)
	s.mux.HandleFunc("DELETE /api/resumes", s.handleClear)
	s.mux.HandleFunc("GET /api/resumes/size", s.handleSize)
	s.mux.HandleFunc("GET /api/resumes/{uuid}", s.handleGet)
	s.mux.HandleFunc("DELETE /api/resumes/{uuid}", s.handleDelete)
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			s.log.Error("ошибка кодирования json-ответа", zap.Error(err))
		}
	}
}

// GET /api/resumes
func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	all := s.storage.GetAll()
	if all == nil {
		all = []*models.Resume{}
	}
	s.writeJSON(w, http.StatusOK, all)
}

// GET /api/resumes/size
func (s *Server) handleSize(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]int{"size": s.storage.Size()})
}

// GET /api/resumes/{uuid}
func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	resume := s.storage.Get(uuid)
	if resume == nil || resume.UUID == "" {
		s.writeJSON(w, http.StatusNotFound, map[string]string{"error": "резюме не найдено"})
		return
	}
	s.writeJSON(w, http.StatusOK, resume)
}

// POST /api/resumes
func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	var resume models.Resume
	if err := json.NewDecoder(r.Body).Decode(&resume); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "некорректный json"})
		return
	}
	if resume.UUID == "" {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "поле uuid обязательно"})
		return
	}
	s.storage.Save(&resume)
	s.writeJSON(w, http.StatusCreated, resume)
}

// DELETE /api/resumes/{uuid}
func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	s.storage.Delete(uuid)
	s.writeJSON(w, http.StatusNoContent, nil)
}

// DELETE /api/resumes
func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	s.storage.Clear()
	s.writeJSON(w, http.StatusNoContent, nil)
}
