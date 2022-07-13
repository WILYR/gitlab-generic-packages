package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"wilyr/gitlabpacks/resources"
)

type PackageList struct {
	Id           float64
	Name         string
	Version      string
	Package_type string
	Status       string
}

type PackageFiles struct {
	Id         float64
	Package_id float64
	File_name  string
}

func main() {

	var command string
	var packname string
	var packtag string
	var filename string
	var resource string
	var token string
	var ifcert string
	var certpath string
	var certpass string
	var projectid string
	var allowduplicate string
	var dirpath string
	var confpath string
	var force string

	flag.StringVar(&command, "c", "", " - Команда серверу (get,send,delete)")
	flag.StringVar(&packname, "pn", "", " - Имя пакета")
	flag.StringVar(&packtag, "pt", "", " - Версия пакета")
	flag.StringVar(&filename, "fn", "", " - Имя файла в пакете")
	flag.StringVar(&resource, "url", "", " - URL ресурса gitlab. По умолчанию из conf файла")
	flag.StringVar(&token, "t", "", " - Токен авторизации на сервере gitlab. По умолчанию из conf файла")
	flag.StringVar(&ifcert, "ic", "", " - Поддержка пользовательского *.pem сертификата. По умолчанию из conf файла")
	flag.StringVar(&certpath, "cp", "", " - Путь к пользовательскому *.pem сертификату. По умолчанию из conf файла")
	flag.StringVar(&certpass, "pw", "", " -  Пароль от пользовательсткого сертификата. По умолчанию из conf файла")
	flag.StringVar(&projectid, "pi", "", " - ID проекта. По умолчанию из conf файла")
	flag.StringVar(&allowduplicate, "ad", "", " - Поддержка дублирования файлов на уровне запроса. По умолчанию из conf файла")
	flag.StringVar(&dirpath, "dir", "", " - Директория файлового каталога. По умолчанию ищет 'out' в текущем")
	flag.StringVar(&confpath, "conf", "conf/application.properties", " - Путь к conf файлу")
	flag.StringVar(&force, "f", "", "Force delete всех повторений файлов и пакетов. По умолчанию из conf файла")

	flag.Parse()

	if command == "" || packname == "" || packtag == "" || filename == "" {
		fmt.Println("########## Система управления регистром пакетов ##########")
		fmt.Println("#################### Версия: v0.1.5 ######################")
		fmt.Println("##########################################################")
		fmt.Println("Доступные тэги: [-help/-h] [-c] [-pn] [-pt] [-fn] [-url] [-t]")
		fmt.Println("             [-pw] [-f] [-ic] [-cp] [-pi] [-ad] [-dir] [-conf]")
		fmt.Println("ОБЯЗАТЕЛЬНЫЕ параметры: [-c] [-pn] [-pt] [-fn]")
		fmt.Println("Каталоги по умолчанию: \n 'out' - Файловый \n 'conf/application.properties' - Свойства приложения")
		fmt.Println("Разработка: n.kovalev")
		fmt.Println("Source: https://github.com/WILYR/gitlab-generic-packages.git")
		fmt.Println("##########################################################")
		return
	}

	properties, err := resources.ReadPropertiesFile(confpath)
	if err != nil {
		fmt.Println(err)
	}

	if resource != "" {
		fmt.Println("Будет переопределено свойство: resource [" + properties["resource"] + "] ---> [" + resource + "]")
		properties["resource"] = resource
	}
	if token != "" {
		fmt.Println("Будет переопределено свойство: token [" + properties["token"] + "] ---> [" + token + "]")
		properties["token"] = token
	}
	if ifcert != "" {
		fmt.Println("Будет переопределено свойство: ifcert [" + properties["ifcert"] + "] ---> [" + ifcert + "]")
		properties["ifcert"] = ifcert
	}
	if certpath != "" {
		fmt.Println("Будет переопределено свойство: certpath [" + properties["certpath"] + "] ---> [" + certpath + "]")
		properties["certpath"] = certpath
	}
	if certpass != "" {
		fmt.Println("Будет переопределено свойство: certpass [" + properties["certpass"] + "] ---> [" + certpass + "]")
		properties["certpass"] = certpass
	}
	if projectid != "" {
		fmt.Println("Будет переопределено свойство: projectid [" + properties["projectid"] + "] ---> [" + projectid + "]")
		properties["projectid"] = projectid
	}
	if allowduplicate != "" {
		fmt.Println("Будет переопределено свойство: allowduplicate [" + properties["allowduplicate"] + "] ---> [" + allowduplicate + "]")
		properties["allowduplicate"] = allowduplicate
	}
	if force != "" {
		fmt.Println("Будет переопределено свойство: force [" + properties["force"] + "] ---> [" + force + "]")
		properties["force"] = force
	}

	if command == "send" {
		if properties["ifcert"] == "true" {
			sendFile(resources.CertClient(properties["certpath"], properties["certpass"]), properties, packname, packtag, filename, outDir(dirpath), properties["force"])
		} else {
			sendFile(&http.Client{}, properties, packname, packtag, filename, outDir(dirpath), properties["force"])
		}
	} else if command == "delete" {
		if properties["ifcert"] == "true" {
			deleteFile(resources.CertClient(properties["certpath"], properties["certpass"]), properties, packname, packtag, filename, properties["force"])
		} else {
			deleteFile(&http.Client{}, properties, packname, packtag, filename, properties["force"])
		}
	} else if command == "get" {
		if properties["ifcert"] == "true" {
			getFile(resources.CertClient(properties["certpath"], properties["certpass"]), properties, packname, packtag, filename, outDir(dirpath))
		} else {
			getFile(&http.Client{}, properties, packname, packtag, filename, outDir(dirpath))
		}
	}

}

