# Examples of interacting with the wide-area shell (wash)

## Principals
Multiple ways to get data, but consistent language within the tool. i.e. may search for a database by saying type is 'db' or 'database', but the tool will always refer to them by 'database'.

## Examples

```
> ls /					            # Shows configured API types
gcp/
aws/
platform9/
kubernetes/
> cd aws/
> ls					            # Multiple ways of accessing the same resources
groups/
regions/
resources/
> tree groups/			            # Specific to AWS. Likely role or project oriented.
groups
├── developer
│   ├── ec2
│   │   ├── vm-106.puppet.com
│   │   └── vm-107.puppet.com
│   ├── lambda
│   │   └── michael-lambda-17
│   └── s3
│       └── michael-bucket1
└── dujour
    ├── ec2
    │   ├── vm-106.puppet.com
    │   └── vm-107.puppet.com
    ├── lambda
    │   └── michael-lambda-17
    └── s3
        └── michael-bucket1
> cd resources/
> ls					            # Only show types of things where you have resources
ec2
lambda
s3
> cd ec2/
> ls -l                             # Show some details about ownership and categorization
Name                Creator         Groups              Created         Tags
vm-106.puppet.com   michael.smith   developer,dujour    Dec 29 10:41    prod,web
> find / -type compute,storage      # Shows the shortest path to each resource
aws/resources/ec2/vm-106.puppet.com
aws/resources/ec2/vm-107.puppet.com
aws/resources/s3/michael-bucket1
> ls /kubernetes/
namespaces/
> tail -f /aws/resources/ec2/vm-106.puppet.com /aws/resources/ec2/vm-107.puppet.com:/var/log/nginx/access.log /aws/resources/lambda/michael-lambda-17 /aws/resources/s3/michael-bucket1 /kubernetes/gke_shared-k8s_us-west1-a_shared-k8s-dev/dujour-dev/pods/r0raxmg1fg276o05wmmqancki8w-dujour-84c7b497cc-fd7m4
==> /aws/resources/ec2/vm-107.puppet.com:/var/log/syslog <==
Jan  2 23:53:50 pe-master systemd[1]: Starting User Manager for UID 1000...
Jan  2 23:53:50 pe-master systemd[1]: Started Session 25386 of user ubuntu.

==> /aws/resources/ec2/vm-107.puppet.com:/var/log/nginx/access.log <==
10.0.25.192 - - [20/Dec/2018:18:03:58 +0000] "GET /index.html HTTP/1.1" 200 603 "https://vm-107.puppet.com/" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.1 Safari/605.1.15" "-"
10.0.25.192 - - [20/Dec/2018:18:04:59 +0000] "GET /index.html HTTP/1.1" 200 605 "https://vm-107.puppet.com/" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.1 Safari/605.1.15" "-"

==> /aws/resources/lambda/michael-lambda-17 <==
=== puma startup: 2019-01-02 06:25:01 +0000 ===

==> /aws/resources/s3/michael-bucket1
79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be mybucket [06/Feb/2014:00:00:38 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be 3E57427F3EXAMPLE REST.GET.VERSIONING - "GET /mybucket?versioning HTTP/1.1" 200 - 113 - 7 - "-" "S3Console/0.4" -
79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be mybucket [06/Feb/2014:00:00:38 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be 891CE47D2EXAMPLE REST.GET.LOGGING_STATUS - "GET /mybucket?logging HTTP/1.1" 200 - 242 - 11 - "-" "S3Console/0.4" -
79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be mybucket [06/Feb/2014:00:00:38 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be A1206F460EXAMPLE REST.GET.BUCKETPOLICY - "GET /mybucket?policy HTTP/1.1" 404 NoSuchBucketPolicy 297 - 38 - "-" "S3Console/0.4" -
79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be mybucket [06/Feb/2014:00:01:00 +0000] 192.0.2.3 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be 7B4A0FABBEXAMPLE REST.GET.VERSIONING - "GET /mybucket?versioning HTTP/1.1" 200 - 113 - 33 - "-" "S3Console/0.4" -

==> /kubernetes/gke_shared-k8s_us-west1-a_shared-k8s-dev/dujour-dev/pods/r0raxmg1fg276o05wmmqancki8w-dujour-84c7b497cc-fd7m4 <==
2018-12-06 19:23:53,215 INFO  [o.e.j.u.log] Logging initialized @49306ms to org.eclipse.jetty.util.log.Slf4jLog
2018-12-06 19:23:54,335 INFO  [p.t.s.w.jetty9-core] Removing buggy security provider SunPKCS11 version 12
2018-12-06 19:23:58,485 INFO  [p.t.s.w.jetty9-service] Initializing web server(s).
2018-12-06 19:23:58,525 INFO  [p.t.s.s.scheduler-service] Initializing Scheduler Service
2018-12-06 19:23:58,635 INFO  [o.q.i.StdSchedulerFactory] Using default implementation for ThreadExecutor

```

### Ways of slicing things
- Namespace-oriented? Are namespaces universal? Resource group (AWS), project (GCP), namespace (K8s). Azure has namespaces within resource groups.
- Region. It's not really clear what common semantics exist for this, maybe we should revisit it later.
- Users (or subscription id).
- Cloud API.
- Multiple hierarchical views? Hardlinks? Symlinks and cycles?
- What are the types? Compute, storage/volume, database/db. Need a consistent taxonomy, lots of different naming patterns across APIs.

### Real-world Examples
- Dujour deployment: pod, deployment, chart, pubsub, dataflow, bigquery
- GKE: the k8s infra itself
- GCP: build servers for Pipelines
- Pipelines SaaS: AWS VMs

