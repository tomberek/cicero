package service

import (
	"database/sql"
	"github.com/input-output-hk/cicero/src/model"
	"github.com/input-output-hk/cicero/src/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

type WorkflowService interface {
	GetAllByName(string) ([]*model.WorkflowInstance, error)
	GetById(uint64) (model.WorkflowInstance, error)
	Save(pgx.Tx, *model.WorkflowInstance) error
	Update(pgx.Tx, uint64, *model.WorkflowInstance) (sql.Result, error)
}

type WorkflowServiceImpl struct {
	logger             *log.Logger
	workflowRepository repository.WorkflowRepository
}

func NewWorkflowService(db *pgxpool.Pool) WorkflowService {
	return &WorkflowServiceImpl{
		logger:             log.New(os.Stderr, "WorkflowService: ", log.LstdFlags),
		workflowRepository: repository.NewWorkflowRepository(db),
	}
}

func (cmd *WorkflowServiceImpl) GetAllByName(name string) ([]*model.WorkflowInstance, error) {
	log.Printf("Get all Workflows by name %s", name)
	return cmd.workflowRepository.GetAllByName(name)
}

func (cmd *WorkflowServiceImpl) GetById(id uint64) (workflow model.WorkflowInstance, err error) {
	log.Printf("Get Workflow by id %d", id)
	workflow, err = cmd.workflowRepository.GetById(id)
	if err != nil {
		log.Printf("Couldn't select existing workflow for id %d: %s\n", id, err)
	}
	return workflow, err
}

func (cmd *WorkflowServiceImpl) Save(tx pgx.Tx, workflow *model.WorkflowInstance) error {
	log.Printf("Saving new workflow %#v", &workflow)
	err := cmd.workflowRepository.Save(tx, workflow)
	if err != nil {
		log.Printf("Couldn't insert workflow: %s", err.Error())
	} else {
		log.Printf("Created workflow %#v", workflow)
	}
	return err
}

func (cmd *WorkflowServiceImpl) Update(tx pgx.Tx, id uint64, workflow *model.WorkflowInstance) (result sql.Result, err error) {
	log.Printf("Update workflow %#v with id %d", workflow, id)
	result, err = cmd.workflowRepository.Update(tx, id, workflow)
	if err != nil {
		log.Printf("Couldn't update workflow with id: %d, error: %s", id, err.Error())
	} else {
		log.Printf("Updated workflow %#v with id %d", workflow, id)
	}
	return result, err
}
