// Copyright (c) 2019-2020 The Zcash developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or https://www.opensource.org/licenses/mit-license.php .

// Package frontend implements the gRPC handlers called by the wallets.
package frontend

import (
	"context"
	"errors"
	"io"
	"net"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adityapk00/lightwalletd/common"
	"github.com/adityapk00/lightwalletd/walletrpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type latencyCacheEntry struct {
	timeNanos   int64
	lastBlock   uint64
	totalBlocks uint64
}

type lwdStreamer struct {
	cache      *common.BlockCache
	chainName  string
	pingEnable bool
	walletrpc.UnimplementedCompactTxStreamerServer
	latencyCache map[string]*latencyCacheEntry
	latencyMutex sync.RWMutex
}

// NewLwdStreamer constructs a gRPC context.
func NewLwdStreamer(cache *common.BlockCache, chainName string, enablePing bool) (walletrpc.CompactTxStreamerServer, error) {
	return &lwdStreamer{cache: cache, chainName: chainName, pingEnable: enablePing, latencyCache: make(map[string]*latencyCacheEntry), latencyMutex: sync.RWMutex{}}, nil
}

// DarksideStreamer holds the gRPC state for darksidewalletd.
type DarksideStreamer struct {
	cache *common.BlockCache
	walletrpc.UnimplementedDarksideStreamerServer
}

// NewDarksideStreamer constructs a gRPC context for darksidewalletd.
func NewDarksideStreamer(cache *common.BlockCache) (walletrpc.DarksideStreamerServer, error) {
	return &DarksideStreamer{cache: cache}, nil
}

// Test to make sure Address is a single t address
func checkTaddress(taddr string) error {
	match, err := regexp.Match("\\At[a-zA-Z0-9]{34}\\z", []byte(taddr))
	if err != nil || !match {
		return errors.New("Invalid address")
	}
	return nil
}

func (s *lwdStreamer) peerIPFromContext(ctx context.Context) string {
	if xRealIP, ok := metadata.FromIncomingContext(ctx); ok {
		realIP := xRealIP.Get("x-real-ip")
		if len(realIP) > 0 {
			return realIP[0]
		}
	}

	if peerInfo, ok := peer.FromContext(ctx); ok {
		ip, _, err := net.SplitHostPort(peerInfo.Addr.String())
		if err == nil {
			return ip
		}
	}

	return "unknown"
}

func (s *lwdStreamer) dailyActiveBlock(height uint64, peerip string) {
	if height%1152 == 0 {
		common.Log.WithFields(logrus.Fields{
			"method":       "DailyActiveBlock",
			"peer_addr":    peerip,
			"block_height": height,
		}).Info("Service")
	}
}

