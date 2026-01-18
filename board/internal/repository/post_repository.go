package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/pisue/go-playground/board/internal/domain"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

// PostRepository 인터페이스 정의
type PostRepository interface {
	Save(post *domain.Post) error
	FindByID(id uint) (*domain.Post, error)
	FindAll() ([]*domain.Post, error)
	Update(post *domain.Post) error
	Delete(id uint) error
}

// memoryPostRepository 구조체 구현
type memoryPostRepository struct {
	mu     sync.RWMutex
	posts  map[uint]*domain.Post
	nextID uint
}

// NewMemoryPostRepository 생성자 함수
func NewMemoryPostRepository() PostRepository {
	return &memoryPostRepository{
		posts:  make(map[uint]*domain.Post),
		nextID: 1,
	}
}

func (r *memoryPostRepository) Save(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post.ID = r.nextID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	r.posts[post.ID] = post
	r.nextID++

	return nil
}

func (r *memoryPostRepository) FindByID(id uint) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (r *memoryPostRepository) FindAll() ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *memoryPostRepository) Update(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[post.ID]; !ok {
		return ErrPostNotFound
	}

	post.UpdatedAt = time.Now()
	r.posts[post.ID] = post
	return nil
}

func (r *memoryPostRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[id]; !ok {
		return ErrPostNotFound
	}

	delete(r.posts, id)
	return nil
}
