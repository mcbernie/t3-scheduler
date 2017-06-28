package typo3

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type databaseConfiguration struct {
	database string
	host     string
	password string
	port     string
	username string
}

func (d *databaseConfiguration) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.username, d.password, d.host, d.port, d.database)
}

type mailConfiguration struct {
	FromMailAddress       string
	FromMailName          string
	Transport             string
	TransportMboxFile     string
	TransportSMTPServer   string
	TransportSMTPEncrypt  string
	TransportSMTPUsername string
	TransportSMTPPassword string
}

//Typo3 Loading Typo3 Configuration and all other Typo3 related stuff
type Typo3 struct {
	DatabaseConfiguration databaseConfiguration
	MailConfiguration     mailConfiguration
}

//Load Load typo3 configuration
func Load(path string) Typo3 {
	configurationFile := filepath.Join(path, "typo3conf", "LocalConfiguration.php")
	if stat, err := os.Stat(configurationFile); err != nil || stat.IsDir() {
		panic("LocalConfiguration not Found!")
	}

	var configurationContent string

	if byteContent, err := ioutil.ReadFile(configurationFile); err == nil {
		configurationContent = string(byteContent)
	} else {
		panic(err.Error())
	}

	dbConfig := loadDatabaseConfiguration(configurationContent)
	mailConfig := loadEMailConfiguration(configurationContent)
	return Typo3{
		DatabaseConfiguration: dbConfig,
		MailConfiguration:     mailConfig,
	}
}

func getPartOfConfig(key, content string) string {
	re := regexp.MustCompile(fmt.Sprintf("\\'%s\\' => (.*),", key))
	match := re.FindStringSubmatch(content)

	if len(match) < 1 {
		return ""
	}

	value := match[1]
	if value[0] == '\'' {
		value = strings.Replace(value, "'", "", -1)
	}

	return value
}

func getSectionConfig(section, content string) string {
	re := regexp.MustCompile(fmt.Sprintf("(?sU)'%s' => \\[(.*)\\]", section))
	match := re.FindStringSubmatch(content)
	if len(match) < 1 {
		return ""
	}

	return match[1]
}

func loadDatabaseConfiguration(configurationContent string) databaseConfiguration {
	content := getSectionConfig("DB", configurationContent)

	config := databaseConfiguration{}

	config.database = getPartOfConfig("database", content)
	config.host = getPartOfConfig("host", content)
	config.password = getPartOfConfig("password", content)
	config.port = getPartOfConfig("port", content)
	config.username = getPartOfConfig("username", content)

	return config
}

func loadEMailConfiguration(configurationContent string) mailConfiguration {
	content := getSectionConfig("MAIL", configurationContent)

	config := mailConfiguration{}

	config.FromMailAddress = getPartOfConfig("defaultMailFromAddress", content)
	config.FromMailName = getPartOfConfig("defaultMailFromName", content)
	config.Transport = getPartOfConfig("transport", content)
	config.TransportMboxFile = getPartOfConfig("transport_mbox_file", content)

	config.TransportSMTPServer = getPartOfConfig("transport_smtp_server", content)
	config.TransportSMTPEncrypt = getPartOfConfig("transport_smtp_encrypt", content)
	config.TransportSMTPUsername = getPartOfConfig("transport_smtp_username", content)
	config.TransportSMTPPassword = getPartOfConfig("transport_smtp_password", content)

	return config
}
