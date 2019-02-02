# ngo
...is an ncat implementation written in pure go!

# Usage & examples
spawn bash connecting stdin stdout and stderr to 127.0.0.1
* ngo -e "bash -li" 127.0.0.1:1337

Transfering files
* ngo localhost:1337 < file.txt
<br>Other side (listen)
* ngo -l :1337 > file.txt
<br>

For the full help run
* ngo -h

# Installation
Download the binary for your system from the releases

# Building from source
First `go get github.com/UlisseMini/ngo` then cd into the directory and run `make`
<br>then the binaries will be installed into ngo/bin, to install into your $GOBIN do `make depends` then `go install`

# Bugs
* Using the -e option holds up the socket even after the process exits (non pty)
* failure to resize pty (wants \*os.File)
* timeout option not supported with tls
* strange data race with the tests

# TODO
* Allow for giving port and ip in other format (besides ip:port)
* ssl options --ssl-key --ssl-cert --ssl etc (maybe create cert auth for ngo?)
* pty support for windows
* Better argument parsing
