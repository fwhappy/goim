package hall

import (
	"goim/core"
	"net"
	"time"

	"github.com/fwhappy/util"
)

// ListenHandShakeTimeout 监听用户连接之后，有没有及时handshake
// 如果用户执行了连接之后，30秒内都没有handshake成功，则将用户踢下线
func ListenHandShakeTimeout(conn *net.TCPConn, c chan int) {
	// 捕获异常
	defer util.RecoverPanic()

	select {
	case value := <-c:
		core.Logger.Debug("用户握手成功或断开了连接，退出监听:%v,remote:%v", value, conn.RemoteAddr().String())
		break
	case <-time.After(10 * time.Second):
		core.Logger.Debug("用户长时间未handshake成功，断开连接:%s", conn.RemoteAddr().String())
		conn.Close()
		break
	}
}
