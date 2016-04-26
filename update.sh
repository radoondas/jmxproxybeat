#!/usr/bin/env bash
sh ../../elastic/beats/libbeat/scripts-3.4/update.sh jmxproxybeat . ../../elastic/beats/libbeat

# python3 ../../elastic/beats/dev-tools/export_dashboards-fix.py --url http://127.0.0.1:9555 --beat jmxproxybeat --index jmxproxybeat-* --dir etc/kibana5/