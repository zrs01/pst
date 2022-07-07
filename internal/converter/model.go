package converter

type ProgSpec struct {
	Modules []Module `yaml:"modules,omitempty"`
}

type Module struct {
	Name     string    `yaml:"name,omitempty"`
	Features []Feature `yaml:"features,omitempty"`
}
type Feature struct {
	Id        string      `yaml:"id,omitempty"`
	Name      string      `yaml:"name,omitempty"`
	Mode      string      `yaml:"mode,omitempty"`
	Desc      string      `yaml:"desc,omitempty"`
	Env       Env         `yaml:"env,omitempty"`
	Resources []Resource  `yaml:"resources,omitempty"`
	Input     []Input     `yaml:"input,omitempty"`
	Scenarios []Scenario  `yaml:"scenarios,omitempty"`
	Remarks   interface{} `yaml:"remarks,omitempty"`
}
type Env struct {
	Sources   interface{} `yaml:"sources,omitempty"`
	Languages interface{} `yaml:"langs,omitempty"`
}
type Resource struct {
	Name  string      `yaml:"name,omitempty"`
	Usage interface{} `yaml:"usage,omitempty"`
}
type Input struct {
	Name        string      `yaml:"name,omitempty"`
	Fields      interface{} `yaml:"fields,omitempty"`
	Constraints interface{} `yaml:"cons,omitempty"`
	Remarks     interface{} `yaml:"remarks,omitempty"`
}
type Scenario struct {
	Name string   `yaml:"name,omitempty"`
	Desc []string `yaml:"desc,omitempty"`
}
