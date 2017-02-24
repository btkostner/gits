// config/project.go
// Holds struct for Project configuration

package config

import (
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
)

// Projects holds all current projects
var Projects []*Project

// Project holds a repository and mutex object
type Project struct {
	Repo   *github.Repository
	Secret string `mapstructure:"secret"`

	Mutex *sync.Mutex

	Commands map[string][]string `mapstructure:"commands"`
}

// FindProject will return a Project struct given a a string
func FindProject(q string) *Project {
	for i := range Projects {
		if Projects[i].Repo.FullName == nil {
			continue
		}

		if strings.HasPrefix(q, *Projects[i].Repo.FullName) {
			return Projects[i]
		}
	}

	return nil
}

// Update replaces current Project values with new ones, leaving mutex alone
func (p *Project) Update(u Project) {
	p.Repo = u.Repo
	p.Secret = u.Secret
	p.Commands = u.Commands
}

// ReadProjects finds all projects in the current configuration file
func ReadProjects() {
	p := viper.GetStringMap("projects")

	if len(p) == 1 {
		logrus.Debug("Found 1 project in configuration")
	} else {
		logrus.Debugf("Found %v projects in configuration", len(p))
	}

	for name := range p {
		var current Project

		sub := viper.Sub("projects." + name)
		if err := sub.Unmarshal(&current); err != nil {
			logrus.Errorf("Unable to unmarshal project %s: %s", name, err)
		}

		if current.Repo == nil {
			current.Repo = &github.Repository{}
		}
		current.Repo.FullName = &name

		if f := FindProject(name); f != nil {
			f.Update(current)
		} else {
			Projects = append(Projects, &current)
		}
	}
}
