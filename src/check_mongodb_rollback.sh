#!/bin/bash
#
## xiaoxiao
## 20150809
## issue: https://pm.qbox.me/issues/16713

disk_list=`grep dbpath /home/qboxserver/mongodb*/mongodb.conf | awk '{print $NF}' 2>/dev/null`
if [ $? -ne 0 ]; then
    echo -1; exit -1
fi

for disk in $disk_list; do
    rollback_dir="${disk}rollback"
    if [ -d $rollback_dir ]; then
        echo 1; exit 0
    fi
done
echo 0
