package service

import (
	"github.com/pisue/go-playground/board/internal/domain"
	"github.com/pisue/go-playground/board/internal/repository"
)

// PostService 인터페이스 정의
type PostService interface {
	CreatePost(title, content, author string) (*domain.Post, error)
	GetPost(id uint) (*domain.Post, error)
	ListPosts() ([]*domain.Post, error)
	UpdatePost(id uint, title, content string) (*domain.Post, error)
	DeletePost(id uint) error
}

// postService 구조체 구현
type postService struct {
	repo repository.PostRepository
}

// NewPostService 생성자 함수 (의존성 주입)
func NewPostService(repo repository.PostRepository) PostService {
	return &postService{
		repo: repo,
	}
}

func (s *postService) CreatePost(title, content, author string) (*domain.Post, error) {
	post := &domain.Post{
		Title:   title,
		Content: content,
		Author:  author,
	}
	if err := s.repo.Save(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) GetPost(id uint) (*domain.Post, error) {
	return s.repo.FindByID(id)
}

func (s *postService) ListPosts() ([]*domain.Post, error) {
	return s.repo.FindAll()
}

func (s *postService) UpdatePost(id uint, title, content string) (*domain.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	post.Title = title
	post.Content = content

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) DeletePost(id uint) error {
	return s.repo.Delete(id)
}
