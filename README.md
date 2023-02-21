# Devops CLI: A tool to do DevOps in style

A Terminal/Web UI to interact with devops tools & services, intended for devops/platform engineers.

## Inspiration
This project is inspired from the following softwares terraform, K9s, Lens. 
- Devops CLI aims to bring the agility (that extensively comes from using keyboard) & speed offered by K9s.
- Devops CLI aims to bring the extensibility provided by terraform (in form of plugins)
- Not everything can be done on TUI, Something are better show in a web app like Lens. The above 2 things will also be package in a Web UI for better experience.

## This project aims to achieve the following:
- Improve debugging
- Improve development agility

## Demo

https://user-images.githubusercontent.com/24411676/210130161-4c467179-9bff-45fd-97e8-0207bb0506c7.mp4


## Installation
**Linux & Mac**

`curl https://storage.googleapis.com/devops-cli-artifacts/releases/devops/0.3.0/install.sh | bash`

## Usage

**Run Server**

`devops`

**Run Client (TUI)

`devops tui`

## Read [Wiki](https://github.com/sharadregoti/devops-cli/wiki) for detailed documentation 

## Supported Plugins
- Kubernetes
- Helm (WIP)

### Kubernetes Features
- View & search all cluster resources
- Create, Read, Update, Delete any resource
- Describe any resource
- View logs, get shell access of pod
