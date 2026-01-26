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
		return nil, &domain.AppError{
			ServiceError: domain.ErrBadRequest,
			Detail:       err,
			Message:      "게시물이 저장되지 않았습니다.",
		}
	}
	return post, nil
}

func (s *postService) GetPost(id uint) (*domain.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, &domain.AppError{
			ServiceError: domain.ErrNotFound,
			Detail:       err,
			Message:      "해당 게시글을 찾을 수 없습니다.",
		}
	}
	return post, nil
}

func (s *postService) ListPosts() ([]*domain.Post, error) {
	posts, err := s.repo.FindAll()
	if err != nil {
		return nil, &domain.AppError{
			ServiceError: domain.ErrInternalFailure,
			Detail:       err,
			Message:      "게시글 목록을 불러오는 중 오류가 발생했습니다.",
		}
	}
	return posts, nil
}

func (s *postService) UpdatePost(id uint, title, content string) (*domain.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, &domain.AppError{
			ServiceError: domain.ErrNotFound,
			Detail:       err,
			Message:      "해당 게시글이 존재하지 않습니다.",
		}
	}

	post.Title = title
	post.Content = content

	if err := s.repo.Update(post); err != nil {
		return nil, &domain.AppError{
			ServiceError: domain.ErrBadRequest,
			Detail:       err,
			Message:      "게시물이 수정되지 않았습니다.",
		}
	}
	return post, nil
}

func (s *postService) DeletePost(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return &domain.AppError{
			ServiceError: domain.ErrNotFound,
			Detail:       err,
			Message:      "삭제하려는 게시글이 존재하지 않습니다.",
		}
	}

	if err := s.repo.Delete(id); err != nil {
		return &domain.AppError{
			ServiceError: domain.ErrInternalFailure,
			Detail:       err,
			Message:      "게시글 삭제에 실패했습니다.",
		}
	}
	return nil
}
