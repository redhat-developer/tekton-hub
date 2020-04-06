package model

import "github.com/jinzhu/gorm"

func initData(db *gorm.DB) {

	cat := &Catalog{
		Name:       "Catalog",
		Type:       "Official",
		Owner:      "tektoncd",
		URL:        "https://github.com/tektoncd/catalog",
		ContextDir: "/tasks",
	}
	db.Create(cat)

	for _, resource := range initResources {
		db.Model(&cat).Association("Resources").Append(&resource)
	}

	for _, resourceTag := range initResourceTags {
		db.Create(&resourceTag)
	}
}

var initResources = []Resource{

	Resource{
		Name:   "ansible-tower-cli",
		Type:   "task",
		Rating: 0,
		Versions: []ResourceVersion{{
			Description: "Ansible Tower (formerly ‘AWX’) is a web-based solution that makes Ansible even more easy to use for IT teams of all kinds, It provides the tower-cli(Tower-CLI) command line tool that simplifies the tasks of starting jobs, workflow jobs, manage users, projects etc.",
			Version:     "1.0",
			URL:         "https://github.com/tektoncd/catalog/tasks/ansible-tower-cli/1.0",
		},
			{
				Description: "Ansible Tower (formerly ‘AWX’) is a web-based solution that makes Ansible even more easy to use for IT teams of all kinds, It provides the tower-cli(Tower-CLI) command line tool that simplifies the tasks of starting jobs, workflow jobs, manage users, projects etc.",
				Version:     "2.0",
				URL:         "https://github.com/tektoncd/catalog/tasks/ansible-tower-cli/2.0",
			}},
	},
	Resource{
		Name:   "argocd",
		Type:   "task",
		Rating: 0,
		Versions: []ResourceVersion{{
			Description: "This task syncs (deploys) an Argo CD application and waits for it to be healthy. To do so, it requires the address of the Argo CD server and some form of authentication - either a username/password or an authentication token.",
			Version:     "1.0",
			URL:         "https://github.com/tektoncd/catalog/tasks/argocd/1.0",
		},
			{
				Description: "This task syncs (deploys) an Argo CD application and waits for it to be healthy. To do so, it requires the address of the Argo CD server and some form of authentication - either a username/password or an authentication token.",
				Version:     "2.0",
				URL:         "https://github.com/tektoncd/catalog/tasks/argocd/2.0",
			},
			{
				Description: "This task syncs (deploys) an Argo CD application and waits for it to be healthy. To do so, it requires the address of the Argo CD server and some form of authentication - either a username/password or an authentication token.",
				Version:     "3.0",
				URL:         "https://github.com/tektoncd/catalog/tasks/argocd/3.0",
			}},
	},
	Resource{
		Name:   "azure-cli",
		Type:   "task",
		Rating: 0,
		Versions: []ResourceVersion{{
			Description: "This task performs operations on Microsoft Azure resources using az.",
			Version:     "1.0",
			URL:         "https://github.com/tektoncd/catalog/tasks/azure-cli/1.0",
		}},
	},
}

var initResourceTags = []ResourceTag{
	ResourceTag{
		ResourceID: 1,
		TagID:      2,
	},
	ResourceTag{
		ResourceID: 1,
		TagID:      7,
	},
	ResourceTag{
		ResourceID: 1,
		TagID:      9,
	},
	ResourceTag{
		ResourceID: 2,
		TagID:      2,
	},
	ResourceTag{
		ResourceID: 2,
		TagID:      5,
	},
	ResourceTag{
		ResourceID: 3,
		TagID:      8,
	},
}
