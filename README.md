# Installation
Download the binary for your system from the releases

# Building from source
You'll need my [utils packages](https://github.com/UlisseMini/utils) then a simple `go build .` or `go install` should work :D

# Bugs
* Using the -e option holds up the socket even after the process exits (need remote end to send data)

# TODO
* Find a better args parser (you can't mix args with flags with the flag package)
* Allow for giving port and ip in other format (besides ip:port)
