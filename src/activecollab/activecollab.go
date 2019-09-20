package activecollab

type ActiveCollab struct {
	Token string

	AuthenticationPayload
	ReportExportPayload
}

const (
	ApiBaseUri      = "https://ac.codezavod.ru/api/v1"
	AuthenticateUri = ApiBaseUri + "/issue-token"
	ReportExportUri = ApiBaseUri + "/reports/export"
)
