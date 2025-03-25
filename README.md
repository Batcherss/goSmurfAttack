# Use for educational purposes only!
This program is a stress tester
use at your own risk
The author is not responsible, the program uses the MIT license.

# Main

![goSmurfAttack'er](https://img.shields.io/badge/Status-Active-green)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)

## Description

This is a smurf attack program developed in Go. It is designed for attacking network load by sending data packets over ICMP(type:Echo). The program uses several libraries to process packets, interact with the network, and generate random data.



### Important Note
For the program to work correctly, you need to install **Npcap**. Without it, the program may not work properly.

You can download **Npcap** from the [official Npcap website](https://nmap.org/npcap/).

## Dependencies

Before running the program, make sure you have installed the following dependencies:

- **`golang.org/x/sys/windows`** — for working with the Windows API.
- Go standard libraries:
  - `bufio`
  - `bytes`
  - `encoding/binary`
  - `fmt`
  - `log`
  - `math/rand`
  - `net`
  - `os`
  - `strconv`
  - `strings`
  - `time`

To install the `golang.org/x/sys/windows` dependency, use the following command:

```cmd
go get golang.org/x/sys/windows
```
Also, make sure you have Npcap installed. You can download and install it from the official Npcap website.

How to Run
Clone the repository:

```git
git clone https://github.com/yourusername/stress-test-program.git
```

Enter the folder:
```cmd
cd stress-test-program
```

Install the necessary dependencies (if not already installed):
```cmd
go mod tidy
```

Run the program:
```cmd
go run main.go
```

Notes
Make sure you have administrator rights to work with network packets.

The program uses Npcap, so be sure to install it before running the program, otherwise, it may not work correctly.

