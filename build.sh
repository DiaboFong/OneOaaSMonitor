CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build
mkdir -p pkg/oneoaas_web/conf
mkdir -p pkg/oneoaas_web/download
mkdir -p pkg/oneoaas_web/db

DOWNLOAD_DIR="pkg/oneoaas_web/download"
#wget -P ${DOWNLOAD_DIR} http://repo.dev.oneoaas.com/monitor/ee/oneoaas-monitor-ee-linux-amd64-2.0.0-1.x86_64.rpm
wget -P ${DOWNLOAD_DIR} http://repo.dev.oneoaas.com/monitor/ee/oneoaas-monitor-ce-linux-amd64-2.0.0-1.x86_64.rpm
wget -P ${DOWNLOAD_DIR} http://repo.dev.oneoaas.com/monitor/ee/oneoaas-monitor-se-linux-amd64-2.0.0-2.x86_64.rpm
cp -r static pkg/oneoaas_web/static
cp -r views pkg/oneoaas_web/views
cp  oneoaas_web pkg/oneoaas_web
cat > pkg/oneoaas_web/start.sh << EOF
killall -9 oneoaas_web
nohup ./oneoaas_web &
EOF
cat > pkg/oneoaas_web/conf/app.conf << EOF
appname = oneoaas_web
httpport = 4007
runmode = prod
dbtype = sqlite3
sms_user = OneOaaS
sms_key = SqaHdxKJLDKaKstQCRbePVA9nbOx4oLL
sms_tid = 7130

#钉钉webhook
ding_webhook = "https://oapi.dingtalk.com/robot/send?access_token=246550e03149b2df949cc61e5a2847241b2d9d90d674b63c9b6cd775a3dd7eb4"

#database
dbtype = "mysql"
dbuser = "oneoaas_web"
dbpass = "oneoaas_web"
dbhost = "127.0.0.1"
dbport = 3306
dbname = "oneoaas_web"

#email
username_email="monitor.apply@oneoaas.com"
password_email="Jci9^Y~e45f=KBPy#rwH?<"
host_email="smtp.exmail.qq.com"
port_email=465

#site_domain
oneoaas_domain="http://www.oneoaas.com"
EOF

cd pkg
tar zcvf oneoaas_web.tar.gz oneoaas_web
