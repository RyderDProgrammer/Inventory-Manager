package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RyderDProgrammer/Inventory-Manager/internal/models"
	"github.com/RyderDProgrammer/Inventory-Manager/internal/repository"
)

func newHandler() *Handler {
	return NewHandler(repository.NewInventoryRepository())
}

func seedItem(t *testing.T, h *Handler) models.Item {
	t.Helper()
	body, _ := json.Marshal(models.Item{Name: "Widget", SKU: "WGT-1"})
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.CreateItem(rec, req)
	var item models.Item
	if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return item
}

func TestHealthCheck(t *testing.T) {
	h := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.HealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal("response is not valid JSON")
	}
	if body["status"] != "ok" {
		t.Errorf("got status %q, want \"ok\"", body["status"])
	}
}

func TestListItems(t *testing.T) {
	tests := []struct {
		name  string
		seed  int
		wantN int
	}{
		{"empty repo returns array", 0, 0},
		{"seeded repo returns items", 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler()
			for i := 0; i < tt.seed; i++ {
				seedItem(t, h)
			}
			req := httptest.NewRequest(http.MethodGet, "/items", nil)
			rec := httptest.NewRecorder()
			h.ListItems(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("got %d, want 200", rec.Code)
			}
			var items []models.Item
			if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
				t.Fatal("response is not a JSON array")
			}
			if len(items) != tt.wantN {
				t.Errorf("got %d items, want %d", len(items), tt.wantN)
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	tests := []struct {
		name     string
		useReal  bool
		wantCode int
	}{
		{"found", true, http.StatusOK},
		{"not found", false, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler()
			id := "nonexistent-id"
			if tt.useReal {
				id = seedItem(t, h).ID
			}
			req := httptest.NewRequest(http.MethodGet, "/items/"+id, nil)
			req.SetPathValue("id", id)
			rec := httptest.NewRecorder()
			h.GetItem(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("got %d, want %d", rec.Code, tt.wantCode)
			}
			if tt.useReal {
				var item models.Item
				if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
					t.Fatal("response is not a valid item JSON")
				}
				if item.ID != id {
					t.Errorf("got ID %q, want %q", item.ID, id)
				}
			}
		})
	}
}

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{"valid body returns 201 with ID", models.Item{Name: "Gadget", SKU: "GDG-1"}, http.StatusCreated},
		{"missing name returns 400", models.Item{SKU: "GDG-1"}, http.StatusBadRequest},
		{"missing SKU returns 400", models.Item{Name: "Gadget"}, http.StatusBadRequest},
		{"bad JSON returns 400", "not-json{", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler()
			var buf bytes.Buffer
			switch v := tt.body.(type) {
			case string:
				buf.WriteString(v)
			default:
				_ = json.NewEncoder(&buf).Encode(v)
			}
			req := httptest.NewRequest(http.MethodPost, "/items", &buf)
			rec := httptest.NewRecorder()
			h.CreateItem(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("got %d, want %d", rec.Code, tt.wantCode)
			}
			if tt.wantCode == http.StatusCreated {
				var item models.Item
				if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
					t.Fatal("response is not a valid item JSON")
				}
				if item.ID == "" {
					t.Error("expected ID to be assigned, got empty string")
				}
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name     string
		useReal  bool
		wantCode int
	}{
		{"updates existing returns 200", true, http.StatusOK},
		{"not found returns 404", false, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler()
			id := "nonexistent-id"
			if tt.useReal {
				id = seedItem(t, h).ID
			}
			update := models.Item{Name: "Updated", SKU: "UPD-1"}
			body, _ := json.Marshal(update)
			req := httptest.NewRequest(http.MethodPut, "/items/"+id, bytes.NewReader(body))
			req.SetPathValue("id", id)
			rec := httptest.NewRecorder()
			h.UpdateItem(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("got %d, want %d", rec.Code, tt.wantCode)
			}
			if tt.wantCode == http.StatusOK {
				var item models.Item
				if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
					t.Fatal("response is not a valid item JSON")
				}
				if item.Name != update.Name {
					t.Errorf("got Name %q, want %q", item.Name, update.Name)
				}
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name     string
		useReal  bool
		wantCode int
	}{
		{"deletes existing returns 204", true, http.StatusNoContent},
		{"not found returns 404", false, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHandler()
			id := "nonexistent-id"
			if tt.useReal {
				id = seedItem(t, h).ID
			}
			req := httptest.NewRequest(http.MethodDelete, "/items/"+id, nil)
			req.SetPathValue("id", id)
			rec := httptest.NewRecorder()
			h.DeleteItem(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("got %d, want %d", rec.Code, tt.wantCode)
			}
		})
	}
}
