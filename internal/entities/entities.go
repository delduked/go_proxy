package entities

type RedirectEntry struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type UpdateRecordsResponse struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
}
