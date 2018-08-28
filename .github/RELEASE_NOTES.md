# Release Notes
## v0.2.0
- fixed a memory leak bug: https://github.com/komuw/meli/pull/114
- upgraded docker-client from v17.03.2-ce to v18.06.1-ce: https://github.com/komuw/meli/pull/112
- added Go module support: https://github.com/komuw/meli/pull/106
- upgraded docker client from verion v1.13.1 to version v17.03.2-ce : https://github.com/komuw/meli/pull/109
- removed vendor directory and dep files: https://github.com/komuw/meli/pull/107
- we no longer release meli for 386 arch on github releases.  
  We now ONLY release amd64 for darwin, linux and windows


## v0.1.9.8
- added support for dot env(.env) files[1] in the docker-compose file: https://github.com/komuw/meli/pull/102        

ref:          
1. https://docs.docker.com/compose/compose-file/#env_file