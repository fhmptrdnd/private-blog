package repository

import (
	"sync"

	"github.com/fhmptrdnd/weather-api-test-web-based/internal/models"
)

// newmemoryrepo, bikin repo baru dengan closure
// return repository struct yang isinya function-function
func NewMemoryRepo() Repository {
    // state disimpan dalam closure (encapsulation via closure)
    articles := make(map[string]models.Article)
    var mu sync.RWMutex

    return Repository{
        Create: func(a models.Article) error {
            mu.Lock()
            defer mu.Unlock()
            articles[a.ID] = a
            return nil
        },
        Get: func(id string) (models.Article, error) {
            mu.RLock()
            defer mu.RUnlock()
            a, ok := articles[id]
            if !ok {
                return models.Article{}, ErrNotFound
            }
            return a, nil
        },
        Update: func(a models.Article) error {
            mu.Lock()
            defer mu.Unlock()
            if _, ok := articles[a.ID]; !ok {
                return ErrNotFound
            }
            articles[a.ID] = a
            return nil
        },
        Delete: func(id string) error {
            mu.Lock()
            defer mu.Unlock()
            if _, ok := articles[id]; !ok {
                return ErrNotFound
            }
            delete(articles, id)
            return nil
        },
        ListByOwner: func(ownerID string) ([]models.Article, error) {
            mu.RLock()
            defer mu.RUnlock()
            var result []models.Article
            for _, a := range articles {
                if a.OwnerID == ownerID {
                    result = append(result, a)
                }
            }
            return result, nil
        },
    }
}
