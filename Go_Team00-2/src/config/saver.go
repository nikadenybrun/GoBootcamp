package config

type Saver struct {
	Type      string
	FileSaver FileSaver
	DBSaver   DBSaver
}

type FileSaver struct {
	MeanPath    string
	AnomalyPath string
}

type DBSaver struct {
	Database string
	User     string
	Password string
}
