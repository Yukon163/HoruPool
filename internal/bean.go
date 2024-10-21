package internal

type DstUrl string

type SrcConfig struct {
	HTTP3           bool   `yaml:"http3"`
	SSL             bool   `yaml:"ssl"`
	KeyFile         string `yaml:"key_file"`
	CertFile        string `yaml:"cert_file"`
	SrcHost         string `yaml:"src_host"`
	Allow0RTT       bool   `yaml:"allow_0rtt"`
	EnableDatagrams bool   `yaml:"enable_datagrams"`
}

type PortEntry struct {
	Port int                  `yaml:"port"`
	Maps map[SrcConfig]DstUrl `yaml:"maps"`
}

type HoruConfig struct {
	EntryList      []PortEntry `yaml:"entry_list"`
	EnableCompress bool        `yaml:"enable_compress"`
	BrotliLevel    int         `yaml:"brotli_level"`
	GzipLevel      int         `yaml:"gzip_level"`
	ZlibLevel      int         `yaml:"zlib_level"`
}
