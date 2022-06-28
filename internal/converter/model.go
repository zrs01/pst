package converter

type ProgSpec struct {
	Modules []Module `yaml:"module,omitempty"`
}

type Module struct {
	Name     string    `yaml:"name,omitempty"`
	Features []Feature `yaml:"features,omitempty"`
}
type Feature struct {
	Id          string       `yaml:"id,omitempty"`
	Name        string       `yaml:"name,omitempty"`
	Mode        string       `yaml:"mode,omitempty"`
	Desc        string       `yaml:"desc,omitempty"`
	Env         Env          `yaml:"env,omitempty"`
	Tables      []Table      `yaml:"tables,omitempty"`
	Validations []Validation `yaml:"validations,omitempty"`
	Scenarios   []Scenario   `yaml:"scenarios,omitempty"`
}
type Env struct {
	Source string `yaml:"source,omitempty"`
	Lang   string `yaml:"lang,omitempty"`
}
type Table struct {
	Name  string `yaml:"name,omitempty"`
	Usage string `yaml:"usage,omitempty"`
}
type Validation struct {
	Input    string `yaml:"input,omitempty"`
	Validate string `yaml:"validate,omitempty"`
	Remarks  string `yaml:"remarks,omitempty"`
}
type Scenario struct {
	Given string `yaml:"given,omitempty"`
	When  string `yaml:"when,omitempty"`
	Then  string `yaml:"then,omitempty"`
	And   string `yaml:"and,omitempty"`
	But   string `yaml:"but,omitempty"`
}
