FROM logstash

COPY rethinkdb.conf /etc/logstash/conf.d/rethinkdb.conf
# COPY elasticsearch-output.conf /etc/logstash/conf.d/elasticsearch-output.conf

RUN logstash-plugin install logstash-input-rethinkdb

CMD ["-f", "/etc/logstash/conf.d/"]
