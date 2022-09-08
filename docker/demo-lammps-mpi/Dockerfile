ARG flux_sched_tag
FROM fluxrm/flux-sched:${flux_sched_tag}
ARG build_jobs=6
ENV build_jobs=${build_jobs}

# docker build --build-arg flux_sched_tag=focal -t ghcr.io/flux-framework/demo-mpi:focal .

# Potential entrypoint:
#    mpirun --allow-run-as-root --mca orte_launch_agent /opt/view/bin/orted --mca plm_rsh_agent rsh -x PATH -x LD_LIBRARY_PATH -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite

ENV DEBIAN_FRONTEND=noninteractive

# It's good to be root :)
USER root
WORKDIR /opt/
RUN apt-get update && \
    apt-get install -y fftw3-dev fftw3

# install laamps
RUN git clone --depth 1 --branch stable_29Sep2021_update2 https://github.com/lammps/lammps.git /opt/laamps && \
    cd /opt/laamps && \
    mkdir build && \
    cd build && \
    . /etc/profile && \ 
    cmake ../cmake -D PKG_REAXFF=yes -D BUILD_MPI=yes -D PKG_OPT=yes -D FFT=FFTW3 && \
    make -j ${build_jobs} && \
    make install
  
# Target this example
ENV PATH=/root/.local/bin:$PATH
WORKDIR /opt/lammps/examples/reaxff/HNS
