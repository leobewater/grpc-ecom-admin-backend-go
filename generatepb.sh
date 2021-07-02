#!/bin/bash

protoc proto/user.proto --go_out=plugins=grpc:pb 
