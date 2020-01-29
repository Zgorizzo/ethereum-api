package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/INFURA/infra-test-benjamin-mateo/config"
	"github.com/INFURA/infra-test-benjamin-mateo/logger"
	"github.com/gorilla/mux"
)

var s *Server

func init() {

	// load configuration
	config.Load()
	// get an API server
	s = NewServer(logger.Init(true), mux.NewRouter())
	s.loadClient(config.ReadString("NODE_URL"))
	defer s.Logger.Sync()
}

func TestHandlers(t *testing.T) {

	tt := []struct {
		URL            string
		routeVariable  string
		handler        func(w http.ResponseWriter, r *http.Request)
		expectedRes    string
		expectedStatus int
	}{
		{"/call/{from:0x(?:[A-Fa-f0-9]{40})}/{to:0x(?:[A-Fa-f0-9]{40})}/{gas:[0-9]+}/{value:[0-9]+}/{data}", "/call/0x5cf2cefd110e7ce39fb353d123776ab683ef9fee/0xe530441f4f73bDB6DC2fA5aF7c3fC5fD551Ec838/30400/0/0x70a082310000000000000000000000005cf2cbfd110e7ce39fb353d123776ab683ef9feb", s.handleCall, `"0x000000000000000000000000000000000000000000000000000000000016bc50"`, http.StatusOK},
		{"/block/last/height", "/block/last/height", s.handleGetLatestBlockID, `{"lastBlockHeight":`, http.StatusOK},
		{"/block/last", "/block/last", s.handleGetLastBlock(false), `{"number":`, http.StatusOK},
		{"/block/last/full", "/block/last/full", s.handleGetLastBlock(true), `"transactionIndex"`, http.StatusOK},
		{"/block/{height:[0-9]+}", "/block/5681044", s.handleGetBlockByHeight(false), `number":"0x56af94","hash":"0xc9ad7040dee3e49e7e7ab396278b23a738bc00e11edf0ac49520a74f1692c9cc"`, http.StatusOK},
		{"/block/{height:[0-9]+}", "/block/dssd", s.handleGetBlockByHeight(false), "", http.StatusNotFound},
		{"/block/{height:[0-9]+}/full", "/block/5681044/full", s.handleGetBlockByHeight(true), `"transactions":[{"blockHash":"0xc9ad7040dee3e49e7e7ab396278b23a738bc00e11edf0ac49520a74f1692c9cc"`, http.StatusOK},
		{"/block/{hash:0x(?:[A-Fa-f0-9]{64})$}", "/block/0x840eed2d9390ea527089a7c2f2b6dd28ef382c094105d29d22b911e606c4e0bc", s.handleGetBlockByHash(false), `"number":"0xefc87","hash":"0x840eed2d9390ea527089a7c2f2b6dd28ef382c094105d29d22b911e606c4e0bc"`, http.StatusOK},
		{"/block/{hash:0x(?:[A-Fa-f0-9]{64})}/full", "/block/0x840eed2d9390ea527089a7c2f2b6dd28ef382c094105d29d22b911e606c4e0bc/full", s.handleGetBlockByHash(true), `"r":"0x87effd10150f3d09c5ffcb79111c654637ba7a3a84c3e733d51003a029660b86","s":"0x5d68bc045b0a2d35160497c9c94da08b4dc32b8e9fb04d917c7442f1ad903205"`, http.StatusOK},
		{"/transaction/{hash:0x(?:[A-Fa-f0-9]{64})$}", "/transaction/0x37e458fcff2a79f32257776aa67f929187d2ff1f8868092bead0b788d248b9b4", s.handleGetTransactionByHash(), `{"blockHash":"0xf247cc1a2cc1b3af1094674bd191594eeeff1e89e5036b377731758055debd51","blockNumber":"0x75a900"`, http.StatusOK},
		{"/transaction/{hash:0x(?:[A-Fa-f0-9]{64})$}", "/transaction/0xgge458fcff2a79f32257776aa67f929187d2ff1f8868092bead0b788d248b9b4", s.handleGetTransactionByHash(), "", http.StatusNotFound},
		{"/gasprice", "/gasprice", s.handleGetGasPrice, `{"gasPrice":`, http.StatusOK},
		{"/balance/{address:0x(?:[A-Fa-f0-9]{40})$}", "/balance/0x5cf2CBfd110E7Ce39fb353d123776Ab683ef9fEB", s.handleGetBalance, `{"balance":`, http.StatusOK},
		{"/log/{from:0x(?:[A-Fa-f0-9]+)}/{to:0x(?:[A-Fa-f0-9]+)}/{topic}", "/log/9135250/9135260/0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", s.handleGetLogs, "", http.StatusNotFound},
		{"/log/{from:0x(?:[A-Fa-f0-9]+)}/{to:0x(?:[A-Fa-f0-9]+)}/{topic}", "/log/0x8B6492/0x8B649C/0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", s.handleGetLogs, `topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"`, http.StatusOK},
		{"/block/{height:[0-9]+}/transaction/{id:[0-9]+}", "/block/9135267/transaction/4", s.handleGetTransactionByIDInBlockHash, `{"blockHash":"0xafddcde383196ecacf07def76a2572c86d1f992081cb1023c8decef92b7712ee"`, http.StatusOK},
		{"/block/{height:[0-9]+}/transaction/{id:[0-9]+}", "/block/9135267/transaction/99999999999999999999999999999999", s.handleGetTransactionByIDInBlockHash, "", http.StatusBadRequest},
	}

	for _, tc := range tt {
		t.Run(tc.URL, func(t *testing.T) {
			//path := fmt.Sprintf("/metrics/%s", tc.routeVariable)
			req, err := http.NewRequest("GET", tc.routeVariable, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			// Need to create a router that we can pass the request through so that the vars will be added to the context
			router := mux.NewRouter()
			router.HandleFunc(tc.URL, tc.handler)
			router.ServeHTTP(rr, req)

			// In this case, our Metrics returns a non-200 response
			// for a route variable it doesn't know about.
			if rr.Code != tc.expectedStatus {
				t.Errorf(" should have failed on routeVariable %s: got %v want %v",
					tc.routeVariable, rr.Code, tc.expectedStatus)
			}
			if !strings.Contains(rr.Body.String(), tc.expectedRes) {
				t.Errorf(" returned unexpected body: got %v want %v", rr.Body.String(), tc.expectedRes)
			}
		})

	}
}
