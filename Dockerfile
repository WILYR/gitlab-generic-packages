#BASE IMAGE
FROM golang:1.18-alpine

#ENV'S
# COMMAND TYPE
ENV C=""
# FILENAME
ENV FN=""
#PACKAGE TAG
ENV PT=""
#PACKAGE NAME
ENV PN=""

#Workdir
WORKDIR /app

#Add packages
RUN apk add openssl

#Create volume dir's
RUN mkdir cert out

#COPY PROJ
COPY . /app

# RUN PROJ
CMD go run package.go -c $C -pn $PN -pt $PT -fn $FN
