package configdata

import (
	"simiango/data/mssql"
	"github.com/simiancreative/simiango/service"
)

func (s configDataService) Result() (interface{}, error) {
	rows, err := readConfigData()
	if err != nil {
		return nil, err
	}

	return buildResp(rows), nil
}

func readConfigData() error {
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

func buildResp(v []configDataResp) []ConfigDataResp {
	return configDataResp;
	/*
	var configdataresponses []ConfigDataResp
	for _, element := range v {
		//code to transform
		configdatatypes = append(configdatatypes, ConfigDataResp{Tag: tag, PlatformType: platformType, AppHash: appHash, TokenSigningKey: tokenSigningKey, RegistrationExpMilliseconds: registrationExpMilliseconds, TokenExpMinutesseconds: tokenExpMinutesseconds, StatusID: statusID})
	}
	return configdataresponses;
	*/
}

func dbReadConfigData() ([]ConfigDataType, error) {
    ctx := context.Background()

    err := db.PingContext(ctx)
    if err != nil {
        return nil, err
    }

    tsql := fmt.Sprintf("EXEC acacia.usp_configuration_info_sel")
    rows, err := db.QueryContext(ctx, tsql)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

	var configdatatypes []ConfigDataType
    for rows.Next() {
		
		err := rows.Scan(&tag, &platformType, &appHash, &tokenSigningKey, &registrationExpMilliseconds, &tokenExpMinutesseconds, &statusID)
		
        if err != nil {
			return nil, err
        }

        log.Println(tag, platformType)
        configdatatypes = append(configdatatypes, ConfigDataType{Tag: tag, PlatformType: platformType, AppHash: appHash, TokenSigningKey: tokenSigningKey, RegistrationExpMilliseconds: registrationExpMilliseconds, TokenExpMinutesseconds: tokenExpMinutesseconds, StatusID: statusID})
	}
	
	return configdatatypes, nil

}
