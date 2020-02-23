## sshtail
tail remote log by ssh

## build
```
make
```

## run
- Usage
```
./sshtail -h
Usage of ./sshtail:
  -cmd string
    	cmd to collect log (default "tail -q -n +1 -F --max-unchanged-stats=5")
  -debug
    	debug log level
  -file string
    	remote log file name (default "/var/log/vipr/emcvipr-object/dataheadsvc-access.log")
  -passwd string
    	ssh user passwd (default "ChangeMe")
  -server string
    	ssh server ip:port (default "192.168.55.2:22")
  -store string
    	where to store log(fs or redis://127.0.0.1:6379/0) (default "fs")
  -user string
    	ssh user name (default "root")
  -version
    	show version
```
- Start container
```
docker run -d --restart always --network host --name sshtail \
    -e USER=admin \
    -e PASSWD=ChangeMe002 \
    myshare.io:5000/sshtail:1.0.1
```