package conf

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string

	// console
	ConsolePort   int
	ConsolePrompt string = "Leaf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
	ProdMode        bool //是否是生产模式
)
