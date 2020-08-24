package main

import (
	"log"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

/*
(gdb) frame 1
#1  0x00005555555576ec in start (argc=3, argv=0x7fffffffeb48) at fileio.c:408

(gdb) print SCST_USER_REGISTER_DEVICE
$3 = 1084257537


(gdb) print desc
$1 = {

	version_str = 93824992274528,
	license_str = 93824992274524,
	type = 0 '\000',
	sgv_shared = 0 '\000',
	sgv_disable_clustered_pool = 0 '\000',
	sgv_single_alloc_pages = 0,
	sgv_purge_interval = 0,
 	opt = {
		 parse_type = 0 '\000',
		 on_free_cmd_type = 1 '\001',
		 memory_reuse_type = 3 '\003',
		 partial_transfers_type = 0 '\000',
		 partial_len = 0,
		 tst = 1 '\001',
		 tmf_only = 0 '\000',
		 queue_alg = 1 '\001',
		 qerr = 0 '\000',
		 tas = 0 '\000',
		 swp = 0 '\000',
		 d_sense = 0 '\000',
		 has_own_order_mgmt = 0 '\000',
		 ext_copy_remap_supported = 1 '\001'
	},
  block_size = 512,

  enable_pr_cmds_notifications = 0 '\000',
  name = "net3",
  '\000' <repeats 45 times>,
  sgv_name = '\000' <repeats 49 times>}

*/

func registerDevice(name string) (int, error) {
	fd, err := unix.Open("/dev/scst_user", unix.O_RDWR, 0)
	if err != nil {
		log.Fatal("Starting Error; /dev/scst_user -> ", err.Error())
	}

	gplString := [4]byte{'G', 'P', 'L', 0}
	// Haha I'm going to hard code what ever git version I was using at the time, is that okay? haha nvm I wasn't asking
	verString := [100]byte{'3', '.', '5', '.', '0', '-', 'p', 'r', 'e', 'b', '0', '3', '8', '9', '3', 'a', '2', '6', '0',
		'c', '0', '8', 'b', '8', '7', 'c', 'b', 'c', '0', 'a', '0', 'f', '8', 'c', '3', 'b', '3', '4', '7', 'd', '4', '8',
		'1', 'b', 'e', 'e', '1', 'd', '7', 'c', 'b', 'b', 'a', '8', '0', '0', '3', '7', 'f', '3', 'f', '5', '7', '2', 'c',
		'4', '2', 'a', 'f', '5', '3', 'c', '5', 'a', '3', 'e', '5', '6', 'e', 'b', '3', 'c', '4', 'f', '5', '2', '3', 'c', '0', 0x00}

	var nameBytes [50]byte
	copy(nameBytes[:], []byte(name))

	def := raw_scst_user_dev_desc{
		Version_str: &verString[0],
		License_str: &gplString[0],
		scst_user_opt: raw_scst_user_opt{
			on_free_cmd_type:         1,
			memory_reuse_type:        3,
			tst:                      1,
			queue_alg:                1,
			ext_copy_remap_supported: 1,
		},
		block_size: 512,
		name:       nameBytes,
		sgv_name:   [50]byte{0},
	}

	return fd, SCST_USER_REGISTER_DEVICE(fd, &def)
}

func SCST_USER_REGISTER_DEVICE(fd int, def *raw_scst_user_dev_desc) error {
	err := ioctl(fd, 1084257537, uintptr(unsafe.Pointer(def)))
	if *debugLogs {
		log.Printf("eee %v", err)
	}
	return err
}

func SCST_USER_REPLY_AND_GET_CMD(fd int, def *raw_scst_user_get_cmd_preply) error {
	tmp := uintptr(unsafe.Pointer(def))
	for {
		err := ioctl(fd, 3256907013, tmp)
		if *debugLogs {
			log.Printf("ooo %v", err)
		}
		if err != nil {
			if *debugLogs {
				log.Printf("======================================= %v =======================================", err)
			}
			time.Sleep(time.Second)
		}
		if err == errEINTR {
			continue
		}
		return err
	}
}

func SCST_USER_REPLY_MEM_ALLOC(fd int, def *raw_scst_user_alloc_reply) error {
	tmp := uintptr(unsafe.Pointer(def))
	for {
		err := ioctl(fd, 3256907013, tmp)
		if *debugLogs {
			log.Printf("ooo %v", err)
		}
		if err != nil {
			if *debugLogs {
				log.Printf("======================================= %v =======================================", err)
			}
			time.Sleep(time.Second)
		}
		if err == errEINTR {
			continue
		}
		return err
	}
}

func SCST_USER_REPLY_AND_GET_CMD_ON_EXEC(fd int, def *raw_scst_user_get_cmd_scsi_cmd_exec) error {
	err := ioctl(fd, 3256907013, uintptr(unsafe.Pointer(def)))
	if *debugLogs {
		log.Printf("ooo %v", err)
	}
	return err
}

func ioctl(fd int, req uint, arg uintptr) (err error) {
	_, _, e1 := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(arg))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

var (
	errEAGAIN error = syscall.EAGAIN
	errEINVAL error = syscall.EINVAL
	errENOENT error = syscall.ENOENT
	errEINTR  error = syscall.EINTR
)

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case unix.EAGAIN:
		return errEAGAIN
	case unix.EINVAL:
		return errEINVAL
	case unix.ENOENT:
		return errENOENT
	case unix.EINTR:
		return errEINTR
	}
	return e
}
