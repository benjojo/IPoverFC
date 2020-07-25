package main

const (
	SCST_USER_PARSE_STANDARD  = 0
	SCST_USER_PARSE_CALL      = 1
	SCST_USER_PARSE_EXCEPTION = 2
	SCST_USER_MAX_PARSE_OPT   = SCST_USER_PARSE_EXCEPTION

	SCST_USER_ON_FREE_CMD_CALL    = 0
	SCST_USER_ON_FREE_CMD_IGNORE  = 1
	SCST_USER_MAX_ON_FREE_CMD_OPT = SCST_USER_ON_FREE_CMD_IGNORE

	SCST_USER_MEM_NO_REUSE      = 0
	SCST_USER_MEM_REUSE_READ    = 1
	SCST_USER_MEM_REUSE_WRITE   = 2
	SCST_USER_MEM_REUSE_ALL     = 3
	SCST_USER_MAX_MEM_REUSE_OPT = SCST_USER_MEM_REUSE_ALL

	SCST_USER_PARTIAL_TRANSFERS_NOT_SUPPORTED     = 0
	SCST_USER_PARTIAL_TRANSFERS_SUPPORTED_ORDERED = 1
	SCST_USER_PARTIAL_TRANSFERS_SUPPORTED         = 2
	SCST_USER_MAX_PARTIAL_TRANSFERS_OPT           = SCST_USER_PARTIAL_TRANSFERS_SUPPORTED
)

/*
struct scst_user_dev_desc {
	aligned_u64 version_str;
	aligned_u64 license_str;
	uint8_t type;
	uint8_t sgv_shared;
	uint8_t sgv_disable_clustered_pool;
	int32_t sgv_single_alloc_pages;
	int32_t sgv_purge_interval;
	struct scst_user_opt opt;
	uint32_t block_size;
	uint8_t enable_pr_cmds_notifications;
	char name[SCST_MAX_NAME];
	char sgv_name[SCST_MAX_NAME];
};
*/

type raw_scst_user_dev_desc struct {
	Version_str                  uintptr
	License_str                  uintptr
	stype                        uint8
	sgv_shared                   uint8
	sgv_disable_clustered_pool   uint8
	sgv_single_alloc_pages       int32
	sgv_purge_interval           int32
	scst_user_opt                raw_scst_user_opt
	block_size                   int32
	enable_pr_cmds_notifications uint8
	name                         [50]byte
	sgv_name                     [50]byte
	// char name[SCST_MAX_NAME];
	// char sgv_name[SCST_MAX_NAME];
}

/*
(gdb) frame 1
#1  0x00005555555576ec in start (argc=3, argv=0x7fffffffeb48) at fileio.c:408

(gdb) print SCST_USER_REGISTER_DEVICE
$3 = 1084257537


(gdb) print desc
$1 = {
	version_str = 93824992274528, // POINTERS
	license_str = 93824992274524, // POINTERS
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

/*
struct scst_user_opt {
	uint8_t parse_type;
	uint8_t on_free_cmd_type;
	uint8_t memory_reuse_type;
	uint8_t partial_transfers_type;
	int32_t partial_len;

	// SCSI control mode page parameters, see SPC
	uint8_t tst;
	uint8_t tmf_only;
	uint8_t queue_alg;
	uint8_t qerr;
	uint8_t tas;
	uint8_t swp;
	uint8_t d_sense;

	uint8_t has_own_order_mgmt;

	uint8_t ext_copy_remap_supported;
};
*/

type raw_scst_user_opt struct {
	parse_type             uint8
	on_free_cmd_type       uint8
	memory_reuse_type      uint8
	partial_transfers_type uint8
	partial_len            int32

	// SCSI control mode page parameters, see SPC
	tst       uint8
	tmf_only  uint8
	queue_alg uint8
	qerr      uint8
	tas       uint8
	swp       uint8
	d_sense   uint8

	has_own_order_mgmt uint8

	ext_copy_remap_supported uint8
}

type raw_scst_user_get_cmd struct {
	cmd_h   uint32
	subcode uint32
}

type raw_scst_user_get_cmd_preply struct {
	cmd_h   uint32
	subcode uint32
	preply  uintptr         // Pointer to a reply
	padding [16 * 1024]byte // I'm scared of the kernel, I want to not accidently overshoot
}

