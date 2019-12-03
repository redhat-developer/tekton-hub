package upload

// TaskSchema represents schema of a particular task
type TaskSchema struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		TaskRef struct {
			Name string `yaml:"name"`
		} `yaml:"taskRef"`
		Inputs struct {
			Resources []struct {
				Name         string `yaml:"name"`
				ResourceSpec struct {
					Type   string `yaml:"type"`
					Params []struct {
						Name  string `yaml:"name"`
						Value string `yaml:"value"`
					} `yaml:"params"`
				} `yaml:"resourceSpec"`
			} `yaml:"resources"`
			Params []struct {
				Name  string `yaml:"name"`
				Value string `yaml:"value"`
			} `yaml:"params"`
		} `yaml:"inputs"`
	} `yaml:"spec"`
}
