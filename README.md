# Go Binary Generic Repository

A template/Fork-able repository to get going quickly with a go binary. 

### Getting Started

Requires Go...
 - clone or template, then run
 - update the [Makefile](./Makefile) with org, repo, binary name, and git server if applicable, this will update the build args that determine where the latest version is fetched.
 - `make build` to generate the binary in your build/ directory.
 - `./build/<binary> version` should print a verision  
 - manually make a 0.0.0 release. 
 - Make a PR with any change and add a label
    - `Release Patch`
    - `Release Minor`
    - `Release Major`
    
   merge the Pr in and github actions (assuming enabled) should build a new binary with the right version, publish it to a new release.
   
### Features:

 - Version command that checks current version and compares to repository of source codes releases
 - Versioned Releases using semantic versions, powered by github actions and github packages for docker images.
 
 ---
 
### Credits

 - #### [Jenkins-X](https://github.com/jenkins-x/jx/)
 - #### [Github-Jira-Bot](https://github.com/Benbentwo/github-jira-bot/)
