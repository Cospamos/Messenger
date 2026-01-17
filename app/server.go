package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Response[T any] struct {
	Addressee      string
	MessageType    string
	MessageContent []T
}

type Message struct {
	Id      string
	ImgUrl  string
	Content string
}

type Payload struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Client struct { // System obj
	id     string
	recent []time.Time
	mu     sync.Mutex
	closed bool
}

type ClientDTO struct { // Data transder obj
	ID     string
	Name   string
	ImgURL string
}

type ClientDBO struct { // Database obj
	Name   string
	ImgURL string
}

var (
	clients    = make(map[*websocket.Conn]*Client)
	messages   = []Message{}
	clientsMu  sync.Mutex
	messagesMu sync.Mutex
	limit      = 20
	window     = 1 * time.Second
	p          Payload
)


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func broadcastUsers[T any](res Response[T]) {
	clientsMu.Lock()
	clientList := make([]*Client, 0, len(clients))
	conns := make([]*websocket.Conn, 0, len(clients))
	for conn, client := range clients {
		clientList = append(clientList, client)
		conns = append(conns, conn)
	}
	clientsMu.Unlock()

	sem := make(chan struct{}, 10)

	for i, conn := range conns {
		cl := clientList[i]
		sem <- struct{}{}
		go func(c *websocket.Conn, cl *Client) {
			defer func() { <-sem }()

			cl.mu.Lock()
			if cl.closed {
				cl.mu.Unlock()
				return
			}

			res.Addressee = cl.id
			b, _ := json.Marshal(res)
			err := c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				cl.closed = true
				conn.Close()
			}
			cl.mu.Unlock()
		}(conn, cl)
	}

	for j := 0; j < cap(sem); j++ {
		sem <- struct{}{}
	}
}

func sendFullMessagesData(conn *websocket.Conn) {
	res := Response[Message]{
		Addressee:      "",
		MessageType:    "messages",
		MessageContent: []Message{},
	}

	messagesMu.Lock()
	msgsCopy := make([]Message, len(messages))
	copy(msgsCopy, messages)
	messagesMu.Unlock()

	res.MessageContent = msgsCopy

	b, _ := json.Marshal(res)

	err := conn.WriteMessage(websocket.TextMessage, b)

	if err != nil {
		fmt.Println("Write message error:", err)
	}
}

func sendMessagesData() {
	res := Response[Message]{
		Addressee:      "",
		MessageType:    "messages",
		MessageContent: []Message{},
	}

	messagesMu.Lock()
	res.MessageContent = append(res.MessageContent, messages[len(messages)-1])
	messagesMu.Unlock()

	broadcastUsers(res)
}

func sendUsersData() {
	res := Response[ClientDTO]{
		Addressee:      "",
		MessageType:    "clients",
		MessageContent: []ClientDTO{},
	}

	clientsMu.Lock()
	idsCopy := make([]ClientDTO, 0, len(clients))
	var clientDBO ClientDBO
	for _, cl := range clients {
		if err := GetClient(cl.id, &clientDBO); err != nil {
			fmt.Println("[ERROR]GetClient:", err)
			continue
		}

		idsCopy = append(idsCopy, ClientDTO{
			ID:     cl.id,
			Name:   clientDBO.Name,
			ImgURL: clientDBO.ImgURL,
		})
	}
	clientsMu.Unlock()

	res.MessageContent = idsCopy

	broadcastUsers(res)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error: ", err)
	}
	defer conn.Close()

	clientsMu.Lock()

	id := uuid.New().String()
	clients[conn] = &Client{
		id:     id,
		recent: []time.Time{time.Now()},
	}

	if err = AddClient(id, ClientDBO{Name: "", ImgURL: ""}); err != nil {
		fmt.Println("[ERROR]AddClient:", err)
		return
	}
	clientsMu.Unlock()

	sendUsersData()
	sendFullMessagesData(conn)

	for {
		_, msg, err := conn.ReadMessage()
		client := clients[conn]

		var clientDBO ClientDBO
		if err := GetClient(id, &clientDBO); err != nil {
			fmt.Println("[ERROR]GetClient:", err)
			continue
		}

		if err != nil {
			client.mu.Lock()
			if !client.closed {
				client.closed = true
				conn.Close()
			}
			client.mu.Unlock()

			clientsMu.Lock()
			if err := RemoveClient(id); err != nil {
				fmt.Println("[ERROR]RemoveClient:", err)
			}
			delete(clients, conn)
			clientsMu.Unlock()

			if ce, ok := err.(*websocket.CloseError); ok && ce.Code == websocket.CloseGoingAway {
				fmt.Println("Client disconected:", client.id)
			} else {
				fmt.Println("Read error:", err)
			}

			sendUsersData()
			return
		}

		client.recent = []time.Time{time.Now()}
		clientsMu.Lock()
		clients[conn] = client
		clientsMu.Unlock()

		if err = json.Unmarshal(msg, &p); err != nil {
			fmt.Println("JSON marshal error:", err)
			continue
		}

		switch p.Type {
		case "POST":
			messagesMu.Lock()
			messages = append(messages, Message{
				Id:      id,
				ImgUrl:  clientDBO.ImgURL,
				Content: p.Content,
			})
			messagesMu.Unlock()

			sendMessagesData()
		case "PATCH/name":
			res := Response[string]{
 			  Addressee:   client.id,
 			  MessageType: "error",
 			  MessageContent: []string{"The selected nickname is already in use. Please choose another one."},
			}
			if IsNameInDB(p.Content) {
				conn.WriteJSON(res)
				continue
			}
			
			clientsMu.Lock()
			clientDBO.Name = p.Content
			clientsMu.Unlock()

			UpdateClient(client.id, &clientDBO)
			sendUsersData()
		case "PATCH/ico":
			clientsMu.Lock()
			clientDBO.ImgURL = p.Content
			clientsMu.Unlock()

			UpdateClient(client.id, &clientDBO)
			sendUsersData()
		}
	}
}

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h.ServeHTTP(w, r)
	})
}


func RunServer() {
	fmt.Println("Local URL:   http://localhost:8080")
	fmt.Printf("Public URL:  https://zonia-interparliament-nonnormally.ngrok-free.dev\n\n")

	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", noCache(fs))

	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}
