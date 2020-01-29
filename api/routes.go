// Package INFURA REST API.
//
// this is a RESTful APIs in golang wrapping call to an INFURA Node.
//
// Terms Of Service:
//
//     Schemes: http
//     Host: localhost:8000
//     Version: 1.0.0
//     Contact: Benjamin MATEO<bmateo@pm.me>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
// swagger:meta
package api

import "net/http"

// all the routes are defined here
func (s *Server) routes() {

	// swagger:operation GET / root handleRoot
	//
	// Returns a simple json with no call to the eth client
	//
	// If the API is running, an ok status will be returned
	//
	// ---
	// responses:
	//   "200":
	//     description: API is running, an ok status will be returned
	//     schema:
	//      type: object
	//      properties:
	//        ok:
	//          type: boolean
	//      example:
	//        ok: true
	s.router.HandleFunc("/", s.handleRoot()).Methods("GET")

	t := s.router.PathPrefix("/transaction").Subrouter()

	// swagger:operation GET /transactions/{hash} transaction handleGetTransactionByHash
	//
	// Returns information about a transaction for a given hash
	//
	// If the transaction is found, transaction will be returned
	// else Error Not Found (404) will be returned.
	//
	// ---
	// parameters:
	// - name: hash
	//   in: path
	//   description: a string representing the hash (32 bytes) of a transaction
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: transaction is returned
	//     schema:
	//       $ref: '#/definitions/Transaction'
	//   "404":
	//     description: transaction not found
	t.HandleFunc("/{hash:0x(?:[A-Fa-f0-9]{64})$}", s.handleGetTransactionByHash()).Methods("GET")

	b := s.router.PathPrefix("/block").Subrouter()

	// swagger:operation GET /block/last block handleGetLastBlock
	//
	// Returns information about the last block.
	//
	//
	// ---
	// responses:
	//   "200":
	//     description: block is returned
	//     schema:
	//       $ref: '#/definitions/Block'
	b.HandleFunc("/last", s.handleGetLastBlock(false)).Methods("GET")

	// swagger:operation GET /block/last/full block handleGetLastBlockFull
	//
	// Returns information about the last block including all the transactions details contained in the block
	//
	// ---
	// responses:
	//   "200":
	//     description: full block is returned
	//     schema:
	//       $ref: '#/definitions/Block'
	b.HandleFunc("/last/full", s.handleGetLastBlock(true)).Methods("GET")

	// swagger:operation GET /block/last/height block handleGetLatestBlockID
	//
	// Returns the last block height
	//
	// ---
	// responses:
	//   "200":
	//     description: transaction is returned
	//     schema:
	//      type: object
	//      properties:
	//        lastBlockHeight:
	//          type: number
	//      example:
	//         lastBlockHeight: 9206229
	b.HandleFunc("/last/height", s.handleGetLatestBlockID).Methods("GET")

	// swagger:operation GET /block/{hash} block handleGetBlockByHash
	//
	// Returns information about a block by hash .
	//
	// ---
	// parameters:
	// - name: hash
	//   in: path
	//   description: a string representing the hash (32 bytes) of a block
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: block is returned
	//     schema:
	//       $ref: '#/definitions/Block'
	//   "404":
	//     description: block not found
	b.HandleFunc("/{hash:0x(?:[A-Fa-f0-9]{64})$}", s.handleGetBlockByHash(false)).Methods("GET")

	// swagger:operation GET /block/{hash}/full block handleGetBlockByHashfull
	//
	// Returns information about a block by hash including all the transactions details contained in the block.
	//
	// ---
	// parameters:
	// - name: hash
	//   in: path
	//   description: a string representing the hash (32 bytes) of a block
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: block is returned
	//     schema:
	//       $ref: '#/definitions/Block'
	//   "404":
	//     description: block not found
	b.HandleFunc("/{hash:0x(?:[A-Fa-f0-9]{64})}/full", s.handleGetBlockByHash(true)).Methods("GET")

	// swagger:operation GET /block/{height} block handleGetBlockByHeight
	//
	// Returns information about a block by height.
	//
	// ---
	// parameters:
	// - name: height
	//   in: path
	//   description: a number representing the height of a block
	//   type: number
	//   required: true
	// responses:
	//   "200":
	//     description: block is returned
	//     schema:
	//       $ref: '#/definitions/Block'
	//   "404":
	//     description: block not found
	b.HandleFunc("/{height:[0-9]+}", s.handleGetBlockByHeight(false)).Methods("GET")

	// swagger:operation GET /block/{height}/full block handleGetBlockByHeightfull
	//
	// Returns information about a block by height including all the transactions details contained in the block.
	//
	// ---
	// parameters:
	// - name: height
	//   in: path
	//   description: a number representing the height of a block
	//   type: number
	//   required: true
	// responses:
	//   "200":
	//     description: block is returned
	//     schema:
	//       $ref: '#/definitions/Block'
	//   "404":
	//     description: block not found
	b.HandleFunc("/{height:[0-9]+}/full", s.handleGetBlockByHeight(true)).Methods("GET")

	// swagger:operation GET /block/{height}/transaction/{id} block handleGetTransactionByIDInBlockHash
	//
	// Returns information about a transaction by block number and transaction index position.
	//
	// If the transaction is found, transaction will be returned
	// else Error Not Found (404) will be returned.
	//
	// ---
	// parameters:
	// - name: height
	//   in: path
	//   description: an integer block number
	//   type: number
	//   required: true
	// - name: id
	//   in: path
	//   description: an integer representing the position in the block
	//   type: number
	//   required: true
	// responses:
	//   "200":
	//     description: transaction is returned
	//     schema:
	//       $ref: '#/definitions/Transaction'
	//   "404":
	//     description: transaction not found
	b.HandleFunc("/{height:[0-9]+}/transaction/{id:[0-9]+}", s.handleGetTransactionByIDInBlockHash).Methods("GET")

	// swagger:operation GET /gasprice gas handleGetGasPrice
	//
	// Returns the current gas price in wei
	//
	// ---
	// responses:
	//   "200":
	//     description: gasPrice is returned in wei
	//     schema:
	//      type: object
	//      properties:
	//        gasPrice:
	//          type: number
	//      example:
	//         gasPrice: 4000000000
	s.router.HandleFunc("/gasprice", s.handleGetGasPrice).Methods("GET")

	// swagger:operation GET /balance/{address} balance handleGetBalance
	//
	// Returns balance in wei of the given address.
	//
	// If the address is found, balance will be returned
	// else Error Not Found (404) will be returned.
	//
	// ---
	// parameters:
	// - name: address
	//   in: path
	//   description: a string representing the address (20 bytes) to check for balance
	//   type: number
	//   required: true
	// responses:
	//   "200":
	//     description: balance is returned
	//     schema:
	//      type: object
	//      properties:
	//        balance:
	//          type: number
	//      example:
	//       balance: 2381188418352874359
	//   "404":
	//     description: address not found
	s.router.HandleFunc("/balance/{address:0x(?:[A-Fa-f0-9]{40})$}", s.handleGetBalance).Methods("GET")

	// swagger:operation GET /log/{from}/{to}/{topic} log handleGetLogs
	//
	// Returns an array of all logs matching a given filter object.
	//
	// If logs are found, logs will be returned
	// else Error Not Found (404) will be returned.
	//
	// ---
	// parameters:
	// - name: from
	//   in: path
	//   description:  an integer block number encoded in hex, or the string "latest", "earliest" or "pending"
	//   type: string
	//   required: true
	// - name: to
	//   in: path
	//   description:  an integer block number encoded in hex, or the string "latest", "earliest" or "pending"
	//   type: string
	//   required: true
	// - name: topic
	//   in: path
	//   description: Array of 32 Bytes DATA topics. Topics are order-dependent
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: logs are returned
	//     schema:
	//       type: array
	//       items:
	//         $ref: '#/definitions/Log'
	//   "404":
	//     description: logs not found
	s.router.HandleFunc("/log/{from:0x(?:[A-Fa-f0-9]+)}/{to:0x(?:[A-Fa-f0-9]+)}/{topic}", s.handleGetLogs).Methods("GET")

	// swagger:operation GET /call/{from}/{to}/{gas}/{value}/{data} call handleCall
	//
	// Executes a new message call immediately without creating a transaction on the block chain.
	//
	//
	// ---
	// parameters:
	// - name: from
	//   in: path
	//   description:  20 Bytes - The address the transaction is sent from.
	//   type: string
	//   required: true
	// - name: to
	//   in: path
	//   description:  20 Bytes - The address the transaction is directed to.
	//   type: string
	//   required: true
	// - name: gas
	//   in: path
	//   description: Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.
	//   type: string
	//   required: true
	// - name: value
	//   in: path
	//   description:  Integer of the value sent with this transaction
	//   type: string
	//   required: true
	// - name: data
	//   in: path
	//   description:  Hash of the method signature and encoded parameters. For details see Ethereum Contract ABI
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description:  the return value of the executed contract method.
	//     schema:
	//       type: string
	//   "404":
	//     description: logs not found
	s.router.HandleFunc("/call/{from:0x(?:[A-Fa-f0-9]{40})}/{to:0x(?:[A-Fa-f0-9]{40})}/{gas:[0-9]+}/{value:[0-9]+}/{data}", s.handleCall).Methods("GET")

	// swagger:operation GET /describe describe handleGetDescription
	//
	// Returns information about the available api routes.
	//
	//
	// ---
	//
	// responses:
	//   "200":
	//     description: routes are returned
	s.router.HandleFunc("/describe", s.handleGetDescription(s.router)).Methods("GET")

	// This will serve files under http://localhost:8000/swaggerui/<filename>
	s.router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui"))))

}
