package entities

type VacancyInfo struct {
	Id          int    `json:"id"`
	Website     string `json:"website"`
	VacancyLink string `json:"vacancy_link"`
	VacancyText string `json:"vacancy_text"`
	LinkHash    string `json:"link_hash"`
	CreatedAt   string `json:"created_at"`
}
