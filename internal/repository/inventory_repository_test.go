package repository

import (
	"testing"

	"github.com/RyderDProgrammer/Inventory-Manager/internal/models"
)

func seedItem(t *testing.T, repo InventoryRepository) models.Item {
	t.Helper()
	item, err := repo.Create(models.Item{Name: "Widget", SKU: "WGT-1"})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	return item
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name  string
		seed  int
		wantN int
	}{
		{"empty", 0, 0},
		{"one item", 1, 1},
		{"multiple items", 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInventoryRepository()
			for i := 0; i < tt.seed; i++ {
				repo.Create(models.Item{Name: "Item", SKU: "SKU"})
			}
			items, err := repo.GetAll()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(items) != tt.wantN {
				t.Errorf("got %d items, want %d", len(items), tt.wantN)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name    string
		useReal bool
		wantErr error
	}{
		{"found", true, nil},
		{"not found", false, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInventoryRepository()
			id := "nonexistent-id"
			if tt.useReal {
				id = seedItem(t, repo).ID
			}
			_, err := repo.GetByID(id)
			if err != tt.wantErr {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
		item models.Item
	}{
		{"assigns id", models.Item{Name: "Gadget", SKU: "GDG-1"}},
		{"preserves fields", models.Item{Name: "Thing", SKU: "THG-1", Quantity: 5, Price: 9.99}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInventoryRepository()
			created, err := repo.Create(tt.item)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if created.ID == "" {
				t.Error("expected ID to be assigned, got empty string")
			}
			if created.Name != tt.item.Name || created.SKU != tt.item.SKU {
				t.Errorf("fields not preserved: got %+v", created)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		useReal bool
		update  models.Item
		wantErr error
	}{
		{"updates existing", true, models.Item{Name: "Updated", SKU: "UPD-1"}, nil},
		{"not found", false, models.Item{Name: "X", SKU: "Y"}, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInventoryRepository()
			id := "nonexistent-id"
			if tt.useReal {
				id = seedItem(t, repo).ID
			}
			updated, err := repo.Update(id, tt.update)
			if err != tt.wantErr {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
			if err == nil {
				if updated.ID != id {
					t.Errorf("ID changed: got %s, want %s", updated.ID, id)
				}
				if updated.Name != tt.update.Name {
					t.Errorf("Name not updated: got %s, want %s", updated.Name, tt.update.Name)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		useReal bool
		wantErr error
	}{
		{"deletes existing", true, nil},
		{"not found", false, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInventoryRepository()
			id := "nonexistent-id"
			if tt.useReal {
				id = seedItem(t, repo).ID
			}
			err := repo.Delete(id)
			if err != tt.wantErr {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
			if err == nil {
				if _, getErr := repo.GetByID(id); getErr != ErrNotFound {
					t.Error("item still exists after delete")
				}
			}
		})
	}
}
