#!/bin/bash

docker run --rm --name=gitlab-package --net=host \
    -v '/home/wilyr/gitlab-package/cert:/app/cert' \
    -v '/home/wilyr/gitlab-package/conf:/app/conf' \
    -v '/home/wilyr/gitlab-package/out:/app/out' \
    --env C=send --env PN=test_package --env PT=0.0.4 --env FN=total.txt gitlab-package:0.2


