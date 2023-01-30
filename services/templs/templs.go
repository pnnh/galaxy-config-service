package templs

import (
	"html/template"
	"strings"

	"github.com/pnnh/multiverse-server/server/helpers"
)

const templatesDir string = "docker/templates"

type Service struct {
	bootstrap *template.Template
}

func NewService() *Service {
	serv := &Service{}

	return serv
}

func (s *Service) Init() error {
	funcMap := helpers.FuncMap()
	bootstrap, err := template.New("templates").Funcs(funcMap).ParseGlob(templatesDir + "/**/*.mst")
	if err != nil {
		return err
	}
	s.bootstrap = bootstrap
	return nil
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Execute(name string, data map[string]interface{}) (string, error) {
	buf := &strings.Builder{}
	if err := s.bootstrap.ExecuteTemplate(buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
