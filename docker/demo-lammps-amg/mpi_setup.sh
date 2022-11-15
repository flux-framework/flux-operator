#!/bin/bash

[[ $(cat /etc/mpi/hostfile | wc -l) != 0 ]] \
    && (date; echo "Hostfile is ready"; cat /etc/mpi/hostfile) \
    || (date; echo "Hostfile not ready ..."; sleep 10; exit 1) 

for host in $( cat /etc/mpi/hostfile )
do
    # Query DNS for hostname and wait if not resolved yet.
    sleep 0.1
    hname=$( host ${host} )
    if [[ "$hname" == *"not found: 3(NXDOMAIN)"* ]]
    then
        echo "host ${host} not resolved"
        while [[ "$hname" == *"not found: 3(NXDOMAIN)"* ]]
        do
            sleep 2
            hname=$( host ${host} )
        done
    fi

    # The host is resolved
    echo "host ${host} resolved successfully"
    ipaddr=$( echo ${hname} | cut -d ' ' -f4 )
    fqdn=$( echo ${hname} | cut -d ' ' -f1 )
    basehost=$( echo ${host} | sed 's/\..*$//' )
    echo "${ipaddr} ${fqdn} ${basehost}" >> /etc/hosts
done

WCOLL=/etc/mpi/hostfile PDSH_RCMD_TYPE=ssh pdsh -S -t 10 hostname -f 2>/dev/null > tmp
if [[ $? -gt 1 ]]
then
    echo "Error: ssh to workers"
    exit 1
fi

while [[ $( cat tmp | sort -u | wc -l ) != $( cat /etc/mpi/hostfile | wc -l ) ]]
do
    comm -13 <( cat tmp | cut -d: -f1 | sort -t - -k 3 -u ) \
             <( cat /etc/mpi/hostfile | cut -d. -f1 | sort -t - -k 3 -u ) 2>/dev/null > tmp2
    echo "Can't reach hosts: "
    cat tmp2
    WCOLL=tmp2 PDSH_RCMD_TYPE=ssh pdsh -t 5 hostname -f 2>/dev/null >> tmp
done

echo "SSH set up successfully on all worker pods"

WCOLL=/etc/mpi/hostfile PDSH_RCMD_TYPE=ssh pdcp /etc/hosts /etc/hosts
if [[ $? -ne 0 ]]
then
    echo "Error: can't pdcp hosts file to workers"
    exit 1
fi

echo "/etc/hosts copied successfully to all worker pods"

rm -f tmp tmp2 2>/dev/null

