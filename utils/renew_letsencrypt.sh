#!/bin/bash

web_service='nginx'
le_path='/home/ubuntu/letsencrypt'
exp_limit=30;
domain='malice.io'
cert_file="/etc/letsencrypt/live/$domain/fullchain.pem"

exp=$(date -d "`openssl x509 -in $cert_file -text -noout|grep "Not After"|cut -c 25-`" +%s)
datenow=$(date -d "now" +%s)
days_exp=$(echo \( $exp - $datenow \) / 86400 |bc)

echo "Checking expiration date for $domain..."

if [ "$days_exp" -gt "$exp_limit" ] ; then
    echo "The certificate is up to date, no need for renewal ($days_exp days left)."
    exit 0;
else
    echo "The certificate for $domain is about to expire soon. Starting webroot renewal script..."
        $le_path/lletsencrypt-auto certonly -a webroot --agree-tos --renew-by-default --webroot-path=/home/ubuntu/server/public/ -d malice.io -d www.malice.io
    echo "Reloading $web_service"
    /usr/sbin/service $web_service reload
    echo "Renewal process finished for domain $domain"
    exit 0;
fi

# sudo crontab -e
#
# 30 2 * * 1 /home/ubuntu/le-renew.sh >> /var/log/le-renewal.log

# A few notes about the above cron job:
#  - Run every Monday at 2:30 am
#  - Outputs the value of the script at /var/log/le-renewal.log
