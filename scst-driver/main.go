package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGURG)
	go func() {
		for aaa := range c {
			log.Printf("SAVED US FROM EXPLOSION? THANKS I GUESS %v", aaa)
		}
	}()

	go startTap()

	fd, err := registerDevice("net3")
	if err != nil {
		log.Fatalf("Failed to register device: %v",
			err)
	}

	go pollForStuff(fd, "net3")

	fd2, err := registerDevice("net4")
	if err != nil {
		log.Fatalf("Failed to register device: %v",
			err)
	}

	pollForStuff(fd2, "net4")
}

var upcomingBug = false

func trap_me() {
	log.Print("trap")
}

func pollForStuff(fd int, Dirtype string) interface{} {
	def := raw_scst_user_get_cmd_preply{}
	ticker := time.NewTicker(time.Second)
	instance := scstInstance{
		globalOutputBufAlign: -1,
		ticker:               ticker.C,
		logger:               log.New(os.Stdout, fmt.Sprintf("[%s] ", Dirtype), log.Ltime),
	}

	if instance.antiGCBufferStorage == nil {
		instance.antiGCBufferStorage = make(map[int][]byte)
	}

	for {

		instance.logger.Printf("ioctl")

		if upcomingBug {
			trap_me()
		}

		SCST_USER_REPLY_AND_GET_CMD(fd,
			&def)

		log.Printf("Entering Switch")
		def.preply = 0

		switch def.subcode {
		case SCST_USER_EXEC:
			// TODO: This is the real biz
			instance.logger.Printf("SCST_USER_EXEC")
			// processExecCmd(raw_scst_user_get_cmd_preply)
			lol := (*raw_scst_user_get_cmd_scsi_cmd_exec)(unsafe.Pointer(&def))
			// log.Printf("SCST_USER_EXEC -> %#v", lol)

			reply := instance.processExecCmd(lol)
			if reply.pbuf != nil {
				instance.logger.Printf("First byte = %#v", *reply.pbuf)
			}
			def.preply = uintptr(unsafe.Pointer(reply))
			// def = def2
			// def_exec = lol

		case SCST_USER_ALLOC_MEM:
			// TODO:

			//
			instance.logger.Printf("SCST_USER_ALLOC_MEM")

			instance.logger.Printf("The module wishes for more memory sir.")
			instance.buffersMade++

			lol := (*scst_user_scsi_cmd_alloc_mem)(unsafe.Pointer(&def))
			instance.logger.Printf("%#v", lol)

			aaa := make([]byte, 2*(1024*1024))

			finalOutputOffset := alignTheBuffer(uintptr(unsafe.Pointer(&aaa[0])))

			instance.antiGCBufferStorage[instance.buffersMade] = aaa
			instance.currentpbuf = aaa[finalOutputOffset:]
			memReply := raw_scst_user_alloc_reply{
				cmd_h:   lol.cmd_h,
				subcode: lol.subcode,
				preply:  &aaa[finalOutputOffset],
			}
			def.preply = uintptr(unsafe.Pointer(&memReply))

			// raw_scst_user_alloc_reply.preply =
		case SCST_USER_PARSE:
			// TODO:
			instance.logger.Printf("SCST_USER_PARSE")
		case SCST_USER_ON_CACHED_MEM_FREE:
			// TODO:
			instance.logger.Printf("SCST_USER_ON_CACHED_MEM_FREE")
		case SCST_USER_ON_FREE_CMD:
			// TODO:
			instance.logger.Printf("SCST_USER_ON_FREE_CMD")
		case SCST_USER_TASK_MGMT_RECEIVED:
			// TODO:
			instance.logger.Printf("SCST_USER_TASK_MGMT_RECEIVED")
		case SCST_USER_TASK_MGMT_DONE:
			// TODO:
			instance.logger.Printf("SCST_USER_TASK_MGMT_DONE")
		case SCST_USER_ATTACH_SESS:
			// TODO: Apparently we don't need to do anything for this.
			instance.logger.Printf("SCST_USER_ATTACH_SESS")
			lol := (*raw_scst_user_get_cmd_scst_user_sess)(unsafe.Pointer(&def))
			instance.logger.Printf("%#v", lol)
			reply := raw_scst_user_reply_cmd_result{
				cmd_h:   def.cmd_h,
				subcode: def.subcode,
				result:  0,
			}

			def.preply = uintptr(unsafe.Pointer(&reply))
		case SCST_USER_DETACH_SESS:
			// TODO: Apparently this is where the interesting stuff happens???
			instance.logger.Printf("SCST_USER_DETACH_SESS")
			lol := (*raw_scst_user_get_cmd_scst_user_sess)(unsafe.Pointer(&def))
			instance.logger.Printf("%#v", lol)

			reply := raw_scst_user_reply_cmd_result{
				cmd_h:   def.cmd_h,
				subcode: def.subcode,
				result:  0,
			}

			def.preply = uintptr(unsafe.Pointer(&reply))
		}

		// time.Sleep(time.Second)
	}
	/*
	   (gdb) print cmd
	   $1 = {cmd_h = 0, subcode = 0, {preply = 0, sess = {sess_h = 0, lun = 0, threads_num = 0, rd_ {etc etc}
	*/

	/* haha this struct would have just changed behind our backs this is so unholy */

	/*

		   So handling this is a huge PITA because it's a union, the only other interface I've
		   seen that does this is TIPC, thankfully there is golang support from that so we can just
		   steal how they deal with it:

		   https://go.googlesource.com/sys/+/master/unix/syscall_linux.go#1056

		   	switch pp.Addrtype {
				case TIPC_SERVICE_RANGE:
					sa.Addr = (*TIPCServiceRange)(unsafe.Pointer(&pp.Addr))
				case TIPC_SERVICE_ADDR:
					sa.Addr = (*TIPCServiceName)(unsafe.Pointer(&pp.Addr))
				case TIPC_SOCKET_ADDR:
					sa.Addr = (*TIPCSocketAddr)(unsafe.Pointer(&pp.Addr))
				default:
					return nil, EINVAL
			}

	*/

	/*
		   (gdb) print cmd
		   $1 = {
			cmd_h = 0,
		   	subcode = 2182640416,
		   	{
				   preply = 18446614906000855744,
		   		   sess = {
						  sess_h = 18446614906000855744,
		   				  lun = 3,
		   	              threads_num = 0,
		   	              rd_only = 1 '\001',
		                  scsi_transport_version = 2304,
		                  phys_transport_version = 3488,
						  initiator_name = "50:01:43:80:21:de:b5:d6",'\000' <repeats 232 times>,
						  target_name = "50:01:43:80:26:68:8e:5c", '\000' <repeats 232 times>
					},
					<... the rest is garbage beacuse of the union?>}

	*/

	// log.Printf("Holy shit %#v", def)
	// log.Printf("Holy -> %#v", string(def.initiator_name[:]))
	// log.Printf("shit -> %#v", string(def.target_name[:]))

}
