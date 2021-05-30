package config

type RawRules struct {
	Map map[string]string `yaml:"column-mappings" validate:"required"`
}

func (c *RawRules) ToConfig() map[string]string {
	return c.Map
}
