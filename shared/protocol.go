package shared

import (
    "encoding/json"
    "net"
)

type Protocol struct {
    conn net.Conn
}

func NewProtocol(conn net.Conn) *Protocol {
    return &Protocol{conn: conn}
}

func (p *Protocol) SendMessage(msg map[string]interface{}) error {
    data, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    
    _, err = p.conn.Write(append(data, '\n'))
    return err
}

func (p *Protocol) ReadMessage() (map[string]interface{}, error) {
    var msg map[string]interface{}
    decoder := json.NewDecoder(p.conn)
    err := decoder.Decode(&msg)
    return msg, err
}

func (p *Protocol) Close() error {
    return p.conn.Close()
}
