package commands

type Collection struct {
	Page           int `json:page`
	PageCount      int `json:pageCount`
	RecordsPerPage int `json:recordsPerPage`
	RecordCount    int `json:recordCount`
}
