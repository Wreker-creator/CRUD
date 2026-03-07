package request

import (
	"encoding/json"
	"fmt"
	"io"
)

type List []Task

func (l List) Find(id int) (*Task, bool) {

	for i := range l {
		if l[i].ID == id {
			return &l[i], true
		}
	}
	return &Task{}, false

}

func NewList(rdr io.Reader) ([]Task, error) {
	var tasks []Task
	err := json.NewDecoder(rdr).Decode(&tasks)
	if err != nil {
		err = fmt.Errorf("Error while parsing the json %v", err)
	}
	return tasks, err
}
