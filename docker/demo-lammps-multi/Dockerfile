ARG flux_sched_tag=focal
FROM fluxrm/flux-sched:${flux_sched_tag}

# docker build -t 

# It's good to be root :)
USER root
WORKDIR /opt/
ARG build_jobs=6

RUN apt-get update
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install \
    build-essential \
    gfortran \
    autotools-dev \
    autoconf \
    automake \
    cmake \
    git \
    python3 \
    openssh-server \
    openssh-client \
    fftw3-dev fftw3 \
    apt-utils \
    libnuma-dev \
    libgomp1 \
    pdsh \
    dnsutils \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# What we want to install and how we want to install it
# is specified in a manifest file (spack.yaml)
RUN mkdir /opt/spack-environment \
    &&  (echo "spack:" \
    &&   echo "  specs:" \
    &&   echo "  - openmpi@4.1.2" \
    &&   echo "  config:" \
    &&   echo "    install_tree: /opt/software" \
    &&   echo "  view: /opt/view") > /opt/spack-environment/spack.yaml \
    && cd /opt/spack-environment \
    && git clone --single-branch --branch develop https://github.com/spack/spack.git \
    && . spack/share/spack/setup-env.sh \
    && spack env activate . \
    && spack external find openssh \
    && spack external find cmake \
    && spack install --reuse --fail-fast \
    && spack gc -y \
    && find -L /opt/view/* -type f -exec readlink -f '{}' \; | \
    xargs file -i | \
    grep 'charset=binary' | \
    grep 'x-executable\|x-archive\|x-sharedlib' | \
    awk -F: '{print $1}' | xargs strip -s \
    && cd /opt/spack-environment \
    && . spack/share/spack/setup-env.sh \
    && spack env activate --sh -d . >> /etc/profile.d/z10_spack_environment.sh

WORKDIR /root
RUN git clone --single-branch --branch stable_29Sep2021_update2 https://github.com/lammps/lammps.git && \
    cd lammps && \
    mkdir build && \
    cd build && \
    . /etc/profile && \ 
    cmake ../cmake -D PKG_REAXFF=yes -D BUILD_MPI=yes -D PKG_OPT=yes -D FFT=FFTW3 && \
    make -j ${build_jobs} && \
    mkdir -p /var/run/sshd

# Allow OpenSSH to talk to containers without asking for confirmation
# by disabling StrictHostKeyChecking.
# mpi-operator mounts the .ssh folder from a Secret. For that to work, we need
# to disable UserKnownHostsFile to avoid write permissions.
# Disabling StrictModes avoids directory and files read permission checks.
ARG port=22
COPY --chown=root mpi_setup.sh /root/mpi_setup.sh
RUN sed -i "s/[ #]\(.*StrictHostKeyChecking \).*/ \1no/g" /etc/ssh/ssh_config \
    && echo "    UserKnownHostsFile /dev/null" >> /etc/ssh/ssh_config \
    && sed -i "s/[ #]\(.*Port \).*/ \1$port/g" /etc/ssh/ssh_config \
    && sed -i "s/#\(StrictModes \).*/\1no/g" /etc/ssh/sshd_config \
    && sed -i "s/#\(Port \).*/\1$port/g" /etc/ssh/sshd_config && \
    mkdir -p /root/lammps && \
    cp -R /root/lammps/build/lmp /usr/local/bin
WORKDIR /root/lammps/examples/reaxff/HNS
