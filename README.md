# spiderloops

Add web spiders in loops with some data.
Just for fun.

```
sudo mkdir /opt/apps
sudo chown -T username:username /opt/apps
cd /opt/apps
git pull https://github.com/sbroekhoven/spiderloops
cd spiderloops
go build .
cp systemd-service.txt /etc/systemd/system/spiderloops.service # change the username
sudo systemctl start spiderloops.service
sudo systemctl status spiderloops.service
```