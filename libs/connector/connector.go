package connector

import (
	"fmt"
	"net"
	"os"
	"time"
	//"io/ioutil"
	"bytes"
	"errors"
	"strconv"
	"sync"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "1337"
	CONN_TYPE = "tcp"
)

type Connection struct {
	ID         int
	Address    string
	Connection net.Conn
	Timeout    int
	Changed    int
	BufferIn   string
	//BufferOut   string
}

var PackageSize = 403 // length 3 checksum + content
var Connections = make(map[int]Connection)
var ConnectionsMutex = &sync.RWMutex{}
var ConnectionIDs = make(map[int]int)
var Active = true
var MaxID = 0
var CompareBytesSmall = make([]byte, 256)

//var CompareBytesBig   = make([]byte, 0, 8 * PackageSize )

func Listen() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// starting the connection input handler
	// as a seperate thread
	go handleInput()
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
		// create a connection struct
		currTime := int(time.Now().Unix())
		addr := l.Addr()
		tmpConnection := Connection{
			ID:         MaxID,
			Address:    addr.String(),
			Connection: conn,
			Timeout:    5,
			Changed:    currTime,
			BufferIn:   "",
		}
		ConnectionsMutex.Lock()
		Connections[MaxID] = tmpConnection
		ConnectionIDs[MaxID] = 1
		MaxID++
		ConnectionsMutex.Unlock()
	}
}

// Handles incoming requests.
func handleInput() {
	for Active == true {
		if len(Connections) > 0 {
			for id, _ := range ConnectionIDs {
				ConnectionsMutex.Lock()
				conn := Connections[id]
				ConnectionsMutex.Unlock()
				buf := make([]byte, 256) // using small tmo buffer for demonstrating
				_, err := conn.Connection.Read(buf)
				if err != nil {
					handlePackageError(conn)
				}
				//if bytes.Equal(tmp, CompareBytesSmall)  {
				//    break
				//} else {
				//fmt.Println("got ", len(tmp), "bytes: ",tmp)
				//}
				//if len(tmp) > 0 {
				//    buf = append(buf, tmp[:n]...)
				//} else {
				//    break
				//}
				//result, _ := ioutil.ReadAll(conn.Connection)
				//strResult := string(result)
				currTime := int(time.Now().Unix())
				// if there is input we update the connection object
				buf = bytes.Trim(buf, "\x00")
				//bufLen    := len(buf)
				if !bytes.Equal(buf, CompareBytesSmall) {
					strResult := string(buf)
					//fmt.Println("data read:",strResult)
					conn.BufferIn += strResult
					if len(conn.BufferIn) > 3 {
						conn.BufferIn, _ = checkBuffer(conn.BufferIn, conn.ID)
					}
					conn.Changed = currTime
				}
				//fmt.Println(strResult)
				//// Send a response back to person contacting us.
				//conn.Write([]byte("Message received."))

				// check if a connections is longer inactive
				// than the given timeout
				ConnectionsMutex.Lock()
				if conn.Changed+conn.Timeout < currTime {
					conn.Connection.Close()
					delete(Connections, id)
					delete(ConnectionIDs, id)
				} else {
					Connections[id] = conn
				}
				ConnectionsMutex.Unlock()
			}
		} else {
			time.Sleep(100000000) // 0,1 seconds
		}
	}
}

func checkBuffer(buffer string, id int) (string, error) {
	// get checklength

	checkLength := buffer[0:3]
	// convert to int
	intLength, err := strconv.Atoi(checkLength)
	if err != nil {
		// handle error
		//handlePackageError(conn)
		//fmt.Println(err)
		//os.Exit(2)
		return "", errors.New("Package error. Buffer clear")
	}
	fmt.Println(intLength)
	// checkLength needs to be converted
	intLength = PackageSize - intLength
	bufferLength := len(buffer)
	fmt.Println("intLength: ", intLength, " | bufferLenth: ", bufferLength)
	// lets check if we got the whole package allready
	if bufferLength == 3+intLength {
		// we got the full package length
		// lets retrieve the package
		strPackage := buffer[3:]
		fmt.Println("Package: ", strPackage)
		go handlePackage(strPackage, id)
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
	return buffer, nil
}

func handlePackage(data string, id int) {
	fmt.Println("\nRecieved data from API:\n", data, " ", id)
	ConnectionsMutex.Lock()
	Connections[id].Connection.Write([]byte("390yes it worked"))
	ConnectionsMutex.Unlock()
}

func handlePackageError(conn Connection) {

}
