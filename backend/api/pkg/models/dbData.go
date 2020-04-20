package models

import "github.com/jinzhu/gorm"

func initialiseTables(db *gorm.DB) {
	for _, resource := range initResource {
		db.Create(&resource)
	}
	for _, category := range initCategories {
		db.Create(&category)
	}
	for _, tag := range initTags {
		db.Create(&tag)
	}
	for _, resourceTags := range initResourceTag {
		db.Create(&resourceTags)
	}
	for _, githubDetail := range initGithubDetail {
		db.Create(&githubDetail)
	}
	for _, resourceRawPath := range initResourceRawPath {
		db.Create(&resourceRawPath)
	}
	// Add New data to Existing DB
	updateDbData(db)
}

var initResourceTag = []ResourceTag{
	ResourceTag{
		ResourceID: 121,
		TagID:      102,
	},
	ResourceTag{
		ResourceID: 122,
		TagID:      107,
	},
	ResourceTag{
		ResourceID: 123,
		TagID:      102,
	},
	ResourceTag{
		ResourceID: 124,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 125,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 126,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 127,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 128,
		TagID:      110,
	},
	ResourceTag{
		ResourceID: 129,
		TagID:      103,
	},
	ResourceTag{
		ResourceID: 130,
		TagID:      107,
	},
	ResourceTag{
		ResourceID: 131,
		TagID:      101,
	},
	ResourceTag{
		ResourceID: 132,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 133,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 134,
		TagID:      102,
	},
	ResourceTag{
		ResourceID: 135,
		TagID:      107,
	},
	ResourceTag{
		ResourceID: 136,
		TagID:      110,
	},
	ResourceTag{
		ResourceID: 137,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 138,
		TagID:      101,
	},
	ResourceTag{
		ResourceID: 139,
		TagID:      102,
	},
	ResourceTag{
		ResourceID: 140,
		TagID:      101,
	},
	ResourceTag{
		ResourceID: 141,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 142,
		TagID:      102,
	},
	ResourceTag{
		ResourceID: 143,
		TagID:      102,
	},
}

var initCategories = []Category{
	Category{
		ID:   1,
		Name: "Build Tools",
	},
	Category{
		ID:   2,
		Name: "CLIs",
	},
	Category{
		ID:   3,
		Name: "Cloud",
	},
	Category{
		ID:   4,
		Name: "Deploy",
	},
	Category{
		ID:   5,
		Name: "Image Build",
	},
	Category{
		ID:   6,
		Name: "Notification",
	},
	Category{
		ID:   7,
		Name: "Test Framework",
	},
	Category{
		ID:   8,
		Name: "Other",
	},
}

var initTags = []Tag{
	Tag{
		ID:         101,
		Name:       "build-tool",
		CategoryID: 1,
	},
	Tag{
		ID:         102,
		Name:       "cli",
		CategoryID: 2,
	},
	Tag{
		ID:         103,
		Name:       "gcp",
		CategoryID: 3,
	},
	Tag{
		ID:         104,
		Name:       "aws",
		CategoryID: 3,
	},
	Tag{
		ID:         105,
		Name:       "azure",
		CategoryID: 3,
	},
	Tag{
		ID:         106,
		Name:       "cloud",
		CategoryID: 3,
	},
	Tag{
		ID:         107,
		Name:       "deploy",
		CategoryID: 4,
	},
	Tag{
		ID:         108,
		Name:       "image-build",
		CategoryID: 5,
	},
	Tag{
		ID:         109,
		Name:       "notification",
		CategoryID: 6,
	},
	Tag{
		ID:         110,
		Name:       "test",
		CategoryID: 7,
	},
}

