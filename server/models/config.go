package models

import (
	"fmt"
	"time"

	"github.com/pnnh/quantum-go/services/sqlxsvc"
	"github.com/jmoiron/sqlx"
)

type ConfigtTable struct {
	Pk         string    `json:"pk"`      
	Project    string    `json:"project"`  
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
	Key        string    `json:"key"`
	Value      string    `json:"value"`
}
 
type ConfigModel struct {
	ConfigtTable 
}
 

func SelectConfigs(project string) ([]*ConfigModel, error) {
	sqlText := `select pk, project, create_time, update_time, key, value from configs where project = :project;`

	sqlParams := map[string]interface{}{"project": project}
	var sqlResults []*ConfigtTable

	rows, err := sqlxsvc.NamedQuery(sqlText, sqlParams)
	if err != nil {
		return nil, fmt.Errorf("NamedQuery: %w", err)
	}
	if err = sqlx.StructScan(rows, &sqlResults); err != nil {
		return nil, fmt.Errorf("StructScan: %w", err)
	}

	models := make([]*ConfigModel, 0)

	for _, v := range sqlResults {
		models = append(models,  &ConfigModel{ConfigtTable: *v})
	}

	return models, nil
}
