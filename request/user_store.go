package request

type InMemoryUserStore struct {
	Tasks []Task
}

func (i *InMemoryUserStore) GetSpecificTask(id int) (Task, bool) {
	for _, task := range i.Tasks {
		if task.ID == id {
			return task, true
		}
	}
	return Task{}, false
}

func (i *InMemoryUserStore) DeleteTask(id int) bool {
	for index, task := range i.Tasks {
		if task.ID == id {
			i.Tasks = append(i.Tasks[:index], i.Tasks[index+1:]...)
			return true
		}
	}
	return false
}

func (i *InMemoryUserStore) GetAllTasks() []Task {
	return i.Tasks
}

func (i *InMemoryUserStore) AddTask(task Task) {
	i.Tasks = append(i.Tasks, task)
}

func (i *InMemoryUserStore) UpdateTask(id int, task Task) bool {

	for index := range i.Tasks {
		if i.Tasks[index].ID == id {
			i.Tasks[index] = task
			return true
		}
	}
	return false
}
