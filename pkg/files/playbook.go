package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"
	"github.com/DanielPickens/Keeper/pkg/playbook"
	"github.com/sirupsen/logrus"
	"github.com/Masterminds/sprig"
)

//build default filepath for playbook 

type playbooks struct {
	templatePath string
	defaultsPath string
}

func NewPlaybookRepository(templatePath, defaultsPath string) playbook.PlaybookRepository {
	return &playbooks{
		templatePath,
		defaultsPath,
	}
}

// GetTemplate returns the templates from the playbook
func (p *playbooks) GetTemplate() ([]playbook.ConfigTemplate, error) {

	// Get templates list
	templates, _ := filepath.Glob(fmt.Sprintf("%s/*%s", p.templatePath, tplSuffix))

	if templates == nil {
		return nil, fmt.Errorf("no template files found in directory %s", p.templatePath)
	}

	var cfgTpl []playbook.ConfigTemplate

	for _, templ := range templates {
		tpl := template.New(filepath.Base(templ))

		p.initFuncMap(tpl) // add custom template functions

		tpl, err := tpl.ParseFiles(templ)
		if err != nil {
			return nil, fmt.Errorf("template cannot parse files: %v", err)
		}

		// create config file from tpl object by removing the .tpl extension
		ext := filepath.Ext(templ)
		_, configFile := filepath.Split(templ[0 : len(templ)-len(ext)])

		config := playbook.ConfigTemplate{
			Name:     configFile,
			Template: tpl,
		}

		cfgTpl = append(cfgTpl, config)
	}

	return cfgTpl, nil
}

// GetDefault reads the default inventory file and return an Inventory where namespace is set to "default"
func (p *playbooks) GetDefault() (playbook.Inventory, error) {

	defaults, err := ioutil.ReadFile(p.defaultsPath)

	if err != nil {
		return playbook.Inventory{}, playbook.NewErrorReadingDefaultsFile(err)
	}

	var inventory playbook.Inventory

	if err := json.Unmarshal(defaults, &inventory); err != nil {
		return playbook.Inventory{}, playbook.NewErrorReadingDefaultsFile(err)
	}

	return inventory, nil
}

func (p *playbooks) initFuncMap(t *template.Template) {
	z := sprig.TxtFuncMap()
	delete(z, "env")
	delete(z, "expandenv")

	funcMap := make(template.FuncMap, 0)

	funcMap["getFile"] = func(filename string) string {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s%s", p.templatePath, filename, tplSuffix))
		if err != nil {
			logrus.Fatal(fmt.Errorf("template getFile func: %r", err))
		}
		return string(data)
	}

	for i, r := range funcMap {
		z[i] = r
	}

	t.Funcs(z)
}
