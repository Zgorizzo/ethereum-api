package api

import (
	"encoding/json"

	"github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/infra-test-benjamin-mateo/node"

	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Warning each HTTP request gets its own go routine we might have concurrent code running against Server

// respond is response helper
// it can be convenient if we want to easily customize response type
func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.Logger.Warnf("can't encode data:%v err:%v", data, err)
		}
	}
}

// handleRoot return a simple json with no call to the eth client and can be used as a health endpoint
func (s *Server) handleRoot() http.HandlerFunc {
	ret := map[string]bool{"ok": true}
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, ret, http.StatusOK)
	}
}

// handleCall handles calls to eth_call which executes a new message call immediately without creating a transaction on the block chain.
func (s *Server) handleCall(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gas, err := strconv.ParseInt(params["gas"], 10, 64)
	if nok := !s.checkTypeError(w, r, gas, err); nok {
		return
	}
	value, err := strconv.ParseInt(params["value"], 10, 64)
	if nok := !s.checkTypeError(w, r, value, err); nok {
		return
	}
	from, err := eth.NewData(params["from"])
	if nok := !s.checkTypeError(w, r, from, err); nok {
		return
	}
	to, err := eth.NewData(params["to"])
	if nok := !s.checkTypeError(w, r, to, err); nok {
		return
	}
	data, err := eth.NewData(params["data"])
	if nok := !s.checkTypeError(w, r, data, err); nok {
		return
	}
	s.Logger.Infof("calling from:%s to:%s", from, to)
	w.Header().Add("Content-Type", "application/json")

	gasPrice, err := s.client.GetGasPrice(r.Context())
	if err != nil {
		s.Logger.Warn("can't get gas price error: ", err)
		s.respond(w, r, err, http.StatusFailedDependency)
	}

	p := node.CallParams{
		Data:     *data,
		From:     *from,
		Gas:      eth.QuantityFromUInt64(uint64(gas)),
		To:       *to,
		Value:    eth.QuantityFromUInt64(uint64(value)),
		GasPrice: eth.QuantityFromUInt64(uint64(gasPrice)),
	}

	res, err := s.client.CallContract(r.Context(), p)
	if err != nil {
		s.Logger.Warnf("call from:%s to:%s failed err:%s", from, to, err)
		s.respond(w, r, err.Error(), http.StatusNotFound)
	} else {
		s.respond(w, r, res, http.StatusOK)
	}

}

// handleGetTransactionByHash returns transaction by hash
func (s *Server) handleGetTransactionByHash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		hash := params["hash"]
		s.Logger.Infof("Request received to get a transaction by hash: %s", hash)

		w.Header().Add("Content-Type", "application/json")

		t, err := s.client.TransactionByHash(r.Context(), hash)
		if err != nil {
			s.Logger.Warnf("Tx hash does not exist: %s err:%s", hash, err)
			s.respond(w, r, err.Error(), http.StatusNotFound)
		} else {
			s.respond(w, r, t, http.StatusOK)
		}
	}
}

// handleGetBlockByHeight handle the root api
func (s *Server) handleGetBlockByHeight(full bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		height := params["height"]
		h, err := strconv.ParseInt(height, 10, 64)
		if err != nil {
			s.Logger.Infof("%d of type %T", h, h)
			s.respond(w, r, err.Error(), http.StatusBadRequest)
		}

		w.Header().Add("Content-Type", "application/json")
		s.Logger.Infof("Request received to get a block by height: %s full: %v", height, full)

		t, err := s.client.BlockByNumber(r.Context(), uint64(h), full)
		if err != nil {
			s.Logger.Warnf("can't get block height:%s err:%s", height, err)
			s.respond(w, r, err.Error(), http.StatusNotFound)
		} else {
			s.respond(w, r, t, http.StatusOK)
		}

	}
}

// handleGetTransactionByIDInBlockHash handle the root api
func (s *Server) handleGetTransactionByIDInBlockHash(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("get Transaction By ID In Block Hash")
	w.Header().Add("Content-Type", "application/json")
	params := mux.Vars(r)
	height := params["height"]
	h, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		s.Logger.Infof("%d of type %T", h, h)
		s.respond(w, r, err.Error(), http.StatusBadRequest)
	}
	id := params["id"]
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		s.Logger.Infof("%d of type %T", h, h)
		s.respond(w, r, err.Error(), http.StatusBadRequest)
	}

	res, err := s.client.TransactionByBlockNumberAndIndex(r.Context(), uint64(h), uint64(i))
	if err != nil {
		s.Logger.Infof("can't get transaction ID:%v in block height:%v err:%s", i, h, err)
		s.respond(w, r, err.Error(), http.StatusNotFound)
	} else {
		s.respond(w, r, res, http.StatusOK)
	}

}

