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
