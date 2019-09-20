package googlesheets

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/art-es/ac-export/src/logging"
	"github.com/art-es/ac-export/src/path"
	"github.com/art-es/ac-export/src/utils"
)

type GoogleSheets struct {
	ClientID             string `json:"client_id"`
	ClientSecret         string `json:"client_secret"`
	RefreshToken         string `json:"refresh_token"`
	AccessToken          string `json:"access_token"`
	EstimatesSpreadsheet `json:"spreadsheet_estimate"`
	TasksSpreadsheet     `json:"spreadsheet_tasks"`
	SubtasksSpreadsheet  `json:"spreadsheet_subtasks"`
}

type Spreadsheet struct {
	SpreadsheetID string `json:"spreadsheet_id"`
	SheetID       int    `json:"sheet_id"`
}

type EstimatesSpreadsheet struct {
	Spreadsheet
}

type TasksSpreadsheet struct {
	Spreadsheet
}

type SubtasksSpreadsheet struct {
	Spreadsheet
}

type BatchUpdateRequest struct {
	Requests                     map[string]interface{} `json:"requests"`
	IncludeSpreadsheetInResponse bool                   `json:"includeSpreadsheetInResponse"`
}

type PasteDataRequest struct {
	Coordinate Coordinate `json:"coordinate"`
	Type       string     `json:"type"`
	Delimiter  string     `json:"delimiter"`
	Data       string     `json:"data"`
}

type Coordinate struct {
	SheetId     int `json:"sheetId"`
	RowIndex    int `json:"rowIndex"`
	ColumnIndex int `json:"columnIndex"`
}

const (
	RefreshAccessTokenUri = "https://www.googleapis.com/oauth2/v4/token"
	BaseApiUri            = "https://sheets.googleapis.com/v4/spreadsheets"
)

func (gs *GoogleSheets) ImportCsvFile(filepath string, spreadsheet Spreadsheet) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("\nFile %s not found\nmessage: %s", filepath, err)
	}

	payload := getBatchUpdatePayload(spreadsheet.SheetID, data)

	resp := utils.SendRequestWithJsonBody(
		"POST",
		getBatchUpdateUri(spreadsheet.SpreadsheetID),
		payload,
		map[string]string{
			"Authorization": "Bearer " + gs.AccessToken,
		},
	)
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		log.Fatal("authentication status not successful\nstatus code: \n", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("authentication request body reading error\nmessage: %s", err)
	}
	logging.RequestInfo(resp, content)

	fmt.Println("Successful importing")
}

func (gs *GoogleSheets) ImportTasksAndSubtasks() {
	filename := path.ExportedTasks

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("\nFile %s not found\nmessage: %s", filename, err)
	}

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("reading csv file error\nerror: %s", err)
	}

	recordsLen := len(records)
	tasks := make([][]string, recordsLen)
	subtasks := make([][]string, recordsLen)

	for _, line := range records {
		if len(line) < 2 {
			continue
		}

		if line[1] == "Task" {
			tasks = append(tasks, line)
			continue
		}

		if line[1] == "Subtask" {
			subtasks = append(subtasks, line)
		}
	}

	gs.importTasks(tasks)
	gs.importSubtasks(subtasks)
}

func (gs *GoogleSheets) importTasks(records [][]string) {
	saveRecordsTo(records, path.TmpTasks)
	defer removeFile(path.TmpTasks)

	gs.ImportCsvFile(path.TmpTasks, gs.TasksSpreadsheet.Spreadsheet)
}

func (gs *GoogleSheets) importSubtasks(records [][]string) {
	saveRecordsTo(records, path.TmpSubtasks)
	defer removeFile(path.TmpSubtasks)

	gs.ImportCsvFile(path.TmpSubtasks, gs.SubtasksSpreadsheet.Spreadsheet)
}

func saveRecordsTo(records [][]string, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("error creating file %s\n%v", filepath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, line := range records {
		if len(line) < 1 {
			continue
		}

		err = writer.Write(line)
		if err != nil {
			fmt.Println(err)
		}
	}

	writer.Flush()
}

func getBatchUpdatePayload(sheetID int, data []byte) *BatchUpdateRequest {
	return &BatchUpdateRequest{
		Requests: map[string]interface{}{
			"paste_data": PasteDataRequest{
				Coordinate: Coordinate{
					SheetId:     sheetID,
					RowIndex:    0,
					ColumnIndex: 0,
				},
				Type:      "PASTE_NORMAL",
				Delimiter: ",",
				Data:      string(data),
			},
		},
		IncludeSpreadsheetInResponse: true,
	}
}

func getBatchUpdateUri(spreadsheetID string) string {
	return fmt.Sprintf("%s/%s:batchUpdate", BaseApiUri, spreadsheetID)
}

func removeFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		log.Fatalf("error deleting file %s\n%v", filepath, err)
	}
}
