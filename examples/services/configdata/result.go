package configdata

import (
	m "github.com/simiancreative/simiango/data/mssql"
	"github.com/simiancreative/simiango/service"
	"fmt"
)

func (s configDataService) Result() (interface{}, error) {
	rows, err := readConfigData()
	if err != nil {
		return nil, err
	}

	return buildResp(rows), nil
}

func readConfigData() ([]ConfigDataResp, error) {
	rows, err := dbReadConfigData()

	if err == nil {
		return rows, nil
	}

	message := "connectivity_issue"
	return nil, &service.ResultError{
		Field:      "ani",
		Status:     503,
		ErrMessage: message,
		Message:    message,
	}
}

func buildResp(v []ConfigDataResp) []ConfigDataResp {
	return v
}

func dbReadConfigData() ([]ConfigDataResp, error) {
	rows := []ConfigDataResp{}
	err := m.Cx.Select(&rows, "EXEC acacia.usp_configuration_info_sel")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
