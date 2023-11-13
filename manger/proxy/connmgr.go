package proxy

import (
	"fmt"
	"github.com/nothollyhigh/kiss/log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ConnMgr = &ConnManger{
		InNum:         0,
		OutNum:        0,
		SuccessNum:    0,
		FailedNum:     0,
		ClientInSize:  0,
		ClientOutSize: 0,
		ServerInSize:  0,
		ServerOutSize: 0,
	}
	AccountMgr = sync.Map{} // 保存 平台账号(key) 服务地址(value)
	ClientMgr  = sync.Map{} // 保存 客户端地址(key) 服务地址(value)
)

type ConnManger struct {
	// sync.Mutex
	InNum         int64 /* 当前客户端连接数 */
	OutNum        int64 /* 当前服务端连接数 */
	SuccessNum    int64 /* 启动以来隧道成功总数 */
	FailedNum     int64 /* 启动以来隧道失败总数 */
	ClientInSize  int64 /* 启动以来客户端读总流量 */
	ClientOutSize int64 /* 启动以来客户端写总流量 */
	ServerInSize  int64 /* 启动以来服务端读总流量 */
	ServerOutSize int64 /* 启动以来服务端写总流量 */
}

// UpdateInNum /* 更新当前客户端连接数 */
func (mgr *ConnManger) UpdateInNum(delta int64) {
	atomic.AddInt64(&mgr.InNum, delta)
}

// GetInNum /* 获取当前客户端连接数 */
func (mgr *ConnManger) GetInNum() int64 {
	return atomic.LoadInt64(&mgr.InNum)
}

// UpdateOutNum /* 更新当前服务端连接数 */
func (mgr *ConnManger) UpdateOutNum(delta int64) {
	atomic.AddInt64(&mgr.OutNum, delta)
}

// GetOutNum /* 获取当前服务端连接数 */
func (mgr *ConnManger) GetOutNum() int64 {
	return atomic.LoadInt64(&mgr.OutNum)
}

// UpdateSuccessNum /* 更新启动以来隧道成功总数 */
func (mgr *ConnManger) UpdateSuccessNum(delta int64) {
	atomic.AddInt64(&mgr.SuccessNum, delta)
}

// GetSuccessNum /* 获取启动以来隧道成功总数 */
func (mgr *ConnManger) GetSuccessNum() int64 {
	return atomic.LoadInt64(&mgr.SuccessNum)
}

// UpdateFailedNum /* 更新启动以来隧道失败总数 */
func (mgr *ConnManger) UpdateFailedNum(delta int64) {
	atomic.AddInt64(&mgr.FailedNum, delta)
}

// GetdateFailedNum /* 获取启动以来隧道失败总数 */
func (mgr *ConnManger) GetdateFailedNum() int64 {
	return atomic.LoadInt64(&mgr.FailedNum)
}

// UpdateClientInSize /* 更新启动以来客户端读总流量 */
func (mgr *ConnManger) UpdateClientInSize(delta int64) {
	atomic.AddInt64(&mgr.ClientInSize, delta)
}

// GetClientInSize /* 获取启动以来客户端读总流量 */
func (mgr *ConnManger) GetClientInSize() int64 {
	return atomic.LoadInt64(&mgr.ClientInSize)
}

// UpdateClientOutSize /* 更新启动以来客户端写总流量 */
func (mgr *ConnManger) UpdateClientOutSize(delta int64) {
	atomic.AddInt64(&mgr.ClientOutSize, delta)
}

// GetClientOutSize /* 获取启动以来客户端写总流量 */
func (mgr *ConnManger) GetClientOutSize() int64 {
	return atomic.LoadInt64(&mgr.ClientOutSize)
}

// UpdateServerInSize /* 更新启动以来服务端读总流量 */
func (mgr *ConnManger) UpdateServerInSize(delta int64) {
	atomic.AddInt64(&mgr.ServerInSize, delta)
}

// GetServerInSize /* 获取启动以来服务端读总流量 */
func (mgr *ConnManger) GetServerInSize() int64 {
	return atomic.LoadInt64(&mgr.ServerInSize)
}

// UpdateServerOutSize /* 更新启动以来服务端写总流量 */
func (mgr *ConnManger) UpdateServerOutSize(delta int64) {
	atomic.AddInt64(&mgr.ServerOutSize, delta)
}

// GetServerOutSize /* 获取启动以来服务端写总流量 */
func (mgr *ConnManger) GetServerOutSize() int64 {
	return atomic.LoadInt64(&mgr.ServerOutSize)
}

func (mgr *ConnManger) LogDataFlowRecord() string {
	return fmt.Sprintf("ClientReadSize: %d B, ClientWriteSize: %d B, ServerReadSize: %d B, ServerWriteSize: %d B)",
		mgr.GetClientInSize(), mgr.GetClientOutSize(),
		mgr.GetServerInSize(), mgr.GetServerOutSize())
}

func (mgr *ConnManger) StartDataFlowRecord(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			log.Warn(mgr.LogDataFlowRecord())
			/*log.Info("ConnManger DataInfo(CR: %d M, CW: %d M, SR: %d M, SW: %d M)",
			mgr.GetClientInSize()/(1024*1024), mgr.GetClientOutSize()/(1024*1024),
			mgr.GetServerInSize()/(1024*1024), mgr.GetServerOutSize()/(1024*1024))*/
		}
	}()
}