type raw_scst_user_get_cmd_scst_user_sess struct {
	cmd_h                  uint32
	subcode                uint32
	sess_h                 uint64
	lun                    uint64
	threads_num            uint16
	rd_only                uint8
	scsi_transport_version uint16
	phys_transport_version uint16
	initiator_name         [256]byte
	target_name            [256]byte
	padding                [64]byte // I'm scared of the kernel, I want to not accidently overshoot
}

/*
struct scst_user_get_cmd
{
	uint32_t cmd_h;
	uint32_t subcode;
	union {
		uint64_t preply;
		struct scst_user_sess sess;
		struct scst_user_scsi_cmd_parse parse_cmd;
		struct scst_user_scsi_cmd_alloc_mem alloc_cmd;
		struct scst_user_scsi_cmd_exec exec_cmd;
		struct scst_user_scsi_on_free_cmd on_free_cmd;
		struct scst_user_on_cached_mem_free on_cached_mem_free;
		struct scst_user_tm tm_cmd;
	};
}
*/

type raw_scst_user_sess struct {
	sess_h                 uint64
	lun                    uint64
	threads_num            uint16
	rd_only                uint8
	scsi_transport_version uint16
	phys_transport_version uint16
	initiator_name         [50]byte
	target_name            [50]byte
}

/*
struct scst_user_sess
{
	uint64_t sess_h;
	uint64_t lun;
	uint16_t threads_num;
	uint8_t rd_only;
	uint16_t scsi_transport_version;
	uint16_t phys_transport_version;
	char initiator_name[SCST_MAX_NAME];
	char target_name[SCST_MAX_NAME];
},
*/

type raw_scst_user_get_cmd_scsi_cmd_exec struct {
	cmd_h   uint32
	subcode uint32
	// sess_h - corresponding session handler
	sess_h int64
	// cdb - SCSI CDB
	cdb [16]byte
	// cdb_len - SCSI CDB length
	cdb_len uint16
	// lba - LBA of the command, if any
	lba int64
	// data_len - command's data length. Could be different from buen for commands like VERIFY,  which transfer different amount of data, than process, or even none of them
	data_len int64
	// bufflen - command's buffer length
	bufflen int32
	// alloc_len - command's buffer length, which should be allocated, if pbuf is 0 and the command requires data transfer
	alloc_len int32
	// pbuf - pointer to command's data buffer or 0 for SCSI commands without data transfer.
	pbuf uintptr
	// queue_type - SCSI task attribute (queue type)
	queue_type uint8
	// data_direction - command's data ow direction, one of SCST_DATA_* constants
	data_direction uint8
	// partial - species, if the command is a partial subcommand, could have the following OR'ed ags:
	//  SCST_USER_SUBCOMMAND - set if the command is a partial subcommand
	//  SCST_USER_SUBCOMMAND_FINAL - set if the subcommand is a nal one
	partial uint8
	// timeout - CDB execution timeout
	timeout int32
	// p_out_buf - for bidirectional commands pointer on command's OUT, i.e. from initiator to target,
	// data buffer or 0 for SCSI commands without data transfer
	p_out_buf uintptr
	// out_bufflen - for bidirectional commands command's OUT, i.e. from initiator to target, buffer length
	out_bufflen int32
	// sn - command's SN, which might be used for task management
	sn uint32
	// parent_cmd_h - has the same unique value for all partial data transfers subcommands of one original
	// (parent) command
	parent_cmd_h uint32
	// parent_cmd_data_len - for partial data transfers subcommand has the size of the overall data
	// transfer of the original (parent) command
	parent_cmd_data_len int32
	// partial_offset - has offset of the subcommand in the original (parent) command
	partial_offset uint32
	padding        [16 * 32]byte // I'm scared of the kernel, I want to not accidently overshoot
}

/*
struct scst_user_scsi_cmd_exec {
	aligned_u64 sess_h;

	uint8_t cdb[SCST_MAX_CDB_SIZE];
	uint16_t cdb_len;

	aligned_i64 lba;

	aligned_i64 data_len;
	int32_t bufflen;
	int32_t alloc_len;
	aligned_u64 pbuf;
	uint8_t queue_type;
	uint8_t data_direction;
	uint8_t partial;
	int32_t timeout;

	aligned_u64 p_out_buf;
	int32_t out_bufflen;

	uint32_t sn;

	uint32_t parent_cmd_h;
	int32_t parent_cmd_data_len;
	uint32_t partial_offset;
};

*/

