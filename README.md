# Система управления регистром пакетов в проекте гитлаба
## Версия v0.1.5

```bash
########## Система управления регистром пакетов ##########
#################### Версия: v0.1.5 ######################
##########################################################
Доступные тэги: [-help/-h] [-c] [-pn] [-pt] [-fn] [-url] [-t]
             [-pw] [-f] [-ic] [-cp] [-pi] [-ad] [-dir] [-conf]
ОБЯЗАТЕЛЬНЫЕ параметры: [-c] [-pn] [-pt] [-fn]
Каталоги по умолчанию: 
 'out' - Файловый 
 'conf/application.properties' - Свойства приложения
Разработка: n.kovalev
Source: https://github.com/WILYR/gitlab-generic-packages.git
##########################################################
```

[Работает на основе gitlab generic api](https://docs.gitlab.com/ee/user/packages/generic_packages/)   
[Релизы](https://github.com/WILYR/gitlab-generic-packages/releases)


## Требования:

### Raw
1. golang 1.18  
2. .pem сертификат, если подключение происходит с помощью клиентского сертификата.
3. openssl, если на клиентском сертификате есть пароль.

### Docker
1. docker
2. .pem сертификат, если подключение происходит с помощью клиентского сертификата.

## Установка:

Скачиваем image  
```shell
docker pull ghcr.io/wilyr/gitlab-generic-packages:latest ### 0.1.5/main
```
Можно просто собрать самостоятельно
```shell
git clone https://github.com/WILYR/gitlab-generic-packages.git
./docker_build.sh
```
:arrow_double_up: **Для сборки docker image тоже требуется golang.**
## Запуск

### Raw go

```shell
git clone https://github.com/WILYR/gitlab-generic-packages.git
go run main.go -c get -pn test_package -pt latest -fn nice.csv
```
#### Доступные флаги:  
:exclamation: - Обязательные флаги  
:ok: - Необязательные

- [:ok:] ad string - Поддержка дублирования файлов на уровне запроса. По умолчанию из conf файла  
    - [x] true/false
- [:exclamation:] c string  - Команда серверу. 
    - [x] get/send/delete
- [:ok:] conf string - Путь к conf файлу. (default "conf/application.properties")  
    - [x] /home/data/application.properties
- [:ok:] cp string - Путь к пользовательскому *.pem сертификату. По умолчанию из conf файла  
    - [x] /home/data/cert.pem
- [:ok:] pw string - Пароль от пользовательсткого сертификата. По умолчанию из conf файла 
    - [x] password
- [:ok:] dir string - Директория файлового каталога. По умолчанию ищет 'out' в текущем
    - [x] /home/data/out
- [:exclamation:] fn string - Имя файла в пакете
    - [x] filename
- [:ok:] ic string - Поддержка пользовательского *.pem сертификата. По умолчанию из conf файла
    - [x] true/false
- [:ok:] pi string- ID прокта. По умолчанию из conf файла
    - [x] 10
- [:exclamation:] pn string - Имя пакета
    - [x] packagename
- [:exclamation:] pt string - Версия пакета
    - [x] 0.0.1/latest
- [:ok:] t string - Токен авторизации на сервере gitlab. По умолчанию из conf файла
    - [x] asfqaewrgqlero341effeaADAS
- [:ok:] url string - URL ресурса gitlab. По умолчанию из conf файла
    - [x] https://your.gitlab.mrgeng.ru
- [:ok:] f string - Force delete всех повторений файлов или пакетов. По умолчанию из conf файла  
    - [x] false/true:packs/true:files

### Docker

Запускаем скрипт. Его описание ниже.  
```shell
./docker_run.sh
```
Рабочий каталог приложения - /app.   
Внутри скрипта запуска обязательно монтируем вложенные каталоги с сертификатом, свойствами и файловый каталог.  
В env задаем свойства, аналогичные флагам go приложения. 
Каталог ENV свойств:
```Dockerfile
#ENV'S
ENV COMMAND=""
ENV HELP=""
ENV PACKAGENAME=""
ENV PACKAGETAG=""
ENV FILENAME=""
ENV URL=""
ENV TOKEN=""
ENV IFCERT=""
ENV CERTPATH=""
ENV CERTPASS=""
ENV PROJECTID=""
ENV ISALLOWDUPLICATE=""
ENV FILEDIR=""
ENV CONFFILE=""
ENV FORCE=""
``` 
После запуска будет выведен лог приложения и контейнер автоматически удалится. Файлы останутся во вмонтированном каталоге.

```shell
#!/bin/bash

#----ALL ENVS-----

#COMMAND=get/send/delete
#HELP=true
#PACKAGENAME=name
#PACKAGETAG=0.0.1
#FILENAME=name
#URL=https://test.gitlab.ru
#TOKEN=asdgawegabWSWE--342S
#IFCERT=true/false
#CERTPATH=/app/cert/test.pem
#CERTPASS=password
#PROJECTID=10
#ISALLOWDUPLICATE=true/false
#FILEDIR=/app/testdir
#CONFFILE=/app/test.conf/application.properties
#FORCE=true:packs/true:files/false

docker run --rm --name=gitlab-package --net=host \
    -v '/home/wilyr/gitlab-package/cert:/app/cert' \
    -v '/home/wilyr/gitlab-package/conf:/app/conf' \
    -v '/home/wilyr/gitlab-package/out:/app/out' \
        --env COMMAND=get \
        --env PACKAGENAME=test_package \
        --env PACKAGETAG=0.0.4 \
        --env FILENAME=total.txt \
        --env ISALLOWDUPLICATE=false \
        --env PROJECTID=10 \
gitlab-generic-packages:0.1.5

```
## Конфигурация

Все свойства подтягиваются из файла conf/application.properties  
:exclamation: **Флаги(если заданы) переопределяют значения свойств!**

```properties
# Адрес ресурса
# Можно переопределить флагом -url
resource=https://your.gitlab.ru

# Токен авторизации
# Можно переопределить флагом -t
token = glglalswflweggaw323

# Есть ли авторизация по ssl сертификату
# Можно переопределить флагом -ic
ifcert=true

# Путь к сертификату. ПРИ СТАРТЕ В КОНТЕЙНЕРЕ УКАЗЫВАЕТСЯ ПУТЬ ВНУТРИ КОНТЕЙНЕРА, А НЕ НА ЗАПУСКАЕМОЙ МАШИНЕ
# Можно переопределить флагом -cp
certpath=/home/wilyr/ssl/temp.pem

# Пароль от клиентского сертификата. Если его нет, то оставляем поле пустым
# Можно переопределить флагом -pw
certpass=asbswdreb

# ID Проекта
# Можно переопределить флагом -pi
projectid=10

# Разрешено ли дублирование файлов в пакете. Если false, то старый файл будет автоматически удален, а новый записан. Если true, то будет загружена новая версия файла без удаления старой 
# Можно переопределить флагом -ad
allowduplicate=false

# Удаление всех вхождений файлов в проекте. 
# Если true:files, то удаляются все файлы во всех пакетах с заданными параметрами, пакеты остаются 
# Если true:packs - полностью удалит все пакеты со всеми файлами.
# Можно переопределить флагом -f 
force=false
```
Токен авторизации можно сгенерировать в настройках профиля в gitlab  
ID Проекта можно получить на главной странице проекта gitlab.  
