package resources

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type AppConfigProperties map[string]string

func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	fmt.Println("Начинаю загрузку свойств...")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Printf("resource: %v\n", config["resource"])
	fmt.Printf("projectid: %v\n", config["projectid"])
	fmt.Printf("token: %v\n", config["token"])
	fmt.Printf("ifcert: %v\n", config["ifcert"])
	fmt.Printf("certpath: %v\n", config["certpath"])
	fmt.Printf("certpass: %v\n", config["certpass"])
	fmt.Printf("allowduplicate: %v\n", config["allowduplicate"])

	if config["resource"] == "" || config["projectid"] == "" || config["token"] == "" {
		fmt.Println("В application.properties не указан один из обязательных параметров")
	}

	if config["ifcert"] != "true" && config["ifcert"] != "false" {
		fmt.Println("Ошибка заполенения 'ifcert'. Должно быть true/false")
	}

	return config, err
}
