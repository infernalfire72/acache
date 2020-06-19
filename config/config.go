package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/infernalfire72/acache/log"
)

type Config struct {
	Database sqlConf `ctag:"Database"`
}

func Create() {
	cdefault := `# Cache Config
Database=akatsuki
Host=
Username=root
Password=`

	f, err := os.Create("./cache.conf")
	if err != nil {
		log.Error(err)
		return
	}

	defer f.Close()
	f.Write([]byte(cdefault))
	log.Info("Created new Config File")
}

func Load() (*Config, error) {
	c := &Config{
		Database: sqlConf{},
	}

	file, err := os.Open("./cache.conf")
	if err != nil {
		if os.IsNotExist(err) {
			Create()
			os.Exit(1337)
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer file.Close()
	kvp := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}

		kvp[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	if err = scanner.Err(); err != nil {
		log.Error(err)
	}

	getValue := func(key, vdefault string) string {
		v := kvp[key]
		if v != "" {
			return v
		} else {
			return vdefault
		}
	}

	c.Database.Database = getValue("Database", "akatsuki")
	c.Database.Host = getValue("Host", "")
	c.Database.Username = getValue("Username", "root")
	c.Database.Password = getValue("Password", "")

	return c, nil
}
