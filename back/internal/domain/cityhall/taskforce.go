package cityhall

import "context"

type TaskForceExecutor interface {
	Task(ctx context.Context, action func(ctx context.Context) error) error
}

// Taskforce is a bit of a special deputy,
// its role is to act as an overseer of actions of
// multiple deputies that should be contained in one task.
//
// A task consists of read/write ops from different deputies that must remain coherent
// at a given time from a data standpoint.
//
// For example, to build an artifact, we must check if the artifact's requiremnt are met;
// to that end, we must first read the current data to compare (city stats, tiles),
// then we can write if requirements are met.
// During this whole task, relevant data holders must be locked,
// otherwise an other write in those data could lead to data that is irrelevant/outdated
// regarding the task at hand (in this case checking if artifact's requirement are met).
//
// and that is what the Taskforce is for:
//
// garanty of relevant data from the onset of an action to the end of an action.
//
// The implementation of this behavior is left as an exercise for the adapter :)
type TaskForce struct {
	Executor TaskForceExecutor
}

func NewTaskForce(store Store) *TaskForce {
	return &TaskForce{
		Executor: store,
	}
}