var initResource = []Resource{

	Resource{
		ID:          121,
		Name:        "ansible-tower-cli",
		Type:        "task",
		Description: "Ansible Tower (formerly ‘AWX’) is a web-based solution that makes Ansible even more easy to use for IT teams of all kinds, It provides the tower-cli(Tower-CLI) command line tool that simplifies the tasks of starting jobs, workflow jobs, manage users, projects etc.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          122,
		Name:        "argocd",
		Type:        "task",
		Description: "This task syncs (deploys) an Argo CD application and waits for it to be healthy. To do so, it requires the address of the Argo CD server and some form of authentication - either a username/password or an authentication token.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          123,
		Name:        "azure-cli",
		Type:        "task",
		Description: "This task performs operations on Microsoft Azure resources using az.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          124,
		Name:        "buildah",
		Type:        "task",
		Description: "This Task builds source into a container image using Project Atomic's Buildah build tool. It uses Buildah's support for building from Dockerfiles, using its buildah bud command. This command executes the directives in the Dockerfile to assemble a container image, then pushes that image to a container registry.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          125,
		Name:        "buildkit-daemonless",
		Type:        "task",
		Description: "This buildkit-daemonless Task is similar to buildkit but does not need creating Secret, Deployment, and Service resources for setting up the buildkitd daemon cluster.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          126,
		Name:        "buildkit",
		Type:        "task",
		Description: "This Task builds source into a container image using Moby BuildKit.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          127,
		Name:        "buildpacks",
		Type:        "task",
		Description: "This build template builds source into a container image using Cloud Native Buildpacks.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          128,
		Name:        "conftest",
		Type:        "task",
		Description: "These tasks make it possible to use Conftest within your Tekton pipelines. Conftest is a tool for testing configuration files using Open Policy Agent.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          129,
		Name:        "gcloud",
		Type:        "task",
		Description: "This task performs operations on Google Cloud Platform resources using gcloud.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          130,
		Name:        "gke-deploy",
		Type:        "task",
		Description: "This Task deploys an application to a Google Kubernetes Engine cluster using gke-deploy.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          131,
		Name:        "golang",
		Type:        "task",
		Description: "These Tasks are Golang task to build, test and validate Go projects.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          132,
		Name:        "jib-maven",
		Type:        "task",
		Description: "This Task builds Java/Kotlin/Groovy/Scala source into a container image using Google's Jib tool.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          133,
		Name:        "kaniko",
		Type:        "task",
		Description: "This Task builds source into a container image using Google's kaniko tool.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          134,
		Name:        "kn",
		Type:        "task",
		Description: "This Task performs operations on Knative resources (services, revisions, routes) using kn CLI.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          135,
		Name:        "knctl",
		Type:        "task",
		Description: "This Task deploys (or update) a Knative service. It uses knctl for that, and supports only the deploy subcommand as of today.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          136,
		Name:        "kubeval",
		Type:        "task",
		Description: "This task makes it possible to use Kubeval within your Tekton pipelines. Kubeval is a tool used for validating Kubernetes configuration files.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          137,
		Name:        "makisu",
		Type:        "task",
		Description: "This Task builds source into a container image using uber's makisu tool.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          138,
		Name:        "maven",
		Type:        "task",
		Description: "This Task can be used to run a Maven build.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          139,
		Name:        "openshift-client",
		Type:        "task",
		Description: "OpenShift is a Kubernetes distribution from Red Hat which provides oc, the OpenShift CLI that complements kubectl for simplifying deployment and configuration applications on OpenShift.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          140,
		Name:        "openwhisk",
		Type:        "task",
		Description: "This directory contains Tekton Task which can be used to Build and Serve Knative compatible applications (i.e., OpenWhisk Actions) on Kubernetes.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          141,
		Name:        "s2i",
		Type:        "task",
		Description: "Source-to-Image (S2I) is a toolkit and workflow for building reproducible container images from source code. S2I produces images by injecting source code into a base S2I container image and letting the container prepare that source code for execution. The base S2I container images contains the language runtime and build tools needed for building and running the source code.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          142,
		Name:        "terraform-cli",
		Type:        "task",
		Description: "Terraform is an open-source infrastructure as codesoftware tool created by HashiCorp. It enables users to define and provision a datacenter infrastructure using a high-level configuration language known as Hashicorp Configuration Language (HCL), or optionally JSON",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          143,
		Name:        "tkn",
		Type:        "task",
		Description: "This task performs operations on Tekton resources using tkn.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
}

var initGithubDetail = []GithubDetail{
	{122, "tektoncd", "catalog", "argocd/argocd.yaml", "argocd/README.md"},
	{123, "tektoncd", "catalog", "azure-cli/azure_cli.yaml", "azure-cli/README.md"},
	{121, "tektoncd", "catalog", "ansible-tower-cli/ansible-tower-cli-task.yaml", "ansible-tower-cli/README.md"},
	{124, "tektoncd", "catalog", "buildah/buildah.yaml", "buildah/README.md"},
	{125, "tektoncd", "catalog", "buildkit-daemonless/buildkit-daemonless.yaml", "buildkit-daemonless/README.md"},
	{126, "tektoncd", "catalog", "buildkit/task.yaml", "buildkit/README.md"},
	{127, "tektoncd", "catalog", "buildpacks/buildpacks-v3.yaml", "buildpacks/README.md"},
	{128, "tektoncd", "catalog", "conftest/helm-conftest.yaml", "conftest/README.md"},
	{129, "tektoncd", "catalog", "gcloud/gcloud.yaml", "gcloud/README.md"},
	{130, "tektoncd", "catalog", "gke-deploy/gke-deploy.yaml", "gke-deploy/README.md"},
	{131, "tektoncd", "catalog", "golang/tests.yaml", "golang/README.md"},
	{132, "tektoncd", "catalog", "jib-maven/jib-maven.yaml", "jib-maven/README.md"},
	{133, "tektoncd", "catalog", "kaniko/kaniko.yaml", "kaniko/README.md"},
	{134, "tektoncd", "catalog", "kn/kn.yaml", "kn/README.md"},
	{135, "tektoncd", "catalog", "knctl/knctl-deploy.yaml", "knctl/README.md"},
	{136, "tektoncd", "catalog", "kubeval/kubeval.yaml", "kubeval/README.md"},
	{137, "tektoncd", "catalog", "makisu/makisu.yaml", "makisu/README.md"},
	{138, "tektoncd", "catalog", "maven/maven.yaml", "maven/README.md"},
	{139, "tektoncd", "catalog", "openshift-client/openshift-client-task.yaml", "openshift-client/README.md"},
	{140, "tektoncd", "catalog", "openwhisk/service-account.yaml", "openwhisk/README.md"},
	{141, "tektoncd", "catalog", "s2i/s2i.yaml", "s2i/README.md"},
	{142, "tektoncd", "catalog", "terraform-cli/terraform-cli-task.yaml", "terraform-cli/README.md"},
	{143, "tektoncd", "catalog", "tkn/tkn.yaml", "tkn/README.md"},
}

var initResourceRawPath = []ResourceRawPath{
	{121, "https://raw.githubusercontent.com/tektoncd/catalog/master/ansible-tower-cli/ansible-tower-cli-task.yaml", "task"},
	{122, "https://raw.githubusercontent.com/tektoncd/catalog/master/argocd/argocd.yaml", "task"},
	{123, "https://raw.githubusercontent.com/tektoncd/catalog/master/azure-cli/azure_cli.yaml", "task"},
	{124, "https://raw.githubusercontent.com/tektoncd/catalog/master/buildah/buildah.yaml", "task"},
	{125, "https://raw.githubusercontent.com/tektoncd/catalog/master/buildkit-daemonless/buildkit-daemonless.yaml", "task"},
	{126, "https://raw.githubusercontent.com/tektoncd/catalog/master/buildkit/deployment+service.privileged.yaml", "task"},
	{126, "https://raw.githubusercontent.com/tektoncd/catalog/master/buildkit/deployment+service.rootless.yaml", "task"},
	{126, "https://raw.githubusercontent.com/tektoncd/catalog/master/buildkit/task.yaml", "task"},
	{127, "https://raw.githubusercontent.com/tektoncd/catalog/master/buildpacks/buildpacks-v3.yaml	", "task"},
	{128, "https://raw.githubusercontent.com/tektoncd/catalog/master/conftest/conftest.yaml	", "task"},
	{128, "https://raw.githubusercontent.com/tektoncd/catalog/master/conftest/helm-conftest.yaml	", "task"},
	{129, "https://raw.githubusercontent.com/tektoncd/catalog/master/gcloud/gcloud.yaml", "task"},
	{130, "https://raw.githubusercontent.com/tektoncd/catalog/master/gke-deploy/build-push-gke-deploy.yaml", "task"},
	{130, "https://raw.githubusercontent.com/tektoncd/catalog/master/gke-deploy/gke-deploy.yaml	", "task"},
	{131, "https://raw.githubusercontent.com/tektoncd/catalog/master/golang/build.yaml", "task"},
	{131, "https://raw.githubusercontent.com/tektoncd/catalog/master/golang/lint.yaml", "task"},
	{131, "https://raw.githubusercontent.com/tektoncd/catalog/master/golang/tests.yaml", "task"},
	{132, "https://raw.githubusercontent.com/tektoncd/catalog/master/jib-maven/jib-maven.yaml", "task"},
	{133, "https://raw.githubusercontent.com/tektoncd/catalog/master/kaniko/kaniko.yaml", "task"},
	{134, "https://raw.githubusercontent.com/tektoncd/catalog/master/kn/kn-deployer.yaml", "task"},
	{134, "https://raw.githubusercontent.com/tektoncd/catalog/master/kn/kn.yaml", "task"},
	{135, "https://raw.githubusercontent.com/tektoncd/catalog/master/knctl/knctl-deploy.yaml", "task"},
	{136, "https://raw.githubusercontent.com/tektoncd/catalog/master/kubeval/kubeval.yaml", "task"},
	{137, "https://raw.githubusercontent.com/tektoncd/catalog/master/makisu/makisu.yaml", "task"},
	{138, "https://raw.githubusercontent.com/tektoncd/catalog/master/maven/maven.yaml", "task"},
	{139, "https://raw.githubusercontent.com/tektoncd/catalog/master/openshift-client/openshift-client-kubecfg-task.yaml", "task"},
	{139, "https://raw.githubusercontent.com/tektoncd/catalog/master/openshift-client/openshift-client-task.yaml", "task"},
	{140, "https://raw.githubusercontent.com/tektoncd/catalog/master/openwhisk/openwhisk.yaml", "task"},
	{140, "https://raw.githubusercontent.com/tektoncd/catalog/master/openwhisk/service-account.yaml", "task"},
	{141, "https://raw.githubusercontent.com/tektoncd/catalog/master/s2i/s2i.yaml", "task"},
	{142, "https://raw.githubusercontent.com/tektoncd/catalog/master/terraform-cli/terraform-cli-task.yaml", "task"},
	{143, "https://raw.githubusercontent.com/tektoncd/catalog/master/tkn/tkn.yaml", "task"},
}

// Update Db Data with the latest of Catalog
func updateDbData(db *gorm.DB) {
	for _, resource := range updateResource {
		db.Create(&resource)
	}
	for _, githubDetail := range updateGithubDetail {
		db.Create(&githubDetail)
	}
	for _, resourceRawPath := range updateResourceRawPath {
		db.Create(&resourceRawPath)
	}
	for _, resourceTags := range updateResourceTag {
		db.Create(&resourceTags)
	}
}

var updateResourceTag = []ResourceTag{
	ResourceTag{
		ResourceID: 148,
		TagID:      106,
	},
	ResourceTag{
		ResourceID: 149,
		TagID:      106,
	},
	ResourceTag{
		ResourceID: 150,
		TagID:      106,
	},
	ResourceTag{
		ResourceID: 151,
		TagID:      106,
	},
	ResourceTag{
		ResourceID: 152,
		TagID:      108,
	},
	ResourceTag{
		ResourceID: 155,
		TagID:      109,
	},
	ResourceTag{
		ResourceID: 156,
		TagID:      109,
	},
}

var updateResource = []Resource{
	Resource{
		ID:          144,
		Name:        "git-clone",
		Type:        "task",
		Description: "These Tasks are Git tasks to work with repositories used by other tasks in your Pipeline.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          145,
		Name:        "github-add-comment",
		Type:        "task",
		Description: "This Task will add a comment to a pull request or an issue.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          146,
		Name:        "github-close-issue",
		Type:        "task",
		Description: "This Task will close a pull request or an issue.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          147,
		Name:        "github-set-status",
		Type:        "task",
		Description: "This task will set the CI as running and add a link to the openshift console viewer url.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          148,
		Name:        "gcs-create-bucket",
		Type:        "task",
		Description: "These Tasks are for copying to and from GCS buckets from Pipelines.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          149,
		Name:        "gcs-delete-bucket",
		Type:        "task",
		Description: "These Tasks are for copying to and from GCS buckets from Pipelines.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          150,
		Name:        "gcs-upload",
		Type:        "task",
		Description: "These Tasks are for copying to and from GCS buckets from Pipelines.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          151,
		Name:        "gcs-download",
		Type:        "task",
		Description: "These Tasks are for copying to and from GCS buckets from Pipelines.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          152,
		Name:        "jib-gradle",
		Type:        "task",
		Description: "This Task builds Java/Kotlin/Groovy/Scala source into a container image using Google's Jib tool.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          153,
		Name:        "sendmail",
		Type:        "task",
		Description: "This task sends a simple email to receivers via SMTP server.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          154,
		Name:        "pull-request",
		Type:        "task",
		Description: "This Task allows a user to interact with an SCM (source control management) system through an abstracted interface.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          155,
		Name:        "send-to-channel-slack",
		Type:        "task",
		Description: "These tasks post a simple message to a slack channel. ",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
	Resource{
		ID:          156,
		Name:        "send-to-webhook-slack",
		Type:        "task",
		Description: "This task uses Incomming Webhooks of the slack to send the message.",
		Rating:      0,
		Github:      "http://github.com/tektoncd/catalog",
		Tags:        []string{},
		Verified:    false,
	},
}

var updateGithubDetail = []GithubDetail{
	{144, "tektoncd", "catalog", "git/git-clone.yaml", "git/README.md"},
	{145, "tektoncd", "catalog", "github/add_comment.yaml", "github/README.md"},
	{146, "tektoncd", "catalog", "github/close_issue.yaml", "github/README.md"},
	{147, "tektoncd", "catalog", "github/set_status.yaml", "github/README.md"},
	{148, "tektoncd", "catalog", "gcs/gcs-create-bucket.yaml", "gcs/README.md"},
	{149, "tektoncd", "catalog", "gcs/gcs-delete-bucket.yaml", "gcs/README.md"},
	{150, "tektoncd", "catalog", "gcs/gcs-upload.yaml", "gcs/README.md"},
	{151, "tektoncd", "catalog", "gcs/gcs-download.yaml", "gcs/README.md"},
	{152, "tektoncd", "catalog", "jib-gradle/jib-gradle.yaml", "jib-gradle/README.md"},
	{153, "tektoncd", "catalog", "mail/sendmail.yaml", "mail/README.md"},
	{154, "tektoncd", "catalog", "pullrequest/pr.yaml", "pullrequest/README.md"},
	{155, "tektoncd", "catalog", "slackmessage/send-to-channel-slack.yaml", "slackmessage/README.md"},
	{156, "tektoncd", "catalog", "slackmessage/send-to-webhook-slack.yaml", "slackmessage/README.md"},
}

var updateResourceRawPath = []ResourceRawPath{
	{144, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/git/git-clone.yaml", "task"},
	{145, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/github/add_comment.yaml", "task"},
	{146, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/github/close_issue.yaml", "task"},
	{147, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/github/set_status.yaml", "task"},
	{148, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/gcs/gcs-create-bucket.yaml", "task"},
	{149, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/gcs/gcs-delete-bucket.yaml", "task"},
	{150, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/gcs/gcs-upload.yaml", "task"},
	{151, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/gcs/gcs-download.yaml", "task"},
	{152, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/jib-gradle/jib-gradle.yaml", "task"},
	{153, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/mail/sendmail.yaml", "task"},
	{154, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/pullrequest/pr.yaml", "task"},
	{155, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/slackmessage/send-to-channel-slack.yaml", "task"},
	{156, "https://raw.githubusercontent.com/tektoncd/catalog/v1beta1/slackmessage/send-to-webhook-slack.yaml", "task"},
}
