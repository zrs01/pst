package converter

type ProgSpec struct {
	Modules []Module `yaml:"modules,omitempty"`
}

type Module struct {
	Name     string    `yaml:"name,omitempty"`
	Features []Feature `yaml:"features,omitempty"`
}
type Feature struct {
	Id        string     `yaml:"id,omitempty"`
	Name      string     `yaml:"name,omitempty"`
	Mode      string     `yaml:"mode,omitempty"`
	Desc      string     `yaml:"desc,omitempty"`
	Env       Env        `yaml:"env,omitempty"`
	Resources []Resource `yaml:"resources,omitempty"`
	Input     []Input    `yaml:"input,omitempty"`
	Scenarios []Scenario `yaml:"scenarios,omitempty"`
}
type Env struct {
	Sources   []string `yaml:"sources,omitempty"`
	Languages []string `yaml:"langs,omitempty"`
}
type Resource struct {
	Name  string `yaml:"name,omitempty"`
	Usage string `yaml:"usage,omitempty"`
}
type Input struct {
	Name        string   `yaml:"name,omitempty"`
	Fields      []string `yaml:"fields,omitempty"`
	Constraints []string `yaml:"cons,omitempty"`
	Remarks     []string `yaml:"remarks,omitempty"`
}
type Scenario struct {
	Name  string   `yaml:"name,omitempty"`
	Given []string `yaml:"given,omitempty"`
	When  []string `yaml:"when,omitempty"`
	And   []string `yaml:"and,omitempty"`
	But   []string `yaml:"but,omitempty"`
	Then  []string `yaml:"then,omitempty"`
}
