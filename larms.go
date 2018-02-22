//// v.0.2
//// added parallelism with hardcoded level 9
////

//// v.0.3
//// added additional thread for structured output
//// added handler for wrong or inaccessible nodes

//// v.0.4
//// added handler for wrong password
//// added percentage

//// v.0.4.1
//// fixed minor issues w. threads

package main

import (
	"bufio"
	"fmt"
	"github.com/howeyc/gopass"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"os/user"
	"strconv"
	"strings"
	"io/ioutil"
)

type hopConfig struct {
	host      string
	port      int
	sshConfig *ssh.ClientConfig
}

type hop struct {
	config *hopConfig
	client *ssh.Client
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {

	//outChan := make(chan string)

	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}


func r(hst string, pass string, tamount *int, wg *sync.WaitGroup, tmax *int, outChan chan string, ctr *int, mctr *int) {
		//*ctr++
		//if *ctr>*mctr {*mctr=*ctr}
		
		defer wg.Done() 
                
		for ( *tamount>=*tmax ) {time.Sleep (100)}

                *tamount++

///////////////////
		usr, err := user.Current()
		if err != nil { /*return err*/ fmt.Println("userCurrent err") }

		args := os.Args[1:]
		args[len(args)-2] = hst

		hostPorts := []string{}
		runLocalWithForward := ""
		commandToRun := []string{}
		var nextArgFn func(string) error = nil

		for argI, arg := range args {
			naf := nextArgFn
			nextArgFn = nil
			if naf != nil {
				if err := naf(arg); err != nil {
					/*return err*/ fmt.Println("nextArgFn err")
				}
				continue
			}

			atCount := 0

			for _, c := range arg {

				if c == '@' {
					atCount++
				}
			}

			if atCount > 0 {
				hostPorts = append(hostPorts, arg)
			} else {
				commandToRun = args[argI:]
				break
			}
		}

		hopConfigs := []*hopConfig{}
		for _, hostPort := range hostPorts {

			hc := &hopConfig{
				host: "localhost",
				port: 22,
				sshConfig: &ssh.ClientConfig{
					Config: ssh.Config{},
					User:   usr.Username,
					Auth:/*authMethods,*/ []ssh.AuthMethod{ssh.Password(pass)},
				},
			}

			if splitUserHost := strings.SplitN(hostPort, "@", 2); len(splitUserHost) == 2 {
				hc.sshConfig.User = splitUserHost[0]
				hc.host = splitUserHost[1]
			} else if len(splitUserHost) == 1 {
				hc.host = splitUserHost[0]
			}

			if splitHostPort := strings.SplitN(hc.host, ":", 2); len(splitHostPort) == 2 {
				if n, err := strconv.Atoi(splitHostPort[1]); err != nil {
					/*return err*/ fmt.Println("splitHostPort err")
				} else {
					hc.host = splitHostPort[0]
					hc.port = n
				}
			}
			hopConfigs = append(hopConfigs, hc)
		}

		if len(hopConfigs) == 0 {
			/*return*/ fmt.Errorf("no ssh hops")
		}

		if runLocalWithForward == "" && len(commandToRun) == 0 {
			/*return*/ fmt.Errorf("shell not implemented yet")
		}

		//var tcpConn Conn
		//var err error

		hops := []*hop{}
		dialFunc := net.Dial
		//fmt.Println(dialFunc)
		//toutFunc := net.DialTimeout

		for _, hc := range hopConfigs {
			hostAddr := fmt.Sprintf("%s:%d", hc.host, hc.port)

			tcpConn, err := dialFunc("tcp", hostAddr)

			//tcpConn, err = toutFunc/*dialFunc*/("tcp", hostAddr, time.Second*5)
			if err != nil {
				fmt.Println(hc.host+" | dial error: inaccessible or does not exist or timeout expired")

				continue
			}

			defer tcpConn.Close()

			sshConn, chans, reqs, err := ssh.NewClientConn(tcpConn, hostAddr, hc.sshConfig)
			if err != nil {
				fmt.Println(hc.host+" | Connection is broken or ssh password is wrong")
				os.Exit(1)
			}

			client := ssh.NewClient(sshConn, chans, reqs)
			defer client.Close()

			dialFunc = client.Dial

			hops = append(hops, &hop{config: hc, client: client})
		}

		//fmt.Println(len(hops)-1)
		lastClient := hops[len(hops)-1].client

		if runLocalWithForward == "" && (len(hops)-1)>0 { // >0 for missing inaccessible hosts
			session, err := lastClient.NewSession()
			if err != nil {
				fmt.Println("last client new sess.")
				//return err
			}
			defer session.Close()

			session.Stdin = os.Stdin
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr

			if len(commandToRun) > 0 {
				cmdString := strings.Join(commandToRun, " ")

				oldout := session.Stdout
                                r, w, _ := os.Pipe()
                                session.Stdout = w

				//
                                    session.Run(cmdString)
				//

				w.Close()
				out, _ := ioutil.ReadAll(r)
				session.Stdout = oldout
                                //fmt.Print(string(out))

				outChan <- string(out)
				
				//outChan <- *ctr

				//fmt.Println(*ctr)

				/*return*/ 
				//session.Run(cmdString)

			} else {
				/*return*/ fmt.Errorf("enter the command, bro")
			}
		}

///////////////////
        *tamount--
	*ctr++

	//fmt.Print("ctr|"); fmt.Println(*ctr)
	//fmt.Print("mctr|");fmt.Println(*mctr)
	//fmt.Print("wg Size|"); fmt.Println(wg.Size())

}


func progress(ctr *int, mctr*int, wg *sync.WaitGroup) {
	defer wg.Done()	
	for {
		if *mctr>0 {
		fmt.Print("prgrs|"); fmt.Print((*ctr)*100/(*mctr));fmt.Print/*ln*/("%\r") //\r
		time.Sleep(time.Second * 5)
			if (*ctr == *mctr) { fmt.Println("prgrs|100%"); break }
		}
		
	}
}

func _main() error {

	outChan := make (chan string)

        wg := &sync.WaitGroup{}

        var tmax int
        tmax=9

        var threadsamount int
        threadsamount=0

	var ctr int
	ctr=0
	var mctr int
	mctr=0

        //usr, err := user.Current()
	//if err != nil {
	//	return err
	//}

	fmt.Print("Enter your account password: ")
	var pass string
	p, _ := gopass.GetPasswd()
	pass = string(p)

	hostlist, err := readLines(os.Args[len(os.Args)-2])
	if err != nil { log.Fatalf("read hostlist error: %s", err) }
	
	wg.Add(1)
	go func() {     
		    defer wg.Done() 
		    
		    for ctr<mctr {
                        select {
                        case msg := <- outChan:
                                fmt.Print(msg)
                        //case <- time.After(2*time.Second):
                        //        fmt.Println("timeout")
                        //default:
                        }
                    }
        }()
	
	for _, _ = range hostlist { mctr++}
	for _, hst := range hostlist {

	wg.Add(1)
/////////////////

        go r(hst, pass ,  &threadsamount, wg, &tmax, outChan, &ctr, &mctr)

/////////////////

	}

	wg.Add(1)
	go progress(&ctr,&mctr,wg)

	wg.Wait()
	return nil
}
