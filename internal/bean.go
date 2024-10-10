package internal

type SrcUrl string
type DstUrl string

type PortEntry struct {
	Port int               `yaml:"port"`
	Maps map[SrcUrl]DstUrl `yaml:"maps"`
}

type Config struct {
	EntryList []PortEntry `yaml:"entry_list"`
	KeyFile   string      `yaml:"key_file"`
	CertFile  string      `yaml:"cert_file"`
}
