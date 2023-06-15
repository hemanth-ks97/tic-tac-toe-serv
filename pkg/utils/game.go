package main

import (
	"fmt"
)

type Cellstate string;

const(
    X Cellstate = "X"
    O Cellstate = "O"
    Blank Cellstate = "_"
)

type GameState string;

const(
    WX GameState = "Player X Wins"
    WO GameState = "Player O wins"
    D GameState = "Draw"
    IP GameState = "InProg"
)

func isValidMove(game *[3][3]Cellstate, player Cellstate, move string) bool{
    return game[move[0]-48][move[1]-48] == Blank;
}

func updateGame(game *[3][3]Cellstate, player Cellstate, move string) bool{
    if game[move[0]-48][move[1]-48] == Blank{
        game[move[0]-48][move[1]-48] = player; 
        return true;
    }
    return false;
}

func getWinState(player Cellstate) GameState{
    if player == X{
        return WX;
    }else{
        return WO;
    }
}

func checkGameState(game *[3][3]Cellstate, player Cellstate, move_num int) GameState{
    //optimization
    if move_num < 5{
        return IP;
    }
    //draw
    if move_num == 9{
        return D;
    }
    //winstate
    winstate := getWinState(player)
    //check rows
    for i := 0; i < 3; i++{
        count := 0;
        for j := 0; j<3; j++{
            if game[i][j] == player{
                count += 1;
            }
        }
        if count == 3{
            return winstate;
        }
    }
    //check colums
    for i := 0;i<3;i++{
        count := 0;
        for j:=0;j<3;j++{
            if game[j][i] == player{
                count += 1;
            }
        }
        if count == 3{
            return winstate;
        }
    } 
    //check diagnols
    count := 0;
    for i:=0;i<3;i++{
        for j:=0;j<3;j++{
            if i==j && game[i][j] == player{
                count += 1;
            }
        }
    }
    if count == 3{
        return winstate;
    }
    count = 0;
    for i:=0;i<3;i++{
        for j:=0;j<3;j++{
            if i+j==2 && game[i][j] == player{
                count += 1;
            }
        }
    }
    if count == 3{
        return winstate;
    }
    //return IP
    return IP;
}

func switchPlayer(player *Cellstate){
    if *player == X{
        *player = O;
    }else{
        *player = X;
    }
}

func Say(s string){
    fmt.Println(s);
}

func main(){
    fmt.Println("Tic-Tac-Toe");
    game := [3][3]Cellstate{{Blank,Blank,Blank},{Blank,Blank,Blank},{Blank,Blank,Blank}};
    curr_player := X;
    move_num := 0;
    for checkGameState(&game, curr_player, move_num) == IP{
        if move_num != 0{
            switchPlayer(&curr_player);
        }
        var move string
        fmt.Printf("Player %s's turn\n", curr_player);
        fmt.Scan(&move);
        
        for !updateGame(&game, curr_player, move){
            fmt.Println("Invalid move. Try again");
            fmt.Scan(&move);
        }

        fmt.Println(game);
        move_num += 1;
     }
     fmt.Println(checkGameState(&game, curr_player, move_num));
}
