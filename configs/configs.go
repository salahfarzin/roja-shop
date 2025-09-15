package configs

import "os"

type DBConfigs struct {
	URL string `env:"DB_URL,required"`
}

type Log struct {
	Level string `env:"LOG_LEVEL,required"`
	Path  string `env:"LOG_PATH,required"`
}

type Configs struct {
	Env        string `env:"APP_ENV,required"`
	Version    string `env:"APP_VERSION,required"`
	URL        string `env:"APP_URL,required"`
	Port       int32  `env:"APP_PORT,required"`
	UploadPath string `env:"UPLOAD_PATH,required"`
	DB         DBConfigs
	Log        Log
}

func New() *Configs {
	return &Configs{
		DB: DBConfigs{
			URL: func() string {
				if url := os.Getenv("DB_URL"); url != "" {
					return url
				}
				return "database/app.db"
			}(),
		},
		UploadPath: func() string {
			if dir := os.Getenv("UPLOAD_PATH"); dir != "" {
				return dir
			}
			return "storage/uploads"
		}(),
	}
}
