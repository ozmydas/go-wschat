package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var allConnection = make([]*websocket.Conn, 0)

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		// upgrade jalur http jd buat websocket
		conn, _ := upgrader.Upgrade(w, r, nil)
		allConnection = append(allConnection, conn) // append koneksi baru ke list seluruh koneksi

		// informasikan jika ada yg join
		fmt.Printf("New Client Connected : %v\n", conn.RemoteAddr())
		Broadcast(conn, 1, "Joined", false)

		for {
			// baca pesan yg dikirim tiap koneksi
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				// jika ada koneksi yg error
				if strings.Contains(err.Error(), "websocket: close") {
					// hapus koneksi yg keluar

					i := 0 // cari dulu offset yg mau dihapus dari list seluruh koneksi
					for _, eachConn := range allConnection {
						if eachConn == conn {
							break // disini kita sudah dapat offset array
						}

						i++
					}

					allConnection[i] = allConnection[len(allConnection)-1] // copy array terakhir ke offset
					allConnection = allConnection[:len(allConnection)-1]   // delete array terakhir

					// informasikan jika ada keluar
					fmt.Printf("%v Exited\n", conn.RemoteAddr())
					Broadcast(conn, 1, "Exited", false)
					return
				}

				return
			}

			// Write message back to each connection
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
			Broadcast(conn, msgType, string(msg), false)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("running on 8080")
	http.ListenAndServe(":8080", nil)
} // end func

func Broadcast(conn *websocket.Conn, msgType int, msg string, isSendMe bool) {
	for _, eachConn := range allConnection {

		if eachConn == conn {
			if isSendMe {
				continue // jangan kirim msg ke diri sendiri
			}
		}

		eachConn.WriteMessage(msgType, []byte(conn.RemoteAddr().String()+" "+msg))
	}
}
