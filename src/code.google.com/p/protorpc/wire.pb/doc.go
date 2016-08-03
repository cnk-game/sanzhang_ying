// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package wire.pb is the protorpc wire format wrapper.
//
//	0. Frame Format
//	len : uvarint64
//	data: byte[len]
//	
//	1. Client Send Request
//	Send RequestHeader: sendFrame(conn, hdr, len(hdr))
//	Send Request: sendFrame(conn, body, hdr.snappy_compressed_request_len)
//	
//	2. Server Recv Request
//	Recv RequestHeader: recvFrame(conn, hdr, max_hdr_len, 0)
//	Recv Request: recvFrame(conn, body, hdr.snappy_compressed_request_len, 0)
//	
//	3. Server Send Response
//	Send ResponseHeader: sendFrame(conn, hdr, len(hdr))
//	Send Response: sendFrame(conn, body, hdr.snappy_compressed_response_len)
//	
//	4. Client Recv Response
//	Recv ResponseHeader: recvFrame(conn, hdr, max_hdr_len, 0)
//	Recv Response: recvFrame(conn, body, hdr.snappy_compressed_response_len, 0)
//	
//	5. Header Size
//	len(RequestHeader)  < Const.max_header_len.default
//	len(ResponseHeader) < Const.max_header_len.default
package google_protobuf_rpc_wire
