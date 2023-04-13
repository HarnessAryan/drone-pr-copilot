A plugin to Drone PR copilot.

# Usage

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
name: default

steps:
- name: run plugins/drone-pr-copilot plugin
  image: plugins/drone-pr-copilot
```

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t plugins/drone-pr-copilot -f docker/Dockerfile .
```

# Testing

Execute the plugin from your current working directory:

```text
docker run --rm \
  -e DRONE_REPO_NAMESPACE=<namespace> \
  -e DRONE_REPO_NAME=<repo> \
  -e PLUGIN_GITHUB_TOKEN=*** \
  -e PLUGIN_OPENAI_KEY=*** \
  -e DRONE_COMMIT_LINK=https://github.com/<namespace>/repo/pull/<PR_NUMBER> \
  -w /drone/src \
  -v $(pwd):/drone/src \
  d1wilko/drone-pr-copilot
```
