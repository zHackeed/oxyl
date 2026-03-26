package models

type CompanyMemberAdded struct {
	CompanyId string `json:"company_id"`
	UserId    string `json:"user_id"`
}

type CompanyMemberRemoved struct {
	CompanyId string `json:"company_id"`
	UserId    string `json:"user_id"`
}

type CompanyDeletion struct {
	CompanyId string `json:"company_id"`
}
