#CGO_ENABLED=0 go build -o main.sh main.go
#BASE IMAGE

FROM alpine:3.15.2

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


#Workdir
WORKDIR /app

#Add packages
RUN apk add openssl

#Create volume dir's
RUN mkdir cert out

#COPY PROJ
COPY . /app
RUN chmod 740 main.sh

CMD if [ "$COMMAND" = "" ]; then \
        break; \
    else \
        export COMMAND="-c $COMMAND"; \
    fi; \ 

    if [ "$HELP" = "" ]; then \
        break; \
    else \
        export HELP="-h $HELP"; \
    fi; \

    if [ "$PACKAGENAME" = "" ]; then \
        break; \
    else \
        export PACKAGENAME="-pn $PACKAGENAME"; \
    fi; \

    if [ "$PACKAGETAG" = "" ]; then \
        break; \
    else \
        export PACKAGETAG="-pt $PACKAGETAG"; \
    fi; \

    if [ "$FILENAME" = "" ]; then \
        break; \
    else \
        export FILENAME="-fn $FILENAME"; \
    fi; \

    if [ "$URL" = "" ]; then \
        break; \
    else \
        export URL="-url $URL"; \
    fi; \

    if [ "$TOKEN" = "" ]; then \
        break; \
    else \
        export TOKEN="-t $TOKEN"; \
    fi; \

    if [ "$IFCERT" = "" ]; then \
        break; \
    else \
        export IFCERT="-ic $IFCERT"; \
    fi; \

    if [ "$CERTPATH" = "" ]; then \
        break; \
    else \
        export CERTPATH="-cp $CERTPATH"; \
    fi; \

    if [ "$CERTPASS" = "" ]; then \
        break; \
    else \
        export CERTPASS="-pw $CERTPASS"; \
    fi; \

    if [ "$PROJECTID" = "" ]; then \
        break; \
    else \
        export PROJECTID="-pi $PROJECTID"; \
    fi; \

    if [ "$ISALLOWDUPLICATE" = "" ]; then \
        break; \
    else \
        export ISALLOWDUPLICATE="-ad $ISALLOWDUPLICATE"; \
    fi; \

    if [ "$FILEDIR" = "" ]; then \
        break; \
    else \
        export FILEDIR="-dir $FILEDIR"; \
    fi; \

    if [ "$CONFFILE" = "" ]; then \
        break; \
    else \
        export CONFFILE="-conf $CONFFILE"; \
    fi; \

    if [ "$FORCE" = "" ]; then \
        break; \
    else \
        export FORCE="-f $FORCE"; \
    fi; \
    
    ./main.sh $COMMAND $HELP $PACKAGENAME $PACKAGETAG $FILENAME $URL $TOKEN $IFCERT $CERTPATH $CERTPASS $PROJECTID $ISALLOWDUPLICATE $FILEDIR $CONFFILE $FORCE