package sqlstorage

import (
	"crawler/engine"
	"crawler/sqldb"
	"crawler/storage"
	"encoding/json"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

type SqlStorage struct {
	dataDocker  []*storage.DataCell //分批输出结果缓存
	columnNames []sqldb.Field       //标题字段
	db          sqldb.DBer
	Table       map[string]struct{}
	options
}

func New(opts ...Option) (*SqlStorage, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &SqlStorage{}
	s.options = options
	s.Table = make(map[string]struct{})
	var err error
	s.db, err = sqldb.New(
		sqldb.WithLogger(s.logger),
		sqldb.WithConnUrl(s.sqlUrl),
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getFields(cell *storage.DataCell) []sqldb.Field {
	taskName := cell.Data["Task"].(string)
	ruleName := cell.Data["Rule"].(string)
	fields := engine.GetField(taskName, ruleName)

	var columnNames []sqldb.Field

	for _, field := range fields {
		columnNames = append(columnNames, sqldb.Field{
			Title: field,
			Type:  "MEDIUMTEXT",
		})
	}

	columnNames = append(columnNames,
		sqldb.Field{Title: "Url", Type: "VARCHAR(255)"},
		sqldb.Field{Title: "Time", Type: "VARCHAR(255)"},
	)
	return columnNames
}

func (s *SqlStorage) Save(dataCells ...*storage.DataCell) error {
	for _, cell := range dataCells {
		name := cell.GetTaskName()
		if _, ok := s.Table[name]; !ok {
			//创建表
			columnNames := getFields(cell)

			err := s.db.CreateTable(sqldb.TableData{
				TableName:   name,
				ColumnNames: columnNames,
				AutoKey:     true,
			})
			if err != nil {
				s.logger.Error("create table failed", zap.Error(err))
			}
			s.Table[name] = struct{}{}
		}
		if len(s.dataDocker) >= s.BatchCount {
			s.Flush()
		}
		s.dataDocker = append(s.dataDocker, cell)
	}
	return nil
}
func (s *SqlStorage) Flush() error {
	if len(s.dataDocker) == 0 {
		return nil
	}
	args := make([]interface{}, 0)
	for _, dataCell := range s.dataDocker {
		ruleName := dataCell.Data["Rule"].(string)
		taskName := dataCell.Data["Task"].(string)
		fields := engine.GetField(taskName, ruleName)
		data := dataCell.Data["Data"].(map[string]interface{})
		value := []string{}
		for _, field := range fields {
			v := data[field]
			switch v.(type) {
			case nil:
				value = append(value, "")
			case string:
				value = append(value, v.(string))
			default:
				j, err := json.Marshal(v)
				if err != nil {
					value = append(value, "")
				} else {
					value = append(value, string(j))
				}
			}
		}
		value = append(value, dataCell.Data["Url"].(string), dataCell.Data["Time"].(string))
		for _, v := range value {
			args = append(args, v)
		}
	}
	return s.db.Insert(sqldb.TableData{
		TableName:   s.dataDocker[0].GetTableName(),
		ColumnNames: getFields(s.dataDocker[0]),
		Args:        args,
		DataCount:   len(s.dataDocker),
	})
}
