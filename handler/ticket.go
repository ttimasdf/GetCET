package handler

import (
	"encoding/json"
	"net/http"
)

// TicketHandler 用于处理准考证找回的POST请求。
// 参数：province, school, name, cet_type
// 返回值（json）：ticket_number
func TicketHandler(w http.ResponseWriter, req *http.Request) {
	ticketNumber := "233336666"
	response, _ := json.Marshal(map[string]string{"ticket_number": ticketNumber})
	w.Write(response)
}
