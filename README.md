# Tripwire
A nice usefull tool to create iptable-rules which logs all connections to a given port to detect ip-scanner and hacker.

# Install
Run 
```go
go get
go build -o tripwire
```
it was tested with go 1.13 if compiling doesn't work, try using my go version

# Usage
Create a rule to log and <b>allow</b> all connections to port 21 and writes them into /var/log/ftpListener.conf
```#./tripwire -p21 -o ftpListener -a```
<br>
Block incomming connections (instead of accepting them)
```#./tripwire -p21 -o ftpListener```
<br>
Specifies the [loglevel](https://highly.illegal-dark-web-server.xyz/i/qszvm-34l8q-9crda-abi85-b0vhv)
```#./tripwire -p21 -o ftpListener -l5 ```
<br>
Delete log and iptable rules for port 21
```#./tripwire -d -p21```

