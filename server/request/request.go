package request

import (
	"encoding/json"
	"fmt"
	"net"

	templ "tcp_http_server.com/server/template"
	"tcp_http_server.com/server/utils"
)

type Request struct {
	Method                  string
	Route                   string
	Query                   map[string]string
	HttpVersion             string
	Host                    string
	Connection              string
	Sec_ch_ua               string
	Sec_ch_ua_mobile        string
	Sec_ch_ua_platform      string
	Dnt                     string
	UpgradeInsecureRequests string
	UserAgent               string
	Accept                  string
	SecFetchSite            string
	SecFetchMode            string
	SecFetchUser            string
	SecFetchDest            string
	AcceptEncoding          string
	AcceptLanguage          string
	conn                    *net.Conn
}

func CreateRequest(headers map[string]string, conn *net.Conn) Request {
	query := make(map[string]string)
	route := headers["route"]
	for pos, ch := range route {
		if ch == '?' {
			query = parseQuery(route[pos+1:])
			route = route[:pos]
		}
	}
	return Request{
		Host:                    headers["Host"],
		UpgradeInsecureRequests: headers["Upgrade-Insecure-Requests"],
		SecFetchMode:            headers["Sec-Fetch-Mode"],
		SecFetchDest:            headers["Sec-Fetch-Dest"],
		Method:                  headers["method"],
		Route:                   route,
		Query:                   query,
		HttpVersion:             headers["http_version"],
		Connection:              headers["Connection"],
		Sec_ch_ua:               headers["sec_ch_ua"],
		Sec_ch_ua_mobile:        headers["sec_ch_ua_mobile"],
		Sec_ch_ua_platform:      headers["sec_ch_ua_platform"],
		Dnt:                     headers["DNT"],
		Accept:                  headers["Accept"],
		SecFetchUser:            headers["Sec-Fetch-User"],
		AcceptLanguage:          headers["Accept-Language"],
		AcceptEncoding:          headers["Accept-Encoding"],
		UserAgent:               headers["User-Agent"],
		SecFetchSite:            headers["Sec-Fetch-Site"],
		conn:                    conn,
	}
}

func parseQuery(q string) map[string]string {
	queryMap := make(map[string]string)
	queries := utils.Split(q, "&")
	for _, qr := range queries {
		query := utils.Split(qr, "=")
		queryMap[query[0]] = query[1]
	}
	return queryMap
}

func ParseRequest(request []byte) map[string]string {
	lines := utils.Split(string(request), string('\n'))
	requestDet := make(map[string]string)
	requestLineRead := false
	for pos, line := range lines {
		println(pos, " -> ", line)
		cleanLine := utils.StripStart(line)
		if len(cleanLine) == 0 {
			continue
		}
		if !requestLineRead {
			requestDet = parseRequestLine(line)
			requestLineRead = true
			continue
		}
		keyValue := utils.SplitFirst(cleanLine, ":")
		requestDet[keyValue[0]] = keyValue[1]
	}
	return requestDet
}

func parseRequestLine(s string) map[string]string {
	words := utils.Split(s, " ")
	request := make(map[string]string)
	request["method"] = words[0]
	request["route"] = words[1]
	request["http_version"] = words[2]
	return request
}

func (req *Request) Write(s []byte) (int, error) {
	conn := *req.conn
	println(string(s))
	// if _, err := conn.Write([]byte(s)); err != nil {
	// 	return err
	// }
	val, err := conn.Write(s)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (req *Request) Html(html string) (int, error) {
	contentLength := len(html) + 1
	r := []byte("HTTP/1.1 200 OK\r\nConnection: close\r\nContent-Type: text/html\r\nContent-Length: " + fmt.Sprintf("%d", contentLength) + "\r\n\r\n " + html)
	return req.Write(r)
}

func (req *Request) RenderHtml(dir string, data map[string]interface{}) (int, error) {
	content, err := templ.LoadTemplateV2(dir, data)
	if err != nil {
		return 0, err
	}

	contentLength := len(content) + 1
	r := []byte("HTTP/1.1 200 OK\r\nConnection: close\r\nContent-Type: text/html\r\nContent-Length: " + fmt.Sprintf("%d", contentLength) + "\r\n\r\n " + string(content))
	return req.Write(r)
}

func (req *Request) Json(jsonMap map[string]interface{}, status int) (int, error) {
	jsonString, err := json.Marshal(jsonMap)
	if err != nil {
		return 0, err
	}
	contentLength := len(jsonString) + 1
	r := []byte("HTTP/1.1 " + fmt.Sprintf("%d", status) + " OK\r\nConnection: close\r\nContent-Type: text/json\r\nContent-Length: " + fmt.Sprintf("%d", contentLength) + "\r\n\r\n " + string(jsonString))
	return req.Write(r)
}

func (req *Request) Http(status int, content []byte) (int, error) {
	contentLength := len(content) + 1
	r := []byte("HTTP/1.1 " + fmt.Sprintf("%d", status) + " OK\r\nConnection: close\r\nContent-Type: text/css\r\nContent-Length: " + fmt.Sprintf("%d", contentLength) + "\r\n\r\n " + string(content))
	return req.Write(r)
}
