# IN WASH SHELL

echo "This screencast demos the Docker, Kubernetes, AWS and GCP plugins (these are the plugins shipped with Wash). We'll show off each plugin's hierarchy, interact with a few of its resources, then conclude the screencast by summarizing all the things that we just demoed. For more detailed core plugin documentation, including information on how to get a specific one working, please refer to https://puppetlabs.github.io/wash/docs/#core-plugins"

echo "Let's start with the AWS plugin. First, we'll use the stree command to get a high-level overview of its hierarchy."

stree aws

echo "We see from the output that the AWS plugin lets us interact with EC2 instances and S3 objects. The available resources are grouped by the AWS profile. Let's check out one such profile."

ls aws

cd aws/<profile>

echo "Note that we are now in the <profile> profile."

ls

cd resources

ls

cd s3

ls

echo "Notice that everything listed here is an S3 bucket."

cd <bucket_name>

ls

echo "Amazon S3 has no concept of directories or files. Everything is an object, and each object is identified by its key. However, we can impose some hierarchical structure on S3 objects by grouping keys with common prefixes into a specific directory. For example, the S3 objects 'foo/bar' and 'foo/baz' would be represented as files with path 'foo/bar' and path 'foo/baz', where 'foo' is a 'directory'. Thus, everything shown in the above output is either an S3 object prefix ('directory') or an S3 object ('file'). Note that you can use the stree command to print this information in case you ever forget it."

stree

echo "Here, '.' is the current directory, which is an S3 bucket."

echo "Let's navigate through this bucket and cat some S3 objects."

cd <prefix>

cat <object>

echo "Note that this object's key is <prefix>/<object>"

cd ../<other_prefix>

cd <other_prefix_prefix>

cat <object>

echo "Note that this object's key is <other_prefix>/<other_prefix_prefix>/<object>"

echo "Now let's navigate into another bucket and do the same thing."

cd ../../../<some_other_bucket>

cat <object>

echo "Pretty neat huh? Could we do anything else with S3 objects? Let's use winfo to find out."

winfo <object>

echo "Looks like S3 objects only support the 'read' action so we cannot do anything else with them. Oh well. Guess it's time to interact with some EC2 instances. We're going to speed things up a bit here for the sake of brevity."

cd ../../ec2/instances

ls

echo "These are all the EC2 instances associated with the <profile> profile. Let's try exec'ing a command on one."

wexec <ec2_instance> uname

echo "Sweet! Let's try cd'ing into it."

cd <ec2_instance>

ls

echo "Here, console.out represents the EC2 instance's console output. fs is the root directory of the EC2 instance's filesystem. Finally, metadata.json is the EC2 instance's metadata. We're going to show off console.out and fs (metadata.json will be discussed in a separate screencast)."

cat console.out

echo "Nice! Now let's check out fs."

ls fs

echo "That should look familiar, because it is equivalent to ls'ing the EC2 instance's root if you were SSH'd into it. That's right, Wash lets you navigate your EC2 instances' filesystems as if you were logged onto the machine. Let's play around with that a bit."

ls fs/b*

echo "Yup, globbing works."

cd fs/home

ls

echo "These are all the users on the current EC2 instance."

cd ../var/log

ls

echo "This is everything in the current EC2 instance's /var/log directory."

cat <log_file>

tail -f <log_file>

echo "Yes, you can 'cat' and 'tail' the log files on this specific EC2 instance. Let's check out another EC2 instance's fs directory."

cd ../../../<other_ec2_instance>

cd fs

ls

cd /var/log

cat <log_file>

echo "Pretty neat, huh? Now check this out..."

cd ../../../../

tail -f <ec2_instance>/fs/var/log/<log_file> <other_ec2_instance>/fs/var/log/<log_file>

echo "That's right. With Wash, you can tail multiple log files spread out across multiple EC2 instances. In fact, you can tail multiple log files spread out across multiple Docker containers, multiple EC2 instances in multiple profiles, or multiple GCP Compute instances in multiple projects. Here's a little preview of what that might look like."

tail -f ../../../../../docker/containers/<container>/fs/var/log/<log_file> <ec2_instance>/fs/var/log/<log_file> ../../../../<other_profile>/resources/ec2/instances/<ec2_instance>/fs/var/log/<log_file>  ../../../../../gcp/<project>/compute/<instance>/fs/var/log/<log_file> ../../../../../gcp/<other_project>/compute/<instance>/fs/var/log/<log_file>

echo "That's pretty neat huh? As a fun exercise, how would you tail those files without Wash? If you're able to come up with a solution, please tell us about it in the Slack channel!"

echo "Anyways, back to the screencast. Another neat thing about Wash is that if your system shell supports globbing, we can do something like"

tail -f i-0*/fs/var/log/<log_file>

echo "i.e. we can glob over our EC2 instances! This is nothing magical though. All we're doing is globbing a path. It just so happens that the part of our path that we're globbing corresponds to an EC2 instance. This means something like"

