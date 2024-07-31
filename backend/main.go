package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	connectionString, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	mux := initializeRoutes(db, validator.New())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	server.ListenAndServe()
}

type (
	Gizmo struct {
		ID          int64  `json:"id" db:"id"`
		Name        string `json:"name" db:"name" validate:"required"`
		Description string `json:"description" db:"description" validate:"required"`
	}

	Widget struct {
		ID      int64  `json:"id" db:"id"`
		GizmoID int64  `json:"gizmoId" db:"gizmo_id"`
		Name    string `json:"name" db:"name" validate:"required"`
	}

	errorResponse struct {
		Error       string `json:"error"`
		Description string `json:"description"`
	}
)

func initializeRoutes(db *sqlx.DB, validate *validator.Validate) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{"ok": true}
		err := json.NewEncoder(w).Encode(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /gizmos", func(w http.ResponseWriter, r *http.Request) {
		gizmos := []Gizmo{}

		if err := db.SelectContext(r.Context(), &gizmos, "SELECT * FROM gizmos"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to fetch gizmos"})
			return
		}

		json.NewEncoder(w).Encode(gizmos)
	})

	mux.HandleFunc("POST /gizmos", func(w http.ResponseWriter, r *http.Request) {
		gizmo := Gizmo{}

		err := json.NewDecoder(r.Body).Decode(&gizmo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to decode request body"})
			return
		}

		err = validate.Struct(gizmo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid gizmo"})
			return
		}

		err = db.GetContext(r.Context(), &gizmo, "INSERT INTO gizmos (name, description) VALUES ($1, $2) RETURNING *", gizmo.Name, gizmo.Description)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to insert gizmo"})
			return
		}

		json.NewEncoder(w).Encode(gizmo)
	})

	mux.HandleFunc("GET /gizmos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid ID"})
			return
		}

		gizmo := Gizmo{}
		err = db.GetContext(r.Context(), &gizmo, "SELECT * FROM gizmos WHERE id = $1", id)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "gizmo not found"})
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to fetch gizmo"})
			return
		}

		json.NewEncoder(w).Encode(gizmo)
	})

	mux.HandleFunc("DELETE /gizmos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid ID"})
			return
		}

		tx, err := db.BeginTxx(r.Context(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to start transaction"})
			return
		}

		res, err := tx.ExecContext(r.Context(), "DELETE FROM widgets WHERE gizmo_id = $1", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to delete gizmo widgets"})
			return
		}

		res, err = tx.ExecContext(r.Context(), "DELETE FROM gizmos WHERE id = $1", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to delete gizmo"})
			return
		}

		affected, err := res.RowsAffected()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to fetch rows affected"})
			return
		}

		err = tx.Commit()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to commit transaction"})
			return
		}

		if affected == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})

	mux.HandleFunc("GET /gizmos/{id}/widgets", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid ID"})
			return
		}

		exists := false
		err = db.GetContext(r.Context(), &exists, "SELECT EXISTS (SELECT 1 FROM gizmos WHERE id = $1)", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to fetch gizmo"})
			return
		} else if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse{"not found", "gizmo not found"})
			return
		}

		widgets := []Widget{}
		err = db.SelectContext(r.Context(), &widgets, "SELECT * FROM widgets WHERE gizmo_id = $1", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to fetch widgets"})
			return
		}

		json.NewEncoder(w).Encode(widgets)
	})

	mux.HandleFunc("POST /gizmos/{id}/widgets", func(w http.ResponseWriter, r *http.Request) {
		widget := Widget{}
		err := error(nil)
		widget.GizmoID, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid ID"})
			return
		}

		exists := false
		err = db.GetContext(r.Context(), &exists, "SELECT EXISTS (SELECT 1 FROM gizmos WHERE id = $1)", widget.GizmoID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to fetch gizmo"})
			return
		} else if !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse{"not found", "gizmo not found"})
			return
		}

		err = json.NewDecoder(r.Body).Decode(&widget)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to decode request body"})
			return
		}

		err = validate.Struct(widget)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid widget"})
			return
		}

		err = db.GetContext(r.Context(), &widget, "INSERT INTO widgets (gizmo_id, name) VALUES ($1, $2) RETURNING *", widget.GizmoID, widget.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to insert widget"})
			return
		}

		json.NewEncoder(w).Encode(widget)
	})

	mux.HandleFunc("DELETE /gizmos/{gizmoId}/widgets/{id}", func(w http.ResponseWriter, r *http.Request) {
		gizmoID, err := strconv.ParseInt(r.PathValue("gizmoId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid gizmo ID"})
			return
		}

		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "invalid ID"})
			return
		}

		res, err := db.ExecContext(r.Context(), "DELETE FROM widgets WHERE gizmo_id = $1 AND id = $2", gizmoID, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to delete widget"})
			return
		}

		affected, err := res.RowsAffected()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{err.Error(), "failed to fetch rows affected"})
			return
		}

		if affected == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})

	return mux
}
