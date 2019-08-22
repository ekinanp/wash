echo "This screencast introduces Wash, which is a shell environment that's layered on top of the system shell. We'll talk about plugins and actions, then navigate through a specific plugin and interact with some of its resources. We'll conclude the screencast with an overview of the stree command, and a small demo of an external plugin."

echo "First, let's start-up the Wash shell."

./wash

ls

echo "Each thing listed here is a 'directory'. Each 'directory' is a Wash plugin, and each Wash plugin serves as an adapter between a given vendor and Wash. Plugins are the sole mode of interaction between a Wash user and a given vendor's API. For example, the docker directory corresponds to the Docker plugin. It lets Wash users interact with Docker resources like containers and volumes. Similarly, the aws directory corresponds to the AWS plugin. It lets Wash users interact with AWS resources like EC2 instances and S3 objects." 

echo "Everything in Wash is an entry, including resources. Each plugin provides a hierarchical view of all its entries. Thus, navigating through a plugin's API reduces to cd'ing into a bunch of directories. Furthermore, each entry supports a specific set of Wash actions. The Wash actions are all the things you can do on a specific kind of entry. Examples of Wash actions include 'list' (listing the entry's children), 'exec' (executing a command on the entry), and 'read' (reading the entry's content). Note that entries which support 'list' are represented as directories in the shell. All other entries are represented as files."

echo "Let's walkthrough an example to see what all of this looks like in practice. We'll navigate through the Docker plugin and interact with some of its entries."

cd docker

echo "Notice how the prompt changes from 'wash . >' to 'wash docker >'. This is a useful way to track where you're currently at when navigating a given plugin."

ls

echo "From the 'containers' and 'volumes' directories, we see that the Docker plugin lets us interact with Docker containers and volumes. Let's try interacting with some containers."

cd containers

ls

echo "Note that everything listed here is a Docker container. Let's check out a Docker container's supported actions to see what we can do with it."

winfo <container>

echo "It looks like a Docker container supports the 'list' and 'exec' actions. Let's try 'exec'ing a command on a container."

wexec <first_letter><HIT_TAB> uname

echo "Nice! Also notice that you can tab-complete Wash entries. This works because Wash is layered on top of the system shell. Thus, if your system shell supports things like tab-completion and globbing, then those features will also carry-over to Wash."

echo "Anyways, back to our example. We see that a Docker container also supports the 'list' action, so it is modeled as a directory. That means we can cd into it, so let's go ahead and do that."

cd <container>

ls

echo "Everything listed here consists of stuff specific to the given container. For example, 'log' represents the container's log; 'fs' represents the container's root directory. Let's check out the log entry."

winfo log

echo "Notice that log doesn't support the 'list' action. That means it's represented as a file, so we can't cd into it."

cd log

echo "Fails as expected. However, log does support 'read' and 'stream' so we can 'cat'/'tail' its content."

cat log

tail -f log

echo "That's enough navigation for now. Let's go back to Wash's root."

cd ../../../

echo "Wash plugins are self-documenting. You can use the stree command to get a high-level overview of a plugin's hierarchy."

stree docker

echo "Notice that the 'containers' subtree contains zero or more containers (the '[]' denote 'zero or more instances of this thing'). Each container contains a 'log' file, a 'metadata.json' file, and an 'fs' directory. All of this matches what we found when we cd'ed through that subtree."

ls

echo "Wash ships with plugins for Docker, Kubernetes, AWS and GCP. However, you can also extend Wash via the external plugin interface. The puppetwash plugin shown here is an external plugin that lets you interact with PuppetDB (shoutout to timidri for writing it!). Let's check it out."

stree puppetwash

ls puppetwash

echo "As expected from the stree output, this shows all the available PE instances."

ls puppetwash/<pe_instance>/nodes

echo "And this shows all the available PE nodes tied to the given PE instance."

exit
