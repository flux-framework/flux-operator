#!/bin/bash
MATCH_FORMAT=${MATCH_FORMAT:-rv1}
NJOBS=${NJOBS:-10}
NNODES=${NNODES:-6}
printf "MATCH_FORMAT=${MATCH_FORMAT} NJOBS=$NJOBS NODES/JOB=$NNODES\n"

flux module remove sched-fluxion-qmanager
flux module remove sched-fluxion-resource
flux module remove resource
flux config load <<EOF
[sched-fluxion-qmanager]
queue-policy = "easy"
[sched-fluxion-resource]
match-format = "$MATCH_FORMAT"

[resource]
noverify = true
norestrict = true

[queues.offline]
requires = ["offline"]

[queues.online]
requires = ["online"]

[[resource.config]]
hosts = "flux-sample[0-3]"
properties = ["online"]

[[resource.config]]
hosts = "flux-sample[0-3],burst[0-99]"
cores = "0-3"

[[resource.config]]
hosts = "burst[0-99]"
properties = ["offline"]
cores = "4-103"
EOF
flux config get | jq '."sched-fluxion-resource"'
# monitor-force-up removed here
flux module load resource noverify
flux module load sched-fluxion-resource
flux module load sched-fluxion-qmanager
flux queue start --all --quiet
flux resource list
t0=$(date +%s.%N)

# These are fake jobs
flux submit -N$NNODES --cc=1-$NJOBS --queue=offline \
    --setattr=exec.test.run_duration=1ms \
    --quiet --wait hostname

# These are real jobs (2 nodes each)
flux submit -N1 --queue=online hostname

ELAPSED=$(echo $(date +%s.%N) - $t0 | bc -l)
THROUGHPUT=$(echo $NJOBS/$ELAPSED | bc -l)
R_SIZE=$(flux job info $(flux job last) R | wc -c)
OBJ_COUNT=$(flux module stats content-sqlite | jq .object_count)
DB_SIZE=$(flux module stats content-sqlite | jq .dbfile_size)

printf "%-12s %5d %4d %8.2f %8.2f %12d %12d %12d\n" \
        $MATCH_FORMAT $NJOBS $NNODES $ELAPSED $THROUGHPUT \
        $R_SIZE $OBJ_COUNT $DB_SIZE

# Get the last job
jobid=$(flux job last)
flux job attach ${jobid}

flux jobs -a
flux jobs -a --json | jq .jobs[0]