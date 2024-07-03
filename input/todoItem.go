package input

import (
	"errors"
	myDao "github.com/wayne011872/myTodoBackend/dao"
)

type TodoItemInput struct {
	*myDao.TodoItem `json:",inline"`
}

func (ti *TodoItemInput) Validate() error {
	if ti.ID == 0 {
		return errors.New("missing id")
	}
	return nil
}