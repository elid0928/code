#!/bin/bash

# 配置丢包100%
tc qdisc add dev ens4 root netem loss 100%
sleep 100

# 删除丢包
tc qdisc del dev ens4 root netem loss 100%

# 配置延迟 100ms
tc qdisc add dev ens4 root netem delay 100ms
sleep 10
# 删除延迟
tc qdisc del dev ens4 root netem delay 100ms