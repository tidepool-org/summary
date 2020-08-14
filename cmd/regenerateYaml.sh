#!/bin/bash
pip3 install PyYAML==5.1

oapi-codegen  -generate=server summary.v1.yaml > api/gen_server.go
oapi-codegen  -generate=types summary.v1.yaml > api/gen_types.go
oapi-codegen  -generate=spec summary.v1.yaml > api/gen_spec.go


sed  -i .bak 's/package Summary/package api/' api/gen_types.go; rm api/gen_types.go.bak
sed  -i .bak 's/package Summary/package api/' api/gen_spec.go; rm api/gen_spec.go.bak
sed  -i .bak 's/package Summary/package api/' api/gen_server.go; rm api/gen_server.go.bak

python3 cmd/createPolicyfile.py
