IPoverFC
===

![FC card](https://blog.benjojo.co.uk/asset/CtSNenJ5eB)

A hacked up fake SCSI device that is backing a TUN/TAP device. Meaning you can send ethernet frames down a SCSI device.

Powered by SCST, meaning this could be expanded over iSCSI, Fiber Channel, or whatever else SCST supports.

See blog post: https://blog.benjojo.co.uk/post/ip-over-fibre-channel-hack

This code must be compiled on a go version lower than go.1.13.16, I used go.1.13.15, anything higher will break the SCST driver due to changes in how go handles IO events (event signals will break SCST communication)
