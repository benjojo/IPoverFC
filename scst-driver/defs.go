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
