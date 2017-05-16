Changelog
---------

-	Moving away from `blacktop/elastic-stack` in favor of one container per service
-	Malice will now wait 20 secs for `blacktop/elasticsearch` to start before giving up
-	Malice will check if elasticsearch fails to start if it is because you don't have enough memory to run it
-	Added initial ability to upgrade `~/.malice/config/config.toml` when a new version comes out

### Fixes

-	[x] fix plugin communication to ES
-	[x] wait for ES to fully start (not just a dumb 10 sec wait)
-	[x] monitor logs to output important info as to why it might not start (not enough RAM etc)
-	[x] add `~/.malice/logs/elastic.log` to catch ES errors
-	[x] deactivate `malice/javascript`

Plugin list
-----------

```
nsrl  
shadow-server  
fileinfo  
yara  
avast  
avg  
bitdefender  
clamav  
comodo  
fprot  
fsecure  
sophos  
floss
```
