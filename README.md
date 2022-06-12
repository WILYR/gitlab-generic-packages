# Система управления регистром пакетов в проекте гитлаба

[Работает на основе gitlab generic api](https://docs.gitlab.com/ee/user/packages/generic_packages/)   
## Требования:

### Raw
1. golang 1.18  
2. .pem сертификат, если подключение происходит с помощью клиентского сертификата.
3. openssl, если на клиентском сертификате есть пароль.

### Docker
1. docker
2. .pem сертификат, если подключение происходит с помощью клиентского сертификата.

## Установка:

```shell 
git clone https://github.com/WILYR/gitlab-generic-packages.git
```
Если хотим в докере, то можно просто скачать последний docker image из репозитория или собрать самостоятельно  
```shell
docker pull ghcr.io/wilyr/gitlab-package:0.1.1
## Или
docker build --label=gitlab-packages -t gitlab-package:0.1.1 .
```

## Запуск

### Raw go

```shell
go run package.go -c get -pn test_package -pt latest -fn nice.csv
```
Доступные флаги:  

- '-c' - Команда серверу (Доступны: get/send/delete). Получить, отправить или удалить файл соответственно. По умолчанию: get
- '-pn'- Имя пакета. По умолчанию: test_package
- '-pt'- Tag пакета. По умолчанию: latest
- '-fn'- Имя файла. По умолчанию text.txt

🔥_**Файловый каталог - 'out/'. Если выполнять команду get, он создается автоматически и туда записыватся файл. Если выполнять команду send, то предварительно необходимо создать этот каталог в корне и скопировать туда нужный для отправки файл. На данный момент поиск будет осуществляться только в этом каталоге. ИМЯ ФАЙЛА ВСЕГДА УКАЗЫВАТЬ БЕЗ FILEPATH, НАПРИМЕР '-fn total.txt'. А ВОТ ТАК НЕЛЬЗЯ ~~'-fn out/total.txt'~~**_🔥  

### Docker

Запускаем скрипт. Его описание ниже.  
```shell
./run.sh
```
Рабочий каталог приложения - /app.   
Внутри скрипта запуска обязательно монтируем вложенные каталоги с серфификатом, свойствами и файловый каталог 'out/'.  
В env задаем свойства, аналогичные флагам go приложения.  
После запуска будет выведен лог приложения и контейнер автоматически удалится. Файлы останутся во вмонтированном каталоге.

```shell
#!/bin/bash

docker run --rm --name=gitlab-package --net=host \
    -v '/home/wilyr/gitlab-package/cert:/app/cert' \
    -v '/home/wilyr/gitlab-package/conf:/app/conf' \
    -v '/home/wilyr/gitlab-package/out:/app/out' \
    --env C=send --env PN=test_package --env PT=latest --env FN=total.txt gitlab-package:0.2
```
## Конфигурация

Все свойства подтягиваются из файла conf/application.properties

```properties
# Адрес ресурса.
resource=https://your.gitlab.ru

#Токен авторизации
token = glAGSEhrheat-xHsWJmgwegwegesdf

#Есть ли авторизация по ssl сертификату
ifcert=true

#Путь к сертификату. ПРИ СТАРТЕ В КОНТЕЙНЕРЕ УКАЗЫВАЕТСЯ ПУТЬ ВНУТРИ КОНТЕЙНЕРА, А НЕ НА ЗАПУСКАЕМОЙ МАШИНЕ
certpath=/home/wilyr/ssl/temp.pem

#Пароль от клиентского сертификата. Если его нет, то оставляем поле пустым
certpass=EMasdgewFWE

#ID Проекта
projectid=10

#Разрешено ли дублирование файлов в пакете. Если false, то старый файл будет автоматически удален, а новый записан. Если true, то будет загружена новая версия файла без удаления старой 
allowduplicate=false
```
Токен авторизации можно сгенерировать в [настройках профиля](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)  
ID Проекта можно получить на главной странице вашего проекта.   


