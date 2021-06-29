package mapping

type RawMapping struct {
	Map map[string]string `yaml:"column-mappings" validate:"required"`
}

func (c *RawMapping) ToConfig() map[string]string {
	return c.Map
}
