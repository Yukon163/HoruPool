package internal

type SrcUrl string
type DstUrl string

type PointList struct {
	Port   int               `yaml:"port"`
	Points map[SrcUrl]DstUrl `yaml:"points"`
}

type Config struct {
	PointLists []PointList `yaml:"point_lists"`
	KeyFile    string      `yaml:"key_file"`
	CertFile   string      `yaml:"cert_file"`
}
