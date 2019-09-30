---
title: Wash
---
Have you ever had to:

<details>
<summary>List all your AWS EC2 instances or Kubernetes pods?</summary>
<pre>aws ec2 describe-instances --profile foo --query 'Reservations[].Instances[].InstanceId' --output text</pre>
<pre>kubectl get pods --all-namespaces`</pre>
</details>
<details>
<summary>Read/cat a GCP Compute instance's console output, or an AWS S3 object's content?</summary>
<pre>gcloud compute instances get-serial-port-output foo</pre>
<pre>aws s3api get-object content.txt --profile foo --bucket bar --key baz && cat content.txt && rm content.txt</pre>
</details>
<details>
<summary>Exec a command on a Kubernetes pod or GCP Compute Instance?</summary>
<pre>kubectl exec foo uname</pre>
<pre>gcloud compute ssh foo --command uname</pre>
</details>
<details>
<summary>Find all AWS EC2 instances with a particular tag, or Docker containers/Kubernetes pods/GCP Compute instances with a specific label?</summary>
<pre>aws ec2 describe-instances --profile foo --query 'Reservations[].Instances[].InstanceId' --filters Name=tag-key,Values=owner --output text</pre>
<pre>docker ps --filter “label=owner”</pre>
</details>

Does it bother you that each of those is a bespoke, cryptic incantation of various vendor-specific tools? It's a lot of commands you have to use, applications you need to install, and DSLs you have to learn just to do some pretty basic tasks. In Wash, these basic tasks are simple. 

* Listing stuff is as easy as `ls` or `find`
* Reading stuff is as easy as `cat`'ing a file
* Execing a command is as easy as `wexec`
* Finding stuff is as easy as `find`

And this is only scratching the surface of Wash's capabilities. Check out the screencast below

<script id="asciicast-mX8Mwa75rr1bJePLi3OnIOkJK" src="https://asciinema.org/a/mX8Mwa75rr1bJePLi3OnIOkJK.js" async></script>

and the [tutorials](tutorials) to learn more.
