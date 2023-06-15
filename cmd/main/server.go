package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	//"github.com/hemanth-ks97/tic-tac-toe-serv/pkg/utils"
)

var PORT = flag.String("port", "8080", "Defines the port on the localhost to start the server");

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func hello(w http.ResponseWriter, r *http.Request){
    if r.Method != "GET"{
        http.Error(w, "Method not supported", http.StatusNotFound);
        return;
    }

    fmt.Fprintf(w, "Hello\n");
}

//echo websocket handler
func echo(w http.ResponseWriter, r *http.Request){
    upgrader.CheckOrigin = func(r *http.Request) bool {
        return true;
    }
    conn, err := upgrader.Upgrade(w,r, nil);
    if err != nil{
        log.Println("upgrade:", err);
        return;
    }

    defer conn.Close();

    for{
        messageType, payload, err := conn.ReadMessage();
        if err != nil{
            log.Println("Read:", err);
        }

        if messageType == websocket.TextMessage {
            log.Println("Received:", string(payload));
        }

        err = conn.WriteMessage(messageType, payload);
        if err != nil{
            log.Println("Write:", err);
        }
    }
}

func setuproutes(){
    http.HandleFunc("/hello", hello);
    http.HandleFunc("/start", start);
    http.HandleFunc("/echows", echo);
}

func start(w http.ResponseWriter, r *http.Request){}

func main(){
    flag.Parse();
    setuproutes();
    //utils.Say("Hello from utils!");

    fmt.Println("HTTP SERVER: starting at", *PORT);
    log.Fatal(http.ListenAndServe(":" + *PORT, nil));
}
