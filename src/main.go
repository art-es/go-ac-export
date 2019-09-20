package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/art-es/ac-export/src/path"

	"github.com/art-es/ac-export/src/activecollab"
	"github.com/art-es/ac-export/src/googlesheets"
)

type ConfigurationSpecification struct {
	ActiveCollabAuth         activecollab.AuthenticationPayload `json:"ac_auth"`
	ActiveCollabReportExport activecollab.ReportExportPayload   `json:"ac_report_export"`
	GoogleSheets             googlesheets.GoogleSheets          `json:"google_sheets"`
}

var config *ConfigurationSpecification

func main() {
	initConfig()

	ac := &activecollab.ActiveCollab{
		AuthenticationPayload: config.ActiveCollabAuth,
		ReportExportPayload:   config.ActiveCollabReportExport,
	}

	ac.Authenticate()
	ac.ExportEstimates()
	ac.ExportTasks()

	gs := config.GoogleSheets
	gs.RefreshAccessToken()
	gs.ImportCsvFile(path.ExportedEstimates, gs.EstimatesSpreadsheet.Spreadsheet)
	gs.ImportTasksAndSubtasks()
}

func initConfig() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("config file not found\nmessage: %s", err)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("config file can not be decode\nmessage: %s", err)
	}
}
