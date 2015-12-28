FROM ubuntu:14.04
MAINTAINER blacktop, https://github.com/blacktop

RUN apt-get update \
  && apt-get install -y unzip

WORKDIR /malice

# Copy over files
ADD https://github.com/blacktop/go-malice/archive/master.zip /
RUN unzip /master.zip && rm /master.zip
ADD config_docker.json /
ADD docker-entry.sh /

RUN chmod +x /docker-entry.sh
ENTRYPOINT /docker-entry.sh

VOLUME ["/malice/data"]

EXPOSE 80 443 9200 9300

# CMD ["/usr/bin/supervisord"]
