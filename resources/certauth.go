package resources

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func CertClient(filepath string, certpassword string) *http.Client {
	certFile := flag.String("certfile", filepath, "trusted CA certificate")
	flag.Parse()

	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", *certFile)
	}

	var clientCert tls.Certificate
	decryptKey := "temp.pem"

	if len(certpassword) != 0 {
		//Снимаю пароль с private key и формирую новый .pem файл.

		cmd := exec.Command("openssl", "rsa", "-in", filepath, "-out", decryptKey, "-passin", "pass:"+certpassword)
		cmd.Dir = "out/"

		output, err := cmd.CombinedOutput()
		log.Println(string(output))

		clientCert, err = tls.LoadX509KeyPair(filepath, "out/"+decryptKey)
		if err != nil {
			log.Fatal(err)
		}

		//Удаляю временные файлы
		log.Println("Удаление временных файлов.....")

		temp, err := os.Open("out/" + decryptKey)
		if err != nil {
			log.Println(err)
		}
		defer temp.Close()

		os.Remove(temp.Name())

	} else {
		clientCert, err = tls.LoadX509KeyPair(filepath, filepath)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Создание http клиента
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{clientCert},
				RootCAs:            certPool,
				InsecureSkipVerify: true,
			},
		},
	}

	return client

}
