# IN WASH SHELL

echo "This screencast talks about entry attributes and metadata, then demos Wash's find command. The find command's useful for filtering out entries that satisfy a certain set of properties, some of which may come from their attributes or metadata. For example, you can use find to enumerate all your Amazon S3/GCP Storage buckets, list all Docker containers/Kubernetes pods that were created within the last 24 hours, or extract all your running AWS EC2 instances/GCP compute instances that contain a specific label. So, let's get started."

echo "All Wash entries are completely described by their metadata. This means that the metadata contains everything you'll ever need to know about the entry. This includes things like creation time, launch time, start time, VPC ID, operating system, expiration date, etc. You can use the 'meta' command to view an entry's metadata. Let's take a look at some examples."

meta docker/containers/<container>

echo "We see that a Docker container is described by quite a bit of properties. For example, the bottom part of the output shows us that it has a 'Platform', a 'ProcessLabel', a 'ResolvConfPath', 'RestartCount', 'SizeRootFs', 'SizeRw', and 'State'. If we scroll up in the output, we see some additional properties like 'Labels', 'Mounts', 'HostConfig', etc."

meta aws/<profile>/resources/ec2/instances/<instance>

echo "Here, we see that an EC2 instance is described by its 'Tags', 'VpcId', 'State', 'SecurityGroups', etc."

meta gcp/<project>/storage/<bucket>

echo "Here, we see that a GCP Storage bucket is described by its 'StorageClass', 'Name', 'Location', 'Labels', 'Etag', etc."

meta docker

echo "Yes, even non-resource entries like the docker directory have metadata. For those entries, the metadata will typically be an empty object."

echo "Some properties will be common across many different kinds of entries (e.g. creation time, last modified time). These common properties are referred to as 'entry attributes.' Examples of valid entry attributes include the creation time (crtime), the last modified time (mtime), and the content size (size). Every entry also includes a subset of its metadata via a special 'meta' attribute. This attribute is typically set to the raw response returned by the plugin's API when one makes a request to a 'list'-like endpoint. For example, a Docker container's 'meta' attribute corresponds to a container JSON object; this JSON object's included in the JSON array returned by Docker's /containers/json endpoint."

echo "You can use winfo to view an entry's attributes."

winfo docker/containers/<container>

echo "Here, we see that this container's crtime, mtime, atime, and ctime are all set to the date '2019-08-21T21:30:07-07:00'. We also see the value of its 'meta' attribute (notice that it is indeed a subset of the full metadata)."

echo "You can also set the 'meta' command's '--attribute' option to view an entry's meta attribute."

meta --attribute docker/containers/<container>

echo "That wraps up our discussion on attributes and metadata. If you'd like to read more about them, please refer to https://puppetlabs.github.io/wash/docs/#attributes-metadata."

echo "Now on to Wash's find command. The find command recursively descends a given path, printing out all of its subchildren."

find docker

echo "If we let find run indefinitely, then the output would contain every entry in the Docker plugin. We prematurely stopped it here because find was descending into the <container> container's root directory, and enumerating root directories takes a long time."

echo "We can use the maxdepth option to limit find's recursion."

find docker -maxdepth 1

find docker -maxdepth 2

echo "We can use the mindepth option to only print entries that are at a minimum depth."

find docker -mindepth 1

echo "Notice that the 'docker' entry's depth is 0, so it did not get printed."

find docker -mindepth 2 -maxdepth 3

echo "Here, we're using both the mindepth and maxdepth options."

echo "find also takes multiple paths."

find aws docker -maxdepth 1

echo "find's biggest advantage is that you can construct a predicate (filter) using its expression syntax. find will take your constructed predicate and print out only the satisfying entries."

echo "We'll give a brief overview of find's expression syntax. You can type 'find --help syntax' for a more detailed description."

find --help syntax

echo "That's a lot of information. We recommend piping it into a terminal pager like 'less' so you can more easily read it."

find --help syntax | less

