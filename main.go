package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/divan/gorilla-xmlrpc/xml"
	"github.com/gorilla/rpc"
	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/parser"
)

type LogService struct {
	ems map[event.ComponentType]*event.EventManager
}

func NewLogService() *LogService {
	ls := &LogService{ems: make(map[event.ComponentType]*event.EventManager)}
	tps := []event.ComponentType{
		event.ComponentTiDB,
		event.ComponentTiKV,
		event.ComponentPD,
		event.ComponentTiFlash,
	}
	for _, tp := range tps {
		em, err := event.NewEventManager(tp)
		if err != nil {
			panic(err)
		}
		ls.ems[tp] = em
	}
	return ls
}

func (h *LogService) ID(
	r *http.Request,
	args *struct {
		Component string
		Log       []byte
	},
	reply *struct {
		ID string
	}) error {

	ct, err := event.GetComponentType(args.Component)
	if err != nil {
		return err
	}
	ls, err := parser.ParseFromString(string(args.Log))
	if err != nil {
		return err
	}
	if len(ls) == 0 {
		return errors.New("no valid logs provided")
	}
	l := ls[0]
	em := h.ems[ct]
	rule := em.GetRuleByLog(l)
	if rule != nil {
		reply.ID = strconv.Itoa(int(rule.ID))
	} else {
		reply.ID = "0"
	}
	return nil
}

type LogField struct {
	Name  string
	Value []byte
}

func (h *LogService) Parse(
	r *http.Request,
	args *struct {
		Component string
		Log       []byte
	},
	reply *struct {
		ID       string
		Name     string
		Level    string
		DateTime string
		File     string
		Line     string
		Message  []byte
		Fields   []LogField
	}) error {

	ct, err := event.GetComponentType(args.Component)
	if err != nil {
		return err
	}
	ls, err := parser.ParseFromString(string(args.Log))
	if err != nil {
		return err
	}
	if len(ls) == 0 {
		return errors.New("no valid logs provided")
	}
	l := ls[0]
	em := h.ems[ct]
	rule := em.GetRuleByLog(l)
	if rule != nil {
		reply.ID = strconv.Itoa(int(rule.ID))
		reply.Name = rule.Name
	} else {
		reply.ID = "0"
		reply.Name = ""
	}
	reply.Level = string(l.Header.Level)
	reply.DateTime = l.Header.DateTime.Format(time.RFC3339)
	reply.File = l.Header.File
	reply.Line = strconv.Itoa(int(l.Header.Line))
	reply.Message = []byte(l.Message)
	reply.Fields = make([]LogField, 0)
	for _, f := range l.Fields {
		reply.Fields = append(reply.Fields, LogField{
			Name:  f.Name,
			Value: []byte(f.Value),
		})
	}
	return nil
}

func main() {
	RPC := rpc.NewServer()
	xmlrpcCodec := xml.NewCodec()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	RPC.RegisterService(NewLogService(), "")
	http.Handle("/RPC2", RPC)

	port := os.Args[1]
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}
