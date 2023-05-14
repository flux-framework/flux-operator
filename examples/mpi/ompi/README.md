# OpenMPI Example

You should be able to create a MiniKube cluster, install the operator with creating the namespace:

```bash
$ minikube start
$ kubectl create namespace flux-operator
$ kubectl apply -f ../../dist/flux-operator.yaml
```

You might want to pre-pull the container:

```bash
$ minikube ssh docker pull ghcr.io/rse-ops/ompi:flux-sched-focal
```

And then create the MiniCluster:

```bash
$ kubectl create -f minicluster.yaml
```

And watch the example run!

```bash
$ kubectl logs -n flux-operator flux-sample-0-5gjqt -f
```

A successful run will show four MPI ranks (and mpich is really vocal huh?)...

```console
broker.info[0]: rc1.0: running /etc/flux/rc1.d/02-cron
broker.info[0]: rc1.0: /etc/flux/rc1 Exited (rc=0) 0.6s
broker.info[0]: rc1-success: init->quorum 0.602697s
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: online: flux-sample-[0-3] (ranks 0-3)
broker.info[0]: quorum-full: quorum->run 0.361781s
Hello, world!  I am 0 of 4(MPICH Version:       3.3a2
MPICH Release date:     Sun Nov 13 09:12:11 MST 2016
MPICH Device:   ch3:nemesis
MPICH configure:        --build=x86_64-linux-gnu --prefix=/usr --includedir=${prefix}/include --mandir=${prefix}/share/man --infodir=${prefix}/share/info --sysconfdir=/etc --localstatedir=/var --disable-silent-rules --libdir=${prefix}/lib/x86_64-linux-gnu --libexecdir=${prefix}/lib/x86_64-linux-gnu --disable-maintainer-mode --disable-dependency-tracking --with-libfabric --enable-shared --prefix=/usr --enable-fortran=all --disable-rpath --disable-wrapper-rpath --sysconfdir=/etc/mpich --libdir=/usr/lib/x86_64-linux-gnu --includedir=/usr/include/mpich --docdir=/usr/share/doc/mpich --with-hwloc-prefix=system --enable-checkpointing --with-hydra-ckpointlib=blcr CPPFLAGS= CFLAGS= CXXFLAGS= FFLAGS= FCFLAGS=
MPICH CC:       gcc  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security  -O2
MPICH CXX:      g++  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security -O2
MPICH F77:      gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
MPICH FC:       gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
, 1297)
Hello, world!  I am 2 of 4(MPICH Version:       3.3a2
MPICH Release date:     Sun Nov 13 09:12:11 MST 2016
MPICH Device:   ch3:nemesis
MPICH configure:        --build=x86_64-linux-gnu --prefix=/usr --includedir=${prefix}/include --mandir=${prefix}/share/man --infodir=${prefix}/share/info --sysconfdir=/etc --localstatedir=/var --disable-silent-rules --libdir=${prefix}/lib/x86_64-linux-gnu --libexecdir=${prefix}/lib/x86_64-linux-gnu --disable-maintainer-mode --disable-dependency-tracking --with-libfabric --enable-shared --prefix=/usr --enable-fortran=all --disable-rpath --disable-wrapper-rpath --sysconfdir=/etc/mpich --libdir=/usr/lib/x86_64-linux-gnu --includedir=/usr/include/mpich --docdir=/usr/share/doc/mpich --with-hwloc-prefix=system --enable-checkpointing --with-hydra-ckpointlib=blcr CPPFLAGS= CFLAGS= CXXFLAGS= FFLAGS= FCFLAGS=
MPICH CC:       gcc  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security  -O2
MPICH CXX:      g++  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security -O2
MPICH F77:      gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
MPICH FC:       gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
, 1297)
Hello, world!  I am 3 of 4(MPICH Version:       3.3a2
MPICH Release date:     Sun Nov 13 09:12:11 MST 2016
MPICH Device:   ch3:nemesis
MPICH configure:        --build=x86_64-linux-gnu --prefix=/usr --includedir=${prefix}/include --mandir=${prefix}/share/man --infodir=${prefix}/share/info --sysconfdir=/etc --localstatedir=/var --disable-silent-rules --libdir=${prefix}/lib/x86_64-linux-gnu --libexecdir=${prefix}/lib/x86_64-linux-gnu --disable-maintainer-mode --disable-dependency-tracking --with-libfabric --enable-shared --prefix=/usr --enable-fortran=all --disable-rpath --disable-wrapper-rpath --sysconfdir=/etc/mpich --libdir=/usr/lib/x86_64-linux-gnu --includedir=/usr/include/mpich --docdir=/usr/share/doc/mpich --with-hwloc-prefix=system --enable-checkpointing --with-hydra-ckpointlib=blcr CPPFLAGS= CFLAGS= CXXFLAGS= FFLAGS= FCFLAGS=
MPICH CC:       gcc  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security  -O2
MPICH CXX:      g++  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security -O2
Hello, world!  I am 1 of 4(MPICH Version:       3.3a2
MPICH F77:      gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
MPICH Release date:     Sun Nov 13 09:12:11 MST 2016
MPICH FC:       gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
MPICH Device:   ch3:nemesis
, 1297)
MPICH configure:        --build=x86_64-linux-gnu --prefix=/usr --includedir=${prefix}/include --mandir=${prefix}/share/man --infodir=${prefix}/share/info --sysconfdir=/etc --localstatedir=/var --disable-silent-rules --libdir=${prefix}/lib/x86_64-linux-gnu --libexecdir=${prefix}/lib/x86_64-linux-gnu --disable-maintainer-mode --disable-dependency-tracking --with-libfabric --enable-shared --prefix=/usr --enable-fortran=all --disable-rpath --disable-wrapper-rpath --sysconfdir=/etc/mpich --libdir=/usr/lib/x86_64-linux-gnu --includedir=/usr/include/mpich --docdir=/usr/share/doc/mpich --with-hwloc-prefix=system --enable-checkpointing --with-hydra-ckpointlib=blcr CPPFLAGS= CFLAGS= CXXFLAGS= FFLAGS= FCFLAGS=
MPICH CC:       gcc  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security  -O2
MPICH CXX:      g++  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -Wformat -Werror=format-security -O2
MPICH F77:      gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
MPICH FC:       gfortran  -g -O2 -fdebug-prefix-map=/build/mpich-O9at2o/mpich-3.3~a2=. -fstack-protector-strong -O2
, 1297)
broker.info[0]: rc2.0: flux submit -N 4 -n 4 --quiet --watch ./hello_cxx Exited (rc=0) 0.4s
broker.info[0]: rc2-success: run->cleanup 0.380367s
broker.info[0]: cleanup.0: flux queue stop --quiet --all --nocheckpoint Exited (rc=0) 0.1s
broker.info[0]: cleanup.1: flux cancel --user=all --quiet --states RUN Exited (rc=0) 0.1s
broker.info[0]: cleanup.2: flux queue idle --quiet Exited (rc=0) 0.1s
broker.info[0]: cleanup-success: cleanup->shutdown 0.264937s
broker.info[0]: children-complete: shutdown->finalize 62.0603ms
broker.info[0]: rc3.0: running /etc/flux/rc3.d/01-sched-fluxion
broker.info[0]: rc3.0: /etc/flux/rc3 Exited (rc=0) 0.2s
broker.info[0]: rc3-success: finalize->goodbye 0.217901s
broker.info[0]: goodbye: goodbye->exit 0.028526ms
```

And the job will be completed.

```bash
kubectl get -n flux-operator pods
```
```console
NAME                  READY   STATUS      RESTARTS   AGE
flux-sample-0-flg28   0/1     Completed   0          9m39s
flux-sample-1-fplvv   0/1     Completed   0          9m39s
flux-sample-2-7bltz   0/1     Completed   0          9m39s
flux-sample-3-p8mtj   0/1     Completed   0          9m39s
```