const (
	SCST_USER_EXEC               = 3228070659
	SCST_USER_ALLOC_MEM          = 3223876354
	SCST_USER_PARSE              = 3226497793
	SCST_USER_ON_CACHED_MEM_FREE = 2148037381
	SCST_USER_ON_FREE_CMD        = 2148561668
	SCST_USER_TASK_MGMT_RECEIVED = 3222827784
	SCST_USER_TASK_MGMT_DONE     = 3222827785
	SCST_USER_ATTACH_SESS        = 2182640416
	SCST_USER_DETACH_SESS        = 2182640417
)

// struct scst_user_reply_cmd {
// 	uint32_t cmd_h;
// 	uint32_t subcode;
// 	union {
// 		int32_t result;
// 		struct scst_user_scsi_cmd_reply_parse parse_reply;
// 		struct scst_user_scsi_cmd_reply_alloc_mem alloc_reply;
// 		struct scst_user_scsi_cmd_reply_exec exec_reply;
// 		struct scst_user_ext_copy_reply_remap remap_reply;
// 	};
// };

type raw_scst_user_reply_cmd_result struct {
	cmd_h   uint32
	subcode uint32
	result  int32
}

/*
struct scst_user_scsi_cmd_reply_exec {
	int32_t resp_data_len;
	aligned_u64 pbuf;

#define SCST_EXEC_REPLY_BACKGROUND	0
#define SCST_EXEC_REPLY_COMPLETED	1
#define SCST_EXEC_REPLY_DO_WRITE_SAME	2
	uint8_t reply_type;

	uint8_t status;
	union {
		struct {
			uint8_t sense_len;
			aligned_u64 psense_buffer;
		};
		struct {
			uint16_t ws_descriptors_len;
			aligned_u64 ws_descriptors;
		};
	};
};
*/

const (
	SCST_EXEC_REPLY_BACKGROUND    = 0
	SCST_EXEC_REPLY_COMPLETED     = 1
	SCST_EXEC_REPLY_DO_WRITE_SAME = 2
)

type raw_scst_user_reply_cmd_exec_reply struct {
	cmd_h         uint32
	subcode       uint32
	resp_data_len int32
	fake          int32
	// do I need to add a fake int32 here to align?
	pbuf       uintptr
	reply_type uint8
	status     uint8
}

type raw_scst_user_reply_cmd_exec_reply_sense struct {
	cmd_h         uint32
	subcode       uint32
	resp_data_len int32
	// fake          int32
	// do I need to add a fake int32 here to align?
	pbuf          uintptr
	reply_type    uint8
	status        uint8
	sense_len     uint8
	psense_buffer uintptr

	/*
		union {
			struct {
				uint8_t sense_len;
				aligned_u64 psense_buffer;
			};
		};
	*/
}

type raw_scst_user_reply_cmd_exec_reply_descriptors struct {
	cmd_h         uint32
	subcode       uint32
	resp_data_len int32
	// fake          int32
	// do I need to add a fake int32 here to align?
	pbuf               uintptr
	reply_type         uint8
	status             uint8
	ws_descriptors_len uint16
	ws_descriptors     uintptr

	/*
			struct {
				uint16_t ws_descriptors_len;
				aligned_u64 ws_descriptors;
			};
		};
	*/
}

const (
	DEVICE_TYPE_DISK           = 0x00
	DEVICE_TYPE_TAPE           = 0x01
	DEVICE_TYPE_PROCESSOR      = 0x03 /* HP scanners use this */
	DEVICE_TYPE_WORM           = 0x04 /* Treated as ROM by our system */
	DEVICE_TYPE_ROM            = 0x05
	DEVICE_TYPE_SCANNER        = 0x06
	DEVICE_TYPE_MOD            = 0x07 /* Magneto-optical disk treated as TYPE_DISK */
	DEVICE_TYPE_MEDIUM_CHANGER = 0x08
	DEVICE_TYPE_ENCLOSURE      = 0x0d /* Enclosure Services Device */
	DEVICE_TYPE_NO_LUN         = 0x7f
)
