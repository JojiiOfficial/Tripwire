# Tripwire
A nice usefull tool to create iptable-rules which logs all connections to a given port to detect ip-scanner and hacker.
You can use it in combination with the [ScanBanServer](https://github.com/JojiiOfficial/ScanBanServer) and [Triplink](https://github.com/JojiiOfficial/triplink) to create a network to collect and block internet scanner

# Install
Run 
```go
go get
go build -o tripwire
```
it was tested with go 1.13. If compiling doesn't work, try using go1.13

# Usage
Show <b>help</b>
<br>```#./tripwire -h```
<br><br>Create a rule to log and <b>allow</b> all connections to port 21 and write them into /var/log/ftpListener.conf
<br>```#./tripwire add -p21 -o ftpListener -a```
<br><br>
...<b>Block</b> incomming connections (instead of accepting them)
<br>```#./tripwire add -p21 -o ftpListener```
<br><br>
...Specifies the <b>[loglevel](https://highly.illegal-dark-web-server.xyz/i/qszvm-34l8q-9crda-abi85-b0vhv)</b>
<br>```#./tripwire add -p21 -o ftpListener -l5 ```
<br><br>
<b>Delete</b> log and iptable rules for port 21
<br>```#./tripwire delete -p21 -o ftpListener```
<br><br>
<b>List</b> all tripwire configurations
<br>```#./tripwire list```
