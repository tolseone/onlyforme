package dau

import (
	"context"
	"sync"
	"time"
)

type Service struct {
	mu          sync.RWMutex
	currentDay  string
	currentData map[int]map[int]struct{} // authorID -> userIDs set
	prevData    map[int]map[int]struct{} // authorID -> userIDs set from previous day
}

type EventRequest struct {
	UserID   int `json:"user_id"`
	AuthorID int `json:"author_id"`
}

func (s *Service) Event(ctx context.Context, request EventRequest) error {
	// Check if request is valid
	if request.UserID == 0 || request.AuthorID == 0 {
		return nil // or return error if needed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if day has changed and rotate data if needed
	s.checkDayChange()

	// Initialize author map if not exists
	if _, ok := s.currentData[request.AuthorID]; !ok {
		s.currentData[request.AuthorID] = make(map[int]struct{})
	}

	// Add user to author's set
	s.currentData[request.AuthorID][request.UserID] = struct{}{}

	return nil
}

func (s *Service) Dau(ctx context.Context, authorIDs []int) ([]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Use previous day's data
	data := s.prevData

	result := make([]int, len(authorIDs))
	for i, authorID := range authorIDs {
		if users, ok := data[authorID]; ok {
			result[i] = len(users)
		} else {
			result[i] = 0
		}
	}

	return result, nil
}

func NewService() *Service {
	return &Service{
		currentData: make(map[int]map[int]struct{}),
		prevData:    make(map[int]map[int]struct{}),
		currentDay:  time.Now().Format("2006-01-02"),
	}
}

// checkDayChange checks if day has changed and rotates data if needed
func (s *Service) checkDayChange() {
	currentDay := time.Now().Format("2006-01-02")
	if currentDay != s.currentDay {
		// Day has changed, rotate data
		s.prevData = s.currentData
		s.currentData = make(map[int]map[int]struct{})
		s.currentDay = currentDay
	}
}
