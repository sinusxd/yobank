package postgres

import (
	"context"
	"strconv"

	"gorm.io/gorm"
	"yobank/domain"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {
	t := &Task{}
	t.FromDomain(task)

	result := tr.db.Create(t)
	if result.Error != nil {
		return result.Error
	}

	task.ID = t.ID
	return nil
}

func (tr *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	result := tr.db.Where("user_id = ?", uint(uid)).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	domainTasks := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		domainTasks[i] = *task.ToDomain()
	}

	return domainTasks, nil
}
