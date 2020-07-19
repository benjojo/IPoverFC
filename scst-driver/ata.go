package main

import (
	"log"
	"unsafe"
)

func processExecCmd(in *raw_scst_user_get_cmd_scsi_cmd_exec) uintptr {
	log.Printf("------> Timeout %d \n%#v", in.timeout, in)
	ATAopCode := in.cdb[0]

	var emptyByte [1]byte

	reply := raw_scst_user_reply_cmd_exec_reply_sense{
		cmd_h:         in.cmd_h,
		subcode:       in.subcode,
		reply_type:    SCST_EXEC_REPLY_COMPLETED,
		resp_data_len: 0,
		pbuf:          uintptr(unsafe.Pointer(&emptyByte)),
		status:        SAM_STAT_GOOD,
		psense_buffer: uintptr(unsafe.Pointer(&emptyByte)),
		sense_len:     0,
	}

	switch ATAopCode {
	case ATA_TEST_UNIT_READY:
		/*
			(gdb) print *cmd
			$4 = {
			  sess_h = 18446614906040681408,
			  cdb = '\000' <repeats 15 times>,
			  cdb_len = 6,
			  lba = 0,
			  data_len = 0,
			  bufflen = 0,
			  alloc_len = 0,
			  pbuf = 0,
			  queue_type = 3 '\003',
			  data_direction = 4 '\004',
			  partial = 0 '\000',
			  timeout = 10,
			  p_out_buf = 0,
			  out_bufflen = 0,
			  sn = 0,
			  parent_cmd_h = 0,
			  parent_cmd_data_len = 0,
			  partial_offset = 0
			}
			(gdb) print *reply
			$5 = {
				resp_data_len = 0,
				pbuf = 0,
				reply_type = 1 '\001',
				status = 0 '\000',
				{
					{
						sense_len = 0 '\000',
						psense_buffer = 0
					},
					{
						ws_descriptors_len = 0,
						ws_descriptors = 0
					}
				}
			}

		*/
		// Do nothing???
	default:
		log.Printf("Unsupported ATA opcode: %d / %x", ATAopCode, ATAopCode)
	}

	return uintptr(unsafe.Pointer(&reply)) // lol total segfault bait
}

const (
	ATA_TEST_UNIT_READY       = 0x00
	ATA_REZERO_UNIT           = 0x01
	ATA_REQUEST_SENSE         = 0x03
	ATA_FORMAT_UNIT           = 0x04
	ATA_READ_BLOCK_LIMITS     = 0x05
	ATA_REASSIGN_BLOCKS       = 0x07
	ATA_READ_6                = 0x08
	ATA_WRITE_6               = 0x0a
	ATA_SEEK_6                = 0x0b
	ATA_READ_REVERSE          = 0x0f
	ATA_WRITE_FILEMARKS       = 0x10
	ATA_SPACE                 = 0x11
	ATA_INQUIRY               = 0x12
	ATA_RECOVER_BUFFERED_DATA = 0x14
	ATA_MODE_SELECT           = 0x15
	ATA_RESERVE               = 0x16
	ATA_RELEASE               = 0x17
	ATA_COPY                  = 0x18
	ATA_ERASE                 = 0x19
	ATA_MODE_SENSE            = 0x1a
	ATA_START_STOP            = 0x1b
	ATA_RECEIVE_DIAGNOSTIC    = 0x1c
	ATA_SEND_DIAGNOSTIC       = 0x1d
	ATA_ALLOW_MEDIUM_REMOVAL  = 0x1e

	ATA_SET_WINDOW             = 0x24
	ATA_READ_CAPACITY          = 0x25
	ATA_READ_10                = 0x28
	ATA_WRITE_10               = 0x2a
	ATA_SEEK_10                = 0x2b
	ATA_WRITE_VERIFY           = 0x2e
	ATA_VERIFY                 = 0x2f
	ATA_SEARCH_HIGH            = 0x30
	ATA_SEARCH_EQUAL           = 0x31
	ATA_SEARCH_LOW             = 0x32
	ATA_SET_LIMITS             = 0x33
	ATA_PRE_FETCH              = 0x34
	ATA_READ_POSITION          = 0x34
	ATA_SYNCHRONIZE_CACHE      = 0x35
	ATA_LOCK_UNLOCK_CACHE      = 0x36
	ATA_READ_DEFECT_DATA       = 0x37
	ATA_MEDIUM_SCAN            = 0x38
	ATA_COMPARE                = 0x39
	ATA_COPY_VERIFY            = 0x3a
	ATA_WRITE_BUFFER           = 0x3b
	ATA_READ_BUFFER            = 0x3c
	ATA_UPDATE_BLOCK           = 0x3d
	ATA_READ_LONG              = 0x3e
	ATA_WRITE_LONG             = 0x3f
	ATA_CHANGE_DEFINITION      = 0x40
	ATA_WRITE_SAME             = 0x41
	ATA_READ_TOC               = 0x43
	ATA_LOG_SELECT             = 0x4c
	ATA_LOG_SENSE              = 0x4d
	ATA_MODE_SELECT_10         = 0x55
	ATA_RESERVE_10             = 0x56
	ATA_RELEASE_10             = 0x57
	ATA_MODE_SENSE_10          = 0x5a
	ATA_PERSISTENT_RESERVE_IN  = 0x5e
	ATA_PERSISTENT_RESERVE_OUT = 0x5f
	ATA_MOVE_MEDIUM            = 0xa5
	ATA_READ_12                = 0xa8
	ATA_WRITE_12               = 0xaa
	ATA_WRITE_VERIFY_12        = 0xae
	ATA_SEARCH_HIGH_12         = 0xb0
	ATA_SEARCH_EQUAL_12        = 0xb1
	ATA_SEARCH_LOW_12          = 0xb2
	ATA_READ_ELEMENT_STATUS    = 0xb8
	ATA_SEND_VOLUME_TAG        = 0xb6
	ATA_WRITE_LONG_2           = 0xea

	SAM_STAT_GOOD                       = 0x00
	SAM_STAT_CHECK_CONDITION            = 0x02
	SAM_STAT_CONDITION_MET              = 0x04
	SAM_STAT_BUSY                       = 0x08
	SAM_STAT_INTERMEDIATE               = 0x10
	SAM_STAT_INTERMEDIATE_CONDITION_MET = 0x14
	SAM_STAT_RESERVATION_CONFLICT       = 0x18
	SAM_STAT_COMMAND_TERMINATED         = 0x22 /* obsolete in SAM-3 */
	SAM_STAT_TASK_SET_FULL              = 0x28
	SAM_STAT_ACA_ACTIVE                 = 0x30
	SAM_STAT_TASK_ABORTED               = 0x40
)
