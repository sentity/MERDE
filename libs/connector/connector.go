package connector

import (
	"bytes"
	//"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CONN_HOST = ""
	CONN_PORT = "1337"
	CONN_TYPE = "tcp"
)

type Intercom struct {
	In  chan string
	Out chan string
}

type Query struct {
	Type string
}

type Connection struct {
	ID       int
	Address  string
	Changed  int
	Intercom Intercom
}

type Response struct {
}

var PackageSize = 50005 // length 3 checksum + content
var Connections = make(map[int]Connection)
var Timeout = 300
var MaxID = 0

// --------------------------------------------------------------------------------------------------------------------
// # # # # # # # # # # # # # # # # # # # # # # # # # # #  NEW  # # # # # # # # # # # # # # # # # # # # # # # # # # # #
// --------------------------------------------------------------------------------------------------------------------

func Listen() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// create a intercom to retrieve
		Intercom := Intercom{
			In:  make(chan string),
			Out: make(chan string),
		}
		// start connection handler
		go handleConnection(conn, Intercom)
		// create a connection struct
		// to store some metainfo
		// and the intercom
		addr := l.Addr()
		tmpConnection := Connection{
			ID:       MaxID,
			Address:  addr.String(),
			Intercom: Intercom,
		}
		// store the connection meta info
		Connections[MaxID] = tmpConnection
		MaxID++
	}
}

func handleConnection(conn net.Conn, inter Intercom) {
	// Prepare input buffer &  lastChanged timeStamp & authed flag
	bufferIn := ""
	lastChanged := int(time.Now().Unix())
	authed := false
	// go into loop to check for
	// new contents
	for {
		// get current timeStamp
		currTime := int(time.Now().Unix())
		// lets read out data from buffer
		buf := make([]byte, 32768) // ### check for buffer dynamic size based on package max size
		_, err := conn.Read(buf)
		if err != nil {
			handlePackageError(conn, inter) ///# FIX PARAMS
		}
		// now we Trim 0bytes from buffer
		buf = bytes.Trim(buf, "\x00")
		//fmt.Println(len(string(buf)))
		// - -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
		// get the length of the curr buffer
		bufLen := len(string(buf))
		// - -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
		// if theres more than 0 chars in the buffer we gonne check it
		// and reset package string
		strPackage := ""
		if bufLen > 0 {
			strResult := string(buf) // ### check if stringcast can be removed since it has been done before not sure if buff gets modified or just the return stringacastete...
			fmt.Println("data read:", strResult)
			// adding the new data into our bufferIn
			bufferIn += strResult
			// if buffer in is longer than the controle number length we gonne check the buffer
			// if it may contain a complete package to handle
			if len(bufferIn) > 5 {
				bufferIn, strPackage, err = checkBuffer(bufferIn)
				if nil != err {
					// #### handle buffer error
				}
			}
			// update the lastChanged flag
			lastChanged = currTime
		}
		// - -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
		// we check if there is a package
		if "" != strPackage {
			// dispatch if this package goes into auth or
			// normal package handling
			if false == authed {
				check, err := handleAuth(strPackage, conn)
				if err != nil {
					errors.New("Auth failed")
				} else {
					authed = check
				}
			} else {
				handlePackage(strPackage, conn)
			}
		}
		// - -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
		// check if a connections is longer inactive
		// than the given timeout
		if lastChanged+Timeout < currTime {
			fmt.Println("Closing connection")
			return
		}
		time.Sleep(5000) // 0,000005 seconds to not burn the CPU
	}
}

func sendResponse(conn net.Conn, message Response) {
	//strPackage := buildPackage(message)
	//conn.Write([]byte(strPackage))
	// #### maybe return something here xDD
}

func buildPackage(message Response) string {
	//message = json.Marshal(message)
	// ## add check for TO BIG responses
	var strPackage strings.Builder
	messageLength := len(message)
	packageControle := PackageSize - messageLength
	strPackage.WriteString(string(packageControle))
	strPackage.WriteString(message)
	return strPackage.String()
}

func checkBuffer(buffer string) (string, string, error) {
	// - -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
	// get checklength and convert it to int
	var strPackage string
	checkLength := buffer[0:5]
	intLength, err := strconv.Atoi(checkLength)
	if err != nil {
		return "", "", errors.New("Package error. Buffer clear")
	}
	// - -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
	// checkLength needs to be converted based on our checking logic
	intLength = PackageSize - intLength
	// and get the bufferLength
	bufferLength := len(buffer)
	fmt.Println("intLength: ", intLength, " | bufferLenth: ", bufferLength)
	// - -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
	// if the package is complete recieved
	if bufferLength == 5+intLength {
		// lets retrieve the package contents
		strPackage = buffer[5:]

		// if there is more
		buffer = ""
	}
	// to review later
	//if bufferLength + 3 >= checkLength {
	//    // we got the full package length
	//    // lets retrieve the package
	//    strPackage  = conn.BufferIn[3:checkLength]
	//    // if there is more
	//    if bufferLength + 3 > checkLength {
	//        conn.BufferIn = conn.BufferIn[ckeckLength + 3:]
	//    }
	//}
	return buffer, strPackage, nil
}

func handleAuth(data string, conn net.Conn) (bool, error) {
	fmt.Println("\n(Auth)Recieved data from API:\n", data)
	return true, nil
}

func handlePackage(data string, conn net.Conn) {
	fmt.Println("\n(Package)Recieved data from API:\n", data)
}

func handlePackageError(conn net.Conn, inter Intercom) {

}
