# Known Bugs/ F.A.Q.

## Elasticsearch

## Ways elasticsearch can fail

### 1) they didn't run:

```
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
sudo sysctl -w vm.max_map_count=262144
```

#### Detect

**grep** `docker logs -f malice-elastic` for:

```bash
ERROR: [1] bootstrap checks failed
[1]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
```

### 2) they don't have enough disk space

#### Detect

**grep** `docker logs -f malice-elastic` for:

```bash
[2018-11-24T03:40:25,677][WARN ][o.e.c.r.a.DiskThresholdMonitor] [LMOut61] flood stage disk watermark [95%] exceeded on [LMOut61HROuVZ0A63UEExQ][LMOut61][/usr/share/elasticsearch/data/nodes/0] free: 189.5mb[1.9%], all indices on this node will be marked read-only
```

#### Fix

Start elasticsearch before scanning with malice like so:

```bash
$ docker run --init -d \
                      --name malice-elastic\
                      -p 9200:9200 \
                      -e cluster.routing.allocation.disk.threshold_enabled=false \
                      malice/elasticsearch:6.5
```

=OR=

Increase the disk size of your VM :wink:

### 3) they don't have enough RAM

#### Detect

If you have a lot of scan plugins that are failing to write their results to `elasticsearch`

```
time="2018-11-24T14:44:52Z" level=fatal msg="exit status 150" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=avg

time="2018-11-24T14:45:22Z" level=fatal msg="failed to initalize elasticsearch: failed to connect to database: failed to ping elasticsearch: Get http://elasticsearch:9200/: EOF" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=fprot
time="2018-11-24T14:45:22Z" level=fatal msg="failed to initalize elasticsearch: failed to connect to database: failed to ping elasticsearch: Get http://elasticsearch:9200/: EOF" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=avast
time="2018-11-24T14:45:50Z" level=fatal msg="failed to initalize elasticsearch: failed to connect to database: failed to ping elasticsearch: Get http://elasticsearch:9200/: dial tcp 172.17.0.2:9200: connect: no route to host" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=mcafee
time="2018-11-24T14:45:51Z" level=fatal msg="failed to initalize elasticsearch: failed to connect to database: failed to ping elasticsearch: Get http://elasticsearch:9200/: dial tcp 172.17.0.2:9200: connect: no route to host" category=av path= plugin=comodo
time="2018-11-24T14:46:03Z" level=fatal msg="signal: killed" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=bitdefender
time="2018-11-24T14:46:05Z" level=fatal msg="signal: killed" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=escan
time="2018-11-24T14:46:06Z" level=fatal msg="signal: killed" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=clamav
time="2018-11-24T14:46:20Z" level=fatal msg="context deadline exceeded" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=drweb
time="2018-11-24T14:46:22Z" level=fatal msg="failed to initalize elasticsearch: failed to connect to database: failed to ping elasticsearch: Get http://elasticsearch:9200/: dial tcp 172.17.0.2:9200: connect: no route to host" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=fsecure
time="2018-11-24T14:46:24Z" level=fatal msg="failed to initalize elasticsearch: failed to connect to database: failed to ping elasticsearch: Get http://elasticsearch:9200/: dial tcp 172.17.0.2:9200: connect: no route to host" category=av path=/malware/4daf4edbb04383b93094e89636f303bb11ab687636a4d40813e930c213c3b513 plugin=sophos
```

#### Fix

Start elasticsearch before scanning with malice like so:

```bash
$ docker run --init -d \
                      --name malice-elastic\
                      -p 9200:9200 \
                      -e ES_JAVA_OPTS="-Xms2g -Xmx2g \
                      malice/elasticsearch:6.5
```

=OR=

Increase the VMs RAM 😉
