# IN WASH SHELL

echo "This screencast talks about some of the more interesting things you could do with Wash (hence the title 'Advanced Wash')."

echo ""

#This screenWe'll see that Wash makes it easy to (1) count the number of things in your infrastructure, (2) view the running processes on an AWS EC2 instance/Docker container/GCP Compute instance/Kubernetes pod, and (3) Tail some log files/execute a command on all AWS EC2 instances/Docker containers/GCP compute instances/Kubernetes pods with a specified label."

echo "We'll start with Wash ps. "

wps --help

echo "Let's try it out."

wps gcp/<project>/compute/<instance>

wps aws/<profile>/ec2/instances/<instance>

wps docker/containers/<container>

wps gcp/<project>/compute/<instance> aws/<profile>/ec2/instances/<instance>

echo "Nice! Unfortunately, wps is experimental so its features are a bit limited (
