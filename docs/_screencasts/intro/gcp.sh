#doitlive speed: 2
#doitlive prompt: {TTY.CYAN}wash {dir} {r_angle}{TTY.RESET}

cd gcp/Wash
ls

# Compute instances
ls compute
find compute -k '*instance' -m '.status' RUNNING -m '.labels.owner' -exists
cd compute
wexec michael-test-instance uname
cd michael-test-instance
ls
cat console.out
cd fs
ls
find var/log -mtime -1h
cat var/log/messages
tail -f var/log/messages
# Hit Ctrl+C

cd $W/gcp/Wash
ls

# Storage
cd storage
ls
find some-wash-bucket -name '*.sh'
cd some-wash-bucket
ls
cat godoc.sh
cd folder
ls