// GetLatestBlock returns the height of the best chain, according to zcashd.
func (s *lwdStreamer) GetLatestBlock(ctx context.Context, placeholder *walletrpc.ChainSpec) (*walletrpc.BlockID, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// GetTaddressTxids is a streaming RPC that returns transaction IDs that have
// the given transparent address (taddr) as either an input or output.
func (s *lwdStreamer) GetTaddressTxids(addressBlockFilter *walletrpc.TransparentAddressBlockFilter, resp walletrpc.CompactTxStreamer_GetTaddressTxidsServer) error {
	return errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// GetBlock returns the compact block at the requested height. Requesting a
// block by hash is not yet supported.
func (s *lwdStreamer) GetBlock(ctx context.Context, id *walletrpc.BlockID) (*walletrpc.CompactBlock, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// GetTreeState returns the note commitment tree state corresponding to the given block.
// See section 3.7 of the Zcash protocol specification. It returns several other useful
// values also (even though they can be obtained using GetBlock).
// The block can be specified by either height or hash.
func (s *lwdStreamer) GetTreeState(ctx context.Context, id *walletrpc.BlockID) (*walletrpc.TreeState, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// GetTransaction returns the raw transaction bytes that are returned
// by the zcashd 'getrawtransaction' RPC.
func (s *lwdStreamer) GetTransaction(ctx context.Context, txf *walletrpc.TxFilter) (*walletrpc.RawTransaction, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// GetLightdInfo gets the LightWalletD (this server) info, and includes information
// it gets from its backend zcashd.
func (s *lwdStreamer) GetLightdInfo(ctx context.Context, in *walletrpc.Empty) (*walletrpc.LightdInfo, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// SendTransaction forwards raw transaction bytes to a zcashd instance over JSON-RPC
func (s *lwdStreamer) SendTransaction(ctx context.Context, rawtx *walletrpc.RawTransaction) (*walletrpc.SendResponse, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

func getTaddressBalanceZcashdRpc(addressList []string) (*walletrpc.Balance, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// GetTaddressBalance returns the total balance for a list of taddrs
func (s *lwdStreamer) GetTaddressBalance(ctx context.Context, addresses *walletrpc.AddressList) (*walletrpc.Balance, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// GetTaddressBalanceStream returns the total balance for a list of taddrs
func (s *lwdStreamer) GetTaddressBalanceStream(addresses walletrpc.CompactTxStreamer_GetTaddressBalanceStreamServer) error {
	return errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// Key is 32-byte txid (as a 64-character string), data is pointer to compact tx.
var mempoolMap *map[string]*walletrpc.CompactTx
var mempoolList []string

// Last time we pulled a copy of the mempool from zcashd.
var lastMempool time.Time

func (s *lwdStreamer) GetMempoolTx(exclude *walletrpc.Exclude, resp walletrpc.CompactTxStreamer_GetMempoolTxServer) error {
	return errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// Return the subset of items that aren't excluded, but
// if more than one item matches an exclude entry, return
// all those items.
func MempoolFilter(items, exclude []string) []string {
	sort.Slice(items, func(i, j int) bool {
		return items[i] < items[j]
	})
	sort.Slice(exclude, func(i, j int) bool {
		return exclude[i] < exclude[j]
	})
	// Determine how many items match each exclude item.
	nmatches := make([]int, len(exclude))
	// is the exclude string less than the item string?
	lessthan := func(e, i string) bool {
		l := len(e)
		if l > len(i) {
			l = len(i)
		}
		return e < i[0:l]
	}
	ei := 0
	for _, item := range items {
		for ei < len(exclude) && lessthan(exclude[ei], item) {
			ei++
		}
		match := ei < len(exclude) && strings.HasPrefix(item, exclude[ei])
		if match {
			nmatches[ei]++
		}
	}

	// Add each item that isn't uniquely excluded to the results.
	tosend := make([]string, 0)
	ei = 0
	for _, item := range items {
		for ei < len(exclude) && lessthan(exclude[ei], item) {
			ei++
		}
		match := ei < len(exclude) && strings.HasPrefix(item, exclude[ei])
		if !match || nmatches[ei] > 1 {
			tosend = append(tosend, item)
		}
	}
	return tosend
}

func getAddressUtxos(arg *walletrpc.GetAddressUtxosArg, f func(*walletrpc.GetAddressUtxosReply) error) error {
	return errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

func (s *lwdStreamer) GetAddressUtxos(ctx context.Context, arg *walletrpc.GetAddressUtxosArg) (*walletrpc.GetAddressUtxosReplyList, error) {
	return nil, errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

func (s *lwdStreamer) GetAddressUtxosStream(arg *walletrpc.GetAddressUtxosArg, resp walletrpc.CompactTxStreamer_GetAddressUtxosStreamServer) error {
	return errors.New("Your version of Zecwallet Lite is deprecated. Please download the latest version from www.zecwallet.co")
}

// This rpc is used only for testing.
var concurrent int64

func (s *lwdStreamer) Ping(ctx context.Context, in *walletrpc.Duration) (*walletrpc.PingResponse, error) {
	// This gRPC allows the client to create an arbitrary number of
	// concurrent threads, which could run the server out of resources,
	// so only allow if explicitly enabled.
	if !s.pingEnable {
		return nil, errors.New("Ping not enabled, start lightwalletd with --ping-very-insecure")
	}
	var response walletrpc.PingResponse
	response.Entry = atomic.AddInt64(&concurrent, 1)
	time.Sleep(time.Duration(in.IntervalUs) * time.Microsecond)
	response.Exit = atomic.AddInt64(&concurrent, -1)
	return &response, nil
}

// SetMetaState lets the test driver control some GetLightdInfo values.
func (s *DarksideStreamer) Reset(ctx context.Context, ms *walletrpc.DarksideMetaState) (*walletrpc.Empty, error) {
	match, err := regexp.Match("\\A[a-fA-F0-9]+\\z", []byte(ms.BranchID))
	if err != nil || !match {
		return nil, errors.New("Invalid branch ID")
	}

	match, err = regexp.Match("\\A[a-zA-Z0-9]+\\z", []byte(ms.ChainName))
	if err != nil || !match {
		return nil, errors.New("Invalid chain name")
	}
	err = common.DarksideReset(int(ms.SaplingActivation), ms.BranchID, ms.ChainName)
	if err != nil {
		return nil, err
	}
	mempoolMap = nil
	mempoolList = nil
	return &walletrpc.Empty{}, nil
}

// StageBlocksStream accepts a list of blocks from the wallet test code,
// and makes them available to present from the mock zcashd's GetBlock rpc.
func (s *DarksideStreamer) StageBlocksStream(blocks walletrpc.DarksideStreamer_StageBlocksStreamServer) error {
	for {
		b, err := blocks.Recv()
		if err == io.EOF {
			blocks.SendAndClose(&walletrpc.Empty{})
			return nil
		}
		if err != nil {
			return err
		}
		common.DarksideStageBlockStream(b.Block)
	}
}

// StageBlocks loads blocks from the given URL to the staging area.
func (s *DarksideStreamer) StageBlocks(ctx context.Context, u *walletrpc.DarksideBlocksURL) (*walletrpc.Empty, error) {
	if err := common.DarksideStageBlocks(u.Url); err != nil {
		return nil, err
	}
	return &walletrpc.Empty{}, nil
}

// StageBlocksCreate stages a set of synthetic (manufactured on the fly) blocks.
func (s *DarksideStreamer) StageBlocksCreate(ctx context.Context, e *walletrpc.DarksideEmptyBlocks) (*walletrpc.Empty, error) {
	if err := common.DarksideStageBlocksCreate(e.Height, e.Nonce, e.Count); err != nil {
		return nil, err
	}
	return &walletrpc.Empty{}, nil
}

// StageTransactionsStream adds the given transactions to the staging area.
func (s *DarksideStreamer) StageTransactionsStream(tx walletrpc.DarksideStreamer_StageTransactionsStreamServer) error {
	// My current thinking is that this should take a JSON array of {height, txid}, store them,
	// then DarksideAddBlock() would "inject" transactions into blocks as its storing
	// them (remembering to update the header so the block hash changes).
	for {
		transaction, err := tx.Recv()
		if err == io.EOF {
			tx.SendAndClose(&walletrpc.Empty{})
			return nil
		}
		if err != nil {
			return err
		}
		err = common.DarksideStageTransaction(int(transaction.Height), transaction.Data)
		if err != nil {
			return err
		}
	}
}

// StageTransactions loads blocks from the given URL to the staging area.
func (s *DarksideStreamer) StageTransactions(ctx context.Context, u *walletrpc.DarksideTransactionsURL) (*walletrpc.Empty, error) {
	if err := common.DarksideStageTransactionsURL(int(u.Height), u.Url); err != nil {
		return nil, err
	}
	return &walletrpc.Empty{}, nil
}

// ApplyStaged merges all staged transactions into staged blocks and all staged blocks into the active blockchain.
func (s *DarksideStreamer) ApplyStaged(ctx context.Context, h *walletrpc.DarksideHeight) (*walletrpc.Empty, error) {
	return &walletrpc.Empty{}, common.DarksideApplyStaged(int(h.Height))
}

// GetIncomingTransactions returns the transactions that were submitted via SendTransaction().
func (s *DarksideStreamer) GetIncomingTransactions(in *walletrpc.Empty, resp walletrpc.DarksideStreamer_GetIncomingTransactionsServer) error {
	// Get all of the incoming transactions we're received via SendTransaction()
	for _, txBytes := range common.DarksideGetIncomingTransactions() {
		err := resp.Send(&walletrpc.RawTransaction{Data: txBytes, Height: 0})
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearIncomingTransactions empties the incoming transaction list.
func (s *DarksideStreamer) ClearIncomingTransactions(ctx context.Context, e *walletrpc.Empty) (*walletrpc.Empty, error) {
	common.DarksideClearIncomingTransactions()
	return &walletrpc.Empty{}, nil
}
