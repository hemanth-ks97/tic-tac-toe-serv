package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hemanth-ks97/tic-tac-toe-serv/pkg/game_utils"
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
    // allow connections from any origin (unsafe -- need to revisit)
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

        result := evaluate(payload);

        err = conn.WriteMessage(messageType, result);
        if err != nil{
            log.Println("Write:", err);
        }
    }
}

type Game struct{
    Board [3][3]game_utils.Cellstate    `json:"board"`
    Curr_player game_utils.Cellstate    `json:"cur_player"`
    Uid string                          `json:"uid"`
    State game_utils.GameState          `json:"state"`
    Move_num int                        `json:"move_num"`
}

//dictionary of active games
var games = make(map[string]Game); 

type message_in struct {
    GameId string                       `json:"gameid"`
    Move string                         `json:"move"`
}

func evaluate(payload []byte) []byte {
    
    switch req := string(payload); req{
    case "start":
        game := [3][3]game_utils.Cellstate{{game_utils.Blank,game_utils.Blank,game_utils.Blank},{game_utils.Blank,game_utils.Blank,game_utils.Blank},{game_utils.Blank,game_utils.Blank,game_utils.Blank}};
        player := game_utils.X;
        seed := rand.NewSource(time.Now().UnixNano());
        rand := rand.New(seed);
        uid := strconv.FormatInt(rand.Int63(), 16);

        games[uid] = Game{
            Board: game,
            Curr_player: player,
            Uid: uid,
            State: game_utils.IP,
            Move_num: 0,
        }

        jsondata, err := json.Marshal(games[uid]);
        if err != nil{
            log.Println("json.Marshal: ", err);
        }

        return jsondata;
    default:
        //extract message
        var msg message_in;
        err := json.Unmarshal([]byte(payload), &msg);
        if err != nil{
            log.Println("json.Unmarshal:", err);
        }

        //pull game
        game := games[msg.GameId];

        //update board
        valid := game_utils.UpdateGame(&game.Board, game.Curr_player, msg.Move); 
        if !valid {
            //generate and return err
            log.Println("Invalid Move");
            return []byte("Invalid Move");
        }
        game.Move_num += 1;

        //update game state
        game.State = game_utils.CheckGameState(&game.Board, game.Curr_player, game.Move_num);
        fmt.Println(game.State);

        //update player
        game_utils.SwitchPlayer(&game.Curr_player);

        //update game map
        games[msg.GameId] = game

        //game end condition
        if game.State != game_utils.IP{
            delete(games, msg.GameId);
            jsondata, err := json.Marshal(game);
            if err != nil{
                log.Println("json.Marshal:", err);
            }
            return jsondata;
        }

        //create resp (game still in progress)
        jsondata, err := json.Marshal(games[msg.GameId]);
        if err != nil{
            log.Println("json.Marshal:", err);
        }

        //send resp
        return jsondata;
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
    // game_utils.Say("Something I'm giving up on you");

    fmt.Println("HTTP SERVER: starting at", *PORT);
    log.Fatal(http.ListenAndServe(":" + *PORT, nil));
}
