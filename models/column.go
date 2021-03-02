package models

type Column struct {
	ColumnName string `json:"columnName,omitempty" bigquery:"column_name"`
	DataType   string `json:"dataType,omitempty"   bigquery:"data_type"`
}

type Columns []Column

func (c *Columns) ColumnNames() []string {
	var result []string
	for _, n := range *c {
		result = append(result, n.ColumnName)
	}

	return result
}