// handleGetLogs
func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("get Transaction By ID In Block Hash")
	w.Header().Add("Content-Type", "application/json")
	params := mux.Vars(r)
	from, err := eth.NewBlockNumberOrTag(params["from"])
	if nok := !s.checkTypeError(w, r, from, err); nok {
		return
	}

	to, err := eth.NewBlockNumberOrTag(params["to"])
	if nok := !s.checkTypeError(w, r, to, err); nok {
		return
	}

	topic, err := eth.NewTopic(params["topic"])
	if nok := !s.checkTypeError(w, r, topic, err); nok {
		return
	}

	filter := eth.LogFilter{FromBlock: from, ToBlock: to, Topics: [][]eth.Data32{[]eth.Data32{*topic}}}

	res, err := s.client.Logs(r.Context(), filter)
	if err != nil {
		s.Logger.Warnf("can't get log topic:%v from block:%v to %v", topic, from, to)
		s.respond(w, r, err.Error(), http.StatusNotFound)
	} else {
		s.respond(w, r, res, http.StatusOK)
	}

}

// handleGetBlockByHash handle the root api
func (s *Server) handleGetBlockByHash(full bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Infof("get Block by hash full: %v", full)
		w.Header().Add("Content-Type", "application/json")
		params := mux.Vars(r)
		hash := params["hash"]
		res, err := s.client.BlockByHash(r.Context(), hash, full)
		if err != nil {
			s.Logger.Warn("can't get  Block By Hash error: ", err)
			s.respond(w, r, err.Error(), http.StatusNotFound)
		} else {
			s.respond(w, r, res, http.StatusOK)
		}
	}
}

// handleGetLastBlock get the latest Block  validated by the node
func (s *Server) handleGetLastBlock(full bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.Logger.Infof("get last block full:%v", full)
		w.Header().Add("Content-Type", "application/json")

		b, err := s.client.BlockNumber(r.Context())
		if err != nil {
			s.Logger.Warn("BlockNumber error: ", err)
			s.respond(w, r, err.Error(), http.StatusFailedDependency)
		}
		t, err := s.client.BlockByNumber(r.Context(), b, full)
		if err != nil {
			s.Logger.Warnf("can't get block height: %v err:%s", b, err)
			s.respond(w, r, err.Error(), http.StatusNotFound)
		} else {
			s.respond(w, r, t, http.StatusOK)
		}
	}
}

// handleGetLatestBlockID get the latest Block height validated by the node
func (s *Server) handleGetLatestBlockID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("get last block height")
	w.Header().Add("Content-Type", "application/json")

	b, err := s.client.BlockNumber(r.Context())
	if err != nil {
		s.Logger.Warn("can't get Block Number error: ", err)
		s.respond(w, r, err.Error(), http.StatusFailedDependency)
	} else {
		data := struct {
			LastBlockHeight uint64 `json:"lastBlockHeight"`
		}{b}
		s.respond(w, r, data, http.StatusOK)
	}
}

// handleGetGasPrice get the current gas price
func (s *Server) handleGetGasPrice(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("get current gas price")
	w.Header().Add("Content-Type", "application/json")

	b, err := s.client.GetGasPrice(r.Context())
	if err != nil {
		s.Logger.Warn("can't get gas price error: ", err)
		s.respond(w, r, err.Error(), http.StatusFailedDependency)
	} else {
		data := struct {
			GasPrice uint64 `json:"gasPrice"`
		}{b}
		s.respond(w, r, data, http.StatusOK)
	}

}

// handleGetBalance get the current balance
func (s *Server) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("get  balance")
	w.Header().Add("Content-Type", "application/json")
	params := mux.Vars(r)
	address := params["address"]
	b, err := s.client.GetBalance(r.Context(), address)
	if err != nil {
		s.Logger.Warnf("can't get balance for:%s error:%s", address, err)
		s.respond(w, r, err.Error(), http.StatusFailedDependency)
	} else {
		data := struct {
			Balance uint64 `json:"balance"`
		}{b}
		s.respond(w, r, data, http.StatusOK)
	}
}

// handleGetDescription describe all the API routes
func (s *Server) handleGetDescription(rout *mux.Router) http.HandlerFunc {
	s.Logger.Info("get API description")
	var ret []map[string]string

	return func(w http.ResponseWriter, r *http.Request) {
		// an example API handler

		err := rout.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			res := make(map[string]string)
			pathTemplate, err := route.GetPathTemplate()
			if err == nil {
				res["ROUTE"] = pathTemplate
			}
			pathRegexp, err := route.GetPathRegexp()
			if err == nil {
				res["Path regexp"] = pathRegexp
			}
			queriesTemplates, err := route.GetQueriesTemplates()
			if err == nil {
				res["Queries templates"] = strings.Join(queriesTemplates, ",")

			}
			queriesRegexps, err := route.GetQueriesRegexp()
			if err == nil {
				res["Queries regexps"] = strings.Join(queriesRegexps, ",")

			}
			methods, err := route.GetMethods()
			if err == nil {
				res["Methods"] = strings.Join(methods, ",")

			}
			ret = append(ret, res)
			return nil
		})

		if err != nil {
			s.Logger.Warn("Description error: ", err)
		}
		s.respond(w, r, ret, http.StatusOK)
	}

}

func (s *Server) checkTypeError(w http.ResponseWriter, r *http.Request, value interface{}, err error) bool {
	if err != nil {
		s.Logger.Infof("can't get %v to type %T", value, value)
		s.respond(w, r, err, http.StatusBadRequest)
		return false
	}
	return true
}
