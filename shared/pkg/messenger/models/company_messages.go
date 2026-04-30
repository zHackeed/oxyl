package models

type CompanyCreation struct {
	CompanyId string `json:"company_id"`
}

type CompanyMemberAdded struct {
	CompanyId string `json:"company_id"`
	UserId    string `json:"user_id"`
}

type CompanyMemberRemoved struct {
	CompanyId string `json:"company_id"`
	UserId    string `json:"user_id"`
}

type CompanyWebhookCreation struct {
	CompanyId string `json:"company_id"`
	Endpoint  string `json:"endpoint"`
}

type CompanyWebhookDeletion struct {
	CompanyId string `json:"company_id"`
	Endpoint  string `json:"endpoint"`
}

type CompanyDeletion struct {
	CompanyId string `json:"company_id"`
}
