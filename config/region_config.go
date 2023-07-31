package config

type RegionConfig struct {
	Options Options
	Config  Config
	Id      int64
	Start   []byte
	End     []byte
}
