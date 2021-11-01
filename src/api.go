package cicero

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/liftbridge-io/go-liftbridge"
)

type Api struct {
	logger    *log.Logger
	bridge    liftbridge.Client
	evaluator Evaluator
}

func (api *Api) init() {
	if api.logger == nil {
		api.logger = log.New(os.Stderr, "api: ", log.LstdFlags)
	}
}

func (a *Api) WorkflowInstances(name string) ([]*WorkflowInstance, error) {
	instances := []*WorkflowInstance{}

	err := pgxscan.Select(
		context.Background(),
		DB,
		&instances,
		`SELECT * FROM workflow_instances WHERE name = $1`,
		name)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (a *Api) WorkflowInstance(name string, id uint64) (*WorkflowInstance, error) {
	instance := &WorkflowInstance{}

	err := pgxscan.Select(
		context.Background(),
		DB,
		instance,
		`SELECT * FROM workflow_instances WHERE name = $1 AND id = $2`,
		name,
		id)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (a *Api) WorkflowForInstance(wfName string, instanceId *uint64, logger *log.Logger) (WorkflowDefinition, error) {
	if instanceId != nil {
		var def WorkflowDefinition

		instance, err := a.WorkflowInstance(wfName, *instanceId)
		if err != nil {
			return def, err
		}

		def, err = instance.GetDefinition(logger, a.evaluator)
		if err != nil {
			return def, err
		}
		return def, nil
	} else {
		var err error
		def, err := a.Workflow(wfName)
		if err != nil {
			return def, err
		}
		return def, nil
	}
}

// TODO superfluous?
func (a *Api) Workflows() ([]string, error) {
	return a.evaluator.ListWorkflows()
}

// TODO superfluous?
func (a *Api) Workflow(name string) (WorkflowDefinition, error) {
	return a.evaluator.EvaluateWorkflow(name, 0, WorkflowCerts{})
}

func (a *Api) WorkflowStart(name string) error {
	return publish(a.logger, a.bridge, fmt.Sprintf("workflow.%s.start", name), "workflow.*.start", WorkflowCerts{})
}

func (a *Api) Steps() ([]*StepInstance, error) {
	instances := []*StepInstance{}

	err := pgxscan.Select(
		context.Background(),
		DB,
		&instances,
		`SELECT * FROM step_instances`)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (a *Api) Step(id uuid.UUID) (*StepInstance, error) {
	instance := &StepInstance{}

	err := pgxscan.Select(
		context.Background(),
		DB,
		instance,
		`SELECT * FROM step_instances WHERE id = $1`,
		id)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
