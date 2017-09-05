#!/usr/bin/env bash
set -e

# Install packages
sudo apt-get update -y
sudo apt-get install -y curl unzip

# Download Malice into some temporary directory
curl -L "${download_url}" > /tmp/malice.zip

# Unzip it
cd /tmp
sudo unzip malice.zip
sudo mv malice /usr/local/bin
sudo chmod 0755 /usr/local/bin/malice
sudo chown root:root /usr/local/bin/malice

# Setup the configuration
cat <<EOF >/tmp/malice-config
${config}
EOF
sudo mv /tmp/malice-config /usr/local/etc/malice-config.json

# Setup the init script
cat <<EOF >/tmp/upstart
description "Malice server"

start on runlevel [2345]
stop on runlevel [!2345]

respawn

script
  if [ -f "/etc/service/malice" ]; then
    . /etc/service/malice
  fi

  # Make sure to use all our CPUs, because Malice can block a scheduler thread
  export GOMAXPROCS=`nproc`

  exec /usr/local/bin/malice server \
    -config="/usr/local/etc/malice-config.json" \
    \$${MALICE_FLAGS} \
    >>/var/log/malice.log 2>&1
end script
EOF
sudo mv /tmp/upstart /etc/init/malice.conf

# Extra install steps (if any)
${extra-install}

# Start Malice
sudo start malice
