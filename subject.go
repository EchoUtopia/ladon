package ladon


type Subject interface {
	GetID() string
	GetConditions() Conditions
	GetTenant() string
}

type PlainSubject string

func (p PlainSubject) GetID() string {
	return string(p)
}

func (p PlainSubject) GetConditions() Conditions {
	return nil
}

func (p PlainSubject) GetTenant() string {
	return ``
}


type TenantSubject struct {
	ID string
	Conditions Conditions
	Tenant string
}

func (t *TenantSubject) GetID() string {
	return t.ID
}

func (t *TenantSubject) GetConditions() Conditions {
	return t.Conditions
}

func (t *TenantSubject) GetTenant() string {
	return t.Tenant
}