tail -f i-0*/fs/var/log/*.log

echo "also makes sense."

echo "Anyways, that's enough of the AWS plugin. Let's check out GCP."

cd ../../../../../

stree gcp

echo "We see from the output that the GCP plugin lets us interact with Storage objects and Compute instances. The available resources are grouped by the GCP project."

echo "Note that we recommend you do 'stree <plugin>' before checking out <plugin>. This way, you get a high-level overview of the plugin's hierarchy while also familiarizing yourself with the vendor's API."

echo "Anyways, let's start walking through the GCP plugin."

ls gcp

cd gcp/<project>

echo "Note that we are now in the <project> project."

cd storage

ls

echo "Notice that everything listed here is a GCP Storage bucket."

cd <bucket_name>

echo "GCP Storage organizes its objects in a similar manner to Amazon S3. This means that objects with keys foo/bar and foo/baz will be grouped under a 'foo' directory."

cat <object>

echo "Note that this object's key is <object>."

cat <prefix>/<object>

echo "and this object's key is <prefix>/<object>."

echo "Now let's check out some GCP compute instances."

cd ../../compute

ls

echo "These are all the Compute instances associated with the <project> project. Let's try exec'ing a command on one."

wexec <compute_instance> uname

echo "Sweet! Let's try cd'ing into it."

cd <compute_instance>

ls

echo "Notice that everything listed here is the same set of stuff that was listed when we ls'ed into an EC2 instance. Here we see the Compute instance's console output, its metadata.json file, and an fs directory representing its filesystem root. Let's play around with the console output and fs directory."

cat console.out

echo "Nice!"

ls fs

cd fs/var/log

cat <log_file>

tail -f <log_file>

echo "Sweet. And remember, the fs directory lets us tail log files spread out across multiple GCP compute instances."

cd ../../../../

tail -f <compute_instance>/fs/var/log/<log_file> <other_compute_instance>/fs/var/log/<log_file>

echo "Sweet! That about wraps up the GCP plugin. Let's now move on to the Kubernetes plugin."

cd ../../../../../

stree kubernetes

echo "We see from the output that the Kubernetes plugin lets us interact with Persistent Volume Claims and Pods. The available resources are grouped by namespace, and the namespaces are grouped by Kubernetes context."

echo "Let's start walking through the Kubernetes plugin."

cd kubernetes

ls

cd kubernetes/<context>

echo "Note that we are now in the <context> context."

ls

echo "These are all the namespaces associated with that context."

cd <namespace>

ls

cd persistentvolumeclaims

ls

echo "Notice that everything here is a PVC (Persistent Volume Claim)"

cd <pvc>

ls

echo "These are the PVC's top-level directories and files."

cat <file>

tail -f <file>

echo "Yup. We can tail a PVC volume file for updates!"

cd <dir>

cat <file_in_dir>

tail -f <file_in_dir>

echo "Nice! There isn't much else to show in a PVC, so let's check out some pods."

cd ../../../pods

ls

echo "These are all the Kubernetes pods associated with the <namespace> namespace in the <context> context. Let's try exec'ing a command on one."

wexec <pod> uname

echo "Sweet! We unfortunately cannot cd into a pod because it does not implement the 'list' action. However, we are planning on exposing the pod's containers, so please follow https://github.com/puppetlabs/wash/issues/228 if this interests you!"

echo "Anyways, that's it for the Kubernetes plugin. Let's now demo the Docker plugin."

cd ../../../../

stree docker

echo "We see from the output that the Docker plugin lets us interact with volumes and containers. Let's start navigating through the plugin."

cd docker

ls

cd volumes

ls

echo "Notice that everything here is a volume."

cd <volume>

ls

echo "These are the volume's top-level directories and files."

cat <file>

tail -f <file>

echo "Yup. We can tail a volume file for updates!"

cd <dir>

cat <file_in_dir>

tail -f <file_in_dir>

echo "Nice! That's about all we can do in a Docker volume. Let's check out some containers."

cd ../../../containers

ls

echo "Notice that everything listed here is a container. Let's try exec'ing a command on one."

wexec <container> uname

echo "Sweet! Let's try cd'ing into it."

cd <container>

ls

echo "Here, we see the usual metadata.json file and fs directory. The only new entry here is 'log', which represents the container's log. Let's checkout 'log' and fs."

cat log

tail -f log

echo "Nice! Note that you can tail a container's log for updates."

ls fs

cd fs/var/log

cat <log_file>

tail -f <log_file>

echo "Sweet. Remember that these are actual files inside our container. Also remember that the fs directory lets us tail log files spread out across multiple Docker containers."

cd ../../../../

tail -f <container>/fs/var/log/<log_file> <other_container>/fs/var/log/<log_file>

echo "Sweet! That about wraps up the Docker plugin (and all of the shipped Wash plugins). We showed off a lot of stuff in this screencast, so let's do a quick recap of everything that we've done."

echo "(1) We used wexec to execute the uname command on an AWS EC2 instance/GCP compute instance/Kubernetes pod/Docker container."

echo "(2) We showed that viewing an AWS EC2 instance's/GCP compute instance's console output is as easy as cat'ing a regular file."

echo "(3) We navigated through an AWS EC2 instance's/GCP compute instance's/Docker container's filesystem as if we were logged onto the machine. We tailed some log files on each of those machines, and showed that with Wash, tailing multiple log files spread out across multiple AWS EC2 instaces/GCP compute instances/Docker containers is as easy as passing in a bunch of paths to 'tail -f' (remember 'tail -f ../../../../../docker/containers/<container>/fs/var/log/<log_file> <ec2_instance>/fs/var/log/<log_file> ../../../../<other_profile>/resources/ec2/instances/<ec2_instance>/fs/var/log/<log_file>  ../../../../../gcp/<project>/compute/<instance>/fs/var/log/<log_file> ../../../../../gcp/<other_project>/compute/<instance>/fs/var/log/<log_file>' ?)"

echo "(4) We navigated through an AWS S3/GCP Storage bucket by cd'ing into a bunch of 'directories', and cat'ed some AWS S3/GCP Storage objects. Remember that an AWS S3/GCP Storage bucket has no concept of 'directories' and 'files'; the AWS/GCP plugins created this hierarchy for us."

echo "(5) We navigated through a Kubernetes PVC and a Docker volume, and cat'ed/tail'ed some files. (And if you're asking whether it's possible to tail a Docker volume file AND an EC2 instance's log file via 'tail -f', the answer is 'YES'. We leave it up to you to figure out the right invocation)."

echo "Pretty neat huh?"

exit
