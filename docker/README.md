# Docker Container

The following automated builds are provided alongside the Flux Operator:

 - [demo-lammps-mpi](demo-lammps-mpi): A demo mpi container that can be used for a MiniCluster
 
You can use these containers as examples of how you should build your flux container
to use with the operator. Generally we recommend using the flux-sched base
so that the install locations are consistent. This assumes that:

 - `/etc/flux` is used for configuration and general setup
 - `/usr/libexec/flux` has executables like flux-imp, flux-shell
 - `pdsh` is installed, which allows to (in parallel) start the cluster
 - flux-core / flux-sched with flux-security should be installed.
 
If/when needed we can lift some of these constraints, but for now they are 
reasonable.