func outDir(dirPath string) string {
	if dirPath == "" {
		err := os.Mkdir("out", 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Println(err)

		}
		return "out/"
	} else {
		return dirPath
	}
}

func getFile(client *http.Client, properties map[string]string, packname string, packtag string, filename string, dirPath string) {

	fmt.Printf("Будет загружен следующий файл: %v:%v:%v\n", packname, packtag, filename)

	outFile, err := os.Create(dirPath + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()

	req, err := http.NewRequest(
		"GET", properties["resource"]+"/api/v4/projects/"+properties["projectid"]+"/packages/generic/"+packname+"/"+packtag+"/"+filename, nil,
	)
	req.Header.Add("PRIVATE-TOKEN", properties["token"])

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	io.Copy(outFile, resp.Body)

	if responseAnalyse(resp.StatusCode) == 4 {
		os.Remove(outFile.Name())
	}

}

func sendFile(client *http.Client, properties map[string]string, packname string, packtag string, filename string, dirPath string, force string) {

	if properties["allowduplicate"] == "false" {
		fmt.Println("-----------Пропуск дубликатов отключен, предудущая версия файла будет удалена из гитлаба-----------")
		deleteFile(client, properties, packname, packtag, filename, force)
	}

	fmt.Printf("Будет отправлен следующий файл: %v:%v:%v\n", packname, packtag, filename)
	sendedFile, err := os.Open(dirPath + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sendedFile.Close()

	req, err := http.NewRequest("PUT", properties["resource"]+"/api/v4/projects/"+properties["projectid"]+"/packages/generic/"+packname+"/"+packtag+"/"+filename, sendedFile)
	req.Header.Add("PRIVATE-TOKEN", properties["token"])
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseAnalyse(resp.StatusCode)
}

func deleteFile(client *http.Client, properties map[string]string, packname string, packtag string, filename string, force string) {
	packId := findPackageId(client, properties, packname, packtag, filename)
	if force == "true:files" {
		fmt.Printf("Был добавлен флаг '-f', все файлы с именем "+"%v"+" во всех пакетах "+"%v:%v"+" будут удалены \n", filename, packname, packtag)
		for packIndex, _ := range packId {
			fileId := findPackageFile(client, properties, filename, packId[packIndex])
			for fileIndex, _ := range fileId {
				removeFile(client, properties, packId[packIndex], fileId[fileIndex])
			}
		}
	} else if force == "true:packs" {
		fmt.Printf("Был добавлен флаг '-f', все пакеты с именем "+"%v:%v будут удалены \n", packname, packtag)
		for packIndex, _ := range packId {
			removePackage(client, properties, packId[packIndex])
		}
	} else {
		fmt.Printf("Будет удален следующий файл: %v:%v:%v\n", packname, packtag, filename)
		if len(packId) == 0 {
			fmt.Println("Пакет с таким именем и версией не найден")
			return
		}
		fileId := findPackageFile(client, properties, filename, packId[0])
		if len(fileId) == 0 {
			fmt.Println("Файл с таким именем не найден в пакете, удаление отменено")
			return
		}
		removeFile(client, properties, packId[0], fileId[0])
	}
}

func findPackageId(client *http.Client, properties map[string]string, packname string, packtag string, filename string) []int {

	var packId []int
	req, err := http.NewRequest("GET", properties["resource"]+"/api/v4/projects/"+properties["projectid"]+"/packages", nil)
	req.Close = true
	req.Header.Add("PRIVATE-TOKEN", properties["token"])
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return packId
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return packId
	}

	var pkgList []PackageList

	err = json.Unmarshal(body, &pkgList)
	if err != nil {
		fmt.Println(err)
		return packId
	}

	for _, v := range pkgList {
		if v.Name == packname && v.Version == packtag {
			packId = append(packId, int(v.Id))
		}
	}
	fmt.Println("ID найденных пакетов:", packId)
	return packId
}

func findPackageFile(client *http.Client, properties map[string]string, filename string, packId int) []int {

	var fileId []int
	req, err := http.NewRequest("GET", properties["resource"]+"/api/v4/projects/"+properties["projectid"]+"/packages/"+strconv.Itoa(packId)+"/package_files", nil)
	req.Close = true
	req.Header.Add("PRIVATE-TOKEN", properties["token"])
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return fileId
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return fileId
	}

	var pkgFiles []PackageFiles

	err = json.Unmarshal(body, &pkgFiles)
	if err != nil {
		fmt.Println(err)
		return fileId
	}

	for _, v := range pkgFiles {
		if v.File_name == filename {
			fileId = append(fileId, int(v.Id))
		}
	}
	fmt.Println("ID найденных файлов:", fileId)
	return fileId
}

