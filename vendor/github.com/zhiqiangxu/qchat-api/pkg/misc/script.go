package misc

import (
	"io/ioutil"

	"github.com/go-redis/redis"
	"github.com/zhiqiangxu/qrpc"
)

const (
	dump = true
)

// RequireCSScript for require cs script
func RequireCSScript(script string) *redis.Script {
	src := RequireFile(
		"pkg/service/asset/common_const.lua",
		"pkg/service/asset/common_func.lua",
		"pkg/service/asset/common_func_redis_only.lua",
		"pkg/service/asset/model/keys.lua",
		"pkg/service/asset/model/cs_online.lua",
		"pkg/service/asset/model/timer.lua",
		"pkg/service/asset/model/cs_session_expire.lua",
		"pkg/service/asset/model/cs_session_last_serial.lua",
		"pkg/service/asset/model/cs_session_q.lua",
		"pkg/service/asset/model/cs_session.lua",
		"pkg/service/asset/model/user_cs.lua",
		"pkg/service/asset/model/cs_user.lua",
		"pkg/service/asset/model/cs_last_cs.lua",
		"pkg/service/asset/model/signal_room_serial.lua",
		"pkg/service/asset/model/signal_room.lua",
		"pkg/service/asset/model/signal_user.lua",
		"pkg/service/asset/service/cs_auto_pick_session_q.lua",
		"pkg/service/asset/service/cs_end_session.lua",
		"pkg/service/asset/service/inner_login.lua",
		"pkg/service/asset/service/cs_manual_pick_session_q.lua",
		"pkg/service/asset/service/cs_allocate_csuid.lua",
		"pkg/service/asset/service/cs_set_status.lua",
		"pkg/service/asset/service/cs_switch_session.lua",
		"pkg/service/asset/service/signal_stop_dial.lua",
		"pkg/service/asset/service/inner_logout.lua",
		"pkg/service/asset/script/"+script,
	)

	if dump {
		ioutil.WriteFile("/tmp/"+script, qrpc.Slice(src), 0644)
	}
	return redis.NewScript(src)
}

// RequireMsgScript for require msg script
func RequireMsgScript(script string) *redis.Script {
	src := RequireFile(
		"pkg/service/asset/common_const.lua",
		"pkg/service/asset/common_func.lua",
		"pkg/service/asset/model/last_nonce.lua",
		"pkg/service/asset/model/msg_ids.lua",
		"pkg/service/asset/model/msg_snippet.lua",
		"pkg/service/asset/model/msg_thread.lua",
		"pkg/service/asset/model/session.lua",
		"pkg/service/asset/model/keys.lua",
		"pkg/service/asset/model/room.lua",
		"pkg/service/asset/model/room_member.lua",
		"pkg/service/asset/model/room_ban.lua",
		"pkg/service/asset/model/room_gag.lua",
		"pkg/service/asset/model/group.lua",
		"pkg/service/asset/model/group_member.lua",
		"pkg/service/asset/model/group_gag.lua",
		"pkg/service/asset/service/group_drain_delete.lua",
		"pkg/service/asset/service/msg_send_body.lua",
		"pkg/service/asset/script/"+script,
	)
	if dump {
		ioutil.WriteFile("/tmp/"+script, qrpc.Slice(src), 0644)
	}
	return redis.NewScript(src)
}
