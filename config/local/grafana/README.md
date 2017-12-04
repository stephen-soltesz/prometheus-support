# Local Testing

## Installing jsonnet

```
git clone https://github.com/google/jsonnet.git
cd jsonnet/
make
mv jsonnet /usr/local/bin/
```

## Installing grafanalib

Grafanalib depends on python3. On Ubuntu 14.04, I was able to install and run
the generator using:

```
apt-get install -y python3-pip
pip3 install grafanalib
python3 /usr/local/bin/generate-dashboard -o example.json example.dashboard.py
```

## Start grafana with docker

The following mounts local directories for Grafana dashboards so we can use an
upstream Grafana docker image.

Grafana performs aggressive changes on file modes and ownership. As well, it is
extremely conservative with modes for reading files. So, use caution when
locally mounting `/etc/grafana` because Grafana may modify the permissions to be
unsuable after the container exits. Also, be certain that dashboard files and
directory are 644 and 755.

```
cd $HOME/src/github.com/m-lab/prometheus-support/config/local/grafana
```

The first time we run the docker image, docker creates a persistent volume
(transparently saved within the docker root directory, i.e. /var/lib/docker by
default) where the grafana.db and other settings are saved across runs. The
first start up will take about two minutes. Just wait for the server to start
and then Ctrl-C to exit.

```
docker run -it -p 3000:3000 -v /var/lib/grafana --name grafana-storage \
    -e "GF_SERVER_ROOT_URL=http://localhost:3000" \
    -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
    grafana/grafana:latest
```

Now that the grafana storage volume is initialized, we can restart the server
using the persistent volume and mounting the local dashboards directory for
experimentation.

```
docker run -it -p 3000:3000 \
    -v `pwd`/dashboards:/grafana/dashboards \
    -e "GF_DASHBOARDS_JSON_ENABLED=true" \
    -e "GF_DASHBOARDS_JSON_PATH=/grafana/dashboards" \
    -e "GF_SERVER_ROOT_URL=http://localhost:3000"  \
    -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
    --volumes-from grafana-storage \
    grafana/grafana:latest
```

NOTE & TODO: Soon it will be possible to specify datasources in a static
configuration. Until then you must manually add datasources through the grafana
web ui. The persistent volume will save that setting across restarts.

NOTE: permissions on the dashboards files and directory must allow world
reading. Also, if any file or subdirectory cannot be read, grafana fails to load
everything. Ask me how I learned this. (╯°□°）╯︵┻━┻

Next, visit http://localhost:3000 and login with the username/password
"admin"/"secret" (without quotes).
