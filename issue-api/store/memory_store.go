package store

import (
	"errors"
	"sync"
	"time"

	"issue-api/models"
)

type MemoryStore struct {
	Users  map[uint]models.User
	Issues map[uint]models.Issue
	mu     sync.Mutex
	nextID uint
}

func NewStore() *MemoryStore {
	users := map[uint]models.User{
		1: {ID: 1, Name: "김개발"},
		2: {ID: 2, Name: "이디자인"},
		3: {ID: 3, Name: "박기획"},
	}
	return &MemoryStore{
		Users:  users,
		Issues: make(map[uint]models.Issue),
		nextID: 1,
	}
}

var store = NewStore()

func GetStore() *MemoryStore {
	return store
}

func (s *MemoryStore) CreateIssue(issue models.Issue, userId *uint) (models.Issue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	issue.ID = s.nextID
	issue.CreatedAt = now
	issue.UpdatedAt = now

	if userId != nil {
		user, exists := s.Users[*userId]
		if !exists {
			return models.Issue{}, errors.New("존재하지 않는 사용자입니다")
		}
		issue.User = &user
		issue.Status = "IN_PROGRESS"
	} else {
		issue.Status = "PENDING"
	}

	s.Issues[s.nextID] = issue
	s.nextID++
	return issue, nil
}

func (s *MemoryStore) GetAllIssues(status *string) []models.Issue {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Issue
	for _, issue := range s.Issues {
		if status == nil || issue.Status == *status {
			result = append(result, issue)
		}
	}
	return result
}

func (s *MemoryStore) GetIssueByID(id uint) (models.Issue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	issue, exists := s.Issues[id]
	if !exists {
		return models.Issue{}, errors.New("이슈를 찾을 수 없습니다")
	}
	return issue, nil
}

func (s *MemoryStore) UpdateIssue(id uint, update models.Issue, userId *uint) (models.Issue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	issue, exists := s.Issues[id]
	if !exists {
		return models.Issue{}, errors.New("이슈를 찾을 수 없습니다")
	}

	if issue.Status == "COMPLETED" || issue.Status == "CANCELLED" {
		return models.Issue{}, errors.New("완료되거나 취소된 이슈는 수정할 수 없습니다")
	}

	if update.Title != "" {
		issue.Title = update.Title
	}
	if update.Description != "" {
		issue.Description = update.Description
	}

	// 담당자 제거
	if userId != nil && *userId == 0 {
		issue.User = nil
		issue.Status = "PENDING"
	} else if userId != nil {
		user, ok := s.Users[*userId]
		if !ok {
			return models.Issue{}, errors.New("존재하지 않는 사용자입니다")
		}
		issue.User = &user
		if issue.Status == "PENDING" && update.Status == "" {
			issue.Status = "IN_PROGRESS"
		}
	}

	// 상태 업데이트
	validStatus := map[string]bool{
		"PENDING":     true,
		"IN_PROGRESS": true,
		"COMPLETED":   true,
		"CANCELLED":   true,
	}
	if update.Status != "" {
		if !validStatus[update.Status] {
			return models.Issue{}, errors.New("유효하지 않은 상태입니다")
		}
		if issue.User == nil && update.Status != "PENDING" && update.Status != "CANCELLED" {
			return models.Issue{}, errors.New("담당자 없이 이 상태로 변경할 수 없습니다")
		}
		issue.Status = update.Status
	}

	issue.UpdatedAt = time.Now()
	s.Issues[id] = issue
	return issue, nil
}
