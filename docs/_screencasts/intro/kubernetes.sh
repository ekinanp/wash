#doitlive speed: 2
#doitlive prompt: {TTY.CYAN}wash {dir} {r_angle}{TTY.RESET}

cd kubernetes/docker-desktop
ls
cd docker
ls

# Pods
ls pods
find pods -k '*pod' -m '.status.phase' Running -m '.metadata.labels.pod-template-hash' -exists
cd pods
wexec compose-6c67d745f6-ljtwr/compose uname
cat compose-6c67d745f6-ljtwr/compose

# TODO: Add PVCs
