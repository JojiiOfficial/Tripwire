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
<b>Help</b>
<br>```#./tripwire -h```
Create a rule to log and <b>allow</b> all connections to port 21 and writes them into /var/log/ftpListener.conf
<br>```#./tripwire add -p21 -o ftpListener -a```
<br><br>
<b>Block</b> incomming connections (instead of accepting them)
<br>```#./tripwire add -p21 -o ftpListener```
<br><br>
Specifies the <b>[loglevel](https://highly.illegal-dark-web-server.xyz/i/qszvm-34l8q-9crda-abi85-b0vhv)</b>
<br>```#./tripwire add -p21 -o ftpListener -l5 ```
<br><br>
<b>Delete</b> log and iptable rules for port 21
<br>```#./tripwire delete -p21 -o ftpListener```
<br><br>
<b>Listing</b> all tripwires
<br>```#./tripwire list```
