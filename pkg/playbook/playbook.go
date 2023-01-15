package playbook 

import (
	"text/template"
)

func Configtemplate struct {
	Name string
	Template *template.Template

}

//PlaybookService represents the way playbooks are managed
type PlaybookService interface {
	GetDefault() (Inventory, error)
	GetTemplate() ([]Configtemplate, error)
}

type PlaybookRepository interface {
	GetDefault() (Inventory, error)
	GetTemplate() ([]Configtemplate, error)
}

type playbookService structr {
	playbooks PlaybookRepository
}

//newPlaybookService returns a new PlaybookService
func newPlaybookService(playbooks PlaybookRepository) playbookService {
	return &playbookService {
		playbooks: playbooks, 
	}
}

//GetTemplate returns the templates for a playbook
func (ps *playbookService) GetTemplate() ([]Configtemplate, error) {
	return ps.playbooks.GetTemplate()
}

//GetDefault returns the default inventory for a playbook
func (ps *playbookService) GetDefault() (Inventory, error) {
	return ps.playbooks.GetDefault()
}