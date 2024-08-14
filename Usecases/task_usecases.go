package usecases

import "github.com/Hailemari/clean_architecture_task_manager/Domain"

type TaskUseCase struct {
    repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
    return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) GetTasks() ([]domain.Task, error) {
    return uc.repo.GetTasks()
}

func (uc *TaskUseCase) GetTask(id string) (domain.Task, bool, error) {
    return uc.repo.GetTaskByID(id)
}

func (uc *TaskUseCase) AddTask(task domain.Task) error {
    if err := task.Validate(); err != nil {
        return err
    }
    return uc.repo.AddTask(task)
}

func (uc *TaskUseCase) UpdateTask(id string, task domain.Task) error {
    return uc.repo.UpdateTask(id, task)
}

func (uc *TaskUseCase) DeleteTask(id string) error {
    return uc.repo.DeleteTask(id)
}