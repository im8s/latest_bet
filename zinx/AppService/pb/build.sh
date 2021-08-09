#!/bin/bash
protoc --go_out=. msg.proto
protoc --cpp_out=. msg.proto

