package repository

import (
	"errors"
	"sync"

	"github.com/RyderDProgrammer/Inventory-Manager/internal/models"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("item not found")

type InventoryRepository interface {
	GetAll() ([]models.Item, error)
	GetByID(id string) (models.Item, error)
	Create(item models.Item) (models.Item, error)
	Update(id string, item models.Item) (models.Item, error)
	Delete(id string) error
}

type inMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]models.Item
}

func NewInventoryRepository() InventoryRepository {
	return &inMemoryRepository{
		data: make(map[string]models.Item),
	}
}

func (r *inMemoryRepository) GetAll() ([]models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]models.Item, 0, len(r.data))
	for _, item := range r.data {
		result = append(result, item)
	}
	return result, nil
}

func (r *inMemoryRepository) GetByID(id string) (models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.data[id]
	if !ok {
		return models.Item{}, ErrNotFound
	}
	return item, nil
}

func (r *inMemoryRepository) Create(item models.Item) (models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = uuid.New().String()
	r.data[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) Update(id string, item models.Item) (models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return models.Item{}, ErrNotFound
	}
	item.ID = id
	r.data[id] = item
	return item, nil
}

func (r *inMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrNotFound
	}
	delete(r.data, id)
	return nil
}
