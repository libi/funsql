package builder

import "github.com/pkg/errors"

var (
	NoTableNameErr = errors.New("[funsql] no table name")
	WhereBetweenParamErr = errors.New("[funsql] where between value must have 2 element")
	WhereParamErr = errors.New("[funsql] where param error")
)
