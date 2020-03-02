agent-dry:
	helm lint sts-agent
	helm install --dry-run --debug sts-agent ./sts-agent

agent-install:
	helm install --debug sts-agent ./sts-agent

traceagent-dry:
	helm lint sts-trace-agent
	helm install --dry-run --debug sts-trace-agent ./sts-trace-agent

traceagent-install:
	helm install --debug sts-trace-agent ./sts-trace-agent

agent-uninstall:
	helm uninstall sts-agent

sts-simulator-dry:
	helm lint sts-simulator
	helm install --dry-run --debug sts-simulator ./sts-simulator

sts-simulator-install:
	helm install --debug sts-simulator ./sts-simulator

sts-simulator-manifest:
	helm get manifest sts-simulator > sts-simulator.manifest.yml
	helm get manifest sts-simulator

sts-simulator-uninstall:
	helm uninstall sts-simulator

trace-java-demo-dry:
	helm install --dry-run --debug app ./trace-java-demo

trace-java-demo-install:
	helm install --debug app ./trace-java-demo
