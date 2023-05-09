# cfparams

Wrangle CloudFormation parameters.

Reads expected parameter keys from a CloudFormation template, accepts values
from command line or YAML file, emits JSON suitable for `aws cloudformation
*-stack --parameters="..."` commands. Supports the use of default and previous
values.

## Installation Instructions

* If you use [Homebrew](https://brew.sh/), you can simply run: `brew install cultureamp/tap/cfparams` (only in the Culture Amp organization at the moment, sorry).
* Alternatively, if you want the latest version, download the executable directly from the [release page](https://github.com/cultureamp/cfparams/releases)
* Rename and move the `cfparams` executable to your path, e.g: `mv cfparams /usr/local/bin` (or move to your favourite spot and symlink).
* To check whether your installation is successful, run the command `cfparams` and the output would be similar to:
    ```shellsession
    $ cfparams
    CloudFormation template or tags file required
        e.g: --template=cloudformation.yaml
        e.g: --tags=tags-foo.yaml
    ```
* If it still doesn't work you may need to `hash -r` to clear the path cache or restart the shell

## Example use-cases

CloudFormation template excerpt describing an ECS service to be provisioned
onto an existing ECS cluster.

```yaml
# cfn.yaml excerpt
Parameters:
  Greeting:
    Type: String
    Description: greeting message to send
    Default: Hello
  Recipient:
    Description: name of the greeting recipient
    Type: String
  ImageRepo:
    Type: String
    Description: repository of Docker image to run
    Default: "123.dkr.ecr.us-east-1.amazonaws.com/greeting"
  ImageTag:
    Type: String
    Description: tag of Docker image to run
    Default: latest
  Cluster:
    Description: ECS cluster ID to run service on
    Type: String
```

### Creating ECS service CloudFormation stack

Launching the CloudFormation stack for the first time.  Accept some defaults
from the template, specify all other parameters.

```sh
params="$(
  cfparams --template=cfn.yaml --accept-defaults --no-previous \
    Recipient=world ImageTag=v1 Cluster=nanoservices
)"
```

* `--template` loads supported Parameters from a CloudFormation template.
* `--accept-defaults` omits keys that have a default in the CloudFormation
  template.
* `--no-previous` means fail if a key has no default in the template and isn't
  specified on the command line. Without this option, those keys will be
  auto-filled as `"UsePreviousValue": true`.

Resulting JSON:

```json
[
  {"ParameterKey": "Recipient", "ParameterValue": "world"},
  {"ParameterKey": "ImageTag", "ParameterValue": "v1"},
  {"ParameterKey": "Cluster", "ParameterValue": "nanoservices"}
]
```

```sh
aws cloudformation create-stack \
  --stack-name=greeting \
  --template-body=file://cfn.yaml \
  --parameters="$params"
```


### Deploying a new version of the app

Deploying a new version of the app, e.g. from CI. Only `ImageTag` should
change, all other parameters use previous value.

```sh
params="$(cfparams --template=cfn.yaml ImageTag=v2)"
```

Resulting JSON:

```json
[
  {"ParameterKey": "Greeting", "UsePreviousValue": true},
  {"ParameterKey": "Recipient", "UsePreviousValue": true},
  {"ParameterKey": "ImageRepo", "UsePreviousValue": true},
  {"ParameterKey": "ImageTag", "ParameterValue": "v2"},
  {"ParameterKey": "Cluster", "UsePreviousValue": true}
]
```

Update stack:

```sh
aws cloudformation update-stack \
  --stack-name=greeting \
  --use-previous-template \
  --parameters="$params"
```


### Updating the CloudFormation stack

Changing the stack, for example introducing a `FooHost` parameter.

```diff
 # cfn.yaml excerpt
 Parameters:
+  FooHost:
+    Type: String
+    Description: API key to access Foo service
   Greeting:
```

```sh
params="$(cfparams --template=cfn-foohost.yaml FooHost=foo.example.com)"
```

Resulting JSON:

```json
[
  {"ParameterKey": "FooHost", "ParameterValue": "foo.example.com"},
  {"ParameterKey": "Greeting", "UsePreviousValue": true},
  {"ParameterKey": "Recipient", "UsePreviousValue": true},
  {"ParameterKey": "ImageRepo", "UsePreviousValue": true},
  {"ParameterKey": "ImageTag", "UsePreviousValue": true},
  {"ParameterKey": "Cluster", "UsePreviousValue": true}
]
```

```sh
name="greeting-update-$(date +%Y%m%d-%H%M%S)"

aws cloudformation create-change-set \
  --stack-name=greeting \
  --change-set-name="$name" \
  --use-previous-template \
  --parameters="$(cfparams --template=cfn.yaml FooHost=foo.example.com)"

# review Change Set here

aws cloudformation execute-change-set \
  --stack-name=greeting \
  --change-set-name="$name"
```

### Introducing version-controlled parameters files

Now we introduce some version-controlled files to the subset of parameters that
make sense to exist in the codebase. `ImageTag` is not included in this file.

```yaml
# parameters-staging.yaml
FooHost: foo.local
Greeting: Howdy
Recipient: team
Cluster: staging
```

```yaml
# parameters-production.yaml
FooHost: foo.example.com
Greeting: Hello
Recipient: world
Cluster: production
```

```sh
params="$(
  cfparams --template=cfn-foohost.yaml --parameters=parameters-staging.yaml \
    ImageTag=v3 Greeting=Bonjour
)
```

Resulting JSON:

```json
[
  {"ParameterKey": "FooHost", "ParameterValue": "foo.local"},
  {"ParameterKey": "Greeting", "ParameterValue": "Howdy"},
  {"ParameterKey": "Recipient", "ParameterValue": "team"},
  {"ParameterKey": "ImageRepo", "UsePreviousValue": true},
  {"ParameterKey": "ImageTag", "ParameterValue": "v3"},
  {"ParameterKey": "Cluster", "ParameterValue": "staging"}
]
```

### Bonus feature: Stack Tags

CloudFormation stacks can be tagged, and those tags flow into all taggable
resources the stack creates. As with `--parameters`, the `aws cloudformation`
commands expect these in an awkward format. `cfparams --tags file.yaml` helps.

```sh
aws cloudformation create-stack \
  ... \
  --tags "$(cfparams --tags=tags-production.yaml)" \
  ...
```

```yaml
# tags-production.yaml
Name: Widgets as a Service
asset: widget-api
workload: production
```

```sh
cfparams --tags=tags-production.yaml
```

```json
[
  {
    "Key": "Name",
    "Value": "Widgets as a Service"
  },
  {
    "Key": "asset",
    "Value": "widget-api"
  },
  {
    "Key": "workload",
    "Value": "production"
  }
]
```

Tags can also be overridden and added to on the command line, just like parameters can:

```yaml
#tags-test.yaml
asset: buildkite
workload: delivery
buildkite-group: build-test
```

```sh
> cfparams --tags=tags-test.yaml buildkite-group=build-pylons message=extra-tag
```

```json
[
  {
    "Key": "buildkite-group",
    "Value": "build-pylons"
  },
  {
    "Key": "message",
    "Value": "extra-tag"
  },
  {
    "Key": "asset",
    "Value": "buildkite"
  },
  {
    "Key": "workload",
    "Value": "delivery"
  }
]
```
