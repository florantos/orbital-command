package handler

type CreateCrewMemmberRequest struct {
	Name           string   `json:"name"`
	Role           string   `json:"role"`
	Qualifications []string `json:"qualifications"`
}