echo "find's expression syntax consists of primaries (the individual predicates) and operators (the things that combine the predicates together). You can view the available primaries by typing in 'find --help'."

find --help

echo "Here, we see that some of the available primaries are 'action' and 'false'. Let's try them out."

find aws docker -action exec

echo "This only printed out the entries that supported the exec action. For the AWS plugin, these entries were EC2 instances. For the Docker plugin, these entries were Docker containers. Thus, the above invocation printed out all EC2 instances and Docker containers."

find docker -false

echo "This printed out nothing because the predicate always returned false."

echo "find's operators include the logical AND, OR and NOT operators. NOT has the highest precedence, followed by AND then OR."

find docker -action exec -false

find docker -action exec -a -false

echo "-action exec -false == -action exec -a -false == (-action exec) AND -false == false. Thus, the above invocation printed out nothing because the expression evaluated to false." 

find docker -action exec -o -false

echo "-action exec -o -false == (-action exec) OR -false == -action exec. Thus, the above invocation printed out all 'execable' entries in the Docker plugin (so all Docker containers) because the expression evaluated to '-action exec'."

find docker ! -false

echo "! -false == -true == true. Thus, the above invocation was equivalent to 'find docker', so it would have enumerated every entry in the Docker plugin if we let it run indefinitely."

echo "You can also use '()' to group terms together."

find docker \( -action exec -o -true \) -a -false

echo "( -action exec -o -true ) -a -false == -action exec -a -false == false. Thus, the above invocation printed out nothing because the expression evaluated to false."

echo "That wraps up our overview of find's expression syntax. Remember that you can type 'find --help syntax' to access this same information."

find --help syntax

echo "Now let's show off a few more primaries."

find --help

find docker/containers -maxdepth 1 -name '*er'

echo "This printed out all Docker containers whose name ends with an 'er'. Note the use of the maxdepth option to avoid indefinite recursion."

find docker/containers -maxdepth 1 -crtime -24h

echo "This printed out all Docker containers that were created within the last 24 hours."

echo "Notice that the last two invocations filtered on a specific kind of entries (Docker containers). We can write them in a more expressive way via the kind primary."

find docker -kind '*container' -name '*er'

find docker -kind '*container' -crtime -24h

echo "Pretty neat huh? Let's go back to our list of primaries."

find --help

echo "Besides crtime, we see some other primaries that let us filter on entry attributes. For example,"

