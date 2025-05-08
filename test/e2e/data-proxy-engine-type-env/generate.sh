#!/bin/sh

export PRISMA_CLIENT_ENGINE_TYPE=dataproxy
go run github.com/gooferOrm/goofer generate --schema schema.out.prisma
