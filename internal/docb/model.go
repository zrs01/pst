package docb

type ProgSpec struct {
	Modules []Module `yaml:"modules,omitempty"`
}

type Module struct {
	Name     string    `yaml:"name,omitempty"`
	Features []Feature `yaml:"features,omitempty"`
}
type Feature struct {
	Id        interface{} `yaml:"id,omitempty"`
	Name      interface{} `yaml:"name,omitempty"`
	Mode      interface{} `yaml:"mode,omitempty"`
	Desc      interface{} `yaml:"desc,omitempty"`
	Env       Env         `yaml:"env,omitempty"`
	Resources []Resource  `yaml:"resources,omitempty"`
	Screens   []Screen    `yaml:"screens,omitempty"`
	Input     []Input     `yaml:"input,omitempty"`
	Scenarios []Scenario  `yaml:"scenarios,omitempty"`
	Remarks   interface{} `yaml:"remarks,omitempty"`
}
type Env struct {
	Sources   interface{} `yaml:"sources,omitempty"`
	Languages interface{} `yaml:"langs,omitempty"`
}
type Resource struct {
	Name  interface{} `yaml:"name,omitempty"`
	Usage interface{} `yaml:"usage,omitempty"`
}
type Screen struct {
	Id    interface{} `yaml:"id,omitempty"`
	Name  interface{} `yaml:"name,omitempty"`
	Image Image       `yaml:"image,omitempty"`
}
type Input struct {
	Name        interface{} `yaml:"name,omitempty"`
	Fields      interface{} `yaml:"fields,omitempty"`
	Constraints interface{} `yaml:"cons,omitempty"`
	Remarks     interface{} `yaml:"remarks,omitempty"`
}
type Scenario struct {
	Name interface{} `yaml:"name,omitempty"`
	Desc []string    `yaml:"desc,omitempty"`
}
type Image struct {
	File  string `yaml:"file,omitempty"`
	Width int    `yaml:"width,omitempty"`
}
