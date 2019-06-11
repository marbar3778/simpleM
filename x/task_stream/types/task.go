package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Task struct {
	TaskTitle       string           `json:"task_title"`
	TaskDescription string           `json:"task_description"`
	Backers         []sdk.AccAddress `json:"backers"`
	Value           sdk.Coins        `json:"value"`
}

func CreateTask(taskTitle string, taskDescription string, backers []sdk.AccAddress, value sdk.Coins) Task {
	return Task{
		TaskTitle:       taskTitle,
		TaskDescription: taskDescription,
		Backers:         backers,
		Value:           value,
	}
}

func (t Task) String() string {
	return fmt.Sprintf(`
	TaskTitle: %s,
	TaskDescription: %s,
	Backers: %v,
	Value: %v`,
		t.TaskTitle, t.TaskDescription, t.Backers, t.Value)
}

func (t Task) AddBacker(addr sdk.AccAddress) {
	t.Backers = append(t.Backers, addr)
}

func (t Task) AddValue(value sdk.Coins) {
	t.Value.Add(value)
}
