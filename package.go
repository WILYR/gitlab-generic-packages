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

	flag.StringVar(&command, "c", "get", "Command(get/send/delete). Default is 'get'")
	flag.StringVar(&packname, "pn", "test_package", "Package name. Default is 'test_package'")
	flag.StringVar(&packtag, "pt", "latest", "Tag package. Default is 'latest'")
	flag.StringVar(&filename, "fn", "text.txt", "Filename. Default is 'text.txt'")
	flag.Parse()

	properties, err := resources.ReadPropertiesFile("conf/application.properties")
	if err != nil {
		fmt.Println(err)
	} else {
		if command == "send" {
			if properties["ifcert"] == "true" {
				sendFile(resources.CertClient(properties["certpath"], properties["certpass"]), properties, packname, packtag, filename)
			} else {
				sendFile(&http.Client{}, properties, packname, packtag, filename)
			}
		} else if command == "delete" {
			if properties["ifcert"] == "true" {
				deleteFile(resources.CertClient(properties["certpath"], properties["certpass"]), properties, packname, packtag, filename)
			} else {
				deleteFile(&http.Client{}, properties, packname, packtag, filename)
			}
		} else {
			if properties["ifcert"] == "true" {
				getFile(resources.CertClient(properties["certpath"], properties["certpass"]), properties, packname, packtag, filename)
			} else {
				getFile(&http.Client{}, properties, packname, packtag, filename)
			}
		}
	}
}

func getFile(client *http.Client, properties map[string]string, packname string, packtag string, filename string) {

	fmt.Printf("Будет загружен следующий файл: %v:%v:%v\n", packname, packtag, filename)

	err := os.Mkdir("out", 0750)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return
	}

	outFile, err := os.Create("out/" + filename)
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

func sendFile(client *http.Client, properties map[string]string, packname string, packtag string, filename string) {

	if properties["allowduplicate"] == "false" {
		fmt.Println("-----------Пропуск дубликатов отключен, предудущая версия файла будет удалена из гитлаба-----------")
		deleteFile(client, properties, packname, packtag, filename)
	}

	fmt.Printf("Будет отправлен следующий файл: %v:%v:%v\n", packname, packtag, filename)
	sendedFile, err := os.Open("out/" + filename)
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

func deleteFile(client *http.Client, properties map[string]string, packname string, packtag string, filename string) {

	fmt.Printf("Будет удален следующий файл: %v:%v:%v\n", packname, packtag, filename)
	packId := findPackageId(client, properties, packname, packtag, filename)
	if packId == 0 {
		fmt.Println("Пакет с таким именем и версией не найден")
		return
	}
	fileId := findPackageFile(client, properties, filename, packId)
	if fileId == 0 {
		fmt.Println("Файл с таким именем не найден в пакете, удаление отменено")
		return
	}
	removeFile(client, properties, packId, fileId)

}

func findPackageId(client *http.Client, properties map[string]string, packname string, packtag string, filename string) int {

	var packId int
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

	//fmt.Println(pkgList)

	for _, v := range pkgList {
		if v.Name == packname && v.Version == packtag {
			packId = int(v.Id)
			fmt.Println("Ид найденного пакета:", packId)
			break
		}
	}
	return packId
}

func findPackageFile(client *http.Client, properties map[string]string, filename string, packId int) int {

	var fileId int
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

	//fmt.Println(pkgFiles)

	for _, v := range pkgFiles {
		if v.File_name == filename {
			fileId = int(v.Id)
			fmt.Println("Ид найденного файла:", fileId)
			break
		}
	}
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