find docker/containers/*/fs/var/log -size +1k

echo "This printed out every log in every container's /var/log directory whose file size exceeded 1 KB."

echo "The meta primary lets you filter on an entry's metadata property. It is useful when you need to filter on properties that aren't a part of Wash's entry attributes. These could be vendor-specific properties like VPC ID for AWS EC2 instances or an image ID for Docker containers. They could also be other 'common' properties that we haven't yet added to the current set of entry attributes. The latter would include things like tags/labels, and state."

echo "The meta primary is its own DSL, similar in flavor to jq and find's expression syntax. You can type 'find --help meta' to view its full documentation."

find --help meta | less

echo "As you can see, the meta primary's pretty powerful (and we'll likely add more features to it). Thus, we'll illustrate its usage via some examples."

echo "For the first example, we're going to find all running EC2 instances. Here, 'state' is the metadata property that we want to filter on. Since an entry's metadata is an arbitrary JSON object/YAML mapping, we'll need to query a representative EC2 instance's metadata to figure out the 'state' property's structure. Since the meta attribute's fetched in a batch request, we'll query that first to optimize find's metadata filtering."

meta --attribute aws/<profile>/resources/ec2/instances/<instance>

echo "From the output, we see that the 'state' property is a YAML mapping. The state's name is specified in the 'Name' key. From https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html, we see that the desired value for 'Name' is 'running'. Thus, our meta primary query reduces to m['state']['name'] == 'running' (where m is the entry's metadata). Note that the meta primary will properly upcase 'state' => 'State' and 'name' => 'Name' so we don't have to worry about correctly casing the keys. Anyways, here's what the find invocation looks like:"

find aws -k '*ec2*instance' -meta '.state.name' 'running'

echo "For the second example, we're going to find all EC2 instances with the 'owner' tag. Here, 'tags' is the metadata property we want to filter on. Let's check our representative entry's metadata to see what this looks like."

meta --attribute aws/<profile>/resources/ec2/instances/<instance>

echo "From the output, we see that the 'Tags' property is an array of YAML mappings. Each mapping has keys 'Key' and value 'Value'. 'Key' specifies the tag's name; 'Value' specifies the tag's value. Thus, 'return true if the EC2 instance has the owner tag' reduces to 'return true if the tags array contains some tag t such that t['key'] == owner'."

find aws -k '*ec2*instance' -meta '.tags[?]' '.key' 'owner'

echo "Note that we don't care about the tag's value. All we care about is whether the tag exists. If we did care about the tag's value, then the meta primary query would be"

find aws -k '*ec2*instance' -meta '.tags[?]' '.key' 'owner' '.value' 'value' 

echo "which translates to 'return true if the tags array contains some tag t such that t['key'] == owner AND t['value'] == value'."

echo "For the third example, we're going to find all Docker containers that started within the last 24 hours. Since there is no 'start time' entry attribute, this property must come from the entry's metadata. So, let's check out a representative entry's metadata."

meta --attribute docker/containers/<container>

echo "From the output, we see that the meta attribute does not contain our desired property. Thus, we resort to checking the entry's full metadata."

meta docker/containers/<container>

echo "Here, we see that our desired property is m['state']['startedAt']. Thus, our meta primary query reduces to 'm['state']['startedAt'] < 24 hours from the current date'." 

find docker -fullmeta -k '*container' -meta '.state.startedAt' -24h

echo "Note that the fullmeta option tells find to construct the meta primary query on the entry's full metadata instead of its meta attribute. You should be careful when using this option because fullmeta will make O(N) requests, where N is the number of visited entries. This may take a while when N is large. It may also cost some money if the plugin's API is subscription based."

echo "Anyways, these three examples should illustrate the general flow of constructing a meta primary query. This flow consists of the following steps."

echo "(1) Check a representative entry's metadata to get a good idea of how the property's value is structured (e.g. is it a YAML mapping, an array, a primitive value?). Start by checking the meta attribute via 'meta --attribute'. Only check the full metadata via 'meta' if the meta attribute does not contain the desired property. If both the meta attribute and the full metadata do not contain the desired property, then we recommend you contact the plugin author to see if they could include that property in the meta attribute (recommended) or the entry's full metadata."

echo "(2) Once you've figured out the value's structure, you can refer to the meta primary's documentation and examples to construct the right predicate."

echo "That concludes our demo of find. We'll finish this screencast by showing the find invocations that correspond to the three use-cases outlined in the screencast's introduction."

echo "(1) Enumerate all Amazon S3/GCP Storage buckets."

find -k 'aws/*s3*bucket' -o -k 'gcp/*storage*bucket'

echo "(2) List all Docker containers/Kubernetes pods that were created within the last 24 hours."

find \( -k 'docker/*container' -o -k 'kubernetes/*pod' \) -crtime -24h

echo "(3) Extract all running AWS EC2 instances/GCP compute instances that contain the 'owner' tag/label."

find \( -k 'aws/*ec2*instance' -m '.state.name' 'running' -m '.tags[?]' '.key' 'owner' \) -o \( -k 'gcp/*compute*instance' -m '.status' 'RUNNING' -m '.labels.owner' -exists \)

echo "Note that it's probably better to split (3) into two separate find commands for readability. This would look something like"

find -k 'aws/*ec2*instance' -m '.state.name' 'running' -m '.tags[?]' '.key' 'owner'

find -k 'gcp/*compute*instance' -m '.status' 'RUNNING' -m '.labels.owner' -exists 

exit
