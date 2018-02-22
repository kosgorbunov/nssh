larms is nssh narrowed fork and able through jumphost or jumphosts chain go to nodes and do same commands there, similar as ansible do.

2do:
- goroutines for parallelism
- run commands from file to avoid quotas

bugs:
- no exceptions on timeouts (use ctrl+c)
 
samples:  
export GOPATH=~/golang
go get github.com/howeyc/gopass
go get golang.org/x/crypto/ssh
 
disk free root
------------------
time go run ./longarms.go jump_user@jump_host  /.../<list of hosts> " echo -n \`hostname --fqdn\` ; echo -n \" \"; echo \`df -h /\` | awk '{print \$12}' "
 
jump_user@jump_host as well could be jump_user1@jump_host1 jump_user2@jump_host2 jump_user3@jump_host3 to build ssh jumps chain

<list of hosts> file must be structured like in below:

userid@h1
userid@h2
userid@h3