func removeFile(client *http.Client, properties map[string]string, packId int, fileId int) {
	req, err := http.NewRequest("DELETE", properties["resource"]+"/api/v4/projects/"+properties["projectid"]+"/packages/"+strconv.Itoa(packId)+"/package_files/"+strconv.Itoa(fileId), nil)
	req.Close = true
	req.Header.Add("PRIVATE-TOKEN", properties["token"])
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	responseAnalyse(resp.StatusCode)
}

func removePackage(client *http.Client, properties map[string]string, packId int) {
	req, err := http.NewRequest("DELETE", properties["resource"]+"/api/v4/projects/"+properties["projectid"]+"/packages/"+strconv.Itoa(packId), nil)
	req.Close = true
	req.Header.Add("PRIVATE-TOKEN", properties["token"])
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	responseAnalyse(resp.StatusCode)
}

func responseAnalyse(statusCode int) int {
	fmt.Printf("Код ответа сервера: %v\n", statusCode)
	if statusCode/100 == 2 {
		fmt.Println("Команда выполнена")
	} else if statusCode/100 == 4 {
		fmt.Println("Ошибка запроса на стороне клиента")
	} else if statusCode/100 == 5 {
		fmt.Println("Ошибка на стороне сервера")
	}
	return statusCode / 100
}
