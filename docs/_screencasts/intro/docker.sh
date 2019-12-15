#doitlive speed: 2
#doitlive prompt: {TTY.CYAN}wash {dir} {r_angle}{TTY.RESET}

cd docker
ls

# Containers
ls containers
find containers -k '*container' -m '.state' running -m '.labels.com\.docker\.compose\.version' -exists
cd containers
wexec wash_tutorial_redis_1 uname
cd wash_tutorial_redis_1
ls
cat log
cd fs
ls
find var/log -mtime -6w
cat var/log/dpkg.log
tail -f var/log/dpkg.log
# Hit Ctrl+C

cd $W/docker
ls

# Volumes
cd volumes
ls
find wash_tutorial_redis -name '*.aof'
cd wash_tutorial_redis
ls
cat appendonly.aof
