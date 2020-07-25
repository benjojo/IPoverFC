package main

import (
	"log"
	"runtime"
	"unsafe"
)

func processExecCmd(in *raw_scst_user_get_cmd_scsi_cmd_exec) raw_scst_user_reply_cmd_exec_reply {
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
	log.Printf("------> Timeout %d \n%#v", in.timeout, in)
	ATAopCode := in.cdb[0]

	// var emptyByte [1]byte

	reply := raw_scst_user_reply_cmd_exec_reply{
		cmd_h:         in.cmd_h,
		subcode:       in.subcode,
		reply_type:    SCST_EXEC_REPLY_COMPLETED,
		resp_data_len: 0,
		// pbuf:          uintptr(unsafe.Pointer(&emptyByte)),
		pbuf:   0,
		status: SAM_STAT_GOOD,
		// psense_buffer: uintptr(unsafe.Pointer(&emptyByte)),
		// sense_len:     0,
	}

	if in.data_direction == 2 { // READ
		reply.resp_data_len = in.bufflen
	} else {
		log.Printf("FUCK IT@S A WRITE RUN")
	}

	log.Printf("------> Opcode %x", ATAopCode)

	switch ATAopCode {
	case ATA_TEST_UNIT_READY:
		log.Printf("ATA_TEST_UNIT_READY")
		// Do nothing???
	case ATA_INQUIRY:
		log.Printf("ATA_INQUIRY")

		handleATAinquiry(in, &reply)
		// Haha oh my fucking god.

	default:
		log.Printf("Unsupported ATA opcode: %d / %x", ATAopCode, ATAopCode)
	}

	runtime.KeepAlive(reply)
	// return uintptr(unsafe.Pointer(&reply)) // lol total segfault bait
	return reply
}

func handleATAinquiry(in *raw_scst_user_get_cmd_scsi_cmd_exec, reply *raw_scst_user_reply_cmd_exec_reply) {
	var finalOutput = [128]byte{0}
	output := make([]byte, in.bufflen)
	resp_len := 0

	output[0] = DEVICE_TYPE_DISK // sure, i'm a uhhh disk
	// Readers note: I really did consider being a DEVICE_TYPE_SCANNER
	// but I have no real desire to figure out what enumeration is needed
	// for that.
	if (in.cdb[1] & 0x01) > 1 {

		if 0 == in.cdb[2] { /* supported vital product data pages */
			// Aka, "Hi frien, what do you support"
			output[3] = 5
			output[4] = 0x0  /* this page */
			output[5] = 0x80 /* unit serial number */
			output[6] = 0x83 /* device identification */
			output[7] = 0xB0 /* block limits */
			output[8] = 0xB1 /* block device characteristics */
			resp_len = int(uint8(output[3]) + 6)

		} else if 0x80 == in.cdb[2] { /* unit serial number */

		} else if 0x83 == in.cdb[2] { /* device identification */

		} else if 0xB0 == in.cdb[2] { /* Block Limits */

		} else if 0xB1 == in.cdb[2] { /* Block Device Characteristics */

		} else {
			// unsupported
		}

	} else {
		// Really basic stuff:

		if in.cdb[2] != 0 {
			// TRACE_DBG("INQUIRY: Unsupported page %x", cmd->cdb[2]);
			// PRINT_INFO("INQUIRY: Unsupported page %x", cmd->cdb[2]);
			// set_cmd_error(vcmd,
			//     SCST_LOAD_SENSE(scst_sense_invalid_field_in_cdb));
			// goto out;
			log.Printf("FUUUUUUUUUUUUUUUUUUCK Unsupported INQ PAGE")
		}

		output[2] = 6    /* Device complies to SPC-4 */
		output[3] = 0x12 /* HiSup + data in format specified in SPC */
		output[4] = 31   /* n - 4 = 35 - 4 = 31 for full 36 byte data */
		output[6] = 1    /* MultiP 1 */
		output[7] = 2    /* CMDQUE 1, BQue 0 => commands queuing supported */

		copy(output[8:], []byte("xXxBCxXx"))
		/* 8 byte ASCII Vendor Identification of the target - left aligned */
		// memcpy(&buf[8], VENDOR, 8);

		/* 16 byte ASCII Product Identification of the target - left aligned */
		copy(output[16:], []byte("                "))
		copy(output[16:], []byte("YOLO"))
		// memset(&buf[16], ' ', 16);
		// len = min(strlen(dev->name), (size_t)16);
		// memcpy(&buf[16], dev->name, len);

		/* 4 byte ASCII Product Revision Level of the target - left aligned */
		// memcpy(&buf[32], FIO_REV, 4);
		copy(output[16:], []byte("350 "))

		resp_len = int(output[4]) + 5

		// */
	}

	log.Printf("debug: resp_len = %d", resp_len)

	copy(finalOutput[:], output[:])
	in.pbuf = uintptr(unsafe.Pointer(&finalOutput))
	reply.resp_data_len = int32(resp_len)
	runtime.KeepAlive(finalOutput)
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
