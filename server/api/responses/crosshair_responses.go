package responses

type CrosshairResponse struct {
	Status      string `json:"status"`
	CHsOnRecord int    `json:"chs_on_record"`
}
