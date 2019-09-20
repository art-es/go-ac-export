package activecollab

import (
	"log"
	"reflect"

	"github.com/art-es/ac-export/src/path"

	"github.com/art-es/ac-export/src/logging"
	"github.com/art-es/ac-export/src/utils"
	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
)

const (
	EstimateType = "TrackingFilter"
	TasksType    = "AssignmentFilter"
)

type ReportExportPayload struct {
	Name                 string `json:"name"`
	TypeFilter           string `json:"type_filter"`
	JobTypeFilter        string `json:"job_type_filter"`
	TrackedByFilter      string `json:"tracked_by_filter"`
	TrackedOnFilter      string `json:"tracked_on_filter"`
	ProjectFilter        string `json:"project_filter"`
	BillableStatusFilter string `json:"billable_status_filter"`
	GroupBy              string `json:"group_by"`
	Type                 string `json:"type"`
	IncludeAllProjects   string `json:"include_all_projects"`
}

func (ac *ActiveCollab) ExportEstimates() {
	ac.ReportExportPayload.Type = EstimateType
	ac.ReportExport(path.ExportedEstimates)
}

func (ac *ActiveCollab) ExportTasks() {
	ac.ReportExportPayload.Type = TasksType
	ac.ReportExport(path.ExportedTasks)
}

func (ac *ActiveCollab) ReportExport(filename string) {
	query := make(map[string]string)
	for key, value := range structs.Map(ac.ReportExportPayload) {
		if t := reflect.TypeOf(value); t.Kind() != reflect.String {
			continue
		}

		query[strcase.ToSnake(key)] = reflect.ValueOf(value).String()
	}

	headers := map[string]string{"X-Angie-AuthApiToken": ac.Token}

	resp := utils.SendRequestWithQueryBody("GET", ReportExportUri, query, headers)
	logging.RequestInfoWithoutBody(resp)

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		log.Fatal("reports export request status not successful\nstatus code: \n", resp.StatusCode)
	}

	utils.SaveToFile(resp.Body, filename)
}